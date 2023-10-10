package main

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/yankooo/tetris-go/tetris"
)

func main() {

	tg := tetris.NewGame()
	pixelgl.Run(func() {
		tg.Initialize()
		tg.Run()
	})
}
