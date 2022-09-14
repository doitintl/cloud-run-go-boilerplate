package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	require := require.New(t)

	api := API{}

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(err, "could not create request")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.indexHandler)
	handler.ServeHTTP(rr, req)

	wantStatus := 200

	require.Equal(wantStatus, rr.Code, "wrong status code returned")

	res := rr.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	wantBody := "<h1>Hello World!</h1>"

	require.Equal(wantBody, string(data), "wrong body returned")
}
