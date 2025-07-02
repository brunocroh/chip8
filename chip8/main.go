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
	DumpMemory()
	Cycle()
	DrawFlag() bool
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
	video      [2048]uint32 // Display buffer
	opcode     uint16       //Current opcode
}

func NewChip8() Chip8 {
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
		video:      [2048]uint32{}, //64*32
		opcode:     0,
	}
}

func (c *chip8) Init() {
	c.pc = START_ADDRESS
	c.memory[c.pc] = 0xA2
	c.memory[c.pc+1] = 0xF0

	for i, v := range fontset {
		c.memory[FONTSET_START_ADDRESS+i] = v
	}

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

func (c *chip8) Cycle() {
	opcode := c.fetchOpcode()
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x00F {
		case 0x0000:
			fmt.Println("CLEAR")
			break
		case 0x000E:
			fmt.Println("Return from subroutine")
			break
		default:
			fmt.Println("Unknown opcode [0x0000]: 0x", strconv.FormatInt(int64(opcode), 16))
		}
		break
	default:
		// fmt.Println("NOT HANDLED OPCODE: ", strconv.FormatInt(int64(opcode), 16))
	}

	c.pc += 2
	c.updateTimers()
}

func (c *chip8) DrawFlag() bool {
	return false
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
