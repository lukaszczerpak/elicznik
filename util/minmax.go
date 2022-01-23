package util

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Foat64Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func Float64Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
