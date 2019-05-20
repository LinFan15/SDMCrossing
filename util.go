package main

func stringInSlice(s string, slice []string) bool {
	for _, x := range slice {
		if s == x {
			return true
		}
	}
	return false
}

func indexOfStringInSlice(s string, slice []string) int {
	for index, x := range slice {
		if s == x {
			return index
		}
	}
	return -1
}

func delFromSlice(index int, slice []string) []string {
	slice = append(slice[:index], slice[index+1:]...)
	return slice
}
