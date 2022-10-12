package httprouters_test

import (
	"bytes"
	"compress/gzip"
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
	s := storage.NewMapStorage(make(map[string]int), make(map[string]string))
	srvs := services.NewService(s)
	h := handlers.NewHandler(srvs, "")
	r := httprouters.NewRouter(h)
	ts := httptest.NewServer(r)
	defer ts.Close()

	s.PutToStorage(0,"xyz", "https://pkg.go.dev/github.com/stretchr/testify@v1.8.0/assert#Containsf")

	resp1, _ := CreateURLRequest(t, ts, "POST", "/")
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp1.Body.Close()

	resp2, _ := GetShortURLRequest(t, ts, "GET", "/xyz")
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp2.Body.Close()

	resp, _ := GetShortURLRequest(t, ts, "PATCH", "/")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	defer resp.Body.Close()

	resp3, _ := CreateURLRequestJSON(t, ts, "POST", "/api/shorten")
	assert.Equal(t, http.StatusCreated, resp3.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp3.Body.Close()

	resp4, _ := CreateURLRequestCompress(t, ts, "POST", "/")
	assert.Equal(t, http.StatusCreated, resp4.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp4.Body.Close()

	resp5, _ := CreateURLsListWrong(t, ts, "GET", "/api/user/urls")
	assert.Equal(t, http.StatusNoContent, resp5.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp5.Body.Close()

	resp6, _ := CreateURLsList(t, ts, "GET", "/api/user/urls")
	assert.Equal(t, http.StatusOK, resp6.StatusCode)
	//assert.Contains(t, "https://", body)
	defer resp6.Body.Close()

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
		Value:  "00000000b38aaf6c89467a765a15a5d40098d050c80503562bebef1c64ded15cc4fbdaeb",
		MaxAge: 300,
	})

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp, string(respBody)
}
