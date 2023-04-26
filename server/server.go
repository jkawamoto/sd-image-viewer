package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/jkawamoto/sd-image-viewer/image"
	"github.com/jkawamoto/sd-image-viewer/server/models"
	"github.com/jkawamoto/sd-image-viewer/server/restapi"
	"github.com/jkawamoto/sd-image-viewer/server/restapi/operations"
)

const defaultLimit = 20

var gmt = time.FixedZone("GMT", 0)

func NewServer(index bleve.Index, pathPrefix string, logger *log.Logger) (*restapi.Server, error) {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return nil, err
	}

	api := operations.NewSdImageViewerAPI(swaggerSpec)
	api.GetImageHandler = GetImageHandler(pathPrefix, logger)
	api.GetImagesHandler = GetImagesHandler(index, pathPrefix, logger)
	api.Logger = logger.Printf

	server := restapi.NewServer(api)
	server.Port = 8080
	server.KeepAlive = 3 * time.Minute
	server.ReadTimeout = 30 * time.Second
	server.WriteTimeout = 60 * time.Second
	server.ConfigureAPI()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// TODO: returns frontend code.
	}))
	mux.Handle("/api/v1/", server.GetHandler())

	server.SetHandler(withLogger(mux, logger))

	return server, nil
}

func GetImageHandler(pathPrefix string, logger *log.Logger) operations.GetImageHandlerFunc {
	return func(params operations.GetImageParams) middleware.Responder {
		name := filepath.Join(pathPrefix, params.ID)

		info, err := os.Stat(name)
		if os.IsNotExist(err) {
			logger.Printf("requested file doesn't exist: %v", err)
			return operations.NewGetImageDefault(http.StatusNotFound).WithPayload(&models.StandardError{
				Message: swag.String(err.Error()),
			})
		} else if err != nil {
			logger.Printf("failed to stat the requested file: %v", err)
			return operations.NewGetImageDefault(http.StatusInternalServerError).WithPayload(&models.StandardError{
				Message: swag.String(err.Error()),
			})
		}

		if since := swag.StringValue(params.IfModifiedSince); since != "" {
			t, err := time.Parse(time.RFC1123, since)
			if err != nil {
				logger.Printf("failed to parse the given If-Modified-Since header value: %v", err)
			} else if !info.ModTime().After(t) {
				return operations.NewGetImageNotModified()
			}
		}

		f, err := os.Open(name)
		if err != nil {
			logger.Printf("failed to open the requested file: %v", err)
			return operations.NewGetImagesDefault(http.StatusInternalServerError).WithPayload(&models.StandardError{
				Message: swag.String(err.Error()),
			})
		}

		return operations.NewGetImageOK().
			WithPayload(f).
			WithCacheControl("max-age=3600").
			WithLastModified(info.ModTime().In(gmt).Format(time.RFC1123))
	}
}

func GetImagesHandler(index bleve.Index, pathPrefix string, logger *log.Logger) operations.GetImagesHandlerFunc {
	query.SetLog(logger)
	return func(params operations.GetImagesParams) middleware.Responder {
		var queries []query.Query
		if params.Query != nil {
			q := query.NewMatchPhraseQuery(swag.StringValue(params.Query))
			q.FieldVal = "prompt"
			q.Analyzer = image.Analyzer(q.FieldVal)

			queries = append(queries, q)
		}
		if params.Size != nil {
			var min, max *float64
			switch swag.StringValue(params.Size) {
			case "small":
				max = swag.Float64(512 * 768)
			case "medium":
				min = swag.Float64(512 * 768)
				max = swag.Float64(512 * 768 * 4)
			case "large":
				min = swag.Float64(512 * 768 * 4)
			}
			q := query.NewNumericRangeInclusiveQuery(min, max, swag.Bool(false), swag.Bool(true))
			q.FieldVal = "pixel"
			queries = append(queries, q)
		}
		if params.After != nil || params.Before != nil {
			var before, after time.Time
			if params.Before != nil {
				before = time.Time(*params.Before)
			}
			if params.After != nil {
				after = time.Time(*params.After)
			}
			q := query.NewDateRangeInclusiveQuery(after, before, swag.Bool(true), swag.Bool(false))
			q.FieldVal = "creation-time"

			queries = append(queries, q)
		}
		if len(queries) == 0 {
			queries = append(queries, query.NewMatchAllQuery())
		}

		page := int(swag.Int64Value(params.Page))
		limit := defaultLimit
		if params.Limit != nil {
			limit = int(swag.Int64Value(params.Limit))
		}

		req := bleve.NewSearchRequestOptions(query.NewConjunctionQuery(queries), limit, limit*page, false)
		req.Fields = []string{"*"}
		if swag.StringValue(params.Order) == "asc" {
			req.SortBy([]string{"creation-time"})
		} else {
			req.SortBy([]string{"-creation-time"})
		}

		res, err := index.Search(req)
		if err != nil {
			logger.Printf("failed to search images: %v", err)
			return operations.NewGetImagesDefault(http.StatusInternalServerError).WithPayload(&models.StandardError{
				Message: swag.String(err.Error()),
			})
		}

		items := make([]*models.Image, len(res.Hits))
		for i, v := range res.Hits {
			id, err := filepath.Rel(pathPrefix, v.ID)
			if err != nil {
				logger.Printf("failed to get a relative path: %v", err)
				return operations.NewGetImagesDefault(http.StatusInternalServerError).WithPayload(&models.StandardError{
					Message: swag.String(err.Error()),
				})
			}

			items[i] = &models.Image{
				ID:                        swag.String(id),
				Prompt:                    getString(v.Fields, "prompt"),
				NegativePrompt:            getString(v.Fields, "negative-prompt"),
				Checkpoint:                getString(v.Fields, "checkpoint"),
				CreationTime:              strfmt.DateTime(getDateTime(v.Fields, "creation-time")),
				Pixel:                     int64(getInt(v.Fields, "pixel")),
				ImageAdditionalProperties: getMap(v.Fields, "metadata"),
			}
		}

		return operations.NewGetImagesOK().WithPayload(&models.ImageList{
			Items: items,
			Metadata: &models.Metadata{
				CurrentPage: swag.Int64(int64(page)),
				TotalItems:  swag.Int64(int64(res.Total)),
				TotalPages:  swag.Int64(int64(res.Total/uint64(limit)) + 1),
			},
		})
	}
}

func getString(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return v
}

func getInt(m map[string]any, key string) int {
	v, _ := m[key].(int)
	return v
}

func getDateTime(m map[string]any, key string) time.Time {
	v, _ := m[key].(string)
	t, _ := time.Parse(time.RFC3339, v)
	return t
}

func getMap(m map[string]any, key string) map[string]any {
	res := make(map[string]any)
	for k, v := range m {
		if strings.HasPrefix(k, key) {
			k = strings.TrimPrefix(k, key+".")
			res[k] = v
		}
	}
	return res
}

type responseWriter struct {
	http.ResponseWriter
	Code int
}

func (w *responseWriter) WriteHeader(code int) {
	w.Code = code
	w.ResponseWriter.WriteHeader(code)
}

func withLogger(h http.Handler, logger *log.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		w := &responseWriter{
			ResponseWriter: res,
		}
		h.ServeHTTP(w, req)
		logger.Printf("%v %v %v", w.Code, req.Method, req.URL)
	}
}
