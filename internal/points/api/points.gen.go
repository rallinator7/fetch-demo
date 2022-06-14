// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// BalanceList defines model for BalanceList.
type BalanceList struct {
	AdditionalProperties map[string]int `json:"-"`
}

// Error defines model for Error.
type Error struct {
	// error code
	Code int `json:"code"`

	// error message
	Message string `json:"message"`
}

// GivePoints defines model for GivePoints.
type GivePoints struct {
	// payer name
	Payer string `json:"payer"`

	// point total
	Points int `json:"points"`

	// timestamp
	Timestamp string `json:"timestamp"`
}

// PayerPoints defines model for PayerPoints.
type PayerPoints struct {
	// payer name
	Payer string `json:"payer"`

	// point total
	Points int `json:"points"`
}

// SpendPoints defines model for SpendPoints.
type SpendPoints struct {
	// points spent
	Points int `json:"points"`
}

// GivePointsJSONBody defines parameters for GivePoints.
type GivePointsJSONBody = GivePoints

// SpendPointsJSONBody defines parameters for SpendPoints.
type SpendPointsJSONBody = SpendPoints

// GivePointsJSONRequestBody defines body for GivePoints for application/json ContentType.
type GivePointsJSONRequestBody = GivePointsJSONBody

// SpendPointsJSONRequestBody defines body for SpendPoints for application/json ContentType.
type SpendPointsJSONRequestBody = SpendPointsJSONBody

// Getter for additional properties for BalanceList. Returns the specified
// element and whether it was found
func (a BalanceList) Get(fieldName string) (value int, found bool) {
	if a.AdditionalProperties != nil {
		value, found = a.AdditionalProperties[fieldName]
	}
	return
}

// Setter for additional properties for BalanceList
func (a *BalanceList) Set(fieldName string, value int) {
	if a.AdditionalProperties == nil {
		a.AdditionalProperties = make(map[string]int)
	}
	a.AdditionalProperties[fieldName] = value
}

// Override default JSON handling for BalanceList to handle AdditionalProperties
func (a *BalanceList) UnmarshalJSON(b []byte) error {
	object := make(map[string]json.RawMessage)
	err := json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if len(object) != 0 {
		a.AdditionalProperties = make(map[string]int)
		for fieldName, fieldBuf := range object {
			var fieldVal int
			err := json.Unmarshal(fieldBuf, &fieldVal)
			if err != nil {
				return fmt.Errorf("error unmarshaling field %s: %w", fieldName, err)
			}
			a.AdditionalProperties[fieldName] = fieldVal
		}
	}
	return nil
}

// Override default JSON handling for BalanceList to handle AdditionalProperties
func (a BalanceList) MarshalJSON() ([]byte, error) {
	var err error
	object := make(map[string]json.RawMessage)

	for fieldName, field := range a.AdditionalProperties {
		object[fieldName], err = json.Marshal(field)
		if err != nil {
			return nil, fmt.Errorf("error marshaling '%s': %w", fieldName, err)
		}
	}
	return json.Marshal(object)
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Gives a user points from a payer
	// (POST /points/give/{user})
	GivePoints(w http.ResponseWriter, r *http.Request, user string)
	// Spends a users points
	// (POST /points/spend/{user})
	SpendPoints(w http.ResponseWriter, r *http.Request, user string)
	// Returns a users balance of points based on the payer that added them
	// (GET /points/{user})
	DescribeBalance(w http.ResponseWriter, r *http.Request, user string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// GivePoints operation middleware
func (siw *ServerInterfaceWrapper) GivePoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "user" -------------
	var user string

	err = runtime.BindStyledParameter("simple", false, "user", chi.URLParam(r, "user"), &user)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "user", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GivePoints(w, r, user)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// SpendPoints operation middleware
func (siw *ServerInterfaceWrapper) SpendPoints(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "user" -------------
	var user string

	err = runtime.BindStyledParameter("simple", false, "user", chi.URLParam(r, "user"), &user)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "user", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SpendPoints(w, r, user)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// DescribeBalance operation middleware
func (siw *ServerInterfaceWrapper) DescribeBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "user" -------------
	var user string

	err = runtime.BindStyledParameter("simple", false, "user", chi.URLParam(r, "user"), &user)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "user", Err: err})
		return
	}

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DescribeBalance(w, r, user)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/points/give/{user}", wrapper.GivePoints)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/points/spend/{user}", wrapper.SpendPoints)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/points/{user}", wrapper.DescribeBalance)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9xWXWsbSwz9K4PufVxi39uHwj41oaUECjXJY8iDsiPbU+arM7KpMf7vRbO73vVHkhYS",
	"Gvpke6SRjo6O5NlCE1wMnjxnqLeQmyU5LF+v0KJv6IvJLD9Ra8MmeLSzFCIlNlTceBMJajCeaUEJdrsK",
	"PqUUktjigWcTNMmnptwkEyUY1EDirIqtOglWgaOccfHovd68v5o5Gb8oMBJ9X5lEGuo76OL37ve7Cj6b",
	"Nc2C6QpHa7/Oob7bwr+J5lDDP5OBmUlHy2SGG0rdpV11XCAbR5nRxVO0g+k5pIPn/U5gjlOeUBrFeJqt",
	"HCuP7gwxFcR9sKNbcq44MNozrTiC2WbeBxOkt5G8fhTpU0mzypE8/0LWPpsYjJ+HVleesSkiJYfGQg0J",
	"rTUeOaT3HxZydtEEBxUUSmq4NS54dYPWCh+HgC5VRhctqcvZteIlsnLocUFZoVplSlm17HZQKmDDVoKW",
	"ytUtpbVphPc1pdyG/O9iejGVTCGSx2ighnflqIKIvCysTNp4k4VZ02QriXYtbe30CZUoCK811GPpSoiE",
	"jphSLvI90t2SCmplNAhjRRy8HKgQI4xJ5rSiqtsDo/neq/W+dabMV0Fvev6lezJFMVrTFKCTb1kQbEeh",
	"npqsUUmluYdliLVjXDXBOfT6BHSRSo7B51Zx/0+nL4buYO5P4ZWZOdRF8ZnjyvKLgWjX6pn0K08/IjVM",
	"WlHnU0FeOYdp05HXq7cncZ6CU9hCLu69/mQQ9bMCHI/636LAcU1nSC7EvLYGDZPLvyXG/crElHDzLPAe",
	"3VvSZyF+tF77FgyiHOS4oDNq/FiyPVD3YvlDinyl1TN+hp3h1prM6qH1eZPdvSFeJT+0t8ca5r0mHzCT",
	"VsEraUy7RMv/LmpNWg6dpN39DAAA//8qjXu8qwoAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
