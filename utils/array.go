package utils

func ArrayStringIncludes(arr []string, key string) bool {
	for i := range arr {
		if arr[i] == key {
			return true
		}
	}
	return false
}
