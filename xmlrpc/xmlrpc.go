package main

import (
	"bytes"
	// "github.com/divan/gorilla-xmlrpc/xml"
	// "github.com/gorilla/rpc"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type request struct {
	Params []param `xml:"params>param"`
}

type param struct {
	Value value `xml:"value"`
}

type value struct {
	Array  []value  `xml:"array>data>value"`
	Struct []member `xml:"struct>member"`
	String string   `xml:"string"`
}

type member struct {
	Name  string `xml:"name"`
	Value value  `xml:"value"`
}

type webhookRequest struct {
	url string
}

func handleMtSupporteMethods(body string, w http.ResponseWriter) bool {
	matched, _ := regexp.MatchString("mt\\.supportedMethods", body)
	if matched {
		io.WriteString(w, "<?xml version=\"1.0\"?><methodResponse><params><param><value><array><data><value><string>metaWeblog.getRecentPosts</string></value><value><string>metaWeblog.newPost</string></value></data></array></value></param></params></methodResponse>")
	}
	return matched
}

func handleMetaWeblogGetRecentPosts(body string, w http.ResponseWriter) bool {
	matched, _ := regexp.MatchString("metaWeblog\\.getRecentPosts", body)
	if matched {
		io.WriteString(w, "<?xml version=\"1.0\"?><methodResponse><params><param><value><array><data/></array></value></param></params></methodResponse>")
	}
	return matched
}

func handleMetaWeblogNewPost(body string, w http.ResponseWriter) bool {
	matched, _ := regexp.MatchString("metaWeblog\\.newPost", body)

	if matched {
		id := strconv.FormatInt(time.Now().Unix(), 10)
		response := "<?xml version=\"1.0\"?><methodResponse><params><param><value><string>" + id + "</string></value></param></params></methodResponse>"

		// Parse the body and forward on the contents to the URL requested
		forwardRequest(body)

		io.WriteString(w, response)
	}
	return matched
}

func forwardRequest(body string) {
	request, err := parseWebhookRequest(body)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(request)
	}
}

func parseWebhookRequest(rawXml string) (webhookRequest, error) {
	var req request

	decoder := xml.NewDecoder(strings.NewReader(rawXml))
	err := decoder.Decode(&req)
	if err != nil {
		fmt.Println(err.Error())
		return webhookRequest{}, err
	} else {
		username := req.Params[1].Value.String
		password := req.Params[2].Value.String
		fmt.Printf("Username: %s, Password: %s\n", username, password)
		for _, mem := range req.Params[3].Value.Struct {
			name := mem.Name
			if name == "description" {
				desc := mem.Value.String
				fmt.Printf("Description: %s\n", desc)
				return webhookRequest{url: strings.TrimSpace(desc)}, nil
			}
		}
		return webhookRequest{}, errors.New("No description found to parse")
	}
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
				fmt.Println("Not a known method call %s", buf)
				// return error
				io.WriteString(w, "<?xml version=\"1.0\"?><methodResponse><fault><value><struct><member><name>faultCode</name><value><int>-32601</int></value></member><member><name>faultString</name><value><string>server error. requested method not found</string></value></member></struct></value></fault></methodResponse>")
			}
		} else {
			fmt.Println(err)
			io.WriteString(w, "<?xml version=\"1.0\"?><methodResponse><fault><value><struct><member><name>faultCode</name><value><int>-32601</int></value></member><member><name>faultString</name><value><string>server error. requested method not found</string></value></member></struct></value></fault></methodResponse>")
		}
	}

	http.HandleFunc("/xmlrpc.php", handleRequests)

	log.Println("Starting XML-RPC server on localhost:80/xmlrpc.php")
	log.Fatal(http.ListenAndServe(":80", nil))
}
