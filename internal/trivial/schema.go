package trivial

type TrivialRequest struct {
	Q string `json:"q"`
}

type TrivialResponse struct {
	A string `json:"a"`
}

//type WriteRequest struct {
//	K string `json:"k"`
//	V string `json:"v"`
//}