package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	c := client()
	resp, err := c.Get("http://127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Status: ", resp.Status, " Body:", string(b))
}

func client() *http.Client {
	config := &tls.Config{
		// GetClientCertificate:  utils.ClientCertReqFunc("", ""),
		// VerifyPeerCertificate: utils.CertificateChains,
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}
}
