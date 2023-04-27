// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/jkawamoto/sd-image-viewer/server/models"
)

// GetCheckpointsOKCode is the HTTP code returned for type GetCheckpointsOK
const GetCheckpointsOKCode int = 200

/*
GetCheckpointsOK A list of checkpoint names.

swagger:response getCheckpointsOK
*/
type GetCheckpointsOK struct {

	/*
	  In: Body
	*/
	Payload []string `json:"body,omitempty"`
}

// NewGetCheckpointsOK creates GetCheckpointsOK with default headers values
func NewGetCheckpointsOK() *GetCheckpointsOK {

	return &GetCheckpointsOK{}
}

// WithPayload adds the payload to the get checkpoints o k response
func (o *GetCheckpointsOK) WithPayload(payload []string) *GetCheckpointsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get checkpoints o k response
func (o *GetCheckpointsOK) SetPayload(payload []string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCheckpointsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]string, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*
GetCheckpointsDefault Error Response

swagger:response getCheckpointsDefault
*/
type GetCheckpointsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.StandardError `json:"body,omitempty"`
}

// NewGetCheckpointsDefault creates GetCheckpointsDefault with default headers values
func NewGetCheckpointsDefault(code int) *GetCheckpointsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetCheckpointsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get checkpoints default response
func (o *GetCheckpointsDefault) WithStatusCode(code int) *GetCheckpointsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get checkpoints default response
func (o *GetCheckpointsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get checkpoints default response
func (o *GetCheckpointsDefault) WithPayload(payload *models.StandardError) *GetCheckpointsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get checkpoints default response
func (o *GetCheckpointsDefault) SetPayload(payload *models.StandardError) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCheckpointsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}