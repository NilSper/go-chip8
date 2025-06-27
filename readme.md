## CHIP 8 emulator/interpreter


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
