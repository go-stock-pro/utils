package typecasting

import "strconv"

func ConvertIntToString(in int) string {
	return strconv.Itoa(in)
}

func ConvertInt64ToString(in int64) string {
	return strconv.FormatInt(in, 10)
}

func ConvertFloat32ToString(in float32) string {
	return strconv.FormatFloat(float64(in), 'f', -1, 32)
}

func ConvertFloat64ToString(in float64) string {
	return strconv.FormatFloat(in, 'f', -1, 64)
}

func ConvertBoolToString(in bool) string {
	return strconv.FormatBool(in)
}

func ConvertBoolToStringYesNo(in bool) string {
	if in {
		return "true"
	}

	return "false"
}
