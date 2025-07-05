package chip8

import (
	"fmt"
	"strconv"
)

const START_ADDRESS = 0x200
const FONTSET_START_ADDRESS = 0x50
const FONTSET_SIZE = 80

var fontset = [FONTSET_SIZE]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type Chip8 interface {
	LoadRom(rom []byte)
	Init()
	Quit()
	DumpMemory()
	Cycle()
	Clear()
	updateTimers()
}

type chip8 struct {
	register   [16]uint8    //V0-VF registers
	memory     [4096]uint8  //4kb of memory
	index      uint16       // index register
	pc         uint16       // Program counter
	stack      [16]uint16   // Stack for storing retunr address
	sp         uint8        //Stack pointer
	delayTimer uint8        // Delay timer
	soundTimer uint8        // Delay timer
	keypad     [16]uint8    // Keypad state
	Video      [2048]uint32 // Display buffer
	opcode     uint16       //Current opcode
	DrawFlag   bool         //Draw flag
}

func NewChip8() *chip8 {
	return &chip8{
		register:   [16]uint8{},
		memory:     [4096]uint8{},
		index:      0,
		pc:         0,
		stack:      [16]uint16{},
		sp:         0,
		delayTimer: 0,
		soundTimer: 0,
		keypad:     [16]uint8{},
		Video:      [2048]uint32{}, //64*32
		opcode:     0,
		DrawFlag:   false,
	}
}

func (c *chip8) Init() {
	c.pc = START_ADDRESS

	for i, v := range fontset {
		c.memory[FONTSET_START_ADDRESS+i] = v
	}

	c.Clear()
}

func (c *chip8) LoadRom(rom []byte) {
	for i, v := range rom {
		c.memory[START_ADDRESS+i] = v
	}

	fmt.Println("Loaded rom into memory")
}

func (c *chip8) DumpMemory() {
	fmt.Printf("% x\n", c.memory)
}

func (c *chip8) fetchOpcode() uint16 {
	return uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
}

func (c *chip8) Clear() {
	for i := range c.Video {
		c.Video[i] = 0
	}
	c.DrawFlag = true
}

func (c *chip8) Cycle() {
	opcode := c.fetchOpcode()
	fmt.Printf("opcode: hex: 0x%s bin: %s \n", strconv.FormatInt(int64(opcode), 16), strconv.FormatInt(int64(opcode), 2))
	switch opcode & 0xF000 {
	case 0xA000:
		c.index = opcode & 0x0FFF
		c.pc += 2
		break
	case 0x1000:
		c.pc = opcode & 0x0FFF
		break
	case 0x6000:
		c.register[uint8(opcode&0x0F00>>8)] = uint8(opcode & 0x00FF)
		c.pc += 2
		break
	case 0x7000:
		c.register[uint8(opcode&0x0F00>>8)] += uint8(opcode & 0x0FF)
		c.pc += 2
		break
	// Dxyn
	case 0xD000:
		x := uint16(c.register[(opcode&0x0F00)>>8])
		y := uint16(c.register[(opcode&0x00F0)>>4])
		n := opcode & 0x000F

		c.register[0xF] = 0

		for yLine := uint16(0); yLine < n; yLine++ {
			pixel := c.memory[c.index+yLine]
			for xLine := uint16(0); xLine < 8; xLine++ {
				if (pixel & (0x80 >> xLine)) != 0 {
					xPos := (x + xLine) % 64
					yPos := (y + yLine) % 32
					screenPos := xPos + (yPos * 64)
					if c.Video[screenPos] == 1 {
						c.register[0xF] = 1
					}

					c.Video[screenPos] ^= 1
				}

			}

		}
		c.index = opcode & 0x0FFF
		c.DrawFlag = true
		c.pc += 2
		break
	case 0x0000:
		switch opcode & 0x00F {
		case 0x0000:
			c.Clear()
			c.pc += 2
			break
		case 0x000E:
			fmt.Println("Return from subroutine")
			break
		default:
			// c.pc += 2
			fmt.Println("Unknown opcode [0x0000]: 0x", strconv.FormatInt(int64(opcode), 16))
		}
		break
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			fmt.Println("E-9e")
			break
		case 0x00A1:
			fmt.Println("E-A1")
			break
		default:
			fmt.Println("E opcode not found")
		}
		break
	case 0xF000:
		fmt.Println("F")
		break
	default:
		fmt.Println("NOT HANDLED OPCODE: ", strconv.FormatInt(int64(opcode), 16))
	}
	c.updateTimers()
}

func (c *chip8) Quit() {
	fmt.Println("Chip-8 Quit")
}

func (c *chip8) updateTimers() {
	if c.delayTimer > 0 {
		c.delayTimer = c.delayTimer - 1
	}

	if c.soundTimer > 0 {
		c.soundTimer = c.soundTimer - 1
		if c.soundTimer == 1 {
			fmt.Println("BEEP")

		}
	}
}
