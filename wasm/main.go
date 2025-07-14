package main

import (
	"brunocroh/chip8/cpu"
	"fmt"
	"os"
	"syscall/js"
	"time"
)

var keepRunning bool = true
var romLoaded bool = false
var chip8 cpu.Chip8

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

func timersThread(chip8 cpu.Chip8) {
	tickerTimers := time.NewTicker(time.Second / 60)
	for range tickerTimers.C {
		fmt.Println("timers update")
		chip8.UpdateTimers()
	}
}

func loadRomJS(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return nil
	}

	uints8Array := args[0]
	length := uints8Array.Get("length").Int()

	rom := make([]byte, length)

	js.CopyBytesToGo(rom, uints8Array)

	chip8.LoadRom(rom)
	romLoaded = true
	return nil
}

func startJS(this js.Value, args []js.Value) interface{} {
	ticker := time.NewTicker(time.Second / 500)
	defer ticker.Stop()

	go timersThread(chip8)

	for range ticker.C {
		if !keepRunning {
			os.Exit(1)
		}
		chip8.Cycle()
		if chip8.DrawFlag() {
			chip8.SetDrawFlag(true)
			fmt.Println("DRAW")
		}
		// listenKeypad()
	}

	return nil
}
