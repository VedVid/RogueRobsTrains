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

func Controls(k int, r rune, p *Creature, b *Board, c *Creatures, o *Objects) bool {
	/* Function Controls is input handler.
	   It takes integer k (key codes are basically numbers,
	   but creating new "type key int" is not convenient),
	   r (character rune produced by the input),
	   and Creature p (which is player).
	   Controls handle input, then returns integer value that depends
	   if player spent turn by action or not. */
	turnSpent := false
	switch k {
	case blt.TK_UP, blt.TK_KP_8:
		turnSpent = p.MoveOrAttack(0, -1, *b, o, *c)
	case blt.TK_RIGHT, blt.TK_KP_6:
		turnSpent = p.MoveOrAttack(1, 0, *b, o, *c)
	case blt.TK_DOWN, blt.TK_KP_2:
		turnSpent = p.MoveOrAttack(0, 1, *b, o, *c)
	case blt.TK_LEFT, blt.TK_KP_4:
		turnSpent = p.MoveOrAttack(-1, 0, *b, o, *c)
	case blt.TK_HOME, blt.TK_KP_7:
		turnSpent = p.MoveOrAttack(-1, -1, *b, o, *c)
	case blt.TK_PAGEUP, blt.TK_KP_9:
		turnSpent = p.MoveOrAttack(1, -1, *b, o, *c)
	case blt.TK_END, blt.TK_KP_1:
		turnSpent = p.MoveOrAttack(-1, 1, *b, o, *c)
	case blt.TK_PAGEDOWN, blt.TK_KP_3:
		turnSpent = p.MoveOrAttack(1, 1, *b, o, *c)
	case blt.TK_SPACE, blt.TK_KP_5:
		turnSpent = true // Pass a turn.
	default:
		switch r {
		case 'k', 'K', 'w', 'W':
			turnSpent = p.MoveOrAttack(0, -1, *b, o, *c)
		case 'l', 'L', 'd', 'D':
			turnSpent = p.MoveOrAttack(1, 0, *b, o, *c)
		case 'j', 'J', 'x', 'X':
			turnSpent = p.MoveOrAttack(0, 1, *b, o, *c)
		case 'h', 'H', 'a', 'A':
			turnSpent = p.MoveOrAttack(-1, 0, *b, o, *c)
		case 'y', 'Y', 'q', 'Q':
			turnSpent = p.MoveOrAttack(-1, -1, *b, o, *c)
		case 'u', 'U', 'e', 'E':
			turnSpent = p.MoveOrAttack(1, -1, *b, o, *c)
		case 'b', 'B', 'z', 'Z':
			turnSpent = p.MoveOrAttack(-1, 1, *b, o, *c)
		case 'n', 'N', 'c', 'C':
			turnSpent = p.MoveOrAttack(1, 1, *b, o, *c)
		case '.', 's', 'S':
			turnSpent = true // Pass a turn.
		case 'f', 'F':
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
		case 'r', 'R':
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
		case 'i', 'I':
			p.Look(*b, *o, *c) // Looking is free action.
		case 'g', 'G':
			turnSpent = p.PickUp(o)
		case 'p', 'P':
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
		case '1':
			if p.ActiveWeapon != SlotWeaponPrimary {
				if p.Equipment[p.ActiveWeapon].Cock == true {
					p.Equipment[p.ActiveWeapon].Cocked = false
				}
				p.ActiveWeapon = SlotWeaponPrimary
				turnSpent = true
			}
		case '2':
			if p.ActiveWeapon != SlotWeaponSecondary {
				if p.Equipment[p.ActiveWeapon].Cock == true {
					p.Equipment[p.ActiveWeapon].Cocked = false
				}
				p.ActiveWeapon = SlotWeaponSecondary
				turnSpent = true
			}
		case '3':
			if p.ActiveWeapon != SlotWeaponMelee {
				if p.Equipment[p.ActiveWeapon].Cock == true {
					p.Equipment[p.ActiveWeapon].Cocked = false
				}
				p.ActiveWeapon = SlotWeaponMelee
				turnSpent = true
			}
		}
	}
	return turnSpent
}
