package utils

import (
	"bufio"
	"os"
)

func ParseAddressFile(fileName string) []string {
	var addresses []string
	file, _ := os.Open(fileName)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		addresses = append(addresses, scanner.Text())
	}

	return addresses
}
