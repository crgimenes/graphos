package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"crg.eti.br/go/graphos"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	cg  = graphos.New()
	c8  = &chip8{}
	key = makeKeyMap()
	ii  = uint8(0)
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

type keyMap struct {
	ekey   ebiten.Key
	c8key  uint8
	cx, cy int
	x, y   int
	x1, y1 int
	chr    string
}

func makeKeyMap() []keyMap {
	/*
		1 2 3 C
		4 5 6 D
		7 8 9 E
		A 0 B F
	*/
	km := [16]keyMap{}

	keys := []byte{
		0x1, 0x2, 0x3, 0xC,
		0x4, 0x5, 0x6, 0xD,
		0x7, 0x8, 0x9, 0xE,
		0xA, 0x0, 0xB, 0xF,
	}

	ekeys := []ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.KeyC,
		ebiten.Key4, ebiten.Key5, ebiten.Key6, ebiten.KeyD,
		ebiten.Key7, ebiten.Key8, ebiten.Key9, ebiten.KeyE,
		ebiten.KeyA, ebiten.Key0, ebiten.KeyB, ebiten.KeyF,
	}

	xBase, yBase := 500, 300
	x, y := xBase, yBase
	for i := 0; i < 16; i++ {
		x += 30
		if i%4 == 0 {
			x = xBase
			y += 30
		}

		km[i] = keyMap{
			ekey:  ekeys[i],
			c8key: keys[i],
			cx:    x,
			cy:    y,
			x:     x - 10,
			y:     y - 9,
			x1:    x + 18,
			y1:    y + 19,
			chr:   fmt.Sprintf("%X", keys[i]),
		}
	}

	return km[:]
}

func drawKeyboard(g *graphos.Instance) {
	black := graphos.Colors16[0]
	white := graphos.Colors16[0xF]
	xBase, yBase := 500, 300
	g.DrawFilledBox(xBase-11, yBase+20, xBase+109, yBase+140, white)
	for i := 0; i < 16; i++ {
		k := key[i]
		c := key[i].chr
		if c8.keys[key[i].c8key] {
			g.DrawFilledBox(k.x, k.y, k.x1, k.y1, white)
			g.DrawString(c, 0, 0x0F, k.cx, k.cy)
			continue
		}
		g.DrawFilledBox(k.x, k.y, k.x1, k.y1, black)
		g.DrawString(c, 0x0F, 0, k.cx, k.cy)
	}
}

func drawDisplay(g *graphos.Instance, x, y int) {
	black := graphos.Colors16[0]
	white := graphos.Colors16[0xF]

	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			if c8.display[i][j] {

				for k := 0; k < 8; k++ {
					for l := 0; l < 8; l++ {
						g.DrawPix(x+i*8+k, y+j*8+l, white)
					}
				}

				continue
			}

			for k := 0; k < 8; k++ {
				for l := 0; l < 8; l++ {
					g.DrawPix(x+i*8+k, y+j*8+l, black)
				}
			}
		}
	}
}

func input(i *graphos.Instance) {
	for _, v := range key {
		if inpututil.IsKeyJustPressed(v.ekey) {
			c8.keys[v.c8key] = true
			continue
		}

		if inpututil.IsKeyJustReleased(v.ekey) {
			c8.keys[v.c8key] = false
			continue
		}

		x, y := ebiten.CursorPosition()

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if x >= v.x && x <= v.x1 && y >= v.y && y <= v.y1 {
				c8.keys[v.c8key] = true
			}
			continue
		}

		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			//if x >= v.x && x <= v.x1 && y >= v.y && y <= v.y1 {
			c8.keys[v.c8key] = false
			//}
		}
	}
}

var (
	firstRun = true
	black    = graphos.Colors16[0]
	white    = graphos.Colors16[0xF]
)

func update(i *graphos.Instance) error {
	if firstRun {
		i.DrawFilledBox(0, 0, i.Width-1, i.Height-1, black)
		firstRun = false
	}
	i.CurrentColor = white

	//i.Input()

	/*
		runes := i.InputChars()
		if len(runes) > 0 {
			log.Printf("runes: %v\n", string(runes))
		}
	*/

	// i.DrawLine(0, 0, i.Width, i.Height)

	//i.DrawChar('A', 0x0F, 0, 300, 300)
	//i.DrawString("Teste de string", 0x0F, 0, 300, 300)
	//i.DrawString("Teste de string", 0x0, 0x0F, 300, 300+17)

	//i.CurrentColor = 0x0F
	i.DrawFilledBox(4, 4, 64*8+8+4, 32*8+8+4, white)

	drawRegisters(i)
	drawDisplay(i, 8, 8)

	/*
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
	*/
	input(i)

	// TODO: draw stack
	drawKeyboard(i)

	c8.DelayTimerTick()
	c8.SoundTimerTick()

	//c8.Cycle()

	i.UpdateScreen = true

	return nil
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//////////////////////////////////////
	f, err := os.Create("chip8.prof")
	if err != nil {

		fmt.Println(err)
		return

	}
	pprof.StartCPUProfile(f)
	//pprof.WriteHeapProfile(f)
	defer func() {
		pprof.StopCPUProfile()
		f.Close()
	}()
	//////////////////////////////////////

	cg.Width = 800
	cg.Height = 600
	cg.ScreenHandler = update
	cg.Title = "chip8"
	cg.CurrentColor = graphos.Colors16[0x0F]

	c8.InitCharSet()
	c8.ClearDisplay()

	cg.InitSound()
	//c8.SetSoundTimer(50)

	/*
		go func() {
			time.Sleep(10 * time.Second)
			c8.SetSoundTimer(20)
		}()
	*/

	//c8.LoadROM("MAZE")
	c8.LoadROM("INVADERS")
	//c8.LoadROM("PONG2")
	c8.PrintRAM()

	c8.PC = 0x200

	go func() {
		for {
			if cg.Running {
				c8.Cycle()
				continue
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	cg.Run()
}
