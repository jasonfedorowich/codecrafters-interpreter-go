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

func IsString(val any) (bool, string) {
	if v, ok := val.(string); ok {
		return true, v
	}
	return false, ""
}

func ParseStrings(strs ...any) []string {
	var s []string
	for _, str := range strs {
		if ok, v := IsString(str); ok {
			s = append(s, v)
		} else {
			return nil
		}
	}
	return s
}

func ParseNumbers(numbers ...any) []float64 {
	var nums []float64
	for _, n := range numbers {
		if ok, v := IsNumber(n); ok {
			nums = append(nums, v)
		} else {
			return nil
		}
	}
	return nums
}
