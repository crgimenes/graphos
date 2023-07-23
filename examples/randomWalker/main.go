package main

import (
	"math/rand"
	"time"

	"crg.eti.br/go/graphos"
)

type Walker struct {
	X int
	Y int
}

var (
	walker Walker
)

func update(screen *graphos.Instance) error {

	x := random(1, 5)
	switch x {
	case 1:
		walker.X++
	case 2:
		walker.X--
	case 3:
		walker.Y++
	case 4:
		walker.Y--
	}

	//screen.CurrentColor = getNextColor()
	screen.CurrentColor = 15
	screen.DrawPix(walker.X, walker.Y, 0x0F)

	return nil
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	rand.Seed(time.Now().Unix())

	cg := graphos.New()
	cg.Width = 800 / 4
	cg.Height = 600 / 4
	cg.Scale = 4
	cg.ScreenHandler = update
	cg.Title = "Random Walker"
	cg.CurrentColor = 0

	walker = Walker{
		X: cg.Width / 2,
		Y: cg.Height / 2,
	}

	cg.Run()
}
