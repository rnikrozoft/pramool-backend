package privacy

import "strings"

func IsValidThaiNationalID(id string) bool {
	id = strings.TrimSpace(id)
	if len(id) != 13 {
		return false
	}
	for _, ch := range id {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += int(id[i]-'0') * (13 - i)
	}
	checkDigit := (11 - (sum % 11)) % 10
	return checkDigit == int(id[12]-'0')
}
