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
	"fmt"
	"regexp"
)

func handleMtSupporteMethods(body string, w http.ResponseWriter) bool {
	matched, _ := regexp.MatchString("mt\\.supportedMethods", body)
	if matched {
		io.WriteString(w, "<?xml version=\"1.0\"?><methodResponse><params><param><value><array><data><value><string>metaWeblog.getRecentPosts</string></value><value><string>metaWeblog.newPost</string></value></data></array></value></param></params></methodResponse>")
	fmt.Println("Responded to mt.supportedMethods");
	}
	return matched
}

func handleMetaWeblogGetRecentPosts(body string, w http.ResponseWriter) bool {
	matched, _ := regexp.MatchString("metaWeblog\\.getRecentPosts", body)
	if matched {
		io.WriteString(w, "<?xml version=\"1.0\"?><methodResponse><params><param><value><array><data/></array></value></param></params></methodResponse>")
		fmt.Println("Responded to get recent posts");
	}
	return matched
}

func handleMetaWeblogNewPost(body string, w http.ResponseWriter) bool {
	matched, _ := regexp.MatchString("metaWeblog\\.newPost", body)

	if matched {
		fmt.Println(body)
		io.WriteString(w, "I don't handle you yet!")
		fmt.Println("Got new post message");
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
				fmt.Println("Not a known method call %s", buf);
				// return error
				io.WriteString(w, "<?xml version=\"1.0\"?><methodResponse><fault><value><struct><member><name>faultCode</name><value><int>-32601</int></value></member><member><name>faultString</name><value><string>server error. requested method not found</string></value></member></struct></value></fault></methodResponse>")
			}
		} else {
			fmt.Println(err);
			io.WriteString(w, "<?xml version=\"1.0\"?><methodResponse><fault><value><struct><member><name>faultCode</name><value><int>-32601</int></value></member><member><name>faultString</name><value><string>server error. requested method not found</string></value></member></struct></value></fault></methodResponse>")
		}
	}

	http.HandleFunc("/xmlrpc.php", handleRequests)

	log.Println("Starting XML-RPC server on localhost:80/xmlrpc.php")
	log.Fatal(http.ListenAndServe(":80", nil))
}
