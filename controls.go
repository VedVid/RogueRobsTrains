/*
Copyright (c) 2018, Tomasz "VedVid" Nowakowski
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	blt "bearlibterminal"
	"os"
)

const (
	KB_QWERTY = iota
	KB_QWERTZ
	KB_AZERTY
)

var KeyMap map[rune]int

var HardcodedKeys = []int{
	blt.TK_RETURN,
	blt.TK_ENTER,
	blt.TK_ESCAPE,
	blt.TK_BACKSPACE,
	blt.TK_TAB,
	blt.TK_SPACE,
	blt.TK_PAUSE,
	blt.TK_INSERT,
	blt.TK_HOME,
	blt.TK_PAGEUP,
	blt.TK_DELETE,
	blt.TK_END,
	blt.TK_PAGEDOWN,
	blt.TK_RIGHT,
	blt.TK_LEFT,
	blt.TK_DOWN,
	blt.TK_UP,
	blt.TK_KP_DIVIDE,
	blt.TK_KP_MULTIPLY,
	blt.TK_KP_MINUS,
	blt.TK_KP_PLUS,
	blt.TK_KP_ENTER,
	blt.TK_KP_1,
	blt.TK_KP_2,
	blt.TK_KP_3,
	blt.TK_KP_4,
	blt.TK_KP_5,
	blt.TK_KP_6,
	blt.TK_KP_7,
	blt.TK_KP_8,
	blt.TK_KP_9,
	blt.TK_KP_0,
	blt.TK_KP_PERIOD,
}

var QWERTYLayoutRunesToCodes = map[rune]int{
	'q': blt.TK_Q,
	'Q': blt.TK_Q,
	'w': blt.TK_W,
	'W': blt.TK_W,
	'e': blt.TK_E,
	'E': blt.TK_E,
	'r': blt.TK_R,
	'R': blt.TK_R,
	't': blt.TK_T,
	'T': blt.TK_T,
	'y': blt.TK_Y,
	'Y': blt.TK_Y,
	'u': blt.TK_U,
	'U': blt.TK_U,
	'i': blt.TK_I,
	'I': blt.TK_I,
	'o': blt.TK_O,
	'O': blt.TK_O,
	'p': blt.TK_P,
	'P': blt.TK_P,
	'a': blt.TK_A,
	'A': blt.TK_A,
	's': blt.TK_S,
	'S': blt.TK_S,
	'd': blt.TK_D,
	'D': blt.TK_D,
	'f': blt.TK_F,
	'F': blt.TK_F,
	'g': blt.TK_G,
	'G': blt.TK_G,
	'h': blt.TK_H,
	'H': blt.TK_H,
	'j': blt.TK_J,
	'J': blt.TK_J,
	'k': blt.TK_K,
	'K': blt.TK_K,
	'l': blt.TK_L,
	'L': blt.TK_L,
	'z': blt.TK_Z,
	'Z': blt.TK_Z,
	'x': blt.TK_X,
	'X': blt.TK_X,
	'c': blt.TK_C,
	'C': blt.TK_C,
	'v': blt.TK_V,
	'V': blt.TK_V,
	'b': blt.TK_B,
	'B': blt.TK_B,
	'n': blt.TK_N,
	'N': blt.TK_N,
	'm': blt.TK_M,
	'M': blt.TK_M,
	',': blt.TK_COMMA,
	'<': blt.TK_COMMA,
	'.': blt.TK_PERIOD,
	'>': blt.TK_PERIOD,
	'/': blt.TK_SLASH,
	'?': blt.TK_SLASH,
	';': blt.TK_SEMICOLON,
	':': blt.TK_SEMICOLON,
	'\'': blt.TK_APOSTROPHE,
	'"': blt.TK_APOSTROPHE,
	'[': blt.TK_LBRACKET,
	'{': blt.TK_LBRACKET,
	']': blt.TK_RBRACKET,
	'}': blt.TK_RBRACKET,
	'1': blt.TK_1,
	'!': blt.TK_1,
	'2': blt.TK_2,
	'@': blt.TK_2,
	'3': blt.TK_3,
	'#': blt.TK_3,
	'4': blt.TK_4,
	'$': blt.TK_4,
	'5': blt.TK_5,
	'%': blt.TK_5,
	'6': blt.TK_6,
	'^': blt.TK_6,
	'7': blt.TK_7,
	'&': blt.TK_7,
	'8': blt.TK_8,
	'*': blt.TK_8,
	'9': blt.TK_9,
	'(': blt.TK_9,
	'0': blt.TK_0,
	')': blt.TK_0,
	'-': blt.TK_MINUS,
	'_': blt.TK_MINUS,
	'=': blt.TK_EQUALS,
	'+': blt.TK_EQUALS,
}

var QWERTZLayoutRunesToCodes map[rune]int

var AZERTYLayoutRunesToCodes map[rune]int

func Controls(k int, p *Creature, b *Board, c *Creatures, o *Objects) bool {
	/* Function Controls is input handler.
	   It takes integer k (key codes are basically numbers,
	   but creating new "type key int" is not convenient)
	   and Creature p (which is player).
	   Controls handle input, then returns integer value that depends
	   if player spent turn by action or not. */
	turnSpent := false
	switch k {
	case blt.TK_UP, blt.TK_KP_8, blt.TK_K, blt.TK_W:
		turnSpent = p.MoveOrAttack(0, -1, *b, o, *c)
	case blt.TK_RIGHT, blt.TK_KP_6, blt.TK_L, blt.TK_D:
		turnSpent = p.MoveOrAttack(1, 0, *b, o, *c)
	case blt.TK_DOWN, blt.TK_KP_2, blt.TK_J, blt.TK_X:
		turnSpent = p.MoveOrAttack(0, 1, *b, o, *c)
	case blt.TK_LEFT, blt.TK_KP_4, blt.TK_H, blt.TK_A:
		turnSpent = p.MoveOrAttack(-1, 0, *b, o, *c)
	case blt.TK_HOME, blt.TK_KP_7, blt.TK_Y, blt.TK_Q:
		turnSpent = p.MoveOrAttack(-1, -1, *b, o, *c)
	case blt.TK_PAGEUP, blt.TK_KP_9, blt.TK_U, blt.TK_E:
		turnSpent = p.MoveOrAttack(1, -1, *b, o, *c)
	case blt.TK_END, blt.TK_KP_1, blt.TK_B, blt.TK_Z:
		turnSpent = p.MoveOrAttack(-1, 1, *b, o, *c)
	case blt.TK_PAGEDOWN, blt.TK_KP_3, blt.TK_N, blt.TK_C:
		turnSpent = p.MoveOrAttack(1, 1, *b, o, *c)
	case blt.TK_SPACE, blt.TK_KP_5, blt.TK_PERIOD, blt.TK_S:
		turnSpent = true // Pass a turn.

	case blt.TK_F:
		if p.ActiveWeapon != SlotWeaponMelee {
			if p.Equipment[p.ActiveWeapon].AmmoCurrent <= 0 {
				AddMessage("You need to reload " + p.Equipment[p.ActiveWeapon].Name + ".")
			} else {
				if (p.Equipment[p.ActiveWeapon].Cock == true &&
					p.Equipment[p.ActiveWeapon].Cocked == true) ||
					p.Equipment[p.ActiveWeapon].Cock == false {
					turnSpent = p.Target(*b, o, *c)
					if turnSpent == true {
						p.Equipment[p.ActiveWeapon].AmmoCurrent--
						if p.Equipment[p.ActiveWeapon].Cock == true {
							p.Equipment[p.ActiveWeapon].Cocked = false
						}
					}
				} else {
					p.Equipment[p.ActiveWeapon].Cocked = true
					turnSpent = true
					AddMessage("You cocked " + p.Equipment[p.ActiveWeapon].Name + ".")
				}
			}
		} else {
			AddMessage("You are using melee weapon.")
		}
	case blt.TK_R:
		if Config.Reloading == AmmoLimited {
			AddMessage("You do not have more ammo!")
		} else {
			if p.ActiveWeapon != SlotWeaponMelee {
				if p.Equipment[p.ActiveWeapon].Cock == false {
					if p.Equipment[p.ActiveWeapon].AmmoCurrent < p.Equipment[p.ActiveWeapon].AmmoMax {
						p.Equipment[p.ActiveWeapon].AmmoCurrent = p.Equipment[p.ActiveWeapon].AmmoMax
						turnSpent = true
					}
				} else {
					if p.Equipment[p.ActiveWeapon].Cocked == true {
						p.Equipment[p.ActiveWeapon].Cocked = false
						AddMessage("You uncocked " + p.Equipment[p.ActiveWeapon].Name + ".")
						turnSpent = true
					} else {
						if p.Equipment[p.ActiveWeapon].AmmoCurrent < p.Equipment[p.ActiveWeapon].AmmoMax {
							p.Equipment[p.ActiveWeapon].AmmoCurrent++
							turnSpent = true
						}
					}
				}
			}
		}
	case blt.TK_I:
		p.Look(*b, *o, *c) // Looking is free action.
	case blt.TK_G:
		turnSpent = p.PickUp(o)
	case blt.TK_P:
		minX := p.X - 1
		if minX < 0 {
			minX = p.X
		}
		maxX := p.X + 1
		if maxX >= MapSizeX {
			maxX = p.X
		}
		minY := p.Y - 1
		if minY < 0 {
			minY = p.Y
		}
		maxY := p.Y + 1
		if maxY >= MapSizeY {
			maxY = p.Y
		}
		lever := false
		for x := minX; x <= maxX; x++ {
			for y := minY; y <= maxY; y++ {
				if (*b)[x][y].Name == "lever" {
					lever = true
					Config.Score += 10
				}
			}
		}
		if lever == false {
			AddMessage("There is no lever to pull here.")
		}
		if lever == true {
			PrintVictoryScreen()
			DeleteSaves()
			blt.Close()
			os.Exit(0)
		}
	case blt.TK_1:
		if p.ActiveWeapon != SlotWeaponPrimary {
			if p.Equipment[p.ActiveWeapon].Cock == true {
				p.Equipment[p.ActiveWeapon].Cocked = false
			}
			p.ActiveWeapon = SlotWeaponPrimary
			turnSpent = true
		}
	case blt.TK_2:
		if p.ActiveWeapon != SlotWeaponSecondary {
			if p.Equipment[p.ActiveWeapon].Cock == true {
				p.Equipment[p.ActiveWeapon].Cocked = false
			}
			p.ActiveWeapon = SlotWeaponSecondary
			turnSpent = true
		}
	case blt.TK_3:
		if p.ActiveWeapon != SlotWeaponMelee {
			if p.Equipment[p.ActiveWeapon].Cock == true {
				p.Equipment[p.ActiveWeapon].Cocked = false
			}
			p.ActiveWeapon = SlotWeaponMelee
			turnSpent = true
		}
	}
	return turnSpent
}

func ReadInput() int {
	key := blt.Read()
	for _, v := range HardcodedKeys {
		if key == v {
			return v
		}
	}
	var r rune
	if blt.Check(blt.TK_WCHAR) != 0 {
		r = rune(blt.State(blt.TK_WCHAR))
	}
	return KeyMap[r]
}

func InitializeKeyboardLayouts() {
	InitializeQWERTZ()
	InitializeAZERTY()
	switch KeyboardLayout {
	case KB_QWERTY: KeyMap = QWERTYLayoutRunesToCodes
	case KB_QWERTZ: KeyMap = QWERTZLayoutRunesToCodes
	case KB_AZERTY: KeyMap = AZERTYLayoutRunesToCodes
	}
}

func InitializeQWERTZ() {
	QWERTZLayoutRunesToCodes = QWERTYLayoutRunesToCodes
	QWERTZLayoutRunesToCodes['z'] = blt.TK_Y
	QWERTZLayoutRunesToCodes['Z'] = blt.TK_Y
	QWERTZLayoutRunesToCodes['y'] = blt.TK_Z
	QWERTZLayoutRunesToCodes['Y'] = blt.TK_Z
	QWERTZLayoutRunesToCodes[';'] = blt.TK_COMMA
	QWERTZLayoutRunesToCodes[':'] = blt.TK_PERIOD
	QWERTZLayoutRunesToCodes['-'] = blt.TK_SLASH
	QWERTZLayoutRunesToCodes['_'] = blt.TK_SLASH
	QWERTZLayoutRunesToCodes['ö'] = blt.TK_SEMICOLON
	QWERTZLayoutRunesToCodes['Ö'] = blt.TK_SEMICOLON
	QWERTZLayoutRunesToCodes['ä'] = blt.TK_APOSTROPHE
	QWERTZLayoutRunesToCodes['Ä'] = blt.TK_APOSTROPHE
	QWERTZLayoutRunesToCodes['ü'] = blt.TK_LBRACKET
	QWERTZLayoutRunesToCodes['Ü'] = blt.TK_LBRACKET
	QWERTZLayoutRunesToCodes['+'] = blt.TK_RBRACKET
	QWERTZLayoutRunesToCodes['*'] = blt.TK_RBRACKET
	QWERTZLayoutRunesToCodes['"'] = blt.TK_2
	QWERTZLayoutRunesToCodes['§'] = blt.TK_3
	QWERTZLayoutRunesToCodes['&'] = blt.TK_6
	QWERTZLayoutRunesToCodes['/'] = blt.TK_7
	QWERTZLayoutRunesToCodes['('] = blt.TK_8
	QWERTZLayoutRunesToCodes[')'] = blt.TK_9
	QWERTZLayoutRunesToCodes['='] = blt.TK_0
	QWERTZLayoutRunesToCodes['ß'] = blt.TK_MINUS
	QWERTZLayoutRunesToCodes['?'] = blt.TK_MINUS
	QWERTZLayoutRunesToCodes['´'] = blt.TK_EQUALS
	QWERTZLayoutRunesToCodes['`'] = blt.TK_EQUALS
}

func InitializeAZERTY() {
	AZERTYLayoutRunesToCodes = QWERTYLayoutRunesToCodes
	AZERTYLayoutRunesToCodes['a'] = blt.TK_Q
	AZERTYLayoutRunesToCodes['A'] = blt.TK_Q
	AZERTYLayoutRunesToCodes['z'] = blt.TK_W
	AZERTYLayoutRunesToCodes['Z'] = blt.TK_W
	AZERTYLayoutRunesToCodes['q'] = blt.TK_A
	AZERTYLayoutRunesToCodes['Q'] = blt.TK_A
	AZERTYLayoutRunesToCodes['w'] = blt.TK_Z
	AZERTYLayoutRunesToCodes['W'] = blt.TK_Z
	AZERTYLayoutRunesToCodes[','] = blt.TK_M
	AZERTYLayoutRunesToCodes['?'] = blt.TK_M
	AZERTYLayoutRunesToCodes[';'] = blt.TK_COMMA
	AZERTYLayoutRunesToCodes['.'] = blt.TK_COMMA
	AZERTYLayoutRunesToCodes[':'] = blt.TK_PERIOD
	AZERTYLayoutRunesToCodes['/'] = blt.TK_PERIOD
	AZERTYLayoutRunesToCodes['!'] = blt.TK_SLASH
	AZERTYLayoutRunesToCodes['§'] = blt.TK_SLASH
	AZERTYLayoutRunesToCodes['m'] = blt.TK_SEMICOLON
	AZERTYLayoutRunesToCodes['M'] = blt.TK_SEMICOLON
	AZERTYLayoutRunesToCodes['ù'] = blt.TK_APOSTROPHE
	AZERTYLayoutRunesToCodes['%'] = blt.TK_APOSTROPHE
	AZERTYLayoutRunesToCodes['^'] = blt.TK_LBRACKET
	AZERTYLayoutRunesToCodes['¨'] = blt.TK_LBRACKET
	AZERTYLayoutRunesToCodes['$'] = blt.TK_RBRACKET
	AZERTYLayoutRunesToCodes['£'] = blt.TK_RBRACKET
	AZERTYLayoutRunesToCodes['&'] = blt.TK_1
	AZERTYLayoutRunesToCodes['é'] = blt.TK_2
	AZERTYLayoutRunesToCodes['"'] = blt.TK_3
	AZERTYLayoutRunesToCodes['\''] = blt.TK_4
	AZERTYLayoutRunesToCodes['('] = blt.TK_5
	AZERTYLayoutRunesToCodes['-'] = blt.TK_6
	AZERTYLayoutRunesToCodes['è'] = blt.TK_7
	AZERTYLayoutRunesToCodes['_'] = blt.TK_8
	AZERTYLayoutRunesToCodes['ç'] = blt.TK_9
	AZERTYLayoutRunesToCodes['à'] = blt.TK_0
	AZERTYLayoutRunesToCodes[')'] = blt.TK_MINUS
	AZERTYLayoutRunesToCodes['°'] = blt.TK_MINUS
}
