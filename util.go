package main

func stringInSlice(s string, slice []string) bool {
	for _, x := range slice {
		if s == x {
			return true
		}
	}
	return false
}
