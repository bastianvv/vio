package dto

func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
