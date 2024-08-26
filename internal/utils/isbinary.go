package utils

import (
	"os"
	"unicode"
)

func IsBinary(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}

	buffer := make([]byte, 512) // Read first 512 bytes
	n, err := file.Read(buffer)
	if err != nil {
		return false, err
	}

	// Loop through the bytes and check if they're printable
	for i := 0; i < n; i++ {
		if buffer[i] == 0 { // Null byte is a good indicator of a binary file
			return true, nil
		}
		r := rune(buffer[i])
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return true, nil
		}
	}
	return false, nil
}
