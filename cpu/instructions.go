package cpu

import (
	"math"
	"math/rand"
)

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
func (m *instructions) sneVxKk(c *chip8, x uint16, kk uint8) {
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

/*
8xy7 - SUBN Vx, Vy

Set Vx = Vy - Vx, set VF = NOT borrow.

If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
*/
func (m *instructions) subnVxVy(c *chip8, x uint16, y uint16) {
	originalVX := c.register[x]
	c.register[x] = c.register[y] - c.register[x]
	if c.register[y] >= originalVX {
		c.register[0xF] = 1
	} else {
		c.register[0xF] = 0
	}
}

/*
8xyE - SHL Vx {, Vy}

Set Vx = Vx SHL 1.

If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
*/
func (m *instructions) shlVx(c *chip8, x uint16) {
	bit := c.register[x]

	c.register[x] = c.register[x] << 1
	if bit&0x01 == 1 {
		c.register[0xF] = 1
	} else {
		c.register[0xF] = 0
	}
}

/*
9xy0 - SNE Vx, Vy

Skip next instruction if Vx != Vy.

The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
*/
func (m *instructions) sneVxVy(c *chip8, x uint16, y uint16) {
	if c.register[x] != c.register[y] {
		c.pc += 2
	}
}

/*
Annn - LD I, addr

Set I = nnn.

The value of register I is set to nnn.
*/
func (m *instructions) loadIndex(c *chip8, nnn uint16) {
	c.index = nnn
}

/*
Bnnn - JP V0, addr

Jump to location nnn + V0.

The program counter is set to nnn plus the value of V0.
*/
func (m *instructions) jumpV0(c *chip8, nnn uint16) {
	c.pc = nnn + uint16(c.register[0])
}

/*
Cxkk - RND Vx, byte

Set Vx = random byte AND kk.

The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk. The results are stored in Vx. See instruction 8xy2 for more information on AND.
*/
func (m *instructions) randonVxKk(c *chip8, x uint16, kk uint8) {
	random := rand.Intn(255)
	c.register[x] = uint8(random) & kk
}

/*
Dxyn - DRW Vx, Vy, nibble

Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.

The interpreter reads n bytes from memory, starting at the address stored in I. These bytes are then displayed as sprites on screen at coordinates (Vx, Vy). Sprites are XORed onto the existing screen. If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0. If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the opposite side of the screen. See instruction 8xy3 for more information on XOR, and section 2.4, Display, for more information on the Chip-8 screen and sprites.
*/
func (m *instructions) draw(c *chip8, x uint16, y uint16, n uint16) {
	vx := uint16(c.register[x])
	vy := uint16(c.register[y])
	c.register[0xF] = 0

	for yLine := uint16(0); yLine < n; yLine++ {
		pixel := c.memory[c.index+yLine]
		for xLine := uint16(0); xLine < 8; xLine++ {
			if (pixel & (0x80 >> xLine)) != 0 {
				xPos := (vx + xLine) % 64
				yPos := (vy + yLine) % 32
				screenPos := xPos + (yPos * 64)
				if c.Video[screenPos] == 1 {
					c.register[0xF] = 1
				}
				c.Video[screenPos] ^= 1
			}
		}
	}
	c.DrawFlag = true
}

/*
Ex9E - SKP Vx

Skip next instruction if key with the value of Vx is pressed.

Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
*/
func (m *instructions) skpVx(c *chip8, x uint16) {
	if c.keypad[c.register[x]] == 1 {
		c.pc += 2
	}
}

/*
ExA1 - SKNP Vx

Skip next instruction if key with the value of Vx is not pressed.

Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
*/
func (m *instructions) sknpVx(c *chip8, x uint16) {
	if c.keypad[c.register[x]] == 0 {
		c.pc += 2
	}
}

/*
Fx07 - LD Vx, DT

Set Vx = delay timer value.

The value of DT is placed into Vx.
*/
func (m *instructions) ldVxDt(c *chip8, x uint16) {
	c.register[x] = c.delayTimer
}

/*
Fx0A - LD Vx, K

Wait for a key press, store the value of the key in Vx.

All execution stops until a key is pressed, then the value of that key is stored in Vx.
*/
func (m *instructions) ldVxK(c *chip8) {
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
}

/*
Fx15 - LD DT, Vx

Set delay timer = Vx.

DT is set equal to the value of Vx.
*/
func (m *instructions) ldDtVx(c *chip8, x uint16) {
	c.delayTimer = c.register[x]
}

/*
Fx18 - LD ST, Vx

Set sound timer = Vx.

ST is set equal to the value of Vx.
*/
func (m *instructions) ldStVx(c *chip8, x uint16) {
	c.soundTimer = c.register[x]
}

/*
Fx1E - ADD I, Vx

Set I = I + Vx.

The values of I and Vx are added, and the results are stored in I.
*/
func (m *instructions) addIndexVx(c *chip8, x uint16) {
	c.index = c.index + uint16(c.register[x])
}

/*
Fx29 - LD F, Vx

Set I = location of sprite for digit Vx.

The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx. See section 2.4, Display, for more information on the Chip-8 hexadecimal font.
*/
func (m *instructions) ldFVx(c *chip8, x uint16) {
	c.index = uint16(c.register[x])
}

/*
Fx33 - LD B, Vx

Store BCD representation of Vx in memory locations I, I+1, and I+2.

The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
*/
func (m *instructions) ldBVx(c *chip8, x uint16) {
	number := c.register[x]
	c.memory[c.index] = number / 100
	c.memory[c.index+1] = (number % 100) / 10
	c.memory[c.index+2] = (number % 100) % 10
}

/*
Fx55 - LD [I], Vx

Store registers V0 through Vx in memory starting at location I.

The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
*/
func (m *instructions) ldIndexVX(c *chip8, x uint16) {
	for i := uint16(0); i <= x; i++ {
		c.memory[c.index+i] = c.register[i]
	}
	c.index += x + 1
}

/*
Fx65 - LD Vx, [I]

Read registers V0 through Vx from memory starting at location I.

The interpreter reads values from memory starting at location I into registers V0 through Vx.
*/
func (m *instructions) ldVxIndex(c *chip8, x uint16) {
	for i := uint16(0); i <= x; i++ {
		c.memory[c.index+i] = c.register[i]
	}
	c.index += x + 1
}
