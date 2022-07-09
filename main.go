package main

import (
	"fmt"
	"k8s.io/klog"
	"encoding/json"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/test", serveTest)
	mux.HandleFunc("/hello", serveHello)

	klog.Info("Starting Go test api server...")
	klog.Flush()

	if err := http.ListenAndServe(":80", mux); err != nil {
		klog.Errorf("Failed to listen and serve test-api-server: %s", err)
	}
	klog.Info("Terminate Go test api server")
}


func serveHello(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		fmt.Println("You call GET /hello !!")
		SetResponse(res, "You call GET /hello !!", nil, http.StatusOK)
	case http.MethodPost:
		fmt.Println("You call POST /hello !!")
		SetResponse(res, "You call POST /hello !!", nil, http.StatusOK)
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

func SetResponse(res http.ResponseWriter, outString string, outJson interface{}, status int) http.ResponseWriter {

	//set Cors
	// res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	res.Header().Set("Access-Control-Max-Age", "3628800")
	res.Header().Set("Access-Control-Expose-Headers", "Content-Type, X-Requested-With, Accept, Authorization, Referer, User-Agent")

	//set Out
	if outJson != nil {
		res.Header().Set("Content-Type", "application/json")
		js, err := json.Marshal(outJson)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		//set StatusCode
		res.WriteHeader(status)
		res.Write(js)
		return res

	} else {
		//set StatusCode
		res.WriteHeader(status)
		res.Write([]byte(outString))
		return res

	}
}
