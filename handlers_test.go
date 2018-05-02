package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := newService()
	srv := httptest.NewServer(http.HandlerFunc(s.handleSet))
	defer srv.Close()

	req, _ := http.NewRequest("PUT", srv.URL+"/test", bytes.NewBuffer([]byte(
		`{"hello":"world"}`)))

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGet(t *testing.T) {
	s := newService()
	srv := httptest.NewServer(http.HandlerFunc(s.handleSet))
	defer srv.Close()
	testBody := []byte(`{"hello":"world"}`)
	req, _ := http.NewRequest("PUT", srv.URL+"/test", bytes.NewBuffer(testBody))

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	srv = httptest.NewServer(http.HandlerFunc(s.handleGet))
	req, _ = http.NewRequest("GET", srv.URL+"/test", nil)
	res, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, testBody, body)
}

func TestDelete(t *testing.T) {
	s := newService()
	srv := httptest.NewServer(http.HandlerFunc(s.handleSet))
	defer srv.Close()

	testBody := []byte(`{"hello":"world"}`)
	req, _ := http.NewRequest("PUT", srv.URL+"/test", bytes.NewBuffer(testBody))
	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	srv = httptest.NewServer(http.HandlerFunc(s.handleDelete))
	req, _ = http.NewRequest("DELETE", srv.URL+"/test", nil)
	res, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	srv = httptest.NewServer(http.HandlerFunc(s.handleGet))
	req, _ = http.NewRequest("GET", srv.URL+"/test", nil)
	res, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, []byte{}, body)
}
