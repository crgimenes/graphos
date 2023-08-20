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
	screen.CurrentColor = graphos.Colors16[0x0F]
	screen.DrawPix(walker.X, walker.Y, screen.CurrentColor)

	return nil
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	rand.Seed(time.Now().Unix())

	cg := graphos.New()
	cg.Width = 800
	cg.Height = 600
	cg.ScreenHandler = update
	cg.Title = "Random Walker"
	cg.CurrentColor = graphos.Colors16[0x0]

	walker = Walker{
		X: cg.Width / 2,
		Y: cg.Height / 2,
	}

	cg.Run()
}
