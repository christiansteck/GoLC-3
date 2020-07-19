package main

import (
	"bufio"
	"os"
)

func GetChar() (uint16, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	return uint16([]byte(input)[0]), nil
}
