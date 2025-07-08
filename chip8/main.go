package chip8

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
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

var teste []uint16

type Chip8 interface {
	LoadRom(rom []byte)
	Init()
	Quit()
	DumpMemory()
	Cycle()
	Clear()
	updateTimers()
	incrementCounter()
	OnKeyEvent(id uint8, down bool)
}

type chip8 struct {
	register   [16]uint8    // V0-VF registers
	memory     [4096]uint8  // 4kb of memory
	index      uint16       // index register
	pc         uint16       // Program counter
	stack      [16]uint16   // Stack for storing retunr address
	sp         uint8        // Stack pointer
	delayTimer uint8        // Delay timer
	soundTimer uint8        // Delay timer
	keypad     [16]uint8    // Keypad state
	Video      [2048]uint32 // Display buffer
	opcode     uint16       // Current opcode
	DrawFlag   bool         // Draw flag
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

func (c *chip8) incrementCounter() {
	c.pc += 2
}

func (c *chip8) Clear() {
	for i := range c.Video {
		c.Video[i] = 0
	}
	c.DrawFlag = true
}

func uniqueOpcodes(opcode uint16) bool {
	return slices.Contains(teste, opcode)
}

func (c *chip8) Cycle() {
	opcode := c.fetchOpcode()
	if !uniqueOpcodes(opcode) {
		teste = append(teste, opcode)
	}

	nnn := opcode & 0x0FFF
	kk := uint8(opcode & 0x00FF)
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := opcode & 0x000F
	c.incrementCounter()

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x00F {
		// 00E0 - CLS
		case 0x0000:
			c.Clear()
		// 00EE - RET
		case 0x000E:
			c.pc = c.stack[c.sp]
			c.sp -= 1
		}
	// 1nnn - JP addr
	case 0x1000:
		c.pc = nnn
	// 2nnn - CALL addr
	case 0x2000:
		c.sp += 1
		c.stack[c.sp] = c.pc
		c.pc = nnn
	// 3xkk - SE Vx, byte
	case 0x3000:
		if c.register[x] == kk {
			c.pc += 2
		}
	// 4xkk - SNE Vx, byte
	case 0x4000:
		if c.register[x] != kk {
			c.pc += 2
		}
	// 5xy0 - SE Vx, Vy
	case 0x5000:
		if c.register[x] == c.register[y] {
			c.pc += 2
		}
	// 6xkk - LD Vx, byte
	case 0x6000:
		c.register[x] = kk
	case 0x7000:
		c.register[x] += kk
	case 0x8000:
		switch n {
		// 8xy0 - LD Vx, Vy
		case 0:
			c.register[x] = c.register[y]
		// 8xy1 - OR Vx, Vy
		case 1:
			c.register[x] = c.register[x] | c.register[y]
		// 8xy2 - AND Vx, Vy
		case 2:
			c.register[x] = c.register[x] & c.register[y]
		// 8xy3 - XOR Vx, Vy
		case 3:
			c.register[x] = c.register[x] ^ c.register[y]
		// 8xy4 - ADD Vx, Vy
		case 4:
			sum := uint16(c.register[x]) + uint16(c.register[y])

			c.register[x] = uint8(sum)
			if sum > math.MaxUint8 {
				c.register[0xF] = 1
			} else {
				c.register[0xF] = 0
			}

		// 8xy5 - SUB Vx, Vy
		case 5:
			originalVX := c.register[x]
			c.register[x] = originalVX - c.register[y]
			if originalVX >= c.register[y] {
				c.register[0xF] = 1
			} else {
				c.register[0xF] = 0
			}
		// 8xy6 - SHR Vx {, Vy}
		case 6:
			bit := c.register[x]

			c.register[x] = c.register[x] >> 1
			if bit&0x01 == 1 {
				c.register[0xF] = 1
			} else {
				c.register[0xF] = 0
			}

		// 8xy7 - SUBN Vx, Vy
		case 7:
			originalVX := c.register[x]
			c.register[x] = c.register[y] - c.register[x]
			if c.register[y] >= originalVX {
				c.register[0xF] = 1
			} else {
				c.register[0xF] = 0
			}
		// 8xyE - SHL, VX {, Vy}
		case 0xE:
			bit := c.register[x]

			c.register[x] = c.register[x] << 1
			c.register[0xF] = bit & 0x80 >> 7
		}
	case 0x9000:
		if c.register[x] != c.register[y] {
			c.pc += 2
		}
	// Annn
	case 0xA000:
		c.index = nnn
	// Bnnn
	case 0xB000:
		c.pc = nnn + uint16(c.register[0])
	case 0xC000:
		random := rand.Intn(255)
		c.register[x] = uint8(random) & kk
	// Dxyn
	case 0xD000:
		x := uint16(c.register[x])
		y := uint16(c.register[y])
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
		c.index = nnn
		c.DrawFlag = true
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			if c.keypad[c.register[x]] == 1 {
				c.pc += 2
			}
		case 0x00A1:
			if c.keypad[c.register[x]] == 0 {
				c.pc += 2
			}
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007:
			c.register[x] = c.delayTimer
		case 0x000A:
			var keyFound uint8
			for i, v := range c.keypad {
				if v == 1 {
					keyFound = uint8(i)
				}
			}

			if keyFound == 0 {
				c.pc -= 2
			}
		case 0x0015:
			c.delayTimer = c.register[x]
		case 0x0018:
			c.soundTimer = c.register[x]
		case 0x001E:
			c.index = c.index + uint16(c.register[x])
		case 0x0029:
			c.index = uint16(c.register[x])
		case 0x0033:
			number := c.register[x]
			c.memory[c.index] = number / 100
			c.memory[c.index+1] = (number % 100) / 10
			c.memory[c.index+2] = (number % 100) % 10
		case 0x0055:
			for i := uint16(0); i <= x; i++ {
				c.memory[c.index+i] = c.register[i]
			}
			c.index += x + 1
		case 0x0065:
			for i := uint16(0); i <= x; i++ {
				c.register[i] = c.memory[c.index+i]
			}
			c.index += x + 1
		}
	}
	c.updateTimers()
}

func (c *chip8) Quit() {
	fmt.Println("Chip-8 Quit")
	fmt.Println("=====")
	for _, v := range teste {
		fmt.Println(strconv.FormatInt(int64(v), 16))
	}
	fmt.Println("=====")
}

func (c *chip8) OnKeyEvent(key uint8, press uint8) {
	c.keypad[key] = press
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
