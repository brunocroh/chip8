package cpu

import "math"

type Instructions interface {
	NewInstructions() Instructions
}

type instructions struct {
}

func NewInstructions() *instructions {
	return &instructions{}
}

/*
00E0 - CLS

Clear the display.
*/
func (m *instructions) cls(c *chip8) {
	for i := range c.Video {
		c.Video[i] = 0
	}
	c.DrawFlag = true
}

/*
00EE - RET

Return from a subroutine.

The interpreter sets the program counter to the address at the top of the stack, then subtracts 1 from the stack pointer.
*/
func (m *instructions) ret(c *chip8) {
	c.pc = c.stack[c.sp]
	c.sp -= 1
}

/*
1nnn - JP addr

Jump to location nnn.

The interpreter sets the program counter to nnn.
*/
func (m *instructions) jump(c *chip8, nnn uint16) {
	c.pc = nnn
}

/*
2nnn - CALL addr

The interpreter sets the program counter to nnn.
*/
func (m *instructions) callSubroutine(c *chip8, nnn uint16) {
	c.sp += 1
	c.stack[c.sp] = c.pc
	c.pc = nnn
}

/*
3xkk - SE Vx, byte

Skip next instruction if Vx = kk.

The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
*/
func (m *instructions) seVxKk(c *chip8, x uint16, kk uint8) {
	if c.register[x] == kk {
		c.pc += 2
	}
}

/*
4xkk - SNE Vx, byte

Skip next instruction if Vx != kk.

The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
*/
func (m *instructions) SneVxKk(c *chip8, x uint16, kk uint8) {
	if c.register[x] != kk {
		c.pc += 2
	}
}

/*
5xy0 - SE Vx, Vy

Skip next instruction if Vx = Vy.

The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
*/
func (m *instructions) seVxVy(c *chip8, x uint16, y uint16) {
	if c.register[x] == c.register[y] {
		c.pc += 2
	}
}

/*
6xkk - LD Vx, byte

Set Vx = kk.

The interpreter puts the value kk into register Vx.
*/
func (m *instructions) loadVx(c *chip8, x uint16, kk uint8) {
	c.register[x] = kk
}

/*
7xkk - ADD Vx, byte

Set Vx = Vx + kk.

Adds the value kk to the value of register Vx, then stores the result in Vx.
*/
func (m *instructions) addVx(c *chip8, x uint16, kk uint8) {
	c.register[x] += kk
}

/*
8xy0 - LD Vx, Vy

Set Vx = Vy.

Stores the value of register Vy in register Vx.
*/
func (m *instructions) loadVyIntoVx(c *chip8, x uint16, y uint16) {
	c.register[x] = c.register[y]
}

/*
8xy1 - OR Vx, Vy

Set Vx = Vx OR Vy.

Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx. A bitwise OR compares the corrseponding bits from two values, and if either bit is 1, then the same bit in the result is also 1. Otherwise, it is 0.
*/
func (m *instructions) orVxVy(c *chip8, x uint16, y uint16) {
	c.register[x] = c.register[x] | c.register[y]
}

/*
8xy2 - AND Vx, Vy

Set Vx = Vx AND Vy.

Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx. A bitwise AND compares the corrseponding bits from two values, and if both bits are 1, then the same bit in the result is also 1. Otherwise, it is 0.
*/
func (m *instructions) andVxVy(c *chip8, x uint16, y uint16) {
	c.register[x] = c.register[x] & c.register[y]
}

/*
8xy3 - XOR Vx, Vy

Set Vx = Vx XOR Vy.

Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx. An exclusive OR compares the corrseponding bits from two values, and if the bits are not both the same, then the corresponding bit in the result is set to 1. Otherwise, it is 0.
*/
func (m *instructions) xorVxVy(c *chip8, x uint16, y uint16) {
	c.register[x] = c.register[x] ^ c.register[y]
}

/*
8xy4 - ADD Vx, Vy

Set Vx = Vx + Vy, set VF = carry.

The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,) VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
*/
func (m *instructions) addVxVy(c *chip8, x uint16, y uint16) {
	sum := uint16(c.register[x]) + uint16(c.register[y])

	c.register[x] = uint8(sum)
	if sum > math.MaxUint8 {
		c.register[0xF] = 1
	} else {
		c.register[0xF] = 0
	}
}

/*
8xy5 - SUB Vx, Vy

Set Vx = Vx - Vy, set VF = NOT borrow.

If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
*/
func (m *instructions) subVxVy(c *chip8, x uint16, y uint16) {
	originalVX := c.register[x]
	c.register[x] = originalVX - c.register[y]
	if originalVX >= c.register[y] {
		c.register[0xF] = 1
	} else {
		c.register[0xF] = 0
	}
}

/*
8xy6 - SHR Vx

Set Vx = Vx SHR 1.

If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
*/
func (m *instructions) shrVx(c *chip8, x uint16) {
	bit := c.register[x]

	c.register[x] = c.register[x] >> 1
	if bit&0x01 == 1 {
		c.register[0xF] = 1
	} else {
		c.register[0xF] = 0
	}
}
