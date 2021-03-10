package iso639_2

func SliceContainsString(slice []string, v string) bool {
	for _, el := range slice {
		if el == v {
			return true
		}
	}
	return false
}
