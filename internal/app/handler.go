package app

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/Esusu2017/rrs-commons/time"
	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
	"github.com/mmuoDev/location-history-api.git/internal/db"
	"github.com/mmuoDev/location-history-api.git/internal/workflow"
	"github.com/mmuoDev/location-history-api.git/pkg"
	"github.com/pkg/errors"
)

const (
	orderID = "order_id"
)

//ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

//DeleteHistoryHandler returns a http request to delete location history for an order
func DeleteHistoryHandler(deleteLocation db.DeleteHistoryFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := httprouter.ParamsFromContext(r.Context())
		orderId := params.ByName(orderID)
		deleteLocation(orderId)
	}
}

//RetrieveHistoryHandler returns a http request to retrieve location history for an order
func RetrieveHistoryHandler(retrieveLocation db.RetrieveHistoryFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := httprouter.ParamsFromContext(r.Context())
		orderId := params.ByName(orderID)
		var qp pkg.QueryParams
		if err := GetQueryParams(&qp, r); err != nil {
			res := ErrorResponse{Error: err.Error()}
			ServeJSON(res, w, http.StatusInternalServerError)
			return
		}
		retrieve := workflow.RetrieveHistory(retrieveLocation)
		res := retrieve(orderId, qp.Max)
		ServeJSON(res, w, http.StatusOK)
	}
}

//AddLocationHandler returns a http request to add location to an order
func AddLocationHandler(addLocation db.AddLocationFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := httprouter.ParamsFromContext(r.Context())
		orderId := params.ByName(orderID)
		var req pkg.Location
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			res := ErrorResponse{Error: err.Error()}
			ServeJSON(res, w, http.StatusInternalServerError)
			return
		}
		if req.Latitude == 0 || req.Longitude == 0 {
			res := ErrorResponse{Error: "Invalid user input"}
			ServeJSON(res, w, http.StatusBadRequest)
			return
		}
		add := workflow.AddLocation(addLocation)
		add(req, orderId)
		ServeJSON(req, w, http.StatusOK)
	}
}

//ServeJSON returns a JSON response
func ServeJSON(res interface{}, w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	bb, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(bb)
}

var decoder = schema.NewDecoder()

// GetQueryParams maps the query params from an http request into an interface
func GetQueryParams(value interface{}, r *http.Request) error {
	// decoder lookup for values on the json tag, instead of the default schema tag
	decoder.SetAliasTag("json")
	var globalErr error
	// Decoder Register for custom type ISO8601
	decoder.RegisterConverter(time.ISO8601{}, func(input string) reflect.Value {
		ISOTime, errISO := time.NewISO8601(input)

		if errISO != nil {
			globalErr = errors.Wrapf(errISO, "handler - invalid iso time provided")
			return reflect.ValueOf(time.ISO8601{})
		}

		return reflect.ValueOf(ISOTime)
	})

	// Decoder Register for custom type Epoch
	decoder.RegisterConverter(time.Epoch(0), func(input string) reflect.Value {
		ISOTime, errISO := time.NewISO8601(input)

		if errISO != nil {
			globalErr = errors.Wrapf(errISO, "handler - invalid iso time provided")
			return reflect.ValueOf(time.ISO8601{}.ToEpoch())
		}

		return reflect.ValueOf(ISOTime.ToEpoch())
	})

	if err := decoder.Decode(value, r.URL.Query()); err != nil {
		return errors.Wrapf(err, "handler - failed to decode query params")
	}

	if globalErr != nil {
		return globalErr
	}

	return nil
}
