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
)
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

var QWERTYLayoutRunesToCodes = map[[2]rune]int{
	[...]rune{'q', 'Q'}: blt.TK_Q,
	[...]rune{'w', 'W'}: blt.TK_W,
	[...]rune{'e', 'E'}: blt.TK_E,
	[...]rune{'r', 'R'}: blt.TK_R,
	[...]rune{'t', 'T'}: blt.TK_T,
	[...]rune{'y', 'Y'}: blt.TK_Y,
	[...]rune{'u', 'U'}: blt.TK_U,
	[...]rune{'i', 'I'}: blt.TK_I,
	[...]rune{'o', 'O'}: blt.TK_O,
	[...]rune{'p', 'P'}: blt.TK_P,
	[...]rune{'a', 'A'}: blt.TK_A,
	[...]rune{'s', 'S'}: blt.TK_S,
	[...]rune{'d', 'D'}: blt.TK_D,
	[...]rune{'f', 'F'}: blt.TK_F,
	[...]rune{'g', 'G'}: blt.TK_G,
	[...]rune{'h', 'H'}: blt.TK_H,
	[...]rune{'j', 'J'}: blt.TK_J,
	[...]rune{'k', 'K'}: blt.TK_K,
	[...]rune{'l', 'L'}: blt.TK_L,
	[...]rune{'z', 'Z'}: blt.TK_Z,
	[...]rune{'x', 'X'}: blt.TK_X,
	[...]rune{'c', 'C'}: blt.TK_C,
	[...]rune{'v', 'V'}: blt.TK_V,
	[...]rune{'b', 'B'}: blt.TK_B,
	[...]rune{'n', 'N'}: blt.TK_N,
	[...]rune{'m', 'M'}: blt.TK_M,
	[...]rune{',', '<'}: blt.TK_COMMA,
	[...]rune{'.', '>'}: blt.TK_PERIOD,
	[...]rune{';', ':'}: blt.TK_SEMICOLON,
	[...]rune{'\'', '"'}: blt.TK_APOSTROPHE,
	[...]rune{'[', '{'}: blt.TK_LBRACKET,
	[...]rune{']', '}'}: blt.TK_RBRACKET,
	[...]rune{'1', '!'}: blt.TK_1,
	[...]rune{'2', '@'}: blt.TK_2,
	[...]rune{'3', '#'}: blt.TK_3,
	[...]rune{'4', '$'}: blt.TK_4,
	[...]rune{'5', '%'}: blt.TK_5,
	[...]rune{'6', '^'}: blt.TK_6,
	[...]rune{'7', '&'}: blt.TK_7,
	[...]rune{'8', '*'}: blt.TK_8,
	[...]rune{'9', '('}: blt.TK_9,
	[...]rune{'0', ')'}: blt.TK_0,
	[...]rune{'-', '_'}: blt.TK_MINUS,
	[...]rune{'=', '+'}: blt.TK_EQUALS,
}

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
	var keyMap map[[2]rune]int
	switch KeyboardLayout {
	case KB_QWERTY: keyMap = QWERTYLayoutRunesToCodes
	}
	for k, v := range keyMap {
		if k[0] == r || k[1] == r {
			return v
		}
	}
	return -1 //wrong value!
}
