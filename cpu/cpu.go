package cpu

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

type Chip8 struct {
	register     [16]uint8   // V0-VF registers
	memory       [4096]uint8 // 4kb of memory
	index        uint16      // index register
	pc           uint16      // Program counter
	stack        [16]uint16  // Stack for storing retunr address
	sp           uint8       // Stack pointer
	delayTimer   uint8       // Delay timer
	soundTimer   uint8       // Delay timer
	keypad       [16]uint8   // Keypad state
	opcode       uint16      // Current opcode
	instructions *instructions
	drawFlag     bool // Draw flag

	Video [2048]uint32 // Display buffer
}

func (c *Chip8) Reset() {
	c.register = [16]uint8{}
	c.memory = [4096]uint8{}
	c.index = 0
	c.pc = 0
	c.stack = [16]uint16{}
	c.sp = 0
	c.delayTimer = 0
	c.soundTimer = 0
	c.keypad = [16]uint8{}
	c.opcode = 0
	c.instructions = NewInstructions()
	c.Video = [2048]uint32{} //64*32
	c.drawFlag = false

}

func NewChip8() *Chip8 {
	return &Chip8{
		register:     [16]uint8{},
		memory:       [4096]uint8{},
		index:        0,
		pc:           0,
		stack:        [16]uint16{},
		sp:           0,
		delayTimer:   0,
		soundTimer:   0,
		keypad:       [16]uint8{},
		opcode:       0,
		instructions: NewInstructions(),

		Video:    [2048]uint32{}, //64*32
		drawFlag: false,
	}
}

func (c *Chip8) Init() {
	c.pc = START_ADDRESS

	for i, v := range fontset {
		c.memory[FONTSET_START_ADDRESS+i] = v
	}

	c.instructions.cls(c)
}

func (c *Chip8) LoadRom(rom []byte) {
	for i, v := range rom {
		c.memory[START_ADDRESS+i] = v
	}
}

func (c *Chip8) incrementCounter() {
	c.pc += 2
}

func (c *Chip8) Cycle() {
	opcode := c.fetchOpcode()
	c.incrementCounter()
	c.decodeExecute(opcode)
}

func (c *Chip8) OnKeyEvent(key uint8, press uint8) {
	c.keypad[key] = press
}

func (c *Chip8) DrawFlag() bool {
	return c.drawFlag
}

func (c *Chip8) SetDrawFlag(v bool) {
	c.drawFlag = v
}

func (c *Chip8) GetVideo() [2048]uint32 {
	return c.Video
}
