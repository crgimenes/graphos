package main

import (
	"fmt"

	"crg.eti.br/go/graphos"
)

func update(i *graphos.Instance) error {
	i.CurrentColor = 0x00
	i.Clear() // Clear screen (with border)
	i.CurrentColor = 0x0F

	if i.Machine == 0 {
		i.Machine++
		i.Println("          1         2         3         4         5         6         7")
		i.Println("01234567890123456789012345678901234567890123456789012345678901234567890123456789")
		i.Println("terminal v0.01")
		i.Println("https://crg.eti.br")
		i.Println(fmt.Sprintf("Width: %v, Height: %v", i.Width, i.Height))
		i.Println("")

		var c byte
		for ; c < 246; c++ {
			i.PutChar(c)
		}
		i.Println("")

	}

	i.DrawVideoTextMode()

	i.Input()
	return nil
}

func main() {

	cg := graphos.New()
	//cg.Border = 40
	cg.Width = 800
	cg.Height = 600
	cg.Scale = 1
	cg.ScreenHandler = update
	cg.Title = "Graphos - Terminal"
	cg.CurrentColor = 0x0F

	cg.Run()

}
