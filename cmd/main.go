package main

import (
	"brunocroh/chip8/chip8"
	"brunocroh/chip8/utils"
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	ticker := time.NewTicker(time.Second / 700)
	defer ticker.Stop()
	romPath := os.Args[1:]
	fmt.Println("start")

	fmt.Println("Initiliaze rom:", romPath)
	time.Sleep(500 * time.Millisecond)

	rom, err := utils.LoadRom(romPath[0])

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

	for range ticker.C {
		chip8.Cycle()
		if chip8.DrawFlag {
			chip8.DrawFlag = false
			renderer.SetDrawColor(255, 0, 0, 255)
			renderer.Clear()

			// Get the display buffer and render
			start := time.Now()
			for i, v := range chip8.Video {
				if v != 0 {
					renderer.SetDrawColor(255, 255, 255, 255)
				} else {
					renderer.SetDrawColor(0, 0, 0, 255)
				}

				renderer.FillRect(&sdl.Rect{
					Y: int32(i/64) * 16,
					X: int32(i%64) * 16,
					W: 16,
					H: 16,
				})
			}

			renderer.Present()
			duration := time.Since(start)
			fmt.Println("duration of render:", duration)
		}
		listenKeypad(chip8)
	}
}

func listenKeypad(chip8 chip8.Chip8) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch et := event.(type) {
		case *sdl.KeyboardEvent:
			if et.Type == sdl.KEYDOWN || et.Type == sdl.KEYUP {
				var ev uint8

				if et.Type == sdl.KEYDOWN {
					ev = 1
				} else {
					ev = 0
				}

				switch et.Keysym.Sym {
				case sdl.K_1:
					chip8.OnKeyEvent(0x1, ev)
				case sdl.K_2:
					chip8.OnKeyEvent(0x2, ev)
				case sdl.K_3:
					chip8.OnKeyEvent(0x3, ev)
				case sdl.K_4:
					chip8.OnKeyEvent(0xC, ev)
				case sdl.K_q:
					chip8.OnKeyEvent(0x4, ev)
				case sdl.K_w:
					chip8.OnKeyEvent(0x5, ev)
				case sdl.K_e:
					chip8.OnKeyEvent(0x6, ev)
				case sdl.K_r:
					chip8.OnKeyEvent(0xD, ev)
				case sdl.K_a:
					chip8.OnKeyEvent(0x7, ev)
				case sdl.K_s:
					chip8.OnKeyEvent(0x8, ev)
				case sdl.K_d:
					chip8.OnKeyEvent(0x9, ev)
				case sdl.K_f:
					chip8.OnKeyEvent(0xE, ev)
				case sdl.K_z:
					chip8.OnKeyEvent(0xA, ev)
				case sdl.K_x:
					chip8.OnKeyEvent(0x0, ev)
				case sdl.K_c:
					chip8.OnKeyEvent(0xB, ev)
				case sdl.K_v:
					chip8.OnKeyEvent(0xF, ev)
				}
			}
		case *sdl.QuitEvent:
			println("Quit")
		}
	}

}
