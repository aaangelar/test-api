package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"fmt"
	"os"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

type GetUsersRequest struct {
	UserID string `json:"user_id"`
}


type GetUsersResponse struct {
 	Details []User   `json:"details"`
}


//Variable for errLogs
var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
	// ServiceName is the service name
	ServiceName = os.Getenv("SERVICE_NAME")
	// ErrUnableToParse is returned when an error occurred while parsing URI params
	ErrUnableToParse = errors.New("Unable to parse query string")
	// ErrEmptyRFQBatch is returned when an there is a 0 count on rfq_batch_details
	ErrEmptyRFQBatch = errors.New("cannot proccess empty batch")
	// ErrEmailSendFailed is returned when an error occured from AWS related settings
	ErrEmailSendFailed = errors.New("unable to send email")
	// ResponseTime is the response time
	ResponseTime = ""
	// Debug is a switch for debug
	Debug = "0"
)

//STEP1

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a profilesvc server.
func MakeHTTPHandler(ctx context.Context, s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/GetUser/{user_id}").Handler(httptransport.NewServer(
		ctx,
		e.GetUserEndpoint,
		decodeGetUsersRequest,
		encodeResponse,
		options...,
	))

	return r
}

//STEP2 - Decode Request
func decodeGetUsersRequest(ctx context.Context, r *http.Request) (interface{}, error) {

    path := r.RequestURI
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return GetUsersRequest{}, ErrUnableToParse
	}

	vars := mux.Vars(r)
	sku, ok := vars["user_id"]
	if !ok {
		return nil, ErrBadRouting
	}

	return GetUsersRequest{UserID: user_id}, nil
}

//STEP3 - Encode Request
func encodeGetUsersRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(GetUsersRequest)
	user_id := url.QueryEscape(r.UserID)
	req.Method, req.URL.Path = "GET", "/"+ServiceName+"/GetSKUComment/"+user_id
	return encodeRequest(ctx, req, request)
}

//STEP4 - Decode Response
func decodeGetUsersResponse(ctx context.Context, req *http.Request, request interface{})(interface{}, error) {

	var response GetUsersResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}


//STEP5
// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	errMsg := err.Error()
	if Debug == "0" {
		errMsg = ""
	}
	json.NewEncoder(w).Encode(errorResponse{
		ResponseTime:    ResponseTime,
		ResponseStatus:  codeFrom(err),
		ResponseMessage: errMsg,
		Result:          map[string]interface{}{},
	})
}

type errorResponse struct {
	ResponseTime    string      `json:"response_time"`
	ResponseStatus  int         `json:"response_status"`
	ResponseMessage string      `json:"response_message,omitempty"`
	Result          interface{} `json:"result"`
}

func codeFrom(err error) int {
	switch err {
	case nil, SubStatusError, JsonError:
		return http.StatusOK
	case ErrNotFound:
		return http.StatusOK
	case ErrAlreadyExists, ErrInconsistentIDs, ErrEmptyRFQBatch, ErrEmailSendFailed:
		return http.StatusBadRequest
	case ErrCannotConnect:
		return 503
	case ErrSQL:
		return 504
	default:
		return http.StatusInternalServerError
	}
}
