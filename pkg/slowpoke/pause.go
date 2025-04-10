package slowpoke

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	// "os"
)

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"` // Defer parsing
}

type RegisterReq struct {
	Endpoint string		`json:"endpoint"`
}

type PauseReq struct {
	Phase int		`json:"phase"`
}

func updateNeighbor(serv string) bool {
	neighborsLock.RLock()
	if _, exists := neighbors[serv]; exists {
		neighborsLock.RUnlock()
		return true
	}
	neighborsLock.RUnlock()
	neighborsLock.Lock()
	if _, exists := neighbors[serv]; exists {
		neighborsLock.Unlock()
		return true
	}
	serverAddr := fmt.Sprintf("%s.%s.svc.cluster.local:%s", serv, "default", "5550")
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	neighbors[serv] = conn
	neighborsLock.Unlock()
	fmt.Printf("Neighbor %s registered\n", serv)
	return false
}

func handleRegister(regReq RegisterReq) {
	serv := regReq.Endpoint
	updateNeighbor(serv)
}

func requestEndpointRegister(serv string) {
	exists := updateNeighbor(serv)
	if !exists {
		neighborsLock.RLock()
		conn := neighbors[serv]
		neighborsLock.RUnlock()
		regReq := RegisterReq{Endpoint: servName}
		regReqJson, err := json.Marshal(regReq)
		if err != nil {
			panic("marshal register request")
		}
		req := Message{Type: "register", Payload: regReqJson}
		jsonData, err := json.Marshal(req)
		if err != nil {
			panic("marshal request")
		}
		jsonData = append(jsonData, '\n')
		if _, err := conn.Write(jsonData); err != nil {
			panic("Client send error")
		}
	}
}

func requestPause(serv string, phase int) {
	conn := neighbors[serv]
	pauseReq := PauseReq{Phase: phase}
	pauseReqJson, err := json.Marshal(pauseReq)
	if err != nil {
		panic("marshal register request")
	}
	req := Message{Type: "pause", Payload: pauseReqJson}
	jsonData, err := json.Marshal(req)
	if err != nil {
		panic("marshal request")
	}
	jsonData = append(jsonData, '\n')
	if _, err := conn.Write(jsonData); err != nil {
		panic("Client send error")
	}
}

func handlePause(pauseReq PauseReq) {
	seen := true
	var delayToDo int64
	sync_guard.Lock()
	p := pauseReq.Phase
	if p > sleepPhase {
		sleepPhase = p
		delayToDo = accumulatedDelay
		accumulatedDelay = 0
		reqcount = 0
		seen = false
	}
	sync_guard.Unlock()
	// fmt.Printf("phase %v seen %v, delayToDo %v, reqCount %v\n", p, seen, delayToDo, reqcount)
	// os.Stdout.Sync()  
	if !seen {
		neighborsLock.RLock()
		defer neighborsLock.RUnlock()
		for neighbor := range neighbors {
			// fmt.Printf("requesting %v %v\n", neighbor, p)
			// os.Stdout.Sync()  
			requestPause(neighbor, p)
		}
		SlowpokeDoDelay(delayToDo)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// fmt.Printf("New connection from %s\n", conn.RemoteAddr())

	// Create a buffered reader
	reader := bufio.NewReader(conn)

	for {
		// Read until newline (assuming JSON objects are newline-delimited)
		data, err := reader.ReadBytes('\n')
		if err != nil {
			//fmt.Printf("Connection %s closed: %v\n", conn.RemoteAddr(), err)
			return
		}

		// Parse JSON
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			fmt.Printf("Invalid JSON from %s: %v\n", conn.RemoteAddr(), err)
			continue
		}

		switch msg.Type {
		case "register":
			var regReq RegisterReq
			if err := json.Unmarshal(msg.Payload, &regReq); err != nil {
				fmt.Printf("Invalid JSON from %s: %v\n", conn.RemoteAddr(), err)
				continue
			}
			handleRegister(regReq)
		case "pause":
			var pauseReq PauseReq
			if err := json.Unmarshal(msg.Payload, &pauseReq); err != nil {
				fmt.Printf("Invalid JSON from %s: %v\n", conn.RemoteAddr(), err)
				continue
			}
			handlePause(pauseReq)
		default:
			fmt.Printf("Invalid JSON from %s: %v\n", conn.RemoteAddr(), err)
			continue
		}
	}
}

func startControlServer() {
	listener, err := net.Listen("tcp", ":5550")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("TCP JSON server listening on :5550")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept error: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}