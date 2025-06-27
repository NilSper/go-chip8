package main

import (
	"fmt"
	"math/rand/v2"
	"os"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

// SDL lib / screen handling

const winTitle string = "CHIP8 Emulator"
const winWidth, winHeight int32 = 64, 32
const modifier int32 = 10

var sdlWidth = winWidth * modifier
var sdlHeight = winHeight * modifier

func clearScreen(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	renderer.Present()
}

// GAME STATE
var opcode uint16
var memory [4096]uint8

var vRegisters [16]uint8

var indexRegister uint16
var programCounter uint16

// gfx
var screenState [winWidth * winHeight]uint8

var delayTimer uint8
var soundTimer uint8

var stack [16]uint16
var stackPointer uint16

var keys [16]uint8

var drawFlag bool
var clearFlag bool

const fontSetOffset = 80

var fontSet = [80]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func runKeyboardInput() {
	// var window *sdl.Window

	// if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
	// 	return err
	// }
	// defer sdl.Quit()

	// window, err = sdl.CreateWindow("Input", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	// if err != nil {
	// 	return err
	// }
	// defer window.Destroy()

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

func setupGraphics(window **sdl.Window, renderer **sdl.Renderer) {
	var err error
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %s\n", err)
	}

	if *window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, sdlWidth, sdlHeight, sdl.WINDOW_SHOWN); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
	}

	if *renderer, err = sdl.CreateRenderer(*window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
	}
	clearScreen(*renderer)
	sdl.Delay(1000)
}

func drawGraphics(renderer *sdl.Renderer) {
	// fmt.Println("Drawing the picture", screenState[1])
	var white = sdl.Color{R: 255, G: 255, B: 255, A: 255}

	clearScreen(renderer)
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

func main() {
	fmt.Println("Starting chip8 emulation..")

	var path = "c8games/PONG" //"c8games/4-flags.ch8" //

	var window *sdl.Window
	var renderer *sdl.Renderer
	setupGraphics(&window, &renderer)
	defer sdl.Quit()
	defer window.Destroy()
	defer renderer.Destroy()

	// setupInput()

	// TODO: need to setup keyboard input ... as with graphics + check each cycle

	go runKeyboardInput()
	// if err := runKeyboardInput(); err != nil {
	// 	os.Exit(1)
	// }

	initialize()
	loadGame(path)

	// Emulation loop
	for {
		// Emulate one cycle
		emulateCycle()

		if drawFlag {
			drawGraphics(renderer)
			drawFlag = false
		}

		if clearFlag {
			clearScreen(renderer)
			clearFlag = false
		}

		// Store the key press state
		// setKeys()
	}

}

// func setKeys() {
// 	panic("unimplemented")
// }

// func setupInput() {
// 	panic("unimplemented")
// }

func emulateCycle() {
	// Fetch Opcode

	// each memory state is 8 bits. An Opcode consists of 16 bits. Thus, we need to merge
	// to consecutive memory states into the correct opcode
	// Casting both memory states to 16 bits before merging them with bit-wise operations
	opcode = (uint16)(memory[programCounter])<<8 | (uint16)(memory[programCounter+1])

	// Decode + execute Opcode
	// fmt.Printf("PC: %X - opcode: 0x%X\n", programCounter, opcode)
	switch opcode {
	case 0x00E0:
		// clear display
		screenState = [len(screenState)]uint8{}
		clearFlag = true
		programCounter += 2
	case 0x00EE:
		stackPointer--
		programCounter = stack[stackPointer]
		programCounter += 2
	default:
		switch opcode & 0xF000 {
		case 0xA000: // ANNN
			indexRegister = opcode & 0x0FFF
			programCounter += 2
		case 0xB000: // BNNN
			programCounter = (opcode & 0x0FFF) + uint16(vRegisters[0])
		case 0x1000: // 1NNN
			// jump to address NNN --> no need to store current location
			programCounter = opcode & 0x0FFF
			if opcode == 0x1228 {
				sdl.Delay(3000)
				panic("final opcode")
			}
		case 0x2000: // 2NNN
			// call subroutine at NNN
			stack[stackPointer] = programCounter
			stackPointer++
			programCounter = opcode & 0x0FFF
		case 0x3000: // 3XNN
			var xValue = (opcode & 0x0F00) >> 8
			var conditionValue = opcode & 0x00FF
			if vRegisters[xValue] == uint8(conditionValue) {
				programCounter += 4
			} else {
				programCounter += 2
			}
		case 0x4000: // 4XNN
			var xValue = (opcode & 0x0F00) >> 8
			var conditionValue = opcode & 0x00FF
			if vRegisters[xValue] != uint8(conditionValue) {
				programCounter += 4
			} else {
				programCounter += 2
			}
		case 0x5000: // 5XY0
			// check last is 0
			var xValue = (opcode & 0x0F00) >> 8
			var yValue = (opcode & 0x00F0) >> 4
			if vRegisters[xValue] == vRegisters[yValue] {
				programCounter += 4
			} else {
				programCounter += 2
			}
		case 0x6000: // 6XNN
			value := uint8(opcode & 0x00FF)
			vRegisters[(opcode&0x0F00)>>8] = value
			programCounter += 2
		case 0x7000: // 7XNN
			value := uint8(opcode & 0x00FF)
			vRegisters[(opcode&0x0F00)>>8] += value
			programCounter += 2
		case 0x8000:
			var xValue = (opcode & 0x0F00) >> 8
			var yValue = (opcode & 0x00F0) >> 4
			switch opcode & 0x000F {
			case 0x0000:
				vRegisters[xValue] = vRegisters[yValue]
				programCounter += 2
			case 0x0001:
				vRegisters[xValue] |= vRegisters[yValue]
				programCounter += 2
			case 0x0002:
				vRegisters[xValue] &= vRegisters[yValue]
				programCounter += 2
			case 0x0003:
				vRegisters[xValue] ^= vRegisters[yValue]
				programCounter += 2
			case 0x0004:
				var checkValue = uint16(vRegisters[xValue]) + uint16(vRegisters[yValue])
				if checkValue > uint16(vRegisters[xValue]) {
					vRegisters[0xF] = 1
				} else {
					vRegisters[0xF] = 0
				}
				vRegisters[xValue] += vRegisters[yValue]
				programCounter += 2
			case 0x0005:
				if vRegisters[xValue] >= vRegisters[yValue] {
					vRegisters[0xF] = 1
				} else {
					vRegisters[0xF] = 0
				}
				vRegisters[xValue] -= vRegisters[yValue]
				programCounter += 2
			case 0x0007:
				if vRegisters[yValue] >= vRegisters[xValue] {
					vRegisters[0xF] = 1
				} else {
					vRegisters[0xF] = 0
				}
				vRegisters[xValue] = vRegisters[yValue] - vRegisters[xValue]
				programCounter += 2
			case 0x0006:
				// TODO config if Set VX to the value of VY
				fmt.Print("Vx: ", vRegisters[xValue], ", Vy: ", vRegisters[yValue], "\n")
				vRegisters[0xF] = vRegisters[xValue] & 1 // least sig bit
				vRegisters[xValue] >>= 1
				programCounter += 2
			case 0x000E:
				// TODO config if Set VX to the value of VY
				fmt.Print("Vx: ", vRegisters[xValue], ", Vy: ", vRegisters[yValue], "\n")
				vRegisters[0xF] = vRegisters[xValue] >> 7 // most sig bit
				vRegisters[xValue] <<= 1
				programCounter += 2
			}
		case 0x9000: // 9XY0 counter to 5XY0
			var xValue = (opcode & 0x0F00) >> 8
			var yValue = (opcode & 0x00F0) >> 4
			if vRegisters[xValue] != vRegisters[yValue] {
				programCounter += 4
			} else {
				programCounter += 2
			}
		case 0xC000: // CXNN
			var xValue = (opcode & 0x0F00) >> 8
			var operand = uint8(opcode & 0x00FF)
			vRegisters[xValue] = uint8(rand.IntN(255)) & operand
			programCounter += 2
		case 0xD000: // DXYN draw to screen
			// read X and Y at 3. and 4. position -> remove buffer 0s
			var xIndex = int32(vRegisters[(opcode&0x0F00)>>8])
			var yIndex = int32(vRegisters[(opcode&0x00F0)>>4])

			var rowNumber = int32(opcode & 0x000F)

			vRegisters[0xF] = 0 // reset VF TODO understand why
			for yLine := int32(0); yLine < rowNumber; yLine++ {
				pixel := memory[int32(indexRegister)+yLine]
				// fixed size of 8 per row
				for xLine := int32(0); xLine < 8; xLine++ {
					if (pixel & (0x80 >> xLine)) != 0 {
						screenIndex := xIndex + xLine + ((yIndex + yLine) * winWidth)
						if screenIndex > winWidth*winHeight {
							screenIndex = screenIndex % (winWidth * winHeight)
						}
						if screenState[screenIndex] == 1 {
							vRegisters[0xF] = 1
						}
						screenState[screenIndex] ^= 1
					}
				}
			}
			// printScreenState()
			drawFlag = true
			programCounter += 2
		case 0xE000: // Key operations
			var xValue = (opcode & 0x0F00) >> 8
			var lowestNibble = vRegisters[xValue] & 0x0F
			switch opcode & 0x00FF {
			case 0x009E: // key is pressed
				if keys[lowestNibble] == 1 {
					programCounter += 4
				} else {
					programCounter += 2
				}
			case 0x00A1: // key is not pressed
				if keys[lowestNibble] == 0 {
					programCounter += 4
				} else {
					programCounter += 2
				}
			}
		case 0xF000: // FX??
			var xValue = (opcode & 0x0F00) >> 8
			switch opcode & 0x00FF {
			case 0x0007:
				vRegisters[xValue] = delayTimer
				programCounter += 2
			case 0x000A:
				// TODO await key press, halt opcodes but continue timers
				// don't increase programCounter until key is set !
				vRegisters[xValue] = getKey()
				panic("Opcode not finished")
			case 0x0015:
				delayTimer = vRegisters[xValue]
				programCounter += 2
			case 0x0018:
				soundTimer = vRegisters[xValue]
				programCounter += 2
			case 0x001E:
				indexRegister += uint16(vRegisters[xValue])
				programCounter += 2
			case 0x0029:
				indexRegister = spriteAddr(vRegisters[xValue])
				programCounter += 2
			case 0x0033:
				// binary-coded decimal representation
				// value 0x9c -> 156: set each digit to on memory location
				memory[indexRegister] = vRegisters[xValue] / 100
				memory[indexRegister+1] = vRegisters[xValue] % 100 / 10
				memory[indexRegister+2] = vRegisters[xValue] % 100 % 10
				programCounter += 2
			case 0x0055:
				// TODO: old instruction would actual set indexRegister += xValue
				for i := range xValue {
					memory[indexRegister+i] = vRegisters[i]
				}
				programCounter += 2
			case 0x0065:
				for i := range xValue {
					vRegisters[i] = memory[indexRegister+i]
				}
				programCounter += 2
			default:
				fmt.Printf("Unknown 0xF... opcode 0x%X\n", opcode)
				panic("unknown code")
			}
		default:
			fmt.Printf("Unknown opcode 0x%X\n", opcode)
			panic("unknown code")
		}
	}

	// Update timers
	// they have a frequency of 60Hz --> 60 ticks per second
	// add some delay before decreasing the timers - to keep 60Hz
	if delayTimer > 0 {
		delayTimer--
	}

	if soundTimer > 0 {
		if soundTimer == 1 {
			fmt.Println("BEEP")
		}
		soundTimer--
	}

}

func getKey() uint8 {
	return 0
}

func spriteAddr(vx uint8) uint16 {
	var value = uint16(vx & 0x0F)

	// width per font: 5
	return fontSetOffset + value*5
}

// func printScreenState() {

// 	fmt.Print("\n")
// 	var oldRow int32 = 0
// 	for i := range int32(len(screenState)) {
// 		var row = i / winWidth
// 		// var col = i % winWidth

// 		if row != oldRow {
// 			// add line break before each new row
// 			fmt.Print("\n")
// 			oldRow = row
// 		}
// 		fmt.Print(screenState[i])

// 	}

// 	fmt.Print("\n")
// }

func initialize() {
	// Initializes registers and memory once
	fmt.Println("Init the state")

	programCounter = 0x200 // Program counter starts at 0x200
	opcode = 0             // Reset opcode, indexRegister and stackPointer
	indexRegister = 0
	stackPointer = 0

	// Clear display
	screenState = [len(screenState)]uint8{}
	// Clear stack
	stack = [len(stack)]uint16{}
	// Clear registers V0-VF
	vRegisters = [len(vRegisters)]uint8{}
	// Clear memory
	memory = [len(memory)]uint8{}

	// Clear key states
	keys = [len(keys)]uint8{}

	// Load font set
	for i := range len(fontSet) {
		memory[i+fontSetOffset] = fontSet[i]
	}

}

func loadGame(path string) {
	fmt.Println("loading the game: ", path)

	// read entire file content as binary array

	buffer, err := os.ReadFile(path)
	if err != nil {
		// TODO: proper error handling
		panic(err)
	}

	for i := range buffer {
		memory[i+512] = buffer[i]
	}
}
