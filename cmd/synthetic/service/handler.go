package main

import (
	// "errors"
	"context"
	"fmt"
	"net"
	"net/http"

	// "github.com/goccy/go-json"
	pb "github.com/eniac/mucache/pkg/pb"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/eniac/mucache/pkg/synthetic"
	"github.com/eniac/mucache/pkg/utility"
	"google.golang.org/grpc"
)

type endpointHandler struct {
	endpoint *synthetic.Endpoint
}

type Response struct {
	CPUResp     string `json:"cpu_response"`
	NetworkResp string `json:"network_response"`
}

func execTask(request *http.Request, endpoint *synthetic.Endpoint) Response {
	cpuResp := execCPU(endpoint)
	networkResp := execNetwork(request, endpoint)
	networkRespStr := "["
	for key, value := range networkResp {
		networkRespStr += fmt.Sprintf("{%s: %s} ", key, value)
	}
	networkRespStr += "]"
	cpuResp = ""
	networkRespStr = ""
	return Response{CPUResp: cpuResp, NetworkResp: networkRespStr}
}

func (handler endpointHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	slowpoke.SlowpokeCheck(handler.endpoint.Name)
	slowpoke.SlowpokePokerPPDelay()
	// slowpoke.SlowpokeDelay()
	response := execTask(request, handler.endpoint)
	utility.DumpJson(response, writer)
	// f, ok := writer.(http.Flusher)
	// if ok {
	// 	slowpoke.SlowpokeFlushDelay(f)
	// } else {
	// 	slowpoke.SlowpokeDelay()
	// }

	// f, ok := writer.(http.Flusher)
	// if ok {
	// 	f.Flush()
	// 	// slowpoke.SlowpokeFlushDelay(f)
	// } else {
	// 	panic("Flusher not available")
	// }
}

// Launch a HTTP server to serve one or more endpoints
func serverHTTP(endpoints []synthetic.Endpoint) {
	slowpoke.SlowpokeInit()
	fmt.Println("Starting HTTP server")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Connected\n"))
		if err != nil {
			return
		}
	})
	mux.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Heartbeat\n"))
		if err != nil {
			return
		}
	})

	for i := range endpoints {
		mux.Handle(fmt.Sprintf("/%s", endpoints[i].Name), endpointHandler{endpoint: &endpoints[i]})
		fmt.Printf("Endpoint %s registered\n", endpoints[i].Name)
	}

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	slowpokeListener := &slowpoke.SlowpokeListener{listener}

	panic(http.Serve(slowpokeListener, mux))
}

// GPRC SERVER

type grpcServer struct {
	pb.UnimplementedSimpleServer
	endpoints []synthetic.Endpoint
}

func (s *grpcServer) SimpleRPC(ctx context.Context, in *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	target_endpoint := in.Endpoint
	var response Response
	for i := range s.endpoints {
		if s.endpoints[i].Name == target_endpoint {
			slowpoke.SlowpokeCheck(target_endpoint)
			response = execTask(nil, &s.endpoints[i])
		}
	}
	jsonBytes, err := utility.MarshalJson(response)
	var respStr string
	if err != nil {
		respStr = fmt.Sprintf("{ \"error\": \"%s\" }", err.Error())
	} else {
		respStr = string(jsonBytes)
	}
	return &pb.SimpleResponse{Resp: respStr}, nil
}

func serveGRPC(endpoints []synthetic.Endpoint) {
	slowpoke.SlowpokeInit()
	fmt.Println("Starting gRPC server")
	lis, err := net.Listen("tcp", ":5000")
	s := grpc.NewServer()
	pb.RegisterSimpleServer(s, &grpcServer{endpoints: endpoints})
	if err != nil {
		panic(err)
	}
	s.Serve(lis)
}
