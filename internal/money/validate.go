package money

import "fmt"

// ValidatePositiveBaht ensures listing deposits and credit deltas are whole positive baht.
func ValidatePositiveBaht(amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive whole baht")
	}
	return nil
}
