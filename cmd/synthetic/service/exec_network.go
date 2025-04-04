package main

import (
	"net/http"
	"fmt"
	"context"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/eniac/mucache/pkg/synthetic"
	// "github.com/goccy/go-json"
)

func execParallelHTTP(calledServices []synthetic.CalledService) map[string]string {
	// forward requests
	respMap := make(map[string]string)
	var wg sync.WaitGroup
	for _, service := range calledServices {
		wg.Add(1)
		go func(service synthetic.CalledService) {
			defer wg.Done()
			// fmt.Printf("Calling %s (%s)\n", service.Service, service.Endpoint)
			resp := slowpoke.InvokeSynthtic(context.Background(), service.Service, service.Endpoint, "")
			key := fmt.Sprintf("%s (%s)", service.Service, service.Endpoint)
			respMap[key] = fmt.Sprintf("{ \"response\": %s }", resp)
		}(service)
	}
	wg.Wait()
	return respMap
}

func execSequentialHTTP(calledServices []synthetic.CalledService) map[string]string {
	// forward requests
	respMap := make(map[string]string)
	for _, service := range calledServices {
		// fmt.Printf("Calling %s (%s)\n", service.Service, service.Endpoint)
		resp := slowpoke.InvokeSynthtic(context.Background(), service.Service, service.Endpoint, "")
		key := fmt.Sprintf("%s (%s)", service.Service, service.Endpoint)
		respMap[key] = fmt.Sprintf("{ \"response\": %s }", resp)
	}
	return respMap
}

func execNetwork(request *http.Request, endpoint *synthetic.Endpoint) map[string]string {
	if endpoint.NetworkComplexity == nil {
		return map[string]string{"nil": "No network complexity"}
	}
	// calledServices := endpoint.NetworkComplexity.CalledServices
	// // forward requests
	// respMap := make(map[string]string)
	// for _, service := range calledServices {
	// 	// fmt.Printf("Calling %s (%s)\n", service.Service, service.Endpoint)
	// 	resp := slowpoke.InvokeSynthtic(context.Background(), service.Service, service.Endpoint, "")
	// 	key := fmt.Sprintf("%s (%s)", service.Service, service.Endpoint)
	// 	respMap[key] = fmt.Sprintf("{ \"response\": %s }", resp)
	// }
	return respMap


	if endpoint.NetworkComplexity.ForwardRequests == "parallel" {
		return execParallelHTTP(endpoint.NetworkComplexity.CalledServices)
	} else {
		return execSequentialHTTP(endpoint.NetworkComplexity.CalledServices)
	}
}