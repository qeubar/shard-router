package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type endpoint struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
}

type shard struct {
	Lower  int    `json:"lower"`
	Upper  int    `json:"upper"`
	Server string `json:"server"`
}

var shards []shard
var endpoints []endpoint

func routes() http.Handler {
	mux := chi.NewMux()
	// TODO: add tests that all methods are caught
	for _, endpoint := range endpoints {
		for _, method := range endpoint.Methods {
			switch method {
			case http.MethodGet:
				mux.Get(endpoint.Path, router)
			case http.MethodHead:
				mux.Head(endpoint.Path, router)
			case http.MethodPost:
				mux.Post(endpoint.Path, router)
			case http.MethodPut:
				mux.Put(endpoint.Path, router)
			case http.MethodPatch:
				mux.Patch(endpoint.Path, router)
			case http.MethodDelete:
				mux.Delete(endpoint.Path, router)
			case http.MethodConnect:
				mux.Connect(endpoint.Path, router)
			case http.MethodOptions:
				mux.Options(endpoint.Path, router)
			case http.MethodTrace:
				mux.Trace(endpoint.Path, router)
			default:
				panic("invalid HTTP method " + method + " for endpoint " + endpoint.Path)
			}
		}
	}

	return mux
}

func router(w http.ResponseWriter, r *http.Request) {
	shardHost := "localhost.shard0:8080"
	// find shard

	uri := shardHost + r.URL.String()
	fmt.Println(uri)
	req, err := http.NewRequest(r.Method, uri, r.Body)
	if err != nil {
		// TODO log err
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.Header = r.Header
	req.Close = true

	resp, err := httpClient.Do(req)
	if err != nil {
		// TODO log err
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// TODO log err
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}
	w.Write(body)
}
