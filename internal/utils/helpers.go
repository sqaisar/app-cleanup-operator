package utils

func HasFinalizer(finalizers []string, target string) bool {
	for _, f := range finalizers {
		if f == target {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, s string) []string {
	result := []string{}
	for _, item := range slice {
		if item != s {
			result = append(result, item)
		}
	}
	return result
}
