package collector

func IndexInStringArray(str string, arr []string) int {
	for i := 0; i < len(arr); i++ {
		if str == arr[i] {
			return i
		}
	}
	return -1
}
