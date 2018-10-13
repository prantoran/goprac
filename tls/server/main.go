package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func main() {
	// create tcp conn
	// if tls then encrypt conn else skip
	// send http packets on conn

	s := server()
	http.HandleFunc("/", handler)
	s.ListenAndServe()
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handling req")
	w.Write([]byte("yo"))
}

func server() *http.Server {
	tls := &tls.Config{
		// GetCertificate:        utils.CertReqFunc("", ""),
		// VerifyPeerCertificate: utils.CertificateChains,
	}

	return &http.Server{
		Addr:      ":8080",
		TLSConfig: tls,
	}
}
