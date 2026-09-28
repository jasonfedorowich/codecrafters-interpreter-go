package types

func IsBoolean(val any) (bool, bool) {
	if val == nil {
		return true, false
	}
	if v, ok := val.(bool); ok {
		return true, v
	}
	return false, false
}

func IsNumber(val any) (bool, float64) {
	if v, ok := val.(float64); ok {
		return true, v
	}
	return false, 0
}
