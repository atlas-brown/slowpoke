package synthetic

type ConfigMap struct {
	Processes int        `json:"processes"`
	Logging   bool       `json:"logging"`
	Protocol  string     `json:"protocol"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Name              string             `json:"name"`
	ExecutionMode     string             `json:"execution_mode"`
	CpuComplexity     *CpuComplexity     `json:"cpu_complexity,omitempty"`
	NetworkComplexity *NetworkComplexity `json:"network_complexity,omitempty"`
}

type CalledService struct {
	Service             string `json:"service"`
	Port                int    `json:"port"`
	Endpoint            string `json:"endpoint"`
	Protocol            string `json:"protocol"`
	TrafficForwardRatio int    `json:"traffic_forward_ratio"`
	RequestPayloadSize  int    `json:"request_payload_size"`
	Probability         int    `json:"probability"`
}

type CpuComplexity struct {
	ExecutionTime float32 `json:"execution_time"`
	Threads       int     `json:"threads"`
}

type NetworkComplexity struct {
	ForwardRequests     string          `json:"forward_requests"`
	ResponsePayloadSize int             `json:"response_payload_size"`
	CalledServices      []CalledService `json:"called_services"`
}

type Response struct {
	Message  string `json:"message"`
	Endpoint string `json:"endpoint"`
}