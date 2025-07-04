package graphics

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func Init() {
	fmt.Println("init graphics")

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Chip-8", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 640, 320, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	surface.FillRect(nil, 0)

	rect := sdl.Rect{0, 0, 200, 200}
	colour := sdl.Color{R: 255, G: 0, B: 255, A: 255}
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(&rect, pixel)
	window.UpdateSurface()
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent: // NOTE: Please use `*sdl.QuitEvent` for `v0.4.x` (current version).
				println("Quit")
				running = false
				break
			default:
				fmt.Println("evento:", event)
			}
		}

		sdl.Delay(33)
	}
}
