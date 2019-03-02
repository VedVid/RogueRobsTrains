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

func (c *Creature) AttackTarget(t *Creature, o *Objects) {
	/* Method Attack handles damage rolls for combat. Receiver "c" is attacker,
	   argument "t" is target.
	   "att" roll needs to be smaller or equal than weapon's effective range on given
	   distance (c.Equipment.Ranges). */
	att := RandInt(100)
	var dist int
	v, err := NewVector(c.X, c.Y, t.X, t.Y)
	if err != nil {
		fmt.Println(err)
	}
	dist = ComputeVector(v)
	var i int
	if dist <= RangeShort {
		i = 0
	} else if dist <= RangeMedium {
		i = 1
	} else {
		i = 2
	}
	weapon := c.Equipment[c.ActiveWeapon]
	weaponRange := weapon.Ranges[i]
	if att <= weaponRange {
		AddMessage(c.Name + " hits " + t.Name + ".")
		t.TakeDamage(1, o)
	} else {
		AddMessage(c.Name + " misses " + t.Name + ".")
	}
}

func (c *Creature) TakeDamage(dmg int, o *Objects) {
	/* Method TakeDamage has *Creature as receiver and takes damage integer
	   as argument. dmg value is deducted from Creature current HP.
	   If HPCurrent is below zero after taking damage, Creature dies. */
	c.HPCurrent -= dmg
	if c.HPCurrent <= 0 {
		c.Die(o)
	}
}
