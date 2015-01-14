package main

import (
	"bytes"
	// "github.com/divan/gorilla-xmlrpc/xml"
	// "github.com/gorilla/rpc"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	// "strings"
	// "fmt"
	"regexp"
)

func handleMtSupporteMethods(body string, w http.ResponseWriter) bool {
	matched, _ := regexp.MatchString("mt\\.supportedMethods", body)
	if matched {
		io.WriteString(w, "get supported methods")
	}
	return matched
}

func handleMetaWeblogGetRecentPosts(body string, w http.ResponseWriter) bool {
	matched, _ := regexp.MatchString("metaWeblog\\.getRecentPosts", body)
	if matched {
		io.WriteString(w, "get recent posts")
	}
	return matched
}

func handleMetaWeblogNewPost(body string, w http.ResponseWriter) bool {
	matched, _ := regexp.MatchString("metaWeblog\\.newPost", body)
	if matched {
		io.WriteString(w, "proxying to webhook!")
	}
	return matched
}

// Entry point for server
func main() {
	handleRequests := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			buf := bytes.NewBuffer(body).String()

			if handleMtSupporteMethods(buf, w) {
				return
			} else if handleMetaWeblogGetRecentPosts(buf, w) {
				return
			} else if handleMetaWeblogNewPost(buf, w) {
				return
			} else {
				// return error
				io.WriteString(w, "Error - method not found!")
			}
		} else {
			io.WriteString(w, err.Error())
		}
	}

	http.HandleFunc("/xmlrpc.php", handleRequests)

	log.Println("Starting XML-RPC server on localhost:1234/xmlrpc.php")
	log.Fatal(http.ListenAndServe(":1234", nil))
}
