package httprouters_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {

	r := httprouters.NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp1, _ := testRequest(t, ts, "POST", "/")
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp1.Body.Close()

	resp2, _ := testRequest(t, ts, "GET", "/xyz") // string(body1))
	assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp2.Body.Close()

	resp, _ := testRequest(t, ts, "PATCH", "/")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

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
