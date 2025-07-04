package main

import (
	"brunocroh/chip8/chip8"
	"brunocroh/chip8/utils"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	fmt.Println("start")

	rom, err := utils.LoadRom()

	if err != nil {
		fmt.Println("fail to load rom")
		return
	}

	chip8 := chip8.NewChip8()

	chip8.Init()
	defer chip8.Quit()

	chip8.LoadRom(rom)
	chip8.DumpMemory()

	// SDL - Graphics
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Chip-8", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 1024, 512, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	surface.FillRect(nil, 255)

	keepRunning := true
	for keepRunning {
		chip8.Cycle()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				keepRunning = false
				break
			default:
				fmt.Println("evento:", event)
			}
		}

		sdl.Delay(100)
	}
}
