package httprouters_test

import (
	"bytes"
	"compress/gzip"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/service"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/settings"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	s := storage.NewMapStorage(make(map[string]string), make(map[string]string), make(map[string]bool))
	ctx := context.Background()
	// ctx = context.WithValue(ctx, settings.CtxKeyUserID , "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff")
	srvs := services.NewService(s, "http://localhost:8080")
	h := handlers.NewHandler(srvs, "http://localhost:8080")
	r := httprouters.NewRouter(h)
	ts := httptest.NewServer(r)
	defer ts.Close()
	ctx, cancel := context.WithTimeout(ctx, settings.StorageTimeout)
	defer cancel()
	s.StoragePut(ctx, "xyz", "https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf", "5e7cb52e-691d-4f46-bc1c-7ae1616a59ff")

	CreateURLRequest, _ := CreateURLRequest(t, ts, "POST", "/")
	assert.Equal(t, http.StatusCreated, CreateURLRequest.StatusCode)
	//assert.Contains(t, "https://", body)
	defer CreateURLRequest.Body.Close()

	GetShortURLRequestOK, _ := GetShortURLRequest(t, ts, "GET", "/xyz")
	assert.Equal(t, http.StatusOK, GetShortURLRequestOK.StatusCode)
	//assert.Contains(t, "https://", body)
	defer GetShortURLRequestOK.Body.Close()

	GetShortURLRequestBad, _ := GetShortURLRequest(t, ts, "PATCH", "/")
	assert.Equal(t, http.StatusBadRequest, GetShortURLRequestBad.StatusCode)
	defer GetShortURLRequestBad.Body.Close()

	CreateURLRequestJSON, _ := CreateURLRequestJSON(t, ts, "POST", "/api/shorten")
	assert.Equal(t, http.StatusCreated, CreateURLRequestJSON.StatusCode)
	//assert.Contains(t, "https://", body)
	defer CreateURLRequestJSON.Body.Close()

	CreateURLRequestCompress, _ := CreateURLRequestCompress(t, ts, "POST", "/")
	assert.Equal(t, http.StatusCreated, CreateURLRequestCompress.StatusCode)
	//assert.Contains(t, "https://", body)
	defer CreateURLRequestCompress.Body.Close()

	CreateURLsListWrong, _ := CreateURLsListWrong(t, ts, "GET", "/api/user/urls")
	assert.Equal(t, http.StatusNoContent, CreateURLsListWrong.StatusCode)
	//assert.Contains(t, "https://", body)
	defer CreateURLsListWrong.Body.Close()

	CreateURLsList, _ := CreateURLsList(t, ts, "GET", "/api/user/urls")
	assert.Equal(t, http.StatusOK, CreateURLsList.StatusCode)
	//assert.Contains(t, "https://", body)
	defer CreateURLsList.Body.Close()

}

func GetShortURLRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func CreateURLRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader("https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf"))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func CreateURLRequestJSON(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(`{"url":"https://yandex.ru/search/?text=AToi+go&lr=213"}`))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func CreateURLRequestCompress(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte("https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf"))
	w.Close()

	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(b.String()))
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp, string(respBody)
}

func CreateURLsListWrong(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte("https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf"))
	w.Close()

	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(b.String()))
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp, string(respBody)
}

func CreateURLsList(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte("https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf"))
	w.Close()

	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(b.String()))
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")

	req.AddCookie(&http.Cookie{
		Name:   "token",
		Value:  "35653763623532652d363931642d346634362d626331632d3761653136313661353966663fe75fb6b45bd519a5e87f62c5507aff32f4410bed855e9c65628b7b9eee35b6",
		MaxAge: 300,
	})

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp, string(respBody)
}
