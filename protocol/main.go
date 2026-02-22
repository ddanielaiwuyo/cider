package protocol

type UserId uint64

type Request struct {
	UserId UserId `json:"userId"`
	Msg    string `json:"msg"`
}

type Response struct {
	From UserId
	Msg  string
}
