package utils

import (
	"fmt"
	"os"
)

func LoadRom() ([]byte, error) {
	data, err := os.ReadFile("./roms/IBM_LOGO.ch8")

	if err != nil {
		fmt.Println("Fail to read the rom", err)
		return nil, err

	}

	return data, nil
}
