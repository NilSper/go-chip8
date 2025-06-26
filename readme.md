## CHIP 8 emulator/interpreter



debug:
draw the ball:
PC: 274 - opcode: 0xD671


next location draw:
PC: 274 - opcode: 0xD671
PC: 276 - opcode: 0x122A - jump back to beginning

PC: 22A - opcode: 0xA2EA
PC: 22C - opcode: 0xDAB6 - draw
PC: 22E - opcode: 0xDCD6 - draw
PC: 230 - opcode: 0x6001
PC: 232 - opcode: 0xE0A1 - key input
PC: 236 - opcode: 0x6004
PC: 238 - opcode: 0xE0A1 - key input
PC: 23C - opcode: 0x601F
PC: 23E - opcode: 0x8B02
PC: 240 - opcode: 0xDAB6 - draw
PC: 242 - opcode: 0x600C
PC: 244 - opcode: 0xE0A1 - key input
PC: 248 - opcode: 0x600D
PC: 24A - opcode: 0xE0A1 - key input
PC: 24E - opcode: 0x601F
PC: 250 - opcode: 0x8D02
PC: 252 - opcode: 0xDCD6 - draw
PC: 254 - opcode: 0xA2F0
PC: 256 - opcode: 0xD671 - draw
PC: 258 - opcode: 0x8684
PC: 25A - opcode: 0x8794
PC: 25C - opcode: 0x603F
PC: 25E - opcode: 0x8602
PC: 260 - opcode: 0x611F
PC: 262 - opcode: 0x8712
PC: 264 - opcode: 0x4602
PC: 268 - opcode: 0x463F
PC: 26C - opcode: 0x471F
PC: 270 - opcode: 0x4700
PC: 274 - opcode: 0xD671


### get key

use scancodes instead of key string values


Keypad                   Keyboard
+-+-+-+-+                +-+-+-+-+
|1|2|3|C|                |1|2|3|4|
+-+-+-+-+                +-+-+-+-+
|4|5|6|D|                |Q|W|E|R|
+-+-+-+-+       =>       +-+-+-+-+
|7|8|9|E|                |A|S|D|F|
+-+-+-+-+                +-+-+-+-+
|A|0|B|F|                |Z|X|C|V|
+-+-+-+-+                +-+-+-+-+


Scancodes in correct order:
1E -> 1
1F -> 2
20 -> 3
21 -> C
14 -> 4
1A -> 5
8  -> 6
15 -> D
4  -> 7
16 -> 8
7  -> 9
9  -> E
1D -> A
1B -> 0
6  -> B
19 -> F


case 1E: keys[1] = 1
case 1F: keys[2] = 1
case 20: keys[3] = 1
case 21: keys[C] = 1
case 14: keys[4] = 1
case 1A: keys[5] = 1
case 8 : keys[6] = 1
case 15: keys[D] = 1
case 4 : keys[7] = 1
case 16: keys[8] = 1
case 7 : keys[9] = 1
case 9 : keys[E] = 1
case 1D: keys[A] = 1
case 1B: keys[0] = 1
case 6 : keys[B] = 1
case 19: keys[F] = 1




### Wrong opcodes:

according to corax tests:
FX33
FX55
FX65
