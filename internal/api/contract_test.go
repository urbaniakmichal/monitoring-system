package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/stretchr/testify/assert"
)

func AssertResponseMatchesOpenAPI(t *testing.T, req *http.Request, resRec *httptest.ResponseRecorder) {
	// copy body response to read it more than once (trick)
	bodyBytes := resRec.Body.Bytes()

	// Load Swagger
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("../../docs/swagger.json")
	assert.NoError(t, err)

	// Create router and assert path
	router, _ := gorillamux.NewRouter(doc)
	route, pathParams, err := router.FindRoute(req)
	assert.NoError(t, err, "Endpoint is not exist in Swagger")

	// Validate response
	input := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request:    req,
			PathParams: pathParams,
			Route:      route,
		},
		Status: resRec.Code,
		Header: resRec.Header(),
		Body:   io.NopCloser(bytes.NewBuffer(bodyBytes)),
	}

	err = openapi3filter.ValidateResponse(context.Background(), input)
	assert.NoError(t, err, "Response is incorrect according to OpenAPI contract")
}
