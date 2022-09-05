package main_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_main(t *testing.T) {
	r := chi.NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Contains(t, "https://", body)

	resp, _ = testRequest(t, ts, "POST", "/")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Contains(t, "https://", body)

	resp, _ = testRequest(t, ts, "PATCH", "/")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	

}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
