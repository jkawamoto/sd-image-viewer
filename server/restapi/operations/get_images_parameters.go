// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetImagesParams creates a new GetImagesParams object
// with the default values initialized.
func NewGetImagesParams() GetImagesParams {

	var (
		// initialize parameters with default values

		orderDefault = string("desc")
	)

	return GetImagesParams{
		Order: &orderDefault,
	}
}

// GetImagesParams contains all the bound params for the get images operation
// typically these are obtained from a http.Request
//
// swagger:parameters getImages
type GetImagesParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Retrieving images created after the given date time.
	  In: query
	*/
	After *strfmt.DateTime
	/*Retrieving images created before the given date time.
	  In: query
	*/
	Before *strfmt.DateTime
	/*Retrieving images that use the given checkpoint.
	  In: query
	*/
	Checkpoint *string
	/*The number of items one page has at most.
	  In: query
	*/
	Limit *int64
	/*
	  In: query
	  Default: "desc"
	*/
	Order *string
	/*Requesting page number.
	  In: query
	*/
	Page *int64
	/*Search query.
	  In: query
	*/
	Query *string
	/*Retrieving the given sized images.
	  In: query
	*/
	Size *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetImagesParams() beforehand.
func (o *GetImagesParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qAfter, qhkAfter, _ := qs.GetOK("after")
	if err := o.bindAfter(qAfter, qhkAfter, route.Formats); err != nil {
		res = append(res, err)
	}

	qBefore, qhkBefore, _ := qs.GetOK("before")
	if err := o.bindBefore(qBefore, qhkBefore, route.Formats); err != nil {
		res = append(res, err)
	}

	qCheckpoint, qhkCheckpoint, _ := qs.GetOK("checkpoint")
	if err := o.bindCheckpoint(qCheckpoint, qhkCheckpoint, route.Formats); err != nil {
		res = append(res, err)
	}

	qLimit, qhkLimit, _ := qs.GetOK("limit")
	if err := o.bindLimit(qLimit, qhkLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	qOrder, qhkOrder, _ := qs.GetOK("order")
	if err := o.bindOrder(qOrder, qhkOrder, route.Formats); err != nil {
		res = append(res, err)
	}

	qPage, qhkPage, _ := qs.GetOK("page")
	if err := o.bindPage(qPage, qhkPage, route.Formats); err != nil {
		res = append(res, err)
	}

	qQuery, qhkQuery, _ := qs.GetOK("query")
	if err := o.bindQuery(qQuery, qhkQuery, route.Formats); err != nil {
		res = append(res, err)
	}

	qSize, qhkSize, _ := qs.GetOK("size")
	if err := o.bindSize(qSize, qhkSize, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAfter binds and validates parameter After from query.
func (o *GetImagesParams) bindAfter(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	// Format: date-time
	value, err := formats.Parse("date-time", raw)
	if err != nil {
		return errors.InvalidType("after", "query", "strfmt.DateTime", raw)
	}
	o.After = (value.(*strfmt.DateTime))

	if err := o.validateAfter(formats); err != nil {
		return err
	}

	return nil
}

// validateAfter carries on validations for parameter After
func (o *GetImagesParams) validateAfter(formats strfmt.Registry) error {

	if err := validate.FormatOf("after", "query", "date-time", o.After.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindBefore binds and validates parameter Before from query.
func (o *GetImagesParams) bindBefore(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	// Format: date-time
	value, err := formats.Parse("date-time", raw)
	if err != nil {
		return errors.InvalidType("before", "query", "strfmt.DateTime", raw)
	}
	o.Before = (value.(*strfmt.DateTime))

	if err := o.validateBefore(formats); err != nil {
		return err
	}

	return nil
}

// validateBefore carries on validations for parameter Before
func (o *GetImagesParams) validateBefore(formats strfmt.Registry) error {

	if err := validate.FormatOf("before", "query", "date-time", o.Before.String(), formats); err != nil {
		return err
	}
	return nil
}

// bindCheckpoint binds and validates parameter Checkpoint from query.
func (o *GetImagesParams) bindCheckpoint(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Checkpoint = &raw

	return nil
}

// bindLimit binds and validates parameter Limit from query.
func (o *GetImagesParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("limit", "query", "int64", raw)
	}
	o.Limit = &value

	return nil
}

// bindOrder binds and validates parameter Order from query.
func (o *GetImagesParams) bindOrder(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetImagesParams()
		return nil
	}
	o.Order = &raw

	if err := o.validateOrder(formats); err != nil {
		return err
	}

	return nil
}

// validateOrder carries on validations for parameter Order
func (o *GetImagesParams) validateOrder(formats strfmt.Registry) error {

	if err := validate.EnumCase("order", "query", *o.Order, []interface{}{"asc", "desc"}, true); err != nil {
		return err
	}

	return nil
}

// bindPage binds and validates parameter Page from query.
func (o *GetImagesParams) bindPage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("page", "query", "int64", raw)
	}
	o.Page = &value

	return nil
}

// bindQuery binds and validates parameter Query from query.
func (o *GetImagesParams) bindQuery(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Query = &raw

	return nil
}

// bindSize binds and validates parameter Size from query.
func (o *GetImagesParams) bindSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Size = &raw

	if err := o.validateSize(formats); err != nil {
		return err
	}

	return nil
}

// validateSize carries on validations for parameter Size
func (o *GetImagesParams) validateSize(formats strfmt.Registry) error {

	if err := validate.EnumCase("size", "query", *o.Size, []interface{}{"small", "medium", "large"}, true); err != nil {
		return err
	}

	return nil
}
