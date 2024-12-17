package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func APITestAfterLogin(t *testing.T, body any, Method, url, token string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	server := newTestServer(t)
	require.NotEmpty(t, server)

	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		request, err := http.NewRequest(Method, url, bytes.NewReader(data))
		require.NoError(t, err)
		require.NotEmpty(t, request)

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", token)

		server.router.ServeHTTP(recorder, request)
		require.NotEmpty(t, recorder)
	} else {
		request, err := http.NewRequest(Method, url, nil)
		require.NoError(t, err)
		require.NotEmpty(t, request)

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", token)

		server.router.ServeHTTP(recorder, request)
		require.NotEmpty(t, recorder)
	}

	return recorder
}

func APITestBeforeLogin(t *testing.T, body any, Method, url string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	server := newTestServer(t)
	require.NotEmpty(t, server)

	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)
		require.NotEmpty(t, data)

		request, err := http.NewRequest(Method, url, bytes.NewReader(data))
		require.NoError(t, err)
		require.NotEmpty(t, request)

		server.router.ServeHTTP(recorder, request)
		require.NotEmpty(t, recorder)
	} else {
		request, err := http.NewRequest(Method, url, nil)
		require.NoError(t, err)
		require.NotEmpty(t, request)

		server.router.ServeHTTP(recorder, request)
		require.NotEmpty(t, recorder)
	}

	return recorder
}
