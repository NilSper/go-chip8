package main

import (
	"chip8/helpers"
	"fmt"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	fmt.Println("Starting chip8 emulation..")

	var drawFlag bool
	var clearFlag bool

	var path = "c8games/PONG"
	var err error

	var window *sdl.Window
	var renderer *sdl.Renderer
	err = helpers.SetupGraphicsAndAudio(&window, &renderer)
	defer sdl.Quit()
	defer window.Destroy()
	defer renderer.Destroy()
	defer mix.Quit()
	defer mix.CloseAudio()
	if err != nil {
		fmt.Println("program crashed when setting up graphics \n--> ", err)
		return
	}

	helpers.InitializeStates()
	if err = helpers.LoadGame(path); err != nil {
		fmt.Println("program crashed when loading game \n--> ", err)
		return
	}

	go helpers.HandleKeyboardInput()

	// Emulation loop
	for {
		// Emulate one cycle

		if err = helpers.EmulateCycle(&clearFlag, &drawFlag); err != nil {
			fmt.Println("program crashed when emulating a cycle \n--> ", err)
			return
		}

		if drawFlag {
			helpers.DrawGraphics(renderer)
			drawFlag = false
		}

		if clearFlag {
			helpers.ClearScreen(renderer)
			clearFlag = false
		}
	}

}
