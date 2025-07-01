package chip8

import (
	"fmt"
)

const START_ADDRESS = 0x200

type Chip8 interface {
	LoadRom(rom []byte)
	DumpMemory()
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

func (c *chip8) LoadRom(rom []byte) {
	for i, v := range rom {
		c.memory[START_ADDRESS+i] = v
	}

	fmt.Println("Loaded rom into memory")
}

func (c *chip8) DumpMemory() {
	fmt.Println(c.memory)
}
