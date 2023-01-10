// Package shttp, sample of http
package shttp

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"

	"github.com/justinas/alice"
)

func firstMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("first middleware\n"))
		h.ServeHTTP(w, r)
	})
}

func secondMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("second middleware\n"))
		h.ServeHTTP(w, r)
	})
}

func realHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

// startHttpServer a http server sample code, refer TestMiddleware() to see how to use it
func startHttpServer(wg *sync.WaitGroup) *http.Server {
	// make middleware chain, the real handler is the last one. we can add as many as middleware functions as we want
	// for example authCheck, logRequest etc.
	chain := alice.New(firstMiddleware, secondMiddleware).ThenFunc(realHandler)
	mux := http.NewServeMux()
	mux.HandleFunc("/", chain.ServeHTTP)

	srv := &http.Server{Addr: ":27182", Handler: mux}

	go func() {
		defer wg.Done() // let main know we are done cleaning up

		// always returns error. ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return srv
}

// TestMiddleware tests middleware
func TestMiddleware(t *testing.T) {
	// start http server
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := startHttpServer(httpServerExitDone)

	// make request to server
	resp, err := http.Get("http://localhost:27182/api/123")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	t.Log("response Status:", resp.Status)
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	expect := "first middleware\nsecond middleware\nhello world"
	if string(respBody) != expect {
		t.Errorf("expected %q, got %q", expect, string(respBody))
	}
	// shutdown http server
	if err := srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	// wait for goroutine started in startHttpServer() to stop
	httpServerExitDone.Wait()

}
