package main

import (
	"cloud-credential-api-server/billing"
	"fmt"

	"k8s.io/klog"

	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/billing", serveBilling)
	mux.HandleFunc("/test", serveTest)
	mux.HandleFunc("/hello", serveHello)

	klog.Info("Starting Cloud Credential server...")
	klog.Flush()

	if err := http.ListenAndServe(":80", mux); err != nil {
		klog.Errorf("Failed to listen and serve AWS COST server: %s", err)
	}
	klog.Info("Terminate Cloud Credential server")
}

func serveBilling(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		billing.Get(res, req)
	case http.MethodPut:
	case http.MethodOptions:
	default:
		//error
	}
}


func serveHello(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		fmt.Println("You call GET /hello !!")
	case http.MethodPost:
		fmt.Println("You call POST /hello !!")
	case http.MethodOptions:
	default:
		//error
	}
}

func serveTest(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		TestGet(res, req)
	case http.MethodPut:
	case http.MethodOptions:
	default:
		//error
	}
}

func TestGet(res http.ResponseWriter, req *http.Request) {
	queryParams := req.URL.Query()

	keys := queryParams["key"]

	for _, k := range keys {
		fmt.Println(k)
	}
	fmt.Println()
}
