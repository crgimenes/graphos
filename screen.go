package graphos

import (
	"fmt"
	"image"
	"log"
	"strings"

	"crg.eti.br/go/graphos/fonts"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rows          = 25 * 2
	columns       = 80 * 2
	columnsWord   = columns * 2
	totalTextSize = rows * columns * 2
)

type Instance struct {
	videoTextMemory  [totalTextSize]byte
	Border           int
	Height           int
	Width            int
	Scale            float64
	CurrentColor     byte
	UTime            uint64
	UpdateScreen     bool
	img              *image.RGBA
	ScreenHandler    func(*Instance) error
	Title            string
	cursor           int
	cursorBlinkTimer int
	cursorSetBlink   bool
	cursorLine       int
	cursorColumn     int
	Machine          int
	cpx, cpy         int
	Font             struct {
		Height int
		Width  int
		Bitmap []byte
	}
	lastKey struct {
		Time uint64
		Char byte
	}
	noKey bool
}

func New() *Instance {
	var i *Instance
	i = &Instance{}
	i.Width = columns * 9
	i.Height = rows * 16
	i.Scale = 1
	i.ScreenHandler = func(i *Instance) error {
		log.Println("ScreenHandler not defined")
		return nil
	}
	i.Title = "term"
	i.CurrentColor = 0x0F
	i.cursorSetBlink = true
	return i
}

type Colors struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

func (c Colors) RGBA() (r, g, b, a uint32) {
	r = c.R
	g = c.G
	b = c.B
	a = c.A

	return
}

var Colors16 = []Colors{
	{0, 0, 0, 0xFFFF},
	{0, 0, 170, 0xFFFF},
	{0, 170, 0, 0xFFFF},
	{0, 170, 170, 0xFFFF},
	{170, 0, 0, 0xFFFF},
	{170, 0, 170, 0xFFFF},
	{170, 85, 0, 0xFFFF},
	{170, 170, 170, 0xFFFF},
	{85, 85, 85, 0xFFFF},
	{85, 85, 255, 0xFFFF},
	{85, 255, 85, 0xFFFF},
	{85, 255, 255, 0xFFFF},
	{255, 85, 85, 0xFFFF},
	{255, 85, 255, 0xFFFF},
	{255, 255, 85, 0xFFFF},
	{255, 255, 255, 0xFFFF},
}

func (i *Instance) Write(p []byte) (n int, err error) {
	lp := len(p)
	i.Print(string(p))
	return lp, nil
}

func MergeColorCode(b, f byte) byte {
	return f&0xff | b<<4
}

func (i *Instance) Run() {

	i.Font.Bitmap = fonts.Bitmap
	i.Font.Height = 16
	i.Font.Width = 9
	i.img = image.NewRGBA(image.Rect(0, 0, i.Width, i.Height))

	i.Clear()
	i.UpdateScreen = true
	i.clearVideoTextMode()

	err := ebiten.RunGame(i)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *Instance) DrawPix(x, y int, color byte) {
	x += i.Border
	y += i.Border
	if x < i.Border ||
		y < i.Border ||
		x >= i.Width-i.Border ||
		y >= i.Height-i.Border {
		return
	}

	pos := 4*y*i.Width + 4*x
	i.img.Pix[pos] = uint8(Colors16[color].R)
	i.img.Pix[pos+1] = uint8(Colors16[color].G)
	i.img.Pix[pos+2] = uint8(Colors16[color].B)
	i.img.Pix[pos+3] = 0xff
	i.UpdateScreen = true
}

func (i *Instance) DrawChar(index, fgColor, bgColor byte, x, y int) {
	var a uint
	var b uint
	var lColor byte
	for b = 0; b < 16; b++ {
		for a = 0; a < 9; a++ {
			x1 := int(a) + x
			y1 := int(b) + y
			if a == 8 {
				c := bgColor
				if index >= 192 && index <= 223 {
					c = lColor
				}
				i.DrawPix(x1, y1, c)
				continue
			}
			idx := uint(index)*16 + b
			if fonts.Bitmap[idx]&(0x80>>a) != 0 {
				lColor = fgColor
				i.DrawPix(x1, y1, lColor)
				continue
			}
			lColor = bgColor
			i.DrawPix(x1, y1, lColor)
		}
	}
}

func (i *Instance) DrawString(s string, fgColor, bgColor byte, x, y int) {
	for idx := 0; idx < len(s); idx++ {
		i.DrawChar(s[idx], fgColor, bgColor, x, y)
		x += 9
	}
}

/*
	func (i *Instance) Clear() {
		for idx := 0; idx < i.Height*i.Width*4; idx += 4 {
			i.img.Pix[idx] = uint8(Colors16[i.CurrentColor].R)
			i.img.Pix[idx+1] = uint8(Colors16[i.CurrentColor].G)
			i.img.Pix[idx+2] = uint8(Colors16[i.CurrentColor].B)
			i.img.Pix[idx+3] = uint8(Colors16[i.CurrentColor].A)
		}
	}
*/
func (i *Instance) Clear() {
	color := Colors16[i.CurrentColor]
	r := uint8(color.R)
	g := uint8(color.G)
	b := uint8(color.B)
	a := uint8(color.A)

	pix := i.img.Pix
	pixLen := len(pix)

	for idx := 0; idx < pixLen; idx += 4 {
		pix[idx] = r
		pix[idx+1] = g
		pix[idx+2] = b
		pix[idx+3] = a
	}
}

func (i *Instance) DrawCursor(index, fgColor, bgColor byte, x, y int) {
	if i.cursorSetBlink {
		if i.cursorBlinkTimer < 15 {
			fgColor, bgColor = bgColor, fgColor
		}
		i.DrawChar(index, fgColor, bgColor, x, y)
		i.cursorBlinkTimer++
		if i.cursorBlinkTimer > 30 {
			i.cursorBlinkTimer = 0
		}
		return
	}
	i.DrawChar(index, bgColor, fgColor, x, y)
}

func (i *Instance) DrawVideoTextMode() {
	idx := 0
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			color := i.videoTextMemory[idx]
			f := color & 0x0f
			b := color & 0xf0 >> 4
			if idx == i.cursor {
				idx++
				i.DrawCursor(i.videoTextMemory[idx], f, b, c*9, r*16)
			} else {
				idx++
				i.DrawChar(i.videoTextMemory[idx], f, b, c*9, r*16)
			}
			idx++
		}
	}
}

func (i *Instance) clearVideoTextMode() {
	copy(i.videoTextMemory[:], make([]byte, totalTextSize))
	for idx := 0; idx < totalTextSize; idx += 2 {
		i.videoTextMemory[idx] = i.CurrentColor
	}
}

func (i *Instance) moveLineUp() {
	copy(i.videoTextMemory[0:], i.videoTextMemory[columnsWord:])
	copy(i.videoTextMemory[totalTextSize-columnsWord:], make([]byte, columnsWord))
	for idx := totalTextSize - columnsWord; idx < totalTextSize; idx += 2 {
		i.videoTextMemory[idx] = i.CurrentColor
	}
}

func (i *Instance) correctVideoCursor() {
	if i.cursor < 0 {
		i.cursor = 0
	}

	for i.cursor >= totalTextSize {
		i.cursor -= columnsWord
		i.moveLineUp()
	}
}

func (i *Instance) PutChar(c byte) {
	i.videoTextMemory[i.cursor] = i.CurrentColor
	i.videoTextMemory[i.cursor+1] = c
	i.cursor += 2
	i.correctVideoCursor()
}

func (i *Instance) Print(msg string) {
	for idx := 0; idx < len(msg); idx++ {
		c := msg[idx]

		switch c {
		case 13:
			i.cursor += columnsWord
			continue
		case 10:
			aux := i.cursor / columnsWord
			aux = aux * columnsWord
			i.cursor = aux
			continue
		}
		i.PutChar(msg[idx])
	}
}

func (i *Instance) Println(msg string) {
	msg += "\r\n"
	i.Print(msg)
}

func (i *Instance) keyTreatment(c byte, f func(c byte)) {
	if i.noKey || i.lastKey.Char != c || i.lastKey.Time+20 < i.UTime {
		f(c)
		i.noKey = false
		i.lastKey.Char = c
		i.lastKey.Time = i.UTime
	}
}

func (i *Instance) getLine() string {
	aux := i.cursor / columnsWord
	var ret string
	for idx := aux*columnsWord + 1; idx < aux*columnsWord+columnsWord; idx += 2 {
		c := i.videoTextMemory[idx]
		if c == 0 {
			break
		}
		ret += string(i.videoTextMemory[idx])
	}

	ret = strings.TrimSpace(ret)
	return ret
}

func (i *Instance) eval(cmd string) {
	fmt.Println("eval:", cmd)
}

func (i *Instance) InputChars() []rune {
	runes := make([]rune, 0, 16)
	return ebiten.AppendInputChars(runes)
}

func (i *Instance) Input() {
	for c := 'A'; c <= 'Z'; c++ {
		if ebiten.IsKeyPressed(ebiten.Key(c) - 'A' + ebiten.KeyA) {
			i.keyTreatment(byte(c), func(c byte) {
				if ebiten.IsKeyPressed(ebiten.KeyShift) {
					i.PutChar(c)
					return
				}
				i.PutChar(c + 32) // convert to lowercase
			})
			return
		}
	}

	for c := '0'; c <= '9'; c++ {
		if ebiten.IsKeyPressed(ebiten.Key(c) - '0' + ebiten.Key0) {
			i.keyTreatment(byte(c), func(c byte) {
				i.PutChar(c)
			})
			return
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		i.keyTreatment(byte(' '), func(c byte) {
			i.PutChar(c)
		})
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyComma) {
		i.keyTreatment(byte(','), func(c byte) {
			i.PutChar(c)
		})
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		i.keyTreatment(0, func(c byte) {
			i.eval(i.getLine())
			i.cursor += columnsWord
			aux := i.cursor / columnsWord
			aux = aux * columnsWord
			i.cursor = aux
			i.correctVideoCursor()
		})
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		i.keyTreatment(0, func(c byte) {
			i.cursor -= 2
			line := i.cursor / columnsWord
			lineEnd := line*columnsWord + columnsWord
			if i.cursor < 0 {
				i.cursor = 0
			}

			copy(i.videoTextMemory[i.cursor:lineEnd], i.videoTextMemory[i.cursor+2:lineEnd])
			i.videoTextMemory[lineEnd-2] = i.CurrentColor
			i.videoTextMemory[lineEnd-1] = 0

			i.correctVideoCursor()
		})
		return
	}

	/*
	   KeyMinus: -
	   KeyEqual: =
	   KeyLeftBracket: [
	   KeyRightBracket: ]
	   KeyBackslash:
	   KeySemicolon: ;
	   KeyApostrophe: '
	   KeySlash: /
	   KeyGraveAccent: `
	*/

	shift := ebiten.IsKeyPressed(ebiten.KeyShift)

	if ebiten.IsKeyPressed(ebiten.KeyEqual) {
		if shift {
			i.keyTreatment('+', func(c byte) {
				i.PutChar(c)
				println("+")
			})
			return
		}
		i.keyTreatment('=', func(c byte) {
			i.PutChar(c)
			println("=")
		})
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		i.keyTreatment(0, func(c byte) {
			i.cursor -= columnsWord
			i.correctVideoCursor()
		})
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		i.keyTreatment(0, func(c byte) {
			i.cursor += columnsWord
			i.correctVideoCursor()
		})
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		i.keyTreatment(0, func(c byte) {
			i.cursor -= 2
			i.correctVideoCursor()
		})
		return
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		i.keyTreatment(0, func(c byte) {
			i.cursor += 2
			i.correctVideoCursor()
		})
		return
	}

	// When the "left mouse button" is pressed...
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		//ebitenutil.DebugPrint(screen, "You're pressing the 'LEFT' mouse button.")
	}
	// When the "right mouse button" is pressed...
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		//ebitenutil.DebugPrint(screen, "\nYou're pressing the 'RIGHT' mouse button.")
	}
	// When the "middle mouse button" is pressed...
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		//ebitenutil.DebugPrint(screen, "\n\nYou're pressing the 'MIDDLE' mouse button.")
	}

	i.cpx, i.cpy = ebiten.CursorPosition()
	//fmt.Printf("X: %d, Y: %d\n", i.cpx, i.cpy)

	// Display the information with "X: xx, Y: xx" format
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %d, Y: %d", x, y))

	i.noKey = true

}

func (i *Instance) Draw(screen *ebiten.Image) {
	//op := &ebiten.DrawImageOptions{}
	//op.GeoM.Scale(1, 1)
	//screen.DrawImage(i.tmpScreen, op)
	screen.WritePixels(i.img.Pix)
}

func (i *Instance) Layout(outsideWidth, outsideHeight int) (int, int) {
	return i.Width, i.Height
}

func (i *Instance) Update() error {
	if i.ScreenHandler != nil {
		err := i.ScreenHandler(i)
		if err != nil {
			return err
		}
	}

	i.UTime++
	return nil
}
