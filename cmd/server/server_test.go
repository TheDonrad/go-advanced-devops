package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteMetric(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "gauge normal",
			want: want{
				statusCode: 200,
			},
			request: "/update/gauge/Alloc/1",
		},
		{
			name: "gauge bad",
			want: want{
				statusCode: 400,
			},
			request: "/update/gauge/PollCount/1",
		},
		{
			name: "counter normal",
			want: want{
				statusCode: 200,
			},
			request: "/update/counter/PollCount/1",
		},
		{
			name: "counter bad",
			want: want{
				statusCode: 400,
			},
			request: "/update/counter/Alloc/1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(writeMetric)
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			err := result.Body.Close()
			require.NoError(t, err)
		})
	}
}
