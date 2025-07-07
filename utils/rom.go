package utils

import (
	"fmt"
	"os"
)

func LoadRom() ([]byte, error) {
	// data, err := os.ReadFile("./roms/IBM_LOGO.ch8")
	// data, err := os.ReadFile("./roms/tetris.ch8")
	data, err := os.ReadFile("./roms/test_opcode.ch8")

	if err != nil {
		fmt.Println("Fail to read the rom", err)
		return nil, err

	}

	return data, nil
}
