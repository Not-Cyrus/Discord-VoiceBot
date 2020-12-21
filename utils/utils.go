package utils

import (
	"bufio"
	"os"
)

func ReadLines(path string) []string {
	var Lines []string
	file, _ := os.Open(path)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Lines = append(Lines, scanner.Text())
	}
	return Lines
}
