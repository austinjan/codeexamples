package main

import (
	"crypto/tls"
	"net/http"

	// for ACME (Let's Encrypt)
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/acme/autocert"
)

func serveHttp() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello chi"))
	})

	// autocert.Manager communicates with Let's Encrypt to get certificates
	certManager := autocert.Manager{
		HostPolicy: autocert.HostWhitelist("mnms.com.io"),
		Prompt:     autocert.AcceptTOS,
		// how to store certificates, DirCache stores them in the local folder
		Cache: autocert.DirCache("certs"),
	}
	server := &http.Server{
		Addr:    ":3333",
		Handler: r,
		TLSConfig: &tls.Config{
			// use autocert.Manager to get certificates
			GetCertificate: certManager.GetCertificate,
		},
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	server.ListenAndServeTLS("", "")
}

func main() {
	serveHttp()
}
