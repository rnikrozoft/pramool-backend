package exception

type Response struct {
	Code int `json:"code"`
}

func Set(code int) Response {
	return Response{
		Code: code,
	}
}
