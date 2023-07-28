package main

const (
	memorySize        = 4096
	displaySize       = 64 * 32
	displaySizeWidth  = 64
	displaySizeHeight = 32
)

var (
	charSet = [80]uint8{
		0xF0, // ****
		0x90, // *  *
		0x90, // *  *
		0x90, // *  *
		0xF0, // ****

		0x20, //   *
		0x60, //  **
		0x20, //   *
		0x20, //   *
		0x70, //  ***

		0xF0, // ****
		0x10, //    *
		0xF0, // ****
		0x80, // *
		0xF0, // ****

		0xF0, // ****
		0x10, //    *
		0xF0, // ****
		0x10, //    *
		0xF0, // ****

		0x90, // *  *
		0x90, // *  *
		0xF0, // ****
		0x10, //    *
		0x10, //    *

		0xF0, // ****
		0x80, // *
		0xF0, // ****
		0x10, //    *
		0xF0, // ****

		0xF0, // ****
		0x80, // *
		0xF0, // ****
		0x90, // *  *
		0xF0, // ****

		0xF0, // ****
		0x10, //    *
		0x20, //   *
		0x40, //  *
		0x40, //  *

		0xF0, // ****
		0x90, // *  *
		0xF0, // ****
		0x90, // *  *
		0xF0, // ****

		0xF0, // ****
		0x90, // *  *
		0xF0, // ****
		0x10, //    *
		0xF0, // ****

		0xF0, // ****
		0x90, // *  *
		0xF0, // ****
		0x90, // *  *
		0x90, // *  *

		0xE0, // ***
		0x90, // *  *
		0xE0, // ***
		0x90, // *  *
		0xE0, // ***

		0xF0, // ****
		0x80, // *
		0x80, // *
		0x80, // *
		0xF0, // ****

		0xE0, // ***
		0x90, // *  *
		0x90, // *  *
		0x90, // *  *
		0xE0, // ***

		0xF0, // ****
		0x80, // *
		0xF0, // ****
		0x80, // *
		0xF0, // ****

		0xF0, // ****
		0x80, // *
		0xF0, // ****
		0x80, // *
		0x80, // *
	}
)

type chip8 struct {
	memory  [memorySize]uint8
	V       [16]uint8                                 // 16 8-bit registers
	I       uint16                                    // 16-bit register for memory address
	PC      uint16                                    // 16-bit program counter
	SP      uint8                                     // 16-bit stack pointer
	stack   [16]uint16                                // 16 16-bit stack
	delay   uint8                                     // 8-bit delay timer
	sound   uint8                                     // 8-bit sound timer
	display [displaySizeWidth][displaySizeHeight]bool // 64x32 monochrome display
	keys    [16]bool                                  // 16-key hexadecimal keypad

	ModifiedRegister [16]int8 // Used to indicate that a register has been modified, it decrements at each frame and stops at zero.
	ModifiedStack    [16]int8
}

func fixAddr(addr uint16) uint16 {
	if addr >= memorySize {
		addr %= memorySize
	}
	return addr
}

func (c *chip8) MemorySet(addr uint16, value uint8) {
	addr = fixAddr(addr)
	c.memory[addr] = value
}

func (c *chip8) MemoryGet(addr uint16) uint8 {
	addr = fixAddr(addr)
	return c.memory[addr]
}

func (c *chip8) SetV(index uint8, value uint8) {
	c.V[index] = value
	c.ModifiedRegister[index] = 60
}

func (c *chip8) GetV(index uint8) uint8 {
	return c.V[index]
}

func (c *chip8) Push(value uint16) {
	c.stack[c.SP] = value
	c.ModifiedStack[c.SP] = 60
	c.SP++
	// TODO: check for stack overflow
}

func (c *chip8) Pop() uint16 {
	c.SP--
	c.ModifiedStack[c.SP] = 0
	return c.stack[c.SP]
}

func (c *chip8) InitCharSet() {
	for i := 0; i < len(charSet); i++ {
		c.MemorySet(uint16(i), charSet[i])
	}
}

func (c *chip8) DrawPixel(x, y uint8) {
	c.display[x][y] = !c.display[x][y]
}

func (c *chip8) GetPixel(x, y uint8) bool {
	return c.display[x][y]
}

func (c *chip8) SetPixel(x, y uint8, value bool) {
	c.display[x][y] = value
}

func (c *chip8) ClearDisplay() {
	copy(c.display[:], make([][displaySizeHeight]bool, displaySizeWidth))
}

func (c *chip8) DrawSprite(x, y, sprite, size uint8) bool {
	colision := false
	for i := uint8(0); i < size; i++ {
		for j := uint8(0); j < 8; j++ {
			xj := x + j
			if xj >= displaySizeWidth {
				xj %= displaySizeWidth
			}
			yi := y + i
			if yi >= displaySizeHeight {
				yi %= displaySizeHeight
			}
			if c.GetPixel(xj, yi) {
				colision = true
			}
			px := c.GetPixel(xj, yi)
			spx := (c.MemoryGet(uint16(sprite)+uint16(i)) & (0x80 >> j)) != 0

			c.SetPixel(xj, yi, px != spx)

		}
	}
	return colision
}

func (c *chip8) DelayTimer() uint8 {
	return c.delay
}

func (c *chip8) SetDelayTimer(value uint8) {
	c.delay = value
}

func (c *chip8) SoundTimer() uint8 {
	return c.sound
}

func (c *chip8) SetSoundTimer(value uint8) {
	c.sound = value
}

func (c *chip8) DelayTimerTick() {
	if c.delay > 0 {
		c.delay--
	}
}

func (c *chip8) SoundTimerTick() {
	if c.sound > 0 {
		c.sound--
	}
}
