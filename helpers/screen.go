package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

// SDL lib / screen handling

func HandleKeyboardInput() {
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
				fmt.Println("Quit event")
			case *sdl.KeyboardEvent:
				scanCode := t.Keysym.Scancode
				// keyCode := t.Keysym.Sym
				// var key = ""

				var isKeyPressed uint8 = 0
				if t.State == sdl.PRESSED {
					isKeyPressed = 1
				}
				switch scanCode {
				case sdl.SCANCODE_1:
					keys[0x1] = isKeyPressed
				case sdl.SCANCODE_2:
					keys[0x2] = isKeyPressed
				case sdl.SCANCODE_3:
					keys[0x3] = isKeyPressed
				case sdl.SCANCODE_4:
					keys[0xC] = isKeyPressed
				case sdl.SCANCODE_Q:
					keys[0x4] = isKeyPressed
				case sdl.SCANCODE_W:
					keys[0x5] = isKeyPressed
				case sdl.SCANCODE_E:
					keys[0x6] = isKeyPressed
				case sdl.SCANCODE_R:
					keys[0xD] = isKeyPressed
				case sdl.SCANCODE_A:
					keys[0x7] = isKeyPressed
				case sdl.SCANCODE_S:
					keys[0x8] = isKeyPressed
				case sdl.SCANCODE_D:
					keys[0x9] = isKeyPressed
				case sdl.SCANCODE_F:
					keys[0xE] = isKeyPressed
				case sdl.SCANCODE_Z:
					keys[0xA] = isKeyPressed
				case sdl.SCANCODE_X:
					keys[0x0] = isKeyPressed
				case sdl.SCANCODE_C:
					keys[0xB] = isKeyPressed
				case sdl.SCANCODE_V:
					keys[0xF] = isKeyPressed
				}

				// if keyCode < 10000 {
				// 	if key != "" {
				// 		key += " + "
				// 	}

				// 	// If the key is held down, this will fire
				// 	if t.State == sdl.PRESSED {
				// 		key += string(keyCode) + " pressed"
				// 	}

				// }

				// if key != "" {
				// 	fmt.Printf("Scancode: %d ", scanCode)
				// 	fmt.Println(key)
				// }
			}
		}

		sdl.Delay(16)
	}
}

func ClearScreen(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Present()
}

func SetupGraphicsAndAudio(window **sdl.Window, renderer **sdl.Renderer) error {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
		return fmt.Errorf("failed to initialize SDL: \n--> %w", err)
	}

	if *window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, sdlWidth, sdlHeight, sdl.WINDOW_SHOWN); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return fmt.Errorf("failed to create window: \n--> %w", err)
	}

	if *renderer, err = sdl.CreateRenderer(*window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return fmt.Errorf("failed to create renderer: \n--> %w", err)
	}
	ClearScreen(*renderer)

	if err := mix.Init(mix.INIT_MP3); err != nil {
		return fmt.Errorf("failed to init mp3: \n--> %w", err)
	}

	if err := mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		return fmt.Errorf("failed to open audio: \n--> %w", err)
	}
	sdl.Delay(1000)
	return nil
}

func DrawGraphics(renderer *sdl.Renderer) {
	var white = sdl.Color{R: 255, G: 255, B: 255, A: 255}

	ClearScreen(renderer)
	for i := range int32(len(screenState)) {
		if screenState[i] != 0 {
			var row = i / winWidth
			var col = i % winWidth

			// need to multiply with modifier to reach SDL Screen scale
			var x1 = col * modifier
			var y1 = row * modifier

			var x2 = x1 + modifier
			var y2 = y1 + modifier
			gfx.RoundedBoxColor(renderer, x1, y1, x2, y2, 1, white)

		}
	}

	renderer.Present()
	sdl.Delay(16)
}

func playBeep() error {
	if music, err := mix.LoadMUS("helpers/test.wav"); err != nil {
		return fmt.Errorf("failed to open file: \n--> %w", err)
	} else if err := music.Play(1); err != nil {
		return fmt.Errorf("failed to play music: \n--> %w", err)
	} else {
		time.Sleep(200 * time.Millisecond)
		music.Free()
	}
	fmt.Println("BEEP")
	return nil
}
