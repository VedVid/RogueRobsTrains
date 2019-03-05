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

func Controls(k int, p *Creature, b *Board, c *Creatures, o *Objects) bool {
	/* Function Controls is input handler.
	   It takes integer k (key codes are basically numbers,
	   but creating new "type key int" is not convenient)
	   and Creature p (which is player).
	   Controls handle input, then returns integer value that depends
	   if player spent turn by action or not. */
	turnSpent := false
	switch k {
	case blt.TK_UP:
		turnSpent = p.MoveOrAttack(0, -1, *b, o, *c)
	case blt.TK_RIGHT:
		turnSpent = p.MoveOrAttack(1, 0, *b, o, *c)
	case blt.TK_DOWN:
		turnSpent = p.MoveOrAttack(0, 1, *b, o, *c)
	case blt.TK_LEFT:
		turnSpent = p.MoveOrAttack(-1, 0, *b, o, *c)

	case blt.TK_F:
		if p.ActiveWeapon != SlotWeaponMelee {
			if p.Equipment[p.ActiveWeapon].AmmoCurrent <= 0 {
				AddMessage("You need to reload!")
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
					AddMessage("Gun cocked.")
				}
			}
		} else {
			AddMessage("You are using melee weapon.")
		}
	case blt.TK_R:
		if p.ActiveWeapon != SlotWeaponMelee {
			if p.Equipment[p.ActiveWeapon].Cock == false {
				if p.Equipment[p.ActiveWeapon].AmmoCurrent < p.Equipment[p.ActiveWeapon].AmmoMax {
					p.Equipment[p.ActiveWeapon].AmmoCurrent = p.Equipment[p.ActiveWeapon].AmmoMax
					turnSpent = true
				}
			} else {
				if p.Equipment[p.ActiveWeapon].Cocked == true {
					p.Equipment[p.ActiveWeapon].Cocked = false
					AddMessage("Gun uncocked.")
					turnSpent = true
				} else {
					if p.Equipment[p.ActiveWeapon].AmmoCurrent < p.Equipment[p.ActiveWeapon].AmmoMax {
						p.Equipment[p.ActiveWeapon].AmmoCurrent++
					}
				}
			}
		}
	case blt.TK_L:
		p.Look(*b, *o, *c) // Looking is free action.
	case blt.TK_G:
		turnSpent = p.PickUp(o)
	case blt.TK_P:
		minX := p.X-1
		if minX < 0 {
			minX = p.X
		}
		maxX := p.X+1
		if maxX >= MapSizeX {
			maxX = p.X
		}
		minY := p.Y-1
		if minY < 0 {
			minY = p.Y
		}
		maxY := p.Y+1
		if maxY >= MapSizeY {
			maxY = p.Y
		}
		lever := false
		for x := minX; x <= maxX; x++ {
			for y := minY; y <= maxY; y++ {
				if (*b)[x][y].Name == "lever" {
					lever = true
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
			p.ActiveWeapon = SlotWeaponMelee
			turnSpent = true
		}
	}
	return turnSpent
}
