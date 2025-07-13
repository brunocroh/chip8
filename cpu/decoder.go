package cpu

import (
	"math"
	"math/rand"
)

func (c *chip8) fetchOpcode() uint16 {
	return uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
}

func (c *chip8) decodeExecute(opcode uint16) {
	nnn := opcode & 0x0FFF
	kk := uint8(opcode & 0x00FF)
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := opcode & 0x000F

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x00F {
		case 0x0000:
			c.instructions.cls(c)
		case 0x000E:
			c.instructions.ret(c)
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
			if bit&0x01 == 1 {
				c.register[0xF] = 1
			} else {
				c.register[0xF] = 0
			}
		}
	// 9xy0 - SNE Vx, Vy
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
	// Cxkk - RND Vx, byte
	case 0xC000:
		random := rand.Intn(255)
		c.register[x] = uint8(random) & kk
	// Dxyn - DRW Vx, Vy, nibble
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
		c.DrawFlag = true
	case 0xE000:
		switch opcode & 0x00FF {
		// Ex9E - SKP VX
		case 0x009E:
			if c.keypad[c.register[x]] == 1 {
				c.pc += 2
			}
		// ExA1 - SKNP VX
		case 0x00A1:
			if c.keypad[c.register[x]] == 0 {
				c.pc += 2
			}
		}
	case 0xF000:
		switch opcode & 0x00FF {
		// Fx07 - LD Vx, DT
		case 0x0007:
			c.register[x] = c.delayTimer
		// Fx0A - LD Vx, K
		case 0x000A:
			keyFound := false
			for _, v := range c.keypad {
				if v == 1 {
					keyFound = true
					break
				}
			}

			if !keyFound {
				c.pc -= 2
			}
		// Fx15 - LD DT, Vx
		case 0x0015:
			c.delayTimer = c.register[x]
		// Fx18 - LD ST, Vx
		case 0x0018:
			c.soundTimer = c.register[x]
		// Fx1E - ADD I, Vx
		case 0x001E:
			c.index = c.index + uint16(c.register[x])
		// Fx29 - LD F, Vx
		case 0x0029:
			c.index = uint16(c.register[x])
		// Fx33 - LD B, Vx
		case 0x0033:
			number := c.register[x]
			c.memory[c.index] = number / 100
			c.memory[c.index+1] = (number % 100) / 10
			c.memory[c.index+2] = (number % 100) % 10
		// Fx55 - LD [I], Vx
		case 0x0055:
			for i := uint16(0); i <= x; i++ {
				c.memory[c.index+i] = c.register[i]
			}
			c.index += x + 1
		// Fx65 - LD Vx, [I]
		case 0x0065:
			for i := uint16(0); i <= x; i++ {
				c.register[i] = c.memory[c.index+i]
			}
			c.index += x + 1
		}
	}
}
