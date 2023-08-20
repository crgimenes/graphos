package main

import (
	"embed"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	memorySize        = 4096
	displaySize       = 64 * 32
	displaySizeWidth  = 64
	displaySizeHeight = 32
	programStart      = 0x200
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

	//go:embed roms/*
	roms embed.FS
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
	/*
		if addr >= memorySize {
			addr %= memorySize
		}
	*/

	return addr & (memorySize - 1)
}

func (c *chip8) MemorySet(addr uint16, value uint8) {
	addr = fixAddr(addr)
	c.memory[addr] = value
}

func (c *chip8) MemoryGet(addr uint16) uint8 {
	addr = fixAddr(addr)
	return c.memory[addr]
}

func (c *chip8) MemorySet16(addr uint16, value uint16) {
	addr = fixAddr(addr)
	c.memory[addr] = uint8(value >> 8)
	c.memory[addr+1] = uint8(value & 0x00FF)
}

func (c *chip8) MemoryGet16(addr uint16) uint16 {
	addr = fixAddr(addr)
	return uint16(c.memory[addr])<<8 | uint16(c.memory[addr+1])
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

func (c *chip8) DrawSprite(x, y uint8, sprite uint16, size uint8) bool {
	colision := false
	for i := uint8(0); i < size; i++ {
		for j := uint8(0); j < 8; j++ {
			xj := x + j
			xj %= displaySizeWidth
			yi := y + i
			yi %= displaySizeHeight
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
		//cg.Play()
		c.sound--
		//if c.sound == 0 {
		//	cg.Stop()
		//}
	}
}

func (c *chip8) LoadProgram(program []uint8) {
	for i := 0; i < len(program); i++ {
		c.MemorySet(uint16(i)+programStart, program[i])
	}
}

func (c *chip8) PrintProgramFromFile(filename string) {
	program, err := roms.ReadFile(filename)
	if err != nil {
		//log.Println("Error loading program file:", err)
		return
	}
	c.LoadProgram(program)
}

func (c *chip8) LoadROM(filename string) {
	program, err := roms.ReadFile("roms/" + filename)
	if err != nil {
		//log.Println("Error loading ROM file:", err)
		return
	}
	c.LoadProgram(program)
}

func retIfPrintable(c uint8) string {
	if c >= 32 && c <= 126 {
		return fmt.Sprintf("%c", c)
	}
	return "." // non printable
}

func (c *chip8) ClearRAM() {
	copy(c.memory[:], make([]uint8, memorySize))
}

func (c *chip8) PrintRAM() {
	s := "" // char column
	for i := 0; i < len(c.memory); i++ {
		// print address
		if i%16 == 0 {
			fmt.Printf("0x%04X: ", i)
		}

		fmt.Printf("%02X ", c.memory[i])

		s += retIfPrintable(c.memory[i])
		if i%16 == 15 {
			fmt.Printf("| %s\n", s)
			s = ""
		}
	}
}

func rnd(min, max int) int {
	return rand.Intn(max-min) + min
}

/*
var jumpTable = [16]func(*chip8, uint16){
	func(c *chip8, opcode uint16) { // 0x0NNN
	},
	func(c *chip8, opcode uint16) { // 0x00E0
		c.ClearDisplay()
		c.PC += 2
	},
}
*/

// var jumpTable = [16]func(*chip8, uint16)

type jtf func(c *chip8, opcode uint16)

func testeJumpTable(c *chip8, opcode uint16) {
}

var jumpTable = [16]jtf{
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
	testeJumpTable,
}

func (c *chip8) ExecOpcodeJT(opcode uint16) {
}

func (c *chip8) ExecOpcode(opcode uint16) {
	if c.PC >= memorySize {
		c.PC %= memorySize
		//log.Printf("PC overflow: %04X\n", c.PC)
	}

	//jumpTable[(opcode&0xF000)>>12](c, opcode)

	j := [16]func(*chip8, uint16){
		testeJumpTable,
		testeJumpTable,
	}

	fmt.Println(j)

	//////////////////////////////

	////log.Printf("executing opcode: %04X\n", opcode)
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x00FF {
		case 0x00E0: // 00E0 - CLS
			//log.Printf("0x%04X: 0x%04X - CLS\n", c.PC, opcode)
			c.ClearDisplay()
			c.PC += 2
		case 0x00EE: // 00EE - RET
			//log.Printf("0x%04X: 0x%04X - RET\n", c.PC, opcode)
			c.PC = c.Pop()
		default:
			//log.Printf("Address: 0x%04X - Unknown opcode: 0x%04X\n", c.PC, opcode)
			os.Exit(1)
		}
	case 0x1000: // 1nnn - JP addr
		//oldPC := c.PC
		c.PC = opcode & 0x0FFF
		//log.Printf("0x%04X: 0x%04X - JP [0x%04X]\n", //oldPC, opcode, c.PC)
	case 0x2000: // 2nnn - CALL addr
		//oldPC := c.PC
		c.Push(c.PC + 2)
		c.PC = opcode & 0x0FFF
		//log.Printf("0x%04X: 0x%04X - CALL [0x%04X]\n", //oldPC, opcode, c.PC)
	case 0x3000: // 3xkk - SE Vx, byte
		//oldPC := c.PC
		if c.GetV(uint8((opcode&0x0F00)>>8)) == uint8(opcode&0x00FF) {
			c.PC += 4
		} else {
			c.PC += 2
		}
		//log.Printf("0x%04X: 0x%04X - SE V%X, 0x%02X\n", //oldPC, opcode, (opcode&0x0F00)>>8, opcode&0x00FF)
	case 0x4000: // 4xkk - SNE Vx, byte
		//oldPC := c.PC
		if c.GetV(uint8((opcode&0x0F00)>>8)) != uint8(opcode&0x00FF) {
			c.PC += 4
		} else {
			c.PC += 2
		}
		//log.Printf("0x%04X: 0x%04X - SNE V%X, 0x%02X\n", //oldPC, opcode, (opcode&0x0F00)>>8, opcode&0x00FF)
	case 0x5000: // 5xy0 - SE Vx, Vy
		//oldPC := c.PC
		if c.GetV(uint8((opcode&0x0F00)>>8)) == c.GetV(uint8((opcode&0x00F0)>>4)) {
			c.PC += 4
		} else {
			c.PC += 2
		}
		//log.Printf("0x%04X: 0x%04X - SE V%X, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8, (opcode&0x00F0)>>4)
	case 0x6000: // 6xkk - LD Vx, byte
		//oldPC := c.PC
		c.SetV(uint8((opcode&0x0F00)>>8), uint8(opcode&0x00FF))
		c.PC += 2
		//log.Printf("0x%04X: 0x%04X - LD V%X, 0x%02X\n", //oldPC, opcode, (opcode&0x0F00)>>8, opcode&0x00FF)
	case 0x7000: // 7xkk - ADD Vx, byte
		//oldPC := c.PC
		c.SetV(uint8((opcode&0x0F00)>>8), c.GetV(uint8((opcode&0x0F00)>>8))+uint8(opcode&0x00FF))
		c.PC += 2
		//log.Printf("0x%04X: 0x%04X - ADD V%X, 0x%02X\n", //oldPC, opcode, (opcode&0x0F00)>>8, opcode&0x00FF)
	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000: // 8xy0 - LD Vx, Vy
			//oldPC := c.PC
			c.SetV(uint8((opcode&0x0F00)>>8), c.GetV(uint8((opcode&0x00F0)>>4)))
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - LD V%X, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8, (opcode&0x00F0)>>4)
		case 0x0001: // 8xy1 - OR Vx, Vy
			//oldPC := c.PC
			c.SetV(uint8((opcode&0x0F00)>>8), c.GetV(uint8((opcode&0x0F00)>>8))|c.GetV(uint8((opcode&0x00F0)>>4)))
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - OR V%X, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8, (opcode&0x00F0)>>4)
		case 0x0002: // 8xy2 - AND Vx, Vy
			//oldPC := c.PC
			c.SetV(uint8((opcode&0x0F00)>>8), c.GetV(uint8((opcode&0x0F00)>>8))&c.GetV(uint8((opcode&0x00F0)>>4)))
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - AND V%X, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8, (opcode&0x00F0)>>4)
		case 0x0003: // 8xy3 - XOR Vx, Vy
			//oldPC := c.PC
			c.SetV(uint8((opcode&0x0F00)>>8), c.GetV(uint8((opcode&0x0F00)>>8))^c.GetV(uint8((opcode&0x00F0)>>4)))
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - XOR V%X, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8, (opcode&0x00F0)>>4)
		case 0x0004: // 8xy4 - ADD Vx, Vy
			//oldPC := c.PC
			vx := c.GetV(uint8((opcode & 0x0F00) >> 8))
			vy := c.GetV(uint8((opcode & 0x00F0) >> 4))
			c.SetV(uint8((opcode&0x0F00)>>8), vx+vy)
			if vx > 0xFF-vy {
				c.SetV(0xF, 1)
			} else {
				c.SetV(0xF, 0)
			}
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - ADD V%X, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8, (opcode&0x00F0)>>4)
		case 0x0005: // 8xy5 - SUB Vx, Vy
			//oldPC := c.PC
			vx := c.GetV(uint8((opcode & 0x0F00) >> 8))
			vy := c.GetV(uint8((opcode & 0x00F0) >> 4))
			if vx > vy {
				c.SetV(0xF, 1)
			} else {
				c.SetV(0xF, 0)
			}
			c.SetV(uint8((opcode&0x0F00)>>8), vx-vy)
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - SUB V%X, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8, (opcode&0x00F0)>>4)
		case 0x0006: // 8xy6 - SHR Vx {, Vy}
			//oldPC := c.PC
			vx := c.GetV(uint8((opcode & 0x0F00) >> 8))
			c.SetV(0xF, vx&0x1)
			c.SetV(uint8((opcode&0x0F00)>>8), vx>>1)
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - SHR V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		case 0x0007: // 8xy7 - SUBN Vx, Vy
			//oldPC := c.PC
			vx := c.GetV(uint8((opcode & 0x0F00) >> 8))
			vy := c.GetV(uint8((opcode & 0x00F0) >> 4))
			if vy > vx {
				c.SetV(0xF, 1)
			} else {
				c.SetV(0xF, 0)
			}
			c.SetV(uint8((opcode&0x0F00)>>8), vy-vx)
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - SUBN V%X, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8, (opcode&0x00F0)>>4)
		case 0x000E: // 8xyE - SHL Vx {, Vy}
			//oldPC := c.PC
			vx := c.GetV(uint8((opcode & 0x0F00) >> 8))
			c.SetV(0xF, vx>>7)
			c.SetV(uint8((opcode&0x0F00)>>8), vx<<1)
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - SHL V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		default:
			//log.Printf("0x%04X: 0x%04X - Unknown opcode\n", c.PC, opcode)
			os.Exit(1)
		}
	case 0x9000: // 9xy0 - SNE Vx, Vy
		//oldPC := c.PC
		if c.GetV(uint8((opcode&0x0F00)>>8)) != c.GetV(uint8((opcode&0x00F0)>>4)) {
			c.PC += 4
		} else {
			c.PC += 2
		}
		//log.Printf("0x%04X: 0x%04X - SNE V%X, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8, (opcode&0x00F0)>>4)
	case 0xA000: // Annn - LD I, addr
		//oldPC := c.PC
		c.I = opcode & 0x0FFF
		c.PC += 2
		//log.Printf("0x%04X: 0x%04X - LD I, 0x%04X\n", //oldPC, opcode, opcode&0x0FFF)
	case 0xB000: // Bnnn - JP V0, addr
		//oldPC := c.PC
		c.PC = (opcode & 0x0FFF) + uint16(c.GetV(0))
		//log.Printf("0x%04X: 0x%04X - JP V0, 0x%04X\n", //oldPC, opcode, opcode&0x0FFF)
	case 0xC000: // Cxkk - RND Vx, byte
		//oldPC := c.PC
		c.SetV(uint8((opcode&0x0F00)>>8), uint8(opcode&0x00FF)&uint8(rnd(0, 255)))
		c.PC += 2
		//log.Printf("0x%04X: 0x%04X - RND V%X, 0x%02X\n", //oldPC, opcode, (opcode&0x0F00)>>8, opcode&0x00FF)
	case 0xD000: // Dxyn - DRW Vx, Vy, nibble
		//oldPC := c.PC
		x := c.GetV(uint8((opcode & 0x0F00) >> 8))
		y := c.GetV(uint8((opcode & 0x00F0) >> 4))
		n := uint8(opcode & 0x000F)
		c.SetV(0xF, 0)
		if c.DrawSprite(x, y, c.I, n) {
			c.SetV(0xF, 1)
		}
		c.PC += 2
		//log.Printf("0x%04X: 0x%04X - DRW V%X, V%X, 0x%02X\n", //oldPC, opcode, (opcode&0x0F00)>>8, (opcode&0x00F0)>>4, opcode&0x000F)
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E: // Ex9E - SKP Vx
			//oldPC := c.PC
			if c.keys[c.GetV(uint8((opcode&0x0F00)>>8))] {
				c.PC += 4
			} else {
				c.PC += 2
			}
			//log.Printf("0x%04X: 0x%04X - SKP V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		case 0x00A1: // ExA1 - SKNP Vx
			//oldPC := c.PC
			if !c.keys[c.GetV(uint8((opcode&0x0F00)>>8))] {
				c.PC += 4
			} else {
				c.PC += 2
			}
			//log.Printf("0x%04X: 0x%04X - SKNP V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8)

		default:
			//log.Printf("0x%04X: 0x%04X - Unknown opcode\n", c.PC, opcode)
			os.Exit(1)
		}

	case 0xF000:
		switch opcode & 0x00FF {
		case 0x0007: // Fx07 - LD Vx, DT
			//oldPC := c.PC
			c.SetV(uint8((opcode&0x0F00)>>8), c.DelayTimer())
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - LD V%X, DT\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		case 0x000A: // Fx0A - LD Vx, K
			//oldPC := c.PC
			for i := 0; i < len(c.keys); i++ {
				if c.keys[i] {
					c.SetV(uint8((opcode&0x0F00)>>8), uint8(i))
					c.PC += 2
					break
				}
			}
			//log.Printf("0x%04X: 0x%04X - LD V%X, K\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		case 0x0015: // Fx15 - LD DT, Vx
			//oldPC := c.PC
			c.SetDelayTimer(c.GetV(uint8((opcode & 0x0F00) >> 8)))
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - LD DT, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		case 0x0018: // Fx18 - LD ST, Vx
			//oldPC := c.PC
			c.SetSoundTimer(c.GetV(uint8((opcode & 0x0F00) >> 8)))
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - LD ST, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		case 0x001E: // Fx1E - ADD I, Vx
			//oldPC := c.PC
			c.I += uint16(c.GetV(uint8((opcode & 0x0F00) >> 8)))
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - ADD I, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		case 0x0029: // Fx29 - LD F, Vx
			//oldPC := c.PC
			c.I = uint16(c.GetV(uint8((opcode&0x0F00)>>8))) * 5
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - LD F, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		case 0x0033: // Fx33 - LD B, Vx
			//oldPC := c.PC
			c.MemorySet(c.I, c.GetV(uint8((opcode&0x0F00)>>8))/100)
			c.MemorySet(c.I+1, (c.GetV(uint8((opcode&0x0F00)>>8))/10)%10)
			c.MemorySet(c.I+2, (c.GetV(uint8((opcode&0x0F00)>>8))%100)%10)
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - LD B, V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		case 0x0055: // Fx55 - LD [I], Vx
			//oldPC := c.PC
			for i := uint16(0); i <= uint16((opcode&0x0F00)>>8); i++ {
				c.MemorySet(c.I+i, c.GetV(uint8(i)))
			}
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - LD [I], V%X\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		case 0x0065: // Fx65 - LD Vx, [I]
			//oldPC := c.PC
			for i := uint16(0); i <= uint16((opcode&0x0F00)>>8); i++ {
				c.SetV(uint8(i), c.MemoryGet(c.I+i))
			}
			c.PC += 2
			//log.Printf("0x%04X: 0x%04X - LD V%X, [I]\n", //oldPC, opcode, (opcode&0x0F00)>>8)
		default:
			//log.Printf("0x%04X: 0x%04X - Unknown opcode\n", c.PC, opcode)
			os.Exit(1)
		}
	default:
		//log.Printf("0x%04X: 0x%04X - Unknown opcode", c.PC, opcode)
		os.Exit(1)

	}

}

var (
	opcodeCount    = 1
	maxOpcodeCount = 20
	opcode         uint16
)

func (c *chip8) Cycle() {
	beginTime := time.Now()

	for i := opcodeCount; i != 0; i-- {
		opcode = c.MemoryGet16(c.PC)
		//jumpTable[(opcode&0xF000)>>12](c, opcode)

		c.ExecOpcode(opcode)
		opcode = c.MemoryGet16(c.PC)
		c.ExecOpcode(opcode)
		opcode = c.MemoryGet16(c.PC)
		c.ExecOpcode(opcode)
		opcode = c.MemoryGet16(c.PC)
		c.ExecOpcode(opcode)
	}

	deltaTime := time.Since(beginTime)
	if deltaTime < (time.Duration(time.Duration(opcodeCount*4) * time.Millisecond)) {

		//if c.PC&15 == 0 {
		fmt.Printf("deltaTime sleep: %v\n", deltaTime)
		//}

		if deltaTime < 1 {
			if opcodeCount <= maxOpcodeCount {
				opcodeCount++
				log.Printf("opcodeCount: %d\n", opcodeCount)
			}
		}

		time.Sleep(time.Duration(time.Duration(opcodeCount*4)*time.Millisecond) - deltaTime)
	}
}
