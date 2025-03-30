package main

import (
	"encoding/json"
	"io"
	"os"
	"runtime"
	"github.com/eniac/mucache/pkg/synthetic"
	"fmt"
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

	serverHTTP(configMap.Endpoints)

	// TODO: Also support gRPC
	// if synthetic.ConfigMap.Protocol == "http" {
	// 	server.HTTP(synthetic.ConfigMap.Endpoints)
	// } else if synthetic.ConfigMap.Protocol == "grpc" {
	// 	server.GRPC(synthetic.ConfigMap.Endpoints)
	// }
}