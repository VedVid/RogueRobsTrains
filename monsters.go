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
	"errors"
	"fmt"
	"unicode/utf8"

	blt "bearlibterminal"
)

const (
	// Special characters.
	CorpseChar = "%"
)

type Creature struct {
	/* Creatures are living objects that
	   moves, attacks, dies, etc. */
	BasicProperties
	VisibilityProperties
	CollisionProperties
	FighterProperties
	EquipmentComponent
	ActiveWeapon int
}

// Creatures holds every creature on map.
type Creatures []*Creature

var FirstNames = []string{
	"Big", "Fat", "Skinny", "Small", "Short",
	"Sad", "Happy", "Angry", "Hateful", "Nervous", "Good", "Bad",
	"Kind", "Nice", "Mean", "Scurvy", "Surly",
	"Dirty", "Ugly",
	"Blind", "Deaf", "Limping",
	"Fast", "Strong", "Quick", "Heavy", "Slow", "Nimble",
	"Stupid", "Smart", "Bright",
}
var SecondNames = []string{
	"Ab", "Adam", "Andy", "Archie",
	"Ben", "Blake", "Brad", "Brock", "Buck", "Bud",
	"Carl", "Cary", "Chad", "Chas", "Chet", "Chris", "Clay", "Cliff", "Cole", "Colin", "Corey",
	"Dallas", "Dan", "Dave", "Dirk", "Doug", "Drew", "Dwight",
	"Earl", "Ed", "Eric",
	"Floyd", "Frank", "Fred",
	"Garth", "Gary", "Gavin", "George", "Glen", "Grant",
	"Hal", "Hank", "Harry", "Henry", "Hugh",
	"Ian",
	"Jack", "Jackson", "Jacob", "Jake", "James", "Jason", "Jay", "Jeb", "Jeff", "Jerry", "Jesse",
	"Jim", "Joel", "Joey", "John", "Jonathan", "Joe", "Justin",
	"Keith", "Kenny", "Kevin", "Kim", "Kyle",
	"Lanny", "Larry", "Laurence", "Lester", "Lewis", "Louis", "Lucas",
	"Mark", "Marshall", "Martin", "Matt", "Melvin", "Murray",
	"Nate", "Neil", "Noel", "Norman",
	"Oliver", "Oscar", "Oswald", "Otto", "Owen",
	"Patsy", "Paul", "Pete", "Phil",
	"Quentin",
	"Ralph", "Randall", "Ray", "Ricky", "Rob", "Roger", "Ron", "Rufus", "Russell",
	"Samuel", "Sam", "Seth", "Shane", "Shawn", "Simon", "Steve", "Swaine",
	"Terence", "Ted", "Tom", "Thomas", "Tim", "Tobias", "Toby", "Travis", "Trevor", "Troy",
	"Ultan",
	"Victor", "Vince",
	"Waldo", "Walter", "Wayne", "Will", "Winnie", "Winston",
	"Zack",
}

func NewCreature(x, y int, monsterFile string) (*Creature, error) {
	/* NewCreature is function that returns new Creature from
	   json file passed as argument. It replaced old code that
	   was encouraging hardcoding data in go files.
	   Errors returned by json package are not very helpful, and
	   hard to work with, so there is lazy panic for them. */
	var monster = &Creature{}
	err := CreatureFromJson(CreaturesPathJson+monsterFile, monster)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	monster.X, monster.Y = x, y
	var err2 error
	if monster.Layer < 0 {
		txt := LayerError(monster.Layer)
		err2 = errors.New("Creature layer is smaller than 0." + txt)
	}
	if monster.Layer != CreaturesLayer {
		txt := LayerWarning(monster.Layer, CreaturesLayer)
		err2 = errors.New("Creature layer is not equal to CreaturesLayer constant." + txt)
	}
	if monster.X < 0 || monster.X >= MapSizeX || monster.Y < 0 || monster.Y >= MapSizeY {
		txt := CoordsError(monster.X, monster.Y)
		err2 = errors.New("Creature coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(monster.Char) != 1 {
		txt := CharacterLengthError(monster.Char)
		err2 = errors.New("Creature character string length is not equal to 1." + txt)
	}
	if monster.HPMax < 0 {
		txt := InitialHPError(monster.HPMax)
		err2 = errors.New("Creature HPMax is smaller than 0." + txt)
	}
	if monster.Attack < 0 {
		txt := InitialAttackError(monster.Attack)
		err2 = errors.New("Creature attack value is smaller than 0." + txt)
	}
	if monster.Defense < 0 {
		txt := InitialDefenseError(monster.Defense)
		err2 = errors.New("Creature defense value is smaller than 0." + txt)
	}
	if monster.ActiveWeapon < 0 || monster.ActiveWeapon >= SlotMax {
		err2 = errors.New("ActiveWeapon of Creature is out of bounds.")
	}
	if monster.Equipment == nil {
		monster.Equipment = Objects{}
	}
	if monster.Inventory == nil {
		monster.Inventory = Objects{}
	}
	newName := SecondNames[RandInt(len(SecondNames)-1)]
	if RandInt(100) >= 60 {
		newName = "\"" + FirstNames[RandInt(len(FirstNames)-1)] + "\" " + newName
	}
	monster.Name = newName
	return monster, err2
}

func (c *Creature) MoveOrAttack(tx, ty int, b Board, o *Objects, all Creatures) bool {
	/* Method MoveOrAttack decides if Creature will move or attack other Creature;
	   It has *Creature receiver, and takes tx, ty (coords) integers as arguments,
	   and map of current level, and list of all Creatures.
	   Starts by target that is nil, then iterates through Creatures. If there is
	   Creature on targeted tile, that Creature becomes new target for attack.
	   Otherwise, Creature moves to specified Tile.
	   It's supposed to take player as receiver (attack / moving enemies is
	   handled differently - check ai.go and combat.go). */
	var target *Creature
	turnSpent := false
	for i, _ := range all {
		if all[i].X == c.X+tx && all[i].Y == c.Y+ty {
			if all[i].HPCurrent > 0 {
				target = all[i]
				break
			}
		}
	}
	if target != nil {
		if c.ActiveWeapon == SlotWeaponMelee {
			c.AttackTarget(target, o)
			turnSpent = true
		} else {
			AddMessage("You need melee weapon to do it.")
		}
	} else {
		turnSpent = c.Move(tx, ty, b, all)
	}
	return turnSpent
}

func (c *Creature) Move(tx, ty int, b Board, cs Creatures) bool {
	/* Move is method of Creature; it takes target x, y as arguments;
	   check if next move won't put Creature off the screen, then updates
	   Creature coords. */
	turnSpent := false
	newX, newY := c.X+tx, c.Y+ty
	if newX >= 0 &&
		newX <= MapSizeX-1 &&
		newY >= 0 &&
		newY <= MapSizeY-1 {
		if b[newX][newY].Blocked == false &&
			GetAliveCreatureFromTile(newX, newY, cs) == nil {
			c.X = newX
			c.Y = newY
			turnSpent = true
		} else if c == cs[0] && b[newX][newY].Name == "doors to next carriage" && G.Alive <= 0 {
			G.LevelInt++
		}
	}
	return turnSpent
}

func (c *Creature) PickUp(o *Objects) bool {
	/* PickUp is method that has *Creature as receiver
	   and slice of *Object as argument.
	   Creature tries to pick object up.
	   If creature stands on object that is possible to pick,
	   object is added to c's inventory, and removed
	   from "global" slice of objects.
	   Picking objects up takes turn only if it is
	   successful attempt. */
	turnSpent := false
	obj := *o
	var allObjects = Objects{}
	var is = []int{}
	for i := 0; i < len(obj); i++ {
		if obj[i].X == c.X && obj[i].Y == c.Y && obj[i].Pickable == true {
			allObjects = append(allObjects, obj[i])
			is = append(is, i)
		}
	}
	if c.AIType == PlayerAI {
		// Print menu
		var elements = []string{
			"top frame",
			"spacing",
			"ESCAPE",
			"spacing",
			"bottom frame",
		}
		for _, v := range allObjects {
			elements = append(elements, v.Name)
		}
		elements = append(elements, "just something else...")
		sizeY := len(elements)
		startY := (MapSizeY / 2) - (sizeY / 2)
		endY := (MapSizeY / 2) + (sizeY/2 + (sizeY % 2))
		for x := 5; x < MapSizeX-5; x++ {
			for y := startY; y < endY; y++ {
				blt.Layer(MenuLayer)
				blt.Print(x, y, "[color=black]▓[/color]")
				switch y {
				case startY:
					blt.Layer(MenuLayer + 1)
					if x == 5 {
						blt.Print(x, y, "[color=#a0785a]╔[/color]")
					} else if x == MapSizeX-5-1 {
						blt.Print(x, y, "[color=#a0785a]╗[/color]")
					} else {
						blt.Print(x, y, "[color=#a0785a]═[/color]")
					}
				case endY - 1:
					blt.Layer(MenuLayer + 1)
					if x == 5 {
						blt.Print(x, y, "[color=#a0785a]╚[/color]")
					} else if x == MapSizeX-5-1 {
						blt.Print(x, y, "[color=#a0785a]╝[/color]")
					} else {
						blt.Print(x, y, "[color=#a0785a]═[/color]")
					}
				default:
					switch x {
					case 5, MapSizeX - 5 - 1:
						blt.Layer(MenuLayer + 1)
						blt.Print(x, y, "[color=#a0785a]║[/color]")
					}
				}
			}
		} //printing finished
		//print objects
		maxI := 0
		for i, v := range allObjects {
			blt.Layer(MenuLayer + 1)
			weaponStr := ""
			weaponStr = weaponStr + "[color=" + v.Color + "]"
			weaponStr = weaponStr + v.Char + " " + v.Name
			if v.Ranges[0] != 0 && (v.Ranges[1] != 0 || v.Ranges[2] != 0) {
				rangesStr := "([/color]"
				for i, _ := range v.Ranges {
					val := v.Ranges[i]
					if val < 25 {
						rangesStr = rangesStr + "[color=darker red]▁[/color]"
					} else if val < 50 {
						rangesStr = rangesStr + "[color=darker flame]▃[/color]"
					} else if val < 75 {
						rangesStr = rangesStr + "[color=darker yellow]▅[/color]"
					} else {
						rangesStr = rangesStr + "[color=darker green]▇[/color]"
					}
				}
				if v.Cock == true {
					rangesStr = rangesStr + "[color=dark red]" + CockedIcon + "[/color]"
				}
				rangesStr = rangesStr + "[color=" + v.Color + "])[/color]"
				weaponStr = weaponStr + rangesStr
			}
			blt.Print(5+2, startY+2+i, OrderToCharacter(i)+") "+weaponStr)
			maxI++
		} //printing finished
		blt.Print(5+2, startY+2+maxI+1, "Press [[ESCAPE]] to cancel.")
		blt.Refresh()
		var key int
		var ord int
		for {
			key = ReadInput()
			if key == blt.TK_ESCAPE {
				return turnSpent
			}
			ord = KeyToOrder(key)
			if ord < len(allObjects) && ord >= 0 {
				break
			} else {
				continue
			}
		}
		weapon := allObjects[ord]
		wName := "[color=" + weapon.Color + "]" + weapon.Name + "[/color]"
		AddMessage("You found " + wName + ".")
		if weapon.Slot != c.ActiveWeapon {
			if c.Equipment[c.ActiveWeapon].Cock == true {
				c.Equipment[c.ActiveWeapon].Cocked = false
			}
			c.ActiveWeapon = weapon.Slot
		}
		c.DropFromEquipment(&obj, weapon.Slot)
		c.EquipItem(weapon, weapon.Slot)
		copy(obj[is[ord]:], obj[is[ord]+1:])
		obj[len(obj)-1] = nil
		*o = obj[:len(obj)-1]
		turnSpent = true
	} else {
		for i := 0; i < len(obj); i++ {
			if obj[i].X == c.X && obj[i].Y == c.Y && obj[i].Pickable == true {
				if obj[i].Slot != c.ActiveWeapon {
					if c.Equipment[c.ActiveWeapon].Cock == true {
						c.Equipment[c.ActiveWeapon].Cocked = false
					}
					c.ActiveWeapon = obj[i].Slot
				}
				c.DropFromEquipment(&obj, obj[i].Slot)
				c.EquipItem(obj[i], obj[i].Slot)
				copy(obj[i:], obj[i+1:])
				obj[len(obj)-1] = nil
				*o = obj[:len(obj)-1]
				turnSpent = true
				break
			}
		}
	}
	return turnSpent
}

func (c *Creature) DropFromInventory(objects *Objects, index int) bool {
	/* Drop is method that has Creature as receiver and takes
	   "global" list of objects as main argument, and additional
	   integer that is index of item to be dropped from c's Inventory.
	   At first, turnSpent is set to false, to make it true
	   at the end of function. It may be considered as obsolete WET,
	   because 'return true' would be sufficient, but it is
	   a bit more readable now.
	   Objs is dereferenced objects and it is absolutely necessary
	   to do any actions on these objects.
	   Drop do two things:
	   at first, it adds specific item to the game map,
	   then it removes this item from its owner Inventory. */
	turnSpent := false
	objs := *objects
	if c.AIType == PlayerAI {
		oName := "[color=" + c.Inventory[index].Color + "]" + c.Inventory[index].Name + "[/color]"
		AddMessage("You dropped " + oName + ".")
	}
	// Add item to the map.
	object := c.Inventory[index]
	object.X, object.Y = c.X, c.Y
	objs = append(objs, object)
	*objects = objs
	// Then remove item from inventory.
	copy(c.Inventory[index:], c.Inventory[index+1:])
	c.Inventory[len(c.Inventory)-1] = nil
	c.Inventory = c.Inventory[:len(c.Inventory)-1]
	turnSpent = true
	return turnSpent
}

func (c *Creature) DropFromEquipment(objects *Objects, slot int) bool {
	/* DropFromEquipment is method of *Creature that takes "global" objects,
	   and int (as index) as arguments, and returns bool (result depends if
	   action was successful, therefore if took a turn).
	   This function is very similar to DropFromInventory, but is kept
	   due to explicitness.
	   The difference is that Equipment checks Equipment index, not
	   specific object, so additionally checks for nils, and instead of
	   removing item from slice, makes it nil.
	   This behavior is important, because while Inventory is "dynamic"
	   slice, Equipment is supposed to be "fixed size" - slots are present
	   all the time, but the can be empty (ie nil) or occupied (ie object). */
	turnSpent := false
	objs := *objects
	object := c.Equipment[slot]
	if object == nil {
		return turnSpent // turn is not spent because there is no object to drop
	}
	// else {
	if c.AIType == PlayerAI {
		oName := "[color=" + object.Color + "]" + object.Name + "[/color]"
		AddMessage("You removed and dropped " + oName + ".")
	}
	// add item to map
	object.X, object.Y = c.X, c.Y
	objs = append(objs, object)
	*objects = objs
	// then remove from slot
	c.Equipment[slot] = nil
	turnSpent = true
	return turnSpent
}

func (c *Creature) EquipItem(o *Object, slot int) (bool, error) {
	/* EquipItem is method of *Creature that takes *Object and int (that is
	   indicator to index of Equipment slot) as arguments; it returns
	   bool and error.
	   At first, EquipItem checks for errors:
	    - if object to equip exists
	    - if this equipment slot is not occupied
	   then equips item and removes it from inventory. */
	var err error
	if o == nil {
		txt := EquipNilError(c)
		err = errors.New("Creature tried to equip *Object that was nil." + txt)
	}
	if c.Equipment[slot] != nil {
		txt := EquipSlotNotNilError(c, slot)
		err = errors.New("Creature tried to equip item into already occupied slot." + txt)
	}
	if o.Slot != slot {
		txt := EquipWrongSlotError(o.Slot, slot)
		err = errors.New("Creature tried to equip item into wrong slot." + txt)
	}
	turnSpent := false
	// Equip item...
	c.Equipment[slot] = o
	if c.AIType == PlayerAI {
		AddMessage("You equipped [color=" + o.Color + "]" + o.Name + "[/color].")
	}
	turnSpent = true
	return turnSpent, err
}

func (c *Creature) DequipItem(slot int) (bool, error) {
	/* DequipItem is method of Creature. It is called when receiver is about
	   to dequip weapon from "ready" equipment slot.
	   At first, weapon is added to Inventory, then Equipment slot is set to nil. */
	var err error
	if c.Equipment[slot] == nil {
		txt := DequipNilError(c, slot)
		err = errors.New("Creature tried to DequipItem that was nil." + txt)
	}
	if c.AIType == PlayerAI {
		oName := "[color=" + c.Equipment[slot].Color + "]" + c.Equipment[slot].Name + "[/color]"
		AddMessage("You removed " + oName + ".")
	}
	turnSpent := false
	c.Inventory = append(c.Inventory, c.Equipment[slot]) //adding items to inventory should have own function, that will check "bounds" of inventory
	c.Equipment[slot] = nil
	turnSpent = true
	return turnSpent, err
}

func (c *Creature) Die(o *Objects) {
	/* Method Die is called when Creature's HP drops below zero.
	   Die() has *Creature as receiver.
	   Receiver properties changes to fit better to corpse. */
	c.Layer = DeadLayer
	c.Name = "corpse"
	c.Color = "dark red"
	c.ColorDark = "dark red"
	c.Char = CorpseChar
	c.Blocked = false
	c.BlocksSight = false
	c.AIType = NoAI
	c.DropFromEquipment(o, c.ActiveWeapon)
	ZeroLastTarget(c)
	G.Alive--
}

func FindMonsterByXY(x, y int, c Creatures) *Creature {
	/* Function FindMonsterByXY takes desired coords and list
	   of all available creatures. It iterates through this list,
	   and returns nil or creature that occupies specified coords. */
	var monster *Creature
	for i := 0; i < len(c); i++ {
		if x == c[i].X && y == c[i].Y {
			monster = c[i]
			break
		}
	}
	return monster
}
