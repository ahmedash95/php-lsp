package util

func CalculatePercentage(value int, total int) int {
	if total == 0 {
		return 0
	}

	percentage := (float64(value) / float64(total)) * 100.0

	return int(percentage)
}
