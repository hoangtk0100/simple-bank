package util

// Constants for all supported currencies
const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	JPY = "JPY"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, JPY:
		return true
	}

	return false
}
