package main

import (
	"brunocroh/chip8/chip8"
	"brunocroh/chip8/utils"
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	var pixels [2048]uint32
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

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Untitled", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 1024, 512, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	sdlTexture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ARGB8888,
		sdl.TEXTUREACCESS_STREAMING,
		64,
		32,
	)
	if err != nil {
		panic(err)
	}
	defer sdlTexture.Destroy()

	keepRunning := true
	for keepRunning {
		chip8.Cycle()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				keepRunning = false
			}
		}

		if chip8.DrawFlag {
			chip8.DrawFlag = false

			for i, pixel := range chip8.Video {
				pixels[i] = (0x00FFFFFF * pixel) | 0xFF000000
			}

			err := sdlTexture.Update(nil, unsafe.Pointer(&pixels[0]), 64*int(unsafe.Sizeof(uint32(0))))
			if err != nil {
				panic(err)
			}

			renderer.Clear()
			renderer.Copy(sdlTexture, nil, nil)
			renderer.Present()
		}

		// 1700/s
	}
}
