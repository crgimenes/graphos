package main

import (
	"math/rand"
	"time"

	"github.com/crgimenes/graphos"
)

type walker struct {
	X int
	Y int
}

var (
	cg      *graphos.Instance
	walkers []walker
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

	for i := 0; i < len(walkers); i++ {
		x := random(1, 5)
		switch x {
		case 1:
			walkers[i].X++
		case 2:
			walkers[i].X--
		case 3:
			walkers[i].Y++
		case 4:
			walkers[i].Y--
		}

		screen.CurrentColor = getNextColor()
		screen.DrawPix(walkers[i].X, walkers[i].Y)
	}
	return nil
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
	cg.Title = "Random Walker"
	cg.CurrentColor = 0x0

	for i := 0; i < 10; i++ {
		w := walker{
			X: random(10, cg.Width-10),
			Y: random(10, cg.Height-10),
		}
		walkers = append(walkers, w)
	}

	cg.Run()
}
