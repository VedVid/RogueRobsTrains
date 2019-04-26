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
	StrMoveNorthwest = "MOVE_NORTHWEST"
	StrMoveNorth     = "MOVE_NORTH"
	StrMoveNortheast = "MOVE_NORTHEAST"
	StrMoveWest      = "MOVE_WEST"
	StrStandStill    = "STAND_STILL"
	StrMoveEast      = "MOVE_EAST"
	StrMoveSouthwest = "MOVE_SOUTHWEST"
	StrMoveSouth     = "MOVE_SOUTH"
	StrMoveSoutheast = "MOVE_SOUTHEAST"

	StrFire    = "FIRE"
	StrReload  = "RELOAD"
	StrInspect = "INSPECT"
	StrPickup  = "PICKUP"
	StrPull    = "PULL"

	StrPrimary   = "PRIMARY"
	StrSecondary = "SECONDARY"
	StrMelee     = "MELEE"
)

var Actions = []string{
	StrMoveNorthwest,
	StrMoveNorth,
	StrMoveNortheast,
	StrMoveWest,
	StrStandStill,
	StrMoveEast,
	StrMoveSouthwest,
	StrMoveSouth,
	StrMoveSoutheast,
	StrFire,
	StrReload,
	StrInspect,
	StrPickup,
	StrPull,
	StrPrimary,
	StrSecondary,
	StrMelee,
}

var CommandKeys = map[int]string{
	blt.TK_UP:        StrMoveNorth,
	blt.TK_KP_8:      StrMoveNorth,
	blt.TK_K:         StrMoveNorth,
	blt.TK_W:         StrMoveNorth,
	blt.TK_RIGHT:     StrMoveEast,
	blt.TK_KP_6:      StrMoveEast,
	blt.TK_L:         StrMoveEast,
	blt.TK_D:         StrMoveEast,
	blt.TK_DOWN:      StrMoveSouth,
	blt.TK_KP_2:      StrMoveSouth,
	blt.TK_J:         StrMoveSouth,
	blt.TK_X:         StrMoveSouth,
	blt.TK_LEFT:      StrMoveWest,
	blt.TK_KP_4:      StrMoveWest,
	blt.TK_H:         StrMoveWest,
	blt.TK_A:         StrMoveWest,
	blt.TK_HOME:      StrMoveNorthwest,
	blt.TK_KP_7:      StrMoveNorthwest,
	blt.TK_Y:         StrMoveNorthwest,
	blt.TK_Q:         StrMoveNorthwest,
	blt.TK_PAGEUP:    StrMoveNortheast,
	blt.TK_KP_9:      StrMoveNortheast,
	blt.TK_U:         StrMoveNortheast,
	blt.TK_E:         StrMoveNortheast,
	blt.TK_END:       StrMoveSouthwest,
	blt.TK_KP_1:      StrMoveSouthwest,
	blt.TK_B:         StrMoveSouthwest,
	blt.TK_Z:         StrMoveSouthwest,
	blt.TK_PAGEDOWN:  StrMoveSoutheast,
	blt.TK_KP_3:      StrMoveSoutheast,
	blt.TK_N:         StrMoveSoutheast,
	blt.TK_C:         StrMoveSoutheast,
	blt.TK_SPACE:     StrStandStill,
	blt.TK_KP_5:      StrStandStill,
	blt.TK_KP_PERIOD: StrStandStill,
	blt.TK_S:         StrStandStill,
	blt.TK_F:         StrFire,
	blt.TK_R:         StrReload,
	blt.TK_I:         StrInspect,
	blt.TK_G:         StrPickup,
	blt.TK_P:         StrPull,
	blt.TK_1:         StrPrimary,
	blt.TK_2:         StrSecondary,
	blt.TK_3:         StrMelee,
}

var CustomCommandKeys = map[int]string{}

func Command(com string, p *Creature, b *Board, c *Creatures, o *Objects) bool {
	turnSpent := false
	switch com {
	case StrMoveNorthwest:
		turnSpent = p.MoveOrAttack(-1, -1, *b, o, *c)
	case StrMoveNorth:
		turnSpent = p.MoveOrAttack(0, -1, *b, o, *c)
	case StrMoveNortheast:
		turnSpent = p.MoveOrAttack(1, -1, *b, o, *c)
	case StrMoveWest:
		turnSpent = p.MoveOrAttack(-1, 0, *b, o, *c)
	case StrStandStill:
		turnSpent = true
	case StrMoveEast:
		turnSpent = p.MoveOrAttack(1, 0, *b, o, *c)
	case StrMoveSouthwest:
		turnSpent = p.MoveOrAttack(-1, 1, *b, o, *c)
	case StrMoveSouth:
		turnSpent = p.MoveOrAttack(0, 1, *b, o, *c)
	case StrMoveSoutheast:
		turnSpent = p.MoveOrAttack(1, 1, *b, o, *c)

	case StrFire:
		if p.ActiveWeapon != SlotWeaponMelee {
			if p.Equipment[p.ActiveWeapon].AmmoCurrent <= 0 {
				AddMessage("You need to reload [color=" + p.Equipment[p.ActiveWeapon].Color + "]" + p.Equipment[p.ActiveWeapon].Name + "[/color].")
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
					AddMessage("You cocked [color=" + p.Equipment[p.ActiveWeapon].Color + "]" + p.Equipment[p.ActiveWeapon].Name + "[/color].")
				}
			}
		} else {
			AddMessage("You are using melee weapon.")
		}

	case StrReload:
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
						AddMessage("You uncocked [color=" + p.Equipment[p.ActiveWeapon].Color + "]" + p.Equipment[p.ActiveWeapon].Name + "[/color].")
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

	case StrInspect:
		p.Look(*b, *o, *c) // Looking is free action.

	case StrPickup:
		turnSpent = p.PickUp(o)

	case StrPull:
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

	case StrPrimary:
		if p.ActiveWeapon != SlotWeaponPrimary {
			if p.Equipment[p.ActiveWeapon].Cock == true {
				p.Equipment[p.ActiveWeapon].Cocked = false
			}
			p.ActiveWeapon = SlotWeaponPrimary
			turnSpent = true
		}

	case StrSecondary:
		if p.ActiveWeapon != SlotWeaponSecondary {
			if p.Equipment[p.ActiveWeapon].Cock == true {
				p.Equipment[p.ActiveWeapon].Cocked = false
			}
			p.ActiveWeapon = SlotWeaponSecondary
			turnSpent = true
		}

	case StrMelee:
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

func Controls(k int, p *Creature, b *Board, c *Creatures, o *Objects) bool {
	turnSpent := false
	var command string
	if CustomControls == false {
		command = CommandKeys[k]
	} else {
		command = CustomCommandKeys[k]
	}
	turnSpent = Command(command, p, b, c, o)
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
