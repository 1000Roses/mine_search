package models

type Resp struct {
	Status int
	Msg    string
	Detail interface{}
}

type RespLocal struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
