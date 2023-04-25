// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/jkawamoto/sd-image-viewer/server/restapi/operations"
)

//go:generate swagger generate server --target ../../server --name SdImageViewer --spec ../../openapi.yaml --principal interface{}

func configureFlags(api *operations.SdImageViewerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SdImageViewerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.GetImageHandler == nil {
		api.GetImageHandler = operations.GetImageHandlerFunc(func(params operations.GetImageParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetImage has not yet been implemented")
		})
	}
	if api.GetImagesHandler == nil {
		api.GetImagesHandler = operations.GetImagesHandlerFunc(func(params operations.GetImagesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetImages has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}
