package graphos

import (
	"math"
)

func Distance(x0, y0, x1, y1 int) int {
	first := math.Pow(float64(x1-x0), 2)
	second := math.Pow(float64(y1-y0), 2)
	return int(math.Sqrt(first + second))
}

func DistanceManhattan(x0, y0, x1, y1 int) int {
	if x1 < x0 {
		x1, x0 = x0, x1
	}
	if y1 < y0 {
		y1, y0 = y0, y1
	}
	return (x1 - x0) + (y1 - y0)
}

func abs(n int) int {
	y := n >> 63
	return (n ^ y) - y
}

func (p *Instance) DrawLine(x0, y0, x1, y1 int) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := int(-1)
	if x0 < x1 {
		sx = 1
	}
	sy := int(-1)
	if y0 < y1 {
		sy = 1
	}
	e := dx - dy

	for {
		p.DrawPix(int(x0), int(y0), p.CurrentColor)

		if x0 == x1 && y0 == y1 {
			return
		}
		var e2 int = 2 * e
		if e2 > -dy {
			e = e - dy
			x0 = x0 + sx
		}
		if e2 < dx {
			e = e + dx
			y0 = y0 + sy
		}
	}
}

func (p *Instance) DrawBox(x1, y1, x2, y2 int) {

	for y := y1; y <= y2; y++ {
		p.DrawPix(x1, y, p.CurrentColor)
		p.DrawPix(x2, y, p.CurrentColor)
	}
	for x := x1; x <= x2; x++ {
		p.DrawPix(x, y1, p.CurrentColor)
		p.DrawPix(x, y2, p.CurrentColor)
	}
}

func (p *Instance) DrawCircle(x0, y0, radius int) {
	x := radius
	y := 0
	e := 0

	for x >= y {
		p.DrawPix(x0+x, y0+y, p.CurrentColor)
		p.DrawPix(x0+y, y0+x, p.CurrentColor)
		p.DrawPix(x0-y, y0+x, p.CurrentColor)
		p.DrawPix(x0-x, y0+y, p.CurrentColor)
		p.DrawPix(x0-x, y0-y, p.CurrentColor)
		p.DrawPix(x0-y, y0-x, p.CurrentColor)
		p.DrawPix(x0+y, y0-x, p.CurrentColor)
		p.DrawPix(x0+x, y0-y, p.CurrentColor)

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
			p.DrawPix(i, y0+y, p.CurrentColor)
			p.DrawPix(i, y0-y, p.CurrentColor)
		}
		for i := x0 - y; i <= x0+y; i++ {
			p.DrawPix(i, y0+x, p.CurrentColor)
			p.DrawPix(i, y0-x, p.CurrentColor)
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

func (p *Instance) DrawFilledBox(x1, y1, x2, y2 int, color Color) {
	pix := p.img.Pix

	array := make([]byte, 4*(x2-x1+1))
	for i := 0; i < len(array); i += 4 {
		copy(array[i:], color[:])
	}

	x1 = 4 * x1

	for y := y1; y <= y2; y++ {
		pos := p.img.Stride*y + x1
		copy(pix[pos:], array)
	}
}
