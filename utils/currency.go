package utils

const (
	IDR = "IDR"
	USD = "USD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case IDR, USD:
		return true
	}
	return false
}
