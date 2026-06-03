package mask

import "strings"

func BankAccountNumber(num string) string {
	num = strings.TrimSpace(num)
	if num == "" {
		return ""
	}
	if len(num) <= 4 {
		return strings.Repeat("*", len(num))
	}
	return strings.Repeat("*", len(num)-4) + num[len(num)-4:]
}

func NationalID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return ""
	}
	if len(id) <= 4 {
		return strings.Repeat("*", len(id))
	}
	return strings.Repeat("*", len(id)-4) + id[len(id)-4:]
}

func Phone(tel string) string {
	tel = strings.TrimSpace(tel)
	if tel == "" {
		return ""
	}
	if len(tel) <= 4 {
		return strings.Repeat("*", len(tel))
	}
	return tel[:2] + strings.Repeat("*", len(tel)-4) + tel[len(tel)-2:]
}
