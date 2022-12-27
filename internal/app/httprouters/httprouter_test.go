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
	"github.com/dimsonson/go-yp-shortener-url/internal/app/handlers/servicemock"
	"github.com/dimsonson/go-yp-shortener-url/internal/app/httprouters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var UIDCookie = "35356630393565302d363765312d343835642d623262622d3630356636326163613566640269ca7988629be7ae55e705ac7d2b5924b95dfc00c796ea5492e17c2bbcd3f7"

func TestNewRouter(t *testing.T) {
	svs := &servicemock.ServiceMock{}
	base := "http://localhost:8080"
	hPut := handlers.NewPutHandler(svs, base)
	hGet := handlers.NewGetHandler(svs, base)
	hDel := handlers.NewDeleteHandler(svs, base)
	hPing := handlers.NewPingHandler(svs, base)
	r := httprouters.NewRouter(hPut, hGet, hDel, hPing)
	ts := httptest.NewServer(r)
	defer ts.Close()

	GetOK := Get(t, ts, http.MethodGet, "/xyz")
	assert.Equal(t, http.StatusOK, GetOK.StatusCode)

	GetNotOk := Get(t, ts, http.MethodPatch, "/")
	assert.Equal(t, http.StatusMethodNotAllowed, GetNotOk.StatusCode)

	PutOk, body := Put(t, ts, http.MethodPost, "/")
	assert.Equal(t, http.StatusCreated, PutOk.StatusCode)
	assert.Contains(t, body, base)
	defer PutOk.Body.Close()

	PutJSONok, body := PutJSON(t, ts, http.MethodPost, "/api/shorten")
	assert.Equal(t, http.StatusCreated, PutJSONok.StatusCode)
	assert.Contains(t, body, base)
	defer PutJSONok.Body.Close()

	PutGzipOk, body := PutGzip(t, ts, http.MethodPost, "/")
	assert.Equal(t, http.StatusCreated, PutGzipOk.StatusCode)
	assert.Contains(t, body, base)
	defer PutGzipOk.Body.Close()

	GetZipNotOk := GetZip(t, ts, http.MethodGet, "/api/")
	assert.Equal(t, http.StatusNotFound, GetZipNotOk.StatusCode)

	PutBatchOk := PutBatch(t, ts, "GET", "/api/user/urls")
	assert.Equal(t, http.StatusOK, PutBatchOk.StatusCode)
	defer PutBatchOk.Body.Close()

}

func Get(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: UidCookie})
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func Put(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader("https://pkg.go.dev/"))
	req.AddCookie(&http.Cookie{Name: "token", Value: UidCookie})
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()
	return resp, string(respBody)
}

func PutJSON(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(`{"url":"https://pkg.go.dev/io#Reader"}`))
	require.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "token", Value: UidCookie})
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	defer resp.Body.Close()
	return resp, string(respBody)
}

func PutGzip(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write([]byte("https://pkg.go.dev/"))
	if err != nil {
		return nil, ""
	}
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

func GetZip(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write([]byte("https://pkg.go.dev/"))
	if err != nil {
		return nil
	}
	w.Close()
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(b.String()))
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func PutBatch(t *testing.T, ts *httptest.Server, method, path string) (*http.Response) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write([]byte("https://pkg.go.dev/"))
	if err != nil {
		return nil
	}
	w.Close()
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(b.String()))
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")
	req.AddCookie(&http.Cookie{Name: "token", MaxAge: 300})
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}
