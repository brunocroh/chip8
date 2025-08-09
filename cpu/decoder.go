package cpu

func (c *Chip8) fetchOpcode() uint16 {
	return uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
}

func (c *Chip8) decodeExecute(opcode uint16) {
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
	case 0x1000:
		c.instructions.jump(c, nnn)
	case 0x2000:
		c.instructions.callSubroutine(c, nnn)
	case 0x3000:
		c.instructions.seVxKk(c, x, kk)
	case 0x4000:
		c.instructions.sneVxKk(c, x, kk)
	case 0x5000:
		c.instructions.seVxVy(c, x, y)
	case 0x6000:
		c.instructions.loadVx(c, x, kk)
	case 0x7000:
		c.instructions.addVx(c, x, kk)
	case 0x8000:
		switch n {
		case 0x0:
			c.instructions.loadVyIntoVx(c, x, y)
		case 0x1:
			c.instructions.orVxVy(c, x, y)
		case 0x2:
			c.instructions.andVxVy(c, x, y)
		case 0x3:
			c.instructions.xorVxVy(c, x, y)
		case 0x4:
			c.instructions.addVxVy(c, x, y)
		case 0x5:
			c.instructions.subVxVy(c, x, y)
		case 0x6:
			c.instructions.shrVx(c, x)
		case 0x7:
			c.instructions.subnVxVy(c, x, y)
		case 0xE:
			c.instructions.shlVx(c, x)
		}
	case 0x9000:
		c.instructions.sneVxVy(c, x, y)
	case 0xA000:
		c.instructions.loadIndex(c, nnn)
	// Bnnn
	case 0xB000:
		c.instructions.jumpV0(c, nnn)
	case 0xC000:
		c.instructions.randonVxKk(c, x, kk)
	case 0xD000:
		c.instructions.draw(c, x, y, n)
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			c.instructions.skpVx(c, x)
		case 0x00A1:
			c.instructions.sknpVx(c, x)
		}
	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007:
			c.instructions.ldVxDt(c, x)
		case 0x000A:
			c.instructions.ldVxK(c)
		case 0x0015:
			c.instructions.ldDtVx(c, x)
		case 0x0018:
			c.instructions.ldStVx(c, x)
		case 0x001E:
			c.instructions.addIndexVx(c, x)
		case 0x0029:
			c.instructions.ldFVx(c, x)
		case 0x0033:
			c.instructions.ldBVx(c, x)
		case 0x0055:
			c.instructions.ldIndexVX(c, x)
		case 0x0065:
			c.instructions.ldVxIndex(c, x)
		}
	}
}
