package main

import (
	"brunocroh/chip8/chip8"
	"brunocroh/chip8/utils"
	"fmt"
)

func main() {
	fmt.Println("start")

	rom, err := utils.LoadRom()

	if err != nil {
		fmt.Println("fail to load rom")
		return
	}

	chip8 := chip8.NewChip8()

	chip8.LoadRom(rom)
	chip8.DumpMemory()

}
