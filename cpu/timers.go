package cpu

import (
	"fmt"
)

func (c *chip8) UpdateTimers() {
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
