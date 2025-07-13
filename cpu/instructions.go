package cpu

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
