//go:build js && wasm

package main

import (
	"brunocroh/chip8/cpu"
	"fmt"
	"os"
	"syscall/js"
	"time"
	"unsafe"
)

var (
	keepRunning bool = true
	romLoaded   bool = false
	chip8       *cpu.Chip8
	drawFunc    js.Value
)

func main() {
	js.Global().Set("loadRom", js.FuncOf(loadRomJS))
	js.Global().Set("start", js.FuncOf(startJS))

	chip8 = cpu.NewChip8()
	chip8.Init()

	select {}
}

func listenKeypad() {
	fmt.Println("map keyboard")
}

func loadRomJS(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return nil
	}

	uints8Array := args[0]
	length := uints8Array.Get("length").Int()

	rom := make([]byte, length)

	js.CopyBytesToGo(rom, uints8Array)

	chip8.Reset()
	chip8.LoadRom(rom)
	romLoaded = true
	return nil
}

func startJS(this js.Value, args []js.Value) interface{} {
	renderCb := args[0]
	videoMemory := args[1]

	lastTimerUpdate := time.Now()
	timerInterval := time.Second / 60

	lastCycle := time.Now()
	cycleInterval := time.Second / 500

	var emulatorLoop js.Func
	emulatorLoop = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !keepRunning {
			os.Exit(1)
		}

		now := time.Now()

		if now.Sub(lastTimerUpdate) >= timerInterval {
			chip8.UpdateTimers()
			fmt.Println("updateTimer")
			lastTimerUpdate = now
		}

		if now.Sub(lastCycle) >= cycleInterval {
			chip8.Cycle()

			if chip8.DrawFlag() {
				chip8.SetDrawFlag(false)
				fmt.Println("render")
				video := chip8.GetVideo()

				videoBytes := (*[2048 * 4]byte)(unsafe.Pointer(&video[0]))[:2048*4]
				js.CopyBytesToJS(videoMemory, videoBytes)

				renderCb.Invoke()
			}

			lastCycle = now
		}

		// listenKeypad()
		js.Global().Call("requestAnimationFrame", emulatorLoop)
		return nil

	})

	js.Global().Call("requestAnimationFrame", emulatorLoop)
	return nil
}
