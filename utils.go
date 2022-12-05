package go_db

func ArrayInBool(slice []string, val string) bool {
	_, result := ArrayIn(slice, val)
	return result
}

func ArrayIn(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
