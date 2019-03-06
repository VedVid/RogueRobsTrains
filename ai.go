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

import "fmt"

const (
	// Types of AI.
	NoAI = iota
	PlayerAI
	MeleeDumbAI
	MeleePatherAI
	RangedDumbAI
	RangedPatherAI
)

const (
	// Probability of triggering AI
	AITrigger = 92
)

func CreaturesTakeTurn(b Board, c Creatures, o *Objects) {
	/* Function CreaturesTakeTurn is supposed to handle all enemy creatures
	   actions: movement, attacking, etc.
	   It takes Board and Creatures as arguments.
	   Iterates through all Creatures slice, and calls HandleAI function with
	   specific parameters.
	   It skips NoAI and PlayerAI. */
	var ai int
	for _, v := range c {
		ai = v.AIType
		if ai == NoAI || ai == PlayerAI {
			continue
		}
		HandleAI(b, c, o, v)
		TriggerAI(b, c[0], v)
	}
}

func TriggerAI(b Board, p, c *Creature) {
	/* TriggerAI is function that takes Board and two Creatures as arguments.
	   First Creature is supposed to be player, second one - enemy.
	   Enemy with AITriggered set to false will ignore the player existence.
	   AITrigger is probability to notice (and, therefore, switch AITriggered)
	   player if is in monster's FOV. */
	if c.AITriggered == false {
		if IsInFOV(b, p.X, p.Y, c.X, c.Y) == true && RandInt(100) <= AITrigger {
			cName := "[color=" + c.Color + "]" + c.Name + "[/color]"
			AddMessage(cName + " spotted you!")
			c.AITriggered = true
		}
	}
}

func HandleAI(b Board, cs Creatures, o *Objects, c *Creature) {
	/* TODO: This docstring needs update!
	   HandleAI is robust function that takes Board, Creatures, Objects,
	   and specific Creature as arguments. The most notable argument is
	   the last one - behavior of this entity will be decided in function body.
	   Its behavior will be decided regarding to AIType.
	   This function is very big and *wet*, but it is here to stay, for a while,
	   at least. I thought about code duplication removal by introducing one
	   generic function that would take Creature as argument, and - after
	   AIType check - would use proper HandleMeleeDumbAI (etc.) functions; or
	   would start with available weapons check. (One may want to peek at
	   issue #98 in repo - https://github.com/VedVid/RAWIG/issues/98 ).
	   But, on the other hand, ai has so many variations and edge cases that
	   unifying monster's behavior would result in smaller flexibility. */
	ai := c.AIType
	cName := "[color=" + c.Color + "]" + c.Name + "[/color]"
	switch ai {
	case MeleeDumbAI:
		if c.AITriggered == true {
			if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
				c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
			} else {
				c.AttackTarget(cs[0], o)
			}
		} else {
			dx := RandRange(-1, 1)
			dy := RandRange(-1, 1)
			nx := c.X+dx
			ny := c.Y+dy
			if nx < 0 || nx >= MapSizeX {
				nx = c.X
			}
			if ny < 0 || ny >= MapSizeY {
				ny = c.Y
			}
			if (b[nx][ny].BlocksSight == false) ||
				(b[nx][ny].BlocksSight == true && RandInt(100) > 80) {
				c.Move(dx, dy, b, cs)
			}
		}
	case MeleePatherAI:
		// The same set of functions as for DumbAI.
		// Just for clarity.
		if c.AITriggered == true {
			if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
				c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
			} else {
				c.AttackTarget(cs[0], o)
			}
		} else {
			dx := RandRange(-1, 1)
			dy := RandRange(-1, 1)
			nx := c.X+dx
			ny := c.Y+dy
			if nx < 0 || nx >= MapSizeX {
				nx = c.X
			}
			if ny < 0 || ny >= MapSizeY {
				ny = c.Y
			}
			if (b[nx][ny].BlocksSight == false) ||
				(b[nx][ny].BlocksSight == true && RandInt(100) > 80) {
				c.Move(dx, dy, b, cs)
			}
		}
	case RangedDumbAI:
		if c.AITriggered == true {
			if c.Equipment[c.ActiveWeapon] != nil {
				if c.ActiveWeapon == SlotWeaponMelee {
					if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
						// TODO:
						// For now, every ranged skill has range equal to FOVLength-1
						// but it should change in future.
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						c.AttackTarget(cs[0], o)
					}
				} else {
					if c.Equipment[c.ActiveWeapon].AmmoCurrent <= 0 {
						if c.Equipment[c.ActiveWeapon].Cock == false {
							c.Equipment[c.ActiveWeapon].AmmoCurrent = c.Equipment[c.ActiveWeapon].AmmoMax
							if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
								AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
							}
							break
						} else if c.Equipment[c.ActiveWeapon].Cock == true {
							if c.Equipment[c.ActiveWeapon].Cocked == true {
								c.Equipment[c.ActiveWeapon].Cocked = false
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " uncocks his " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							} else {
								c.Equipment[c.ActiveWeapon].AmmoCurrent++
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							}
						}
					} else if c.Equipment[c.ActiveWeapon].AmmoCurrent < c.Equipment[c.ActiveWeapon].AmmoMax &&
						IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == false {
						if c.Equipment[c.ActiveWeapon].Cock == false {
							c.Equipment[c.ActiveWeapon].AmmoCurrent = c.Equipment[c.ActiveWeapon].AmmoMax
							break
						} else {
							if c.Equipment[c.ActiveWeapon].Cocked == true {
								c.Equipment[c.ActiveWeapon].Cocked = false
								break
							} else {
								c.Equipment[c.ActiveWeapon].AmmoCurrent++
								break
							}
						}
					}
					if c.DistanceTo(cs[0].X, cs[0].Y) >= FOVLength-1 ||
						IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == false { // should it use DistanceTo, instead of ComputeVector?
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						// DumbAI will not check if target is valid
						vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
						if err != nil {
							fmt.Println(err)
						}
						_ = ComputeVector(vec)
						_, _, target, _ := ValidateVector(vec, b, cs, *o)
						if target != nil {
							if c.Equipment[c.ActiveWeapon].Cock == false {
								c.AttackTarget(target, o)
								c.Equipment[c.ActiveWeapon].AmmoCurrent--
							} else {
								if c.Equipment[c.ActiveWeapon].Cocked == true {
									c.AttackTarget(target, o)
									c.Equipment[c.ActiveWeapon].AmmoCurrent--
									c.Equipment[c.ActiveWeapon].Cocked = false
								} else {
									c.Equipment[c.ActiveWeapon].Cocked = true
									AddMessage(cName + " cocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
							}
						}
					}
				}
			} else {
				if c.Equipment[SlotWeaponPrimary] != nil {
					// Use primary ranged weapon.
					if c.Equipment[SlotWeaponPrimary].AmmoCurrent <= 0 {
						if c.Equipment[SlotWeaponPrimary].Cock == false {
							c.Equipment[SlotWeaponPrimary].AmmoCurrent = c.Equipment[SlotWeaponPrimary].AmmoMax
							if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
								AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
							}
							break
						} else if c.Equipment[SlotWeaponPrimary].Cock == true {
							if c.Equipment[SlotWeaponPrimary].Cocked == true {
								c.Equipment[SlotWeaponPrimary].Cocked = false
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " uncocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							} else {
								c.Equipment[SlotWeaponPrimary].AmmoCurrent++
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							}
						}
					} else if c.Equipment[SlotWeaponPrimary].AmmoCurrent < c.Equipment[SlotWeaponPrimary].AmmoMax &&
						IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == false {
						if c.Equipment[SlotWeaponPrimary].Cock == false {
							c.Equipment[SlotWeaponPrimary].AmmoCurrent = c.Equipment[SlotWeaponPrimary].AmmoMax
							break
						} else {
							if c.Equipment[SlotWeaponPrimary].Cocked == true {
								c.Equipment[SlotWeaponPrimary].Cocked = false
								break
							} else {
								c.Equipment[SlotWeaponPrimary].AmmoCurrent++
								break
							}
						}
					}
					if c.DistanceTo(cs[0].X, cs[0].Y) >= FOVLength-1 ||
						IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == false {
						// TODO:
						// For now, every ranged skill has range equal to FOVLength-1
						// but it should change in future.
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						// DumbAI will not check if target is valid
						vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
						if err != nil {
							fmt.Println(err)
						}
						_ = ComputeVector(vec)
						_, _, target, _ := ValidateVector(vec, b, cs, *o)
						if target != nil {
							if c.Equipment[SlotWeaponPrimary].Cock == false {
								c.AttackTarget(target, o)
								c.Equipment[SlotWeaponPrimary].AmmoCurrent--
							} else {
								if c.Equipment[SlotWeaponPrimary].Cocked == true {
									c.AttackTarget(target, o)
									c.Equipment[SlotWeaponPrimary].AmmoCurrent--
									c.Equipment[SlotWeaponPrimary].Cocked = false
								} else {
									c.Equipment[SlotWeaponPrimary].Cocked = true
									AddMessage(cName + " cocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
							}
						}
					}
				} else if c.Equipment[SlotWeaponSecondary] != nil {
					// Use secondary ranged weapon.
					if c.Equipment[SlotWeaponSecondary].AmmoCurrent <= 0 {
						if c.Equipment[SlotWeaponSecondary].Cock == false {
							c.Equipment[SlotWeaponSecondary].AmmoCurrent = c.Equipment[SlotWeaponSecondary].AmmoMax
							if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
								AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
							}
							break
						} else if c.Equipment[SlotWeaponSecondary].Cock == true {
							if c.Equipment[SlotWeaponSecondary].Cocked == true {
								c.Equipment[SlotWeaponSecondary].Cocked = false
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " uncocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							} else {
								c.Equipment[SlotWeaponSecondary].AmmoCurrent++
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							}
						}
					} else if c.Equipment[SlotWeaponSecondary].AmmoCurrent < c.Equipment[SlotWeaponSecondary].AmmoMax &&
						IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == false {
						if c.Equipment[SlotWeaponSecondary].Cock == false {
							c.Equipment[SlotWeaponSecondary].AmmoCurrent = c.Equipment[SlotWeaponSecondary].AmmoMax
							break
						} else {
							if c.Equipment[SlotWeaponSecondary].Cocked == true {
								c.Equipment[SlotWeaponSecondary].Cocked = false
								break
							} else {
								c.Equipment[SlotWeaponSecondary].AmmoCurrent++
								break
							}
						}
					}
					if c.DistanceTo(cs[0].X, cs[0].Y) >= FOVLength-1 ||
						IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == false {
						// TODO:
						// For now, every ranged skill has range equal to FOVLength-1
						// but it should change in future.
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						// DumbAI will not check if target is valid
						vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
						if err != nil {
							fmt.Println(err)
						}
						_ = ComputeVector(vec)
						_, _, target, _ := ValidateVector(vec, b, cs, *o)
						if target != nil {
							if c.Equipment[SlotWeaponSecondary].Cock == false {
								c.AttackTarget(target, o)
								c.Equipment[SlotWeaponSecondary].AmmoCurrent--
							} else {
								if c.Equipment[SlotWeaponSecondary].Cocked == true {
									c.AttackTarget(target, o)
									c.Equipment[SlotWeaponSecondary].AmmoCurrent--
									c.Equipment[SlotWeaponSecondary].Cocked = false
								} else {
									c.Equipment[SlotWeaponSecondary].Cocked = true
									AddMessage(cName + " cocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
							}
						}
					}
				} else {
					if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						c.AttackTarget(cs[0], o)
					}
				}
			}
		} else {
			if c.Equipment[c.ActiveWeapon] != nil &&
				c.Equipment[c.ActiveWeapon].AmmoCurrent < c.Equipment[c.ActiveWeapon].AmmoMax {
				if c.Equipment[c.ActiveWeapon].Cock == false {
					c.Equipment[c.ActiveWeapon].AmmoCurrent = c.Equipment[c.ActiveWeapon].AmmoMax
					AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
					break
				} else if c.Equipment[c.ActiveWeapon].Cock == true {
					if c.Equipment[c.ActiveWeapon].Cocked == true {
						c.Equipment[c.ActiveWeapon].Cocked = false
						AddMessage(cName + " uncocks " + c.Equipment[c.ActiveWeapon].Name + ".")
						break
					} else {
						c.Equipment[c.ActiveWeapon].AmmoCurrent++
						AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
						break
					}
				}
			} else {
				dx := RandRange(-1, 1)
				dy := RandRange(-1, 1)
				nx := c.X+dx
				ny := c.Y+dy
				if nx < 0 || nx >= MapSizeX {
					nx = c.X
				}
				if ny < 0 || ny >= MapSizeY {
					ny = c.Y
				}
				if (b[nx][ny].BlocksSight == false) ||
					(b[nx][ny].BlocksSight == true && RandInt(100) > 80) {
					c.Move(dx, dy, b, cs)
				}
			}
		}
	case RangedPatherAI: // It will depend on ranged weapons and equipment implementation
		if c.AITriggered == true {
			if c.Equipment[c.ActiveWeapon] != nil {
				if c.ActiveWeapon == SlotWeaponMelee {
					if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						c.AttackTarget(cs[0], o)
					}
				} else {
					if c.Equipment[c.ActiveWeapon].AmmoCurrent <= 0 {
						if c.Equipment[c.ActiveWeapon].Cock == false {
							c.Equipment[c.ActiveWeapon].AmmoCurrent = c.Equipment[c.ActiveWeapon].AmmoMax
							if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
								AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
							}
							break
						} else if c.Equipment[c.ActiveWeapon].Cock == true {
							if c.Equipment[c.ActiveWeapon].Cocked == true {
								c.Equipment[c.ActiveWeapon].Cocked = false
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " uncocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							} else {
								c.Equipment[c.ActiveWeapon].AmmoCurrent++
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							}
						}
					} else if c.Equipment[c.ActiveWeapon].AmmoCurrent < c.Equipment[c.ActiveWeapon].AmmoMax &&
						IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == false {
						if c.Equipment[c.ActiveWeapon].Cock == false {
							c.Equipment[c.ActiveWeapon].AmmoCurrent = c.Equipment[c.ActiveWeapon].AmmoMax
							break
						} else {
							if c.Equipment[c.ActiveWeapon].Cocked == true {
								c.Equipment[c.ActiveWeapon].Cocked = false
								break
							} else {
								c.Equipment[c.ActiveWeapon].AmmoCurrent++
								break
							}
						}
					}
					bestDistance := FindMaxInSlice(c.Equipment[c.ActiveWeapon].Ranges)
					if c.DistanceTo(cs[0].X, cs[0].Y) >= FOVLength-1 {
						// TODO:
						// For now, every ranged skill has range equal to FOVLength-1
						// but it should change in future.
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else if c.DistanceTo(cs[0].X, cs[0].Y) > bestDistance {
						// If distance between creature and target is bigger than
						// optimal effective range of currently wielded weapon,
						// move towards target.
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
						if err != nil {
							fmt.Println(err)
						}
						_ = ComputeVector(vec)
						_, _, target, _ := ValidateVector(vec, b, cs, *o)
						if target != cs[0] {
							c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
						} else {
							if c.Equipment[c.ActiveWeapon].Cock == false {
								c.AttackTarget(target, o)
								c.Equipment[c.ActiveWeapon].AmmoCurrent--
							} else {
								if c.Equipment[c.ActiveWeapon].Cocked == true {
									c.AttackTarget(target, o)
									c.Equipment[c.ActiveWeapon].Cocked = false
								} else {
									c.Equipment[c.ActiveWeapon].Cocked = true
									AddMessage(cName + " cocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
							}
						}
					}
				}
			} else {
				if c.Equipment[SlotWeaponPrimary] != nil {
					if c.Equipment[SlotWeaponPrimary].AmmoCurrent <= 0 {
						if c.Equipment[SlotWeaponPrimary].Cock == false {
							c.Equipment[SlotWeaponPrimary].AmmoCurrent = c.Equipment[SlotWeaponPrimary].AmmoMax
							if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
								AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
							}
							break
						} else if c.Equipment[SlotWeaponPrimary].Cock == true {
							if c.Equipment[SlotWeaponPrimary].Cocked == true {
								c.Equipment[SlotWeaponPrimary].Cocked = false
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " uncocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							} else {
								c.Equipment[SlotWeaponPrimary].AmmoCurrent++
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							}
						}
					} else if c.Equipment[SlotWeaponPrimary].AmmoCurrent < c.Equipment[SlotWeaponPrimary].AmmoMax &&
						IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == false {
						if c.Equipment[SlotWeaponPrimary].Cock == false {
							c.Equipment[SlotWeaponPrimary].AmmoCurrent = c.Equipment[SlotWeaponPrimary].AmmoMax
							break
						} else {
							if c.Equipment[SlotWeaponPrimary].Cocked == true {
								c.Equipment[SlotWeaponPrimary].Cocked = false
								break
							} else {
								c.Equipment[SlotWeaponPrimary].AmmoCurrent++
								break
							}
						}
					}
					if c.DistanceTo(cs[0].X, cs[0].Y) >= FOVLength-1 {
						// TODO:
						// For now, every ranged skill has range equal to FOVLength-1
						// but it should change in future.
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
						if err != nil {
							fmt.Println(err)
						}
						_ = ComputeVector(vec)
						_, _, target, _ := ValidateVector(vec, b, cs, *o)
						if target != cs[0] {
							c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
						} else {
							if c.Equipment[SlotWeaponPrimary].Cock == false {
								c.AttackTarget(target, o)
								c.Equipment[SlotWeaponPrimary].AmmoCurrent--
							} else {
								if c.Equipment[SlotWeaponPrimary].Cocked == true {
									c.AttackTarget(target, o)
									c.Equipment[SlotWeaponPrimary].Cocked = false
								} else {
									c.Equipment[SlotWeaponPrimary].Cocked = true
									AddMessage(cName + " cocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
							}
						}
					}
				} else if c.Equipment[SlotWeaponSecondary] != nil {
					if c.Equipment[SlotWeaponSecondary].AmmoCurrent <= 0 {
						if c.Equipment[SlotWeaponSecondary].Cock == false {
							c.Equipment[SlotWeaponSecondary].AmmoCurrent = c.Equipment[SlotWeaponSecondary].AmmoMax
							if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
								AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
							}
							break
						} else if c.Equipment[SlotWeaponSecondary].Cock == true {
							if c.Equipment[SlotWeaponSecondary].Cocked == true {
								c.Equipment[SlotWeaponSecondary].Cocked = false
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " uncocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							} else {
								c.Equipment[SlotWeaponSecondary].AmmoCurrent++
								if IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == true {
									AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
								break
							}
						}
					} else if c.Equipment[SlotWeaponSecondary].AmmoCurrent < c.Equipment[SlotWeaponSecondary].AmmoMax &&
						IsInFOV(b, c.X, c.Y, cs[0].X, cs[0].Y) == false {
						if c.Equipment[SlotWeaponSecondary].Cock == false {
							c.Equipment[SlotWeaponSecondary].AmmoCurrent = c.Equipment[SlotWeaponSecondary].AmmoMax
							break
						} else {
							if c.Equipment[SlotWeaponSecondary].Cocked == true {
								c.Equipment[SlotWeaponSecondary].Cocked = false
								break
							} else {
								c.Equipment[SlotWeaponSecondary].AmmoCurrent++
								break
							}
						}
					}
					if c.DistanceTo(cs[0].X, cs[0].Y) >= FOVLength-1 {
						// TODO:
						// For now, every ranged skill has range equal to FOVLength-1
						// but it should change in future.
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						vec, err := NewVector(c.X, c.Y, cs[0].X, cs[0].Y)
						if err != nil {
							fmt.Println(err)
						}
						_ = ComputeVector(vec)
						_, _, target, _ := ValidateVector(vec, b, cs, *o)
						if target != cs[0] {
							c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
						} else {
							if c.Equipment[SlotWeaponSecondary].Cock == false {
								c.AttackTarget(target, o)
								c.Equipment[SlotWeaponSecondary].AmmoCurrent--
							} else {
								if c.Equipment[SlotWeaponSecondary].Cocked == true {
									c.AttackTarget(target, o)
									c.Equipment[SlotWeaponSecondary].Cocked = false
								} else {
									c.Equipment[SlotWeaponSecondary].Cocked = true
									AddMessage(cName + " cocks " + c.Equipment[c.ActiveWeapon].Name + ".")
								}
							}
						}
					}
				} else {
					if c.DistanceTo(cs[0].X, cs[0].Y) > 1 {
						c.MoveTowards(b, cs, cs[0].X, cs[0].Y, ai)
					} else {
						c.AttackTarget(cs[0], o)
					}
				}
			}
		} else {
			if c.Equipment[c.ActiveWeapon] != nil &&
				c.Equipment[c.ActiveWeapon].AmmoCurrent < c.Equipment[c.ActiveWeapon].AmmoMax {
				if c.Equipment[c.ActiveWeapon].Cock == false {
					c.Equipment[c.ActiveWeapon].AmmoCurrent = c.Equipment[c.ActiveWeapon].AmmoMax
					AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
					break
				} else if c.Equipment[c.ActiveWeapon].Cock == true {
					if c.Equipment[c.ActiveWeapon].Cocked == true {
						c.Equipment[c.ActiveWeapon].Cocked = false
						AddMessage(cName + " uncocks " + c.Equipment[c.ActiveWeapon].Name + ".")
						break
					} else {
						c.Equipment[c.ActiveWeapon].AmmoCurrent++
						AddMessage(cName + " reloads " + c.Equipment[c.ActiveWeapon].Name + ".")
						break
					}
				}
			} else {
				dx := RandRange(-1, 1)
				dy := RandRange(-1, 1)
				nx := c.X+dx
				ny := c.Y+dy
				if nx < 0 || nx >= MapSizeX {
					nx = c.X
				}
				if ny < 0 || ny >= MapSizeY {
					ny = c.Y
				}
				if (b[nx][ny].BlocksSight == false) ||
					(b[nx][ny].BlocksSight == true && RandInt(100) > 80) {
					c.Move(dx, dy, b, cs)
				}
			}
		}
	}
}
