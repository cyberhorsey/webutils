package webutils

func ContainsString(slice []string, contains string) bool {
	for _, value := range slice {
		if value == contains {
			return true
		}
	}

	return false
}
