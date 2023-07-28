package main

import (
	"fmt"

	"crg.eti.br/go/graphos"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	c8  = &chip8{}
	key = map[ebiten.Key]uint8{}
)

func drawRegisters(g *graphos.Instance) {
	x, y := 700, 0
	for i := 0; i < 16; i++ {
		y += 17
		s := fmt.Sprintf("V%X: 0x%02X", i, c8.V[i])
		if c8.ModifiedRegister[i] > 0 {
			g.DrawString(s, 0x0, 0xF, x, y)
			c8.ModifiedRegister[i]--
			continue
		}
		g.DrawString(s, 0x0F, 0, x, y)
	}
}

func drawChip8Keys(g *graphos.Instance) {
	/*
		1 2 3 C
		4 5 6 D
		7 8 9 E
		A 0 B F
	*/

	keys := []byte{
		0x1, 0x2, 0x3, 0xC,
		0x4, 0x5, 0x6, 0xD,
		0x7, 0x8, 0x9, 0xE,
		0xA, 0x0, 0xB, 0xF,
	}
	xBase, yBase := 500, 300
	x, y := xBase, yBase
	for i := 0; i < 16; i++ {
		x += 30
		if i%4 == 0 {
			x = xBase
			y += 30
		}
		c := fmt.Sprintf("%X", keys[i])
		if c8.keys[keys[i]] {
			g.DrawString(c, 0, 0x0F, x, y)
			continue
		}
		g.DrawString(c, 0x0F, 0, x, y)

	}
}

func drawDisplay(g *graphos.Instance, x, y int) {

	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			if c8.display[i][j] {

				for k := 0; k < 8; k++ {
					for l := 0; l < 8; l++ {
						g.DrawPix(x+i*8+k, y+j*8+l, 0xF)
					}
				}

				continue
			}

			for k := 0; k < 8; k++ {
				for l := 0; l < 8; l++ {
					g.DrawPix(x+i*8+k, y+j*8+l, 0x0)
				}
			}

		}
	}
}

func input(i *graphos.Instance) {
	for k, v := range key {
		i.InputPressed(k, func(i *graphos.Instance) {
			c8.keys[v] = true
		})

		i.InputReleased(k, func(i *graphos.Instance) {
			c8.keys[v] = false
		})
	}

}

var ii = uint8(0)

func update(i *graphos.Instance) error {
	i.CurrentColor = 0x0
	i.DrawFilledBox(0, 0, i.Width, i.Height)
	i.CurrentColor = 0xF

	//i.Input()

	/*
		runes := i.InputChars()
		if len(runes) > 0 {
			fmt.Printf("runes: %v\n", string(runes))
		}
	*/

	// i.DrawLine(0, 0, i.Width, i.Height)
	//i.CurrentColor = 0x00
	//i.DrawFilledBox(0, 0, 800, 600)

	//i.DrawChar('A', 0x0F, 0, 300, 300)
	//i.DrawString("Teste de string", 0x0F, 0, 300, 300)
	//i.DrawString("Teste de string", 0x0, 0x0F, 300, 300+17)

	//i.CurrentColor = 0x0F
	i.DrawFilledBox(4, 4, 64*8+8+4, 32*8+8+4)

	drawRegisters(i)
	drawDisplay(i, 8, 8)

	for x, y := uint8(0), uint8(0); x < 64; x++ {
		for y = 0; y < 32; y++ {
			c8.SetPixel(x, y, true)
		}
	}

	x := uint8(0)
	for y := uint8(0); y < 32; y++ {
		c8.SetPixel(x, y, false)
		x++
	}

	ii++
	c8.DrawSprite(62, 10, 0, 5)
	c8.DrawSprite(10, 10+ii, 0, 5)
	c8.DrawSprite(10+ii, 20, 0, 5)
	c8.DrawSprite(11+ii, 21+ii, 0, 5)

	// TODO: draw stack
	// TODO: draw display
	drawChip8Keys(i)

	input(i)

	c8.DelayTimerTick()
	c8.SoundTimerTick()
	return nil
}

func main() {

	cg := graphos.New()
	cg.Width = 800
	cg.Height = 600
	cg.Scale = 1
	cg.ScreenHandler = update
	cg.Title = "chip8"
	cg.CurrentColor = 0x0F

	key = map[ebiten.Key]uint8{
		ebiten.Key0: 0x0,
		ebiten.Key1: 0x1,
		ebiten.Key2: 0x2,
		ebiten.Key3: 0x3,
		ebiten.Key4: 0x4,
		ebiten.Key5: 0x5,
		ebiten.Key6: 0x6,
		ebiten.Key7: 0x7,
		ebiten.Key8: 0x8,
		ebiten.Key9: 0x9,
		ebiten.KeyA: 0xA,
		ebiten.KeyB: 0xB,
		ebiten.KeyC: 0xC,
		ebiten.KeyD: 0xD,
		ebiten.KeyE: 0xE,
		ebiten.KeyF: 0xF,
	}

	c8.InitCharSet()
	c8.ClearDisplay()

	cg.Run()
}
