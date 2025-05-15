package util

// boolToStr converts a boolean to "true" or "false"
func BoolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
