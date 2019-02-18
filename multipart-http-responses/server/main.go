// https://peter.bourgon.org/blog/2019/02/12/multipart-http-responses.html

package main

import (
	"log"
	"mime"
	"mime/multipart"
	"net/http"
)

func getValues() [][]byte {
	return [][]byte{[]byte("a bytes in string"), []byte("bbbb")}
}

// handle using multipart, adv: no need to go through a base64 conversion
func handle(w http.ResponseWriter, r *http.Request) {
	mediatype, _, err := mime.ParseMediaType(r.Header.Get("Accept"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	if mediatype != "multipart/form-data" {
		http.Error(w, "set Accept: multipart/form-data", http.StatusMultipleChoices)
		return
	}
	mw := multipart.NewWriter(w)
	w.Header().Set("Content-Type", mw.FormDataContentType())
	for _, value := range getValues() {
		fw, err := mw.CreateFormField("value")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fw.Write(value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := mw.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// response:
// 2019/02/18 13:11:28 Value: a bytes in string
// 2019/02/18 13:11:28 Value: bbbb
