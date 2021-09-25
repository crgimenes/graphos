package main

import (
	"math/rand"
	"time"

	"github.com/crgimenes/graphos"
)

type dot struct {
	X int
	Y int
}

var (
	cg *graphos.Instance

	xAux, yAux int

	dotMain      []dot
	dotUnreached []dot
	dotReached   []dot
)

var color byte

func getNextColor() byte {
	color++
	if color > 15 {
		color = 1
	}
	return color
}

func update(screen *graphos.Instance) error {

	screen.CurrentColor = 0

	screen.Clear()

	screen.CurrentColor = 0xf

	for i := 0; i < len(dotMain); i++ {
		x := random(1, 5)
		switch x {
		case 1:
			dotMain[i].X++
		case 2:
			dotMain[i].X--
		case 3:
			dotMain[i].Y++
		case 4:
			dotMain[i].Y--

		}
	}

	// --=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	for i := 0; i < len(dotMain); i++ {
		if i == 0 {
			dotReached = append(dotReached, dotMain[0])
		} else {
			dotUnreached = append(dotUnreached, dotMain[i])
		}
	}

	var rIndex int
	var uIndex int

	for len(dotUnreached) > 0 {
		var record = 1000000
		for r := 0; r < len(dotReached); r++ {
			for u := 0; u < len(dotUnreached); u++ {
				r1 := dotReached[r]
				u1 := dotUnreached[u]

				d := graphos.Distance(r1.X, r1.Y, u1.X, u1.Y)
				if d < record {
					record = d
					uIndex = u
					rIndex = r
				}
			}
		}

		/*
			screen.DrawFilledCircle(
				dotReached[rIndex].X,
				dotReached[rIndex].Y, 4)
		*/

		screen.DrawLine(
			dotReached[rIndex].X,
			dotReached[rIndex].Y,
			dotUnreached[uIndex].X,
			dotUnreached[uIndex].Y)

		dotReached = append(dotReached, dotUnreached[uIndex])
		dotUnreached = RemoveDot(dotUnreached, uIndex)

	}

	// --=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	for i := 0; i < len(dotMain); i++ {
		screen.CurrentColor = getNextColor()

		screen.DrawFilledCircle(dotMain[i].X, dotMain[i].Y, 4)

		/*
			if xAux == 0 {
				xAux = dotMain[i].X
				yAux = dotMain[i].Y
				screen.DrawFilledCircle(dotMain[i].X, dotMain[i].Y, 4)
			} else {
				screen.DrawFilledCircle(dotMain[i].X, dotMain[i].Y, 4)
				screen.DrawLine(xAux, yAux, dotMain[i].X, dotMain[i].Y)
				xAux = dotMain[i].X
				yAux = dotMain[i].Y
			}
		*/

	}
	dotReached = nil
	return nil
}

func RemoveDot(s []dot, index int) []dot {
	return append(s[:index], s[index+1:]...)
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	rand.Seed(time.Now().Unix())

	cg = graphos.Get()
	cg.Width = 800
	cg.Height = 600
	cg.Scale = 1
	cg.ScreenHandler = update
	cg.Title = "Minimum Spanning Tree - Prim's Algorithm"
	cg.CurrentColor = 0x0

	for i := 0; i < 60; i++ {

		d := dot{
			X: random(10, cg.Width-10),
			Y: random(10, cg.Height-10),
		}
		dotMain = append(dotMain, d)
	}

	cg.Run()
}
