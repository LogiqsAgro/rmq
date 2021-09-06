package cmd

import "os"

func writeError(err error) bool {
	if err != nil {
		os.Stderr.WriteString("ERROR: " + err.Error())
		return true
	}
	return false
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
