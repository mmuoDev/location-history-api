package app_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/mmuoDev/location-history-api.git/internal/app"
	"github.com/mmuoDev/location-history-api.git/pkg"
	"github.com/mmuoDev/location-history-api.git/pkg/db"
	"github.com/stretchr/testify/assert"
)

const (
	orderID = "12345"
)

func storageProvider() *db.Storage {
	return &db.Storage{}
}

func TestDeleteHistoryWorksAsExpected(t *testing.T) {
	deleteHistoryIsInvoked := false
	mockDeleteStorage := func(o *app.OptionalArgs) {
		o.DeleteHistory = func(orderID string) {
			deleteHistoryIsInvoked = true
		}
	}
	opts := []app.Options{
		mockDeleteStorage,
	}
	ap := app.New(storageProvider(), opts...)
	serverURL, cleanUpServer := newTestServer(ap.Handler())
	defer cleanUpServer()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/location/%s", serverURL, orderID), nil)
	client := &http.Client{}
	res, _ := client.Do(req)

	t.Run("Http Status Code is 200", func(t *testing.T) {
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})
	t.Run("Delete from storage is invoked", func(t *testing.T) {
		assert.True(t, deleteHistoryIsInvoked)
	})
}
func TestRetrieveHistoryWorksAsExpected(t *testing.T) {
	retrieveStorageIsInvoked := false
	mockRetrieveStorage := func(o *app.OptionalArgs) {
		o.RetrieveHistory = func(orderID string, max int) []pkg.Location {
			retrieveStorageIsInvoked = true
			var res []pkg.Location
			fileToStruct(filepath.Join("testdata", "retrieve_history_storage_response.json"), &res)
			return res
		}
	}
	opts := []app.Options{
		mockRetrieveStorage,
	}
	ap := app.New(storageProvider(), opts...)
	serverURL, cleanUpServer := newTestServer(ap.Handler())
	defer cleanUpServer()

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/location/%s", serverURL, orderID), nil)
	client := &http.Client{}
	res, _ := client.Do(req)

	t.Run("Http Status Code is 200", func(t *testing.T) {
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})
	t.Run("Retrieve from storage is invoked", func(t *testing.T) {
		assert.True(t, retrieveStorageIsInvoked)
	})
	t.Run("Response Body is as expected", func(t *testing.T) {
		var (
			expectedResponse pkg.RetrieveLocationsResponse
			actualResponse   pkg.RetrieveLocationsResponse
		)
		json.NewDecoder(res.Body).Decode(&actualResponse)
		fileToStruct(filepath.Join("testdata", "retrieve_history_response.json"), &expectedResponse)
		assert.Equal(t, expectedResponse, actualResponse)
	})
}

func TestAddLocationWorksAsExpected(t *testing.T) {
	addToStorageIsInvoked := false

	mockAddAstorage := func(o *app.OptionalArgs) {
		o.AddLocation = func(req pkg.Location, orderID string) {
			addToStorageIsInvoked = true
			t.Run("OrderID is as expected", func(t *testing.T) {
				assert.Equal(t, "12345", orderID)
			})
		}
	}

	opts := []app.Options{
		mockAddAstorage,
	}
	ap := app.New(storageProvider(), opts...)
	serverURL, cleanUpServer := newTestServer(ap.Handler())
	defer cleanUpServer()

	reqPayload, _ := os.Open(filepath.Join("testdata", "add_location_request.json"))
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/location/%s/now", serverURL, orderID), reqPayload)

	client := &http.Client{}
	res, _ := client.Do(req)

	t.Run("Http Status Code is 200", func(t *testing.T) {
		assert.Equal(t, res.StatusCode, http.StatusOK)
	})
	t.Run("Add to Storage is invoked", func(t *testing.T) {
		assert.True(t, addToStorageIsInvoked)
	})
}

func newTestServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)
	return ts.URL, func() { ts.Close() }
}

// fileToStruct reads a json file to a struct
func fileToStruct(filepath string, s interface{}) io.Reader {
	bb, _ := ioutil.ReadFile(filepath)
	json.Unmarshal(bb, s)
	return bytes.NewReader(bb)
}
