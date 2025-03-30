package main

import (
	"errors"
	"fmt"
	"net/http"
	// "github.com/goccy/go-json"
	"github.com/eniac/mucache/pkg/utility"
	"github.com/eniac/mucache/pkg/synthetic"
	"github.com/eniac/mucache/pkg/slowpoke"
)

type endpointHandler struct {
	endpoint *synthetic.Endpoint
}

type Response struct {
	CPUResp  string            `json:"cpu_response"`
	NetworkResp   map[string]string `json:"network_response"`
}

func execTask(request *http.Request, endpoint *synthetic.Endpoint) Response {
	cpuResp := execCPU(endpoint)
	networkResp := execNetwork(request, endpoint)
	return Response{CPUResp: cpuResp, NetworkResp: networkResp}
}

func (handler endpointHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// fmt.Printf("Request to %s\n", handler.endpoint.Name)
	slowpoke.SlowpokeCheck(handler.endpoint.Name)
	response := execTask(request, handler.endpoint)
	// fmt.Printf("Response: %s\n", response)
	// respJSON, err := json.Marshal(response)
	// if err != nil {
	// 	panic(err)
	// }
	utility.DumpJson(response, writer)
}

// Launch a HTTP server to serve one or more endpoints
func serverHTTP(endpoints []synthetic.Endpoint) {
	slowpoke.SlowpokeInit()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Connected\n"))
		if err != nil {
			return
		}
	})

	for i := range endpoints {
		mux.Handle(fmt.Sprintf("/%s", endpoints[i].Name), endpointHandler{endpoint: &endpoints[i]})
		fmt.Printf("Endpoint %s registered\n", endpoints[i].Name)
	}

	err := http.ListenAndServe(":5000", mux)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
