package main

import (
	"goAdvancedTpl/internal/server/handlers"
	"goAdvancedTpl/internal/server/storage"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteMetric(t *testing.T) {

	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	statusCode, _ := testRequest(t, ts, "POST", "/update/gauge/Alloc/1")
	assert.Equal(t, http.StatusOK, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/someBad/PollCount/1")
	assert.Equal(t, http.StatusNotImplemented, statusCode)

	statusCode, _ = testRequest(t, ts, "POST", "/update/counter/PollCount/1")
	assert.Equal(t, http.StatusOK, statusCode)

}

func TestGetMetric(t *testing.T) {

	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	statusCode, _ := testRequest(t, ts, "GET", "/value/gauge/Alloc")
	assert.Equal(t, http.StatusNotFound, statusCode)

}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	return resp.StatusCode, string(respBody)
}

func NewRouter() chi.Router {
	metStorage := storage.NewMetricStorage()
	h := handlers.NewAPIHandler(metStorage, "", "")
	r := chi.NewRouter()

	r.Post("/update/{metricType}/{metricName}/{metricValue}", h.WriteMetric)
	r.Get("/value/{metricType}/{metricName}", h.GetMetric)

	return r
}
