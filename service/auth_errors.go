package service

import "github.com/rnikrozoft/pramool-core/exception"

func errPasswordRequired() error {
	return exception.Unauthorized("กรุณากรอกรหัสผ่าน")
}

func errInvalidCredentials() error {
	return exception.Unauthorized("เบอร์โทรศัพท์ อีเมล หรือรหัสผ่านไม่ถูกต้อง")
}

func errAccountSuspended() error {
	return exception.Forbidden("บัญชีถูกระงับการใช้งาน กรุณาติดต่อผู้ดูแลระบบ")
}
