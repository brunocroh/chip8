//go:build !js && !wasm
// +build !js,!wasm

package main

import (
	"fmt"
	"github.com/brunocroh/chip8/cpu"
	"github.com/brunocroh/chip8/utils"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var keepRunning bool = true

func main() {
	ticker := time.NewTicker(time.Second / 500)
	defer ticker.Stop()
	romPath := os.Args[1:]
	fmt.Println("Initiliaze rom:", romPath)
	time.Sleep(500 * time.Millisecond)

	rom, err := utils.LoadRom(romPath[0])

	if err != nil {
		fmt.Println("Fail to load rom")
		return
	}

	chip8 := cpu.NewChip8()
	chip8.Init()
	chip8.LoadRom(rom)

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

	go timersThread(chip8)

	for range ticker.C {
		if !keepRunning {
			os.Exit(1)
		}
		chip8.Cycle()
		if chip8.DrawFlag() {
			chip8.SetDrawFlag(false)
			renderer.SetDrawColor(255, 0, 0, 255)
			renderer.Clear()

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
		}
		listenKeypad(chip8)
	}
}

func listenKeypad(chip8 *cpu.Chip8) {
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
			keepRunning = false
		}
	}

}

func timersThread(chip8 *cpu.Chip8) {
	tickerTimers := time.NewTicker(time.Second / 60)
	for range tickerTimers.C {
		chip8.UpdateTimers()
	}
}
