package exception

type Exception struct {
	Code   string
	Detail string
}

func Set(code, detail string) Exception {
	return Exception{
		Code:   code,
		Detail: detail,
	}
}
