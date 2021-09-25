package graphos

import "math"

func Distance(x1, y1, x2, y2 int) int {
	first := math.Pow(float64(x2-x1), 2)
	second := math.Pow(float64(y2-y1), 2)
	return int(math.Sqrt(first + second))
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (p *Instance) DrawLine(x0, y0, x1, y1 int) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	var e2 int
	for {
		p.DrawPix(x0, y0)

		if x0 == x1 && y0 == y1 {
			return
		}
		e2 = 2 * err
		if e2 > -dy {
			err = err - dy
			x0 = x0 + sx
		}
		if e2 < dx {
			err = err + dx
			y0 = y0 + sy
		}
	}
}

func (p *Instance) Box(x1, y1, x2, y2 int) {
	for y := y1; y <= y2; y++ {
		p.DrawPix(x1, y)
		p.DrawPix(x2, y)
	}
	for x := x1; x <= x2; x++ {
		p.DrawPix(x, y1)
		p.DrawPix(x, y2)
	}
}

func (p *Instance) Circle(x0, y0, radius int) {
	x := radius
	y := 0
	e := 0

	for x >= y {
		p.DrawPix(x0+x, y0+y)
		p.DrawPix(x0+y, y0+x)
		p.DrawPix(x0-y, y0+x)
		p.DrawPix(x0-x, y0+y)
		p.DrawPix(x0-x, y0-y)
		p.DrawPix(x0-y, y0-x)
		p.DrawPix(x0+y, y0-x)
		p.DrawPix(x0+x, y0-y)

		if e <= 0 {
			y += 1
			e += 2*y + 1
		}
		if e > 0 {
			x -= 1
			e -= 2*x + 1
		}
	}
}

func (p *Instance) DrawFilledCircle(x0, y0, radius int) {
	x := radius
	y := 0
	xChange := 1 - (radius << 1)
	yChange := 0
	radiusError := 0

	for x >= y {
		for i := x0 - x; i <= x0+x; i++ {
			p.DrawPix(i, y0+y)
			p.DrawPix(i, y0-y)
		}
		for i := x0 - y; i <= x0+y; i++ {
			p.DrawPix(i, y0+x)
			p.DrawPix(i, y0-x)
		}

		y++
		radiusError += yChange
		yChange += 2
		if ((radiusError << 1) + xChange) > 0 {
			x--
			radiusError += xChange
			xChange += 2
		}
	}
}

func (p *Instance) DrawChar(index, fgColor, bgColor byte, x, y int) {
	var a, b uint64
	for a = 0; a < 8; a++ {
		for b = 0; b < 8; b++ {
			if p.Font.Bitmap[index][b]&(0x80>>a) != 0 {
				p.CurrentColor = fgColor
				p.DrawPix(int(a)+x, int(b)+y)
				continue
			}
			p.CurrentColor = bgColor
			p.DrawPix(int(a)+x, int(b)+y)
		}
	}
}
