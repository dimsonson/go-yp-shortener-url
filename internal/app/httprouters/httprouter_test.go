package httprouters_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/services"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	storage.DB["xyz"] = "https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf"

	s := storage.NewMapStorage("map")
	srvs := services.NewService(s)
	h := handlers.NewHandler(srvs)
	r := httprouters.NewRouter(h)
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp1, _ := testRequest1(t, ts, "POST", "/")
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp1.Body.Close()

	resp2, _ := testRequest(t, ts, "GET", "/xyz") // string(body1))
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp2.Body.Close()

	resp, _ := testRequest(t, ts, "PATCH", "/")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	resp3, _ := testRequest2(t, ts, "POST", "/api/shorten")
	assert.Equal(t, http.StatusCreated, resp3.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp3.Body.Close()

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

func testRequest1(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader("https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf"))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func testRequest2(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(`{"url":"https://yandex.ru/search/?text=AToi+go&lr=213"}`))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}
