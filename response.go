package douyin_security

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Success() *Response {
	return &Response{
		Code: 200,
		Msg:  "success",
	}
}

func Fail(code int, msg string) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
	}
}
