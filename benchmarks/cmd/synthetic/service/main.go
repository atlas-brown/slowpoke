package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/eniac/mucache/pkg/synthetic"
)

// Load the config map from the CONF environment variable
func LoadConfigMap() (*synthetic.ConfigMap, error) {
	configFilename := os.Getenv("CONF")
	configFile, err := os.Open(configFilename)
	configFileByteValue, _ := io.ReadAll(configFile)

	if err != nil {
		return nil, err
	}

	inputConfig := &synthetic.ConfigMap{}
	err = json.Unmarshal(configFileByteValue, inputConfig)

	if err != nil {
		return nil, err
	}

	return inputConfig, nil
}

func main() {
	configMap, err := LoadConfigMap()
	if err != nil {
		panic(err)
	}

	runtime.GOMAXPROCS(configMap.Processes)

	fmt.Printf("Starting synthetic service with %d processes\n", runtime.GOMAXPROCS(0))

	// TODO: Also support gRPC
	if configMap.Protocol == "http" {
		serverHTTP(configMap.Endpoints)
	} else if configMap.Protocol == "grpc" {
		serveGRPC(configMap.Endpoints)
	} else {
		panic("Unsupported protocol")
	}
}
