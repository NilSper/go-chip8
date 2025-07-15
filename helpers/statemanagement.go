package helpers

import (
	"fmt"
	"math/rand/v2"
	"os"
)

// GAME STATE
var opcode uint16
var memory [4096]uint8

var vRegisters [16]uint8

var indexRegister uint16
var programCounter uint16

var stack [16]uint16
var stackPointer uint16

var keys [16]uint8

var delayTimer uint8
var soundTimer uint8

// Screen state
const winTitle string = "CHIP8 Emulator"
const winWidth, winHeight int32 = 64, 32
const modifier int32 = 10

var screenState [winWidth * winHeight]uint8

var sdlWidth = winWidth * modifier
var sdlHeight = winHeight * modifier

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

func EmulateCycle(clearFlag *bool, drawFlag *bool) error {
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
		*clearFlag = true
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
			default:
				return fmt.Errorf("unknown opcode 0x%X", opcode)
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

			vRegisters[0xF] = 0 // reset VF
			for yLine := int32(0); yLine < rowNumber; yLine++ {
				pixel := memory[int32(indexRegister)+yLine]
				// fixed size of 8 per row
				for xLine := int32(0); xLine < 8; xLine++ {
					if (pixel & (0x80 >> xLine)) != 0 {
						screenIndex := xIndex + xLine + ((yIndex + yLine) * winWidth)
						// Avoid index out of bounds
						if screenIndex >= winHeight*winHeight {
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
			*drawFlag = true
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
			default:
				return fmt.Errorf("unknown opcode 0x%X", opcode)
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
				return fmt.Errorf("unknown 0xF... opcode 0x%X", opcode)
			}
		default:
			return fmt.Errorf("unknown opcode 0x%X", opcode)
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
			go playBeep()
		}
		soundTimer--
	}
	return nil

}

func getKey() uint8 {
	fmt.Println("Wanted to getKey but wasn't implemented properly yet")
	return 0
}

func spriteAddr(vx uint8) uint16 {
	var value = uint16(vx & 0x0F)

	// width per font: 5
	return fontSetOffset + value*5
}

// For debugging only
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

func InitializeStates() {
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

func LoadGame(path string) (err error) {
	fmt.Println("loading the game: ", path)

	// read entire file content as binary array
	buffer, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("loadGame failed while trying to read file %s \n--> %w", path, err)
	}

	for i := range buffer {
		memory[i+512] = buffer[i]
	}

	return nil
}
