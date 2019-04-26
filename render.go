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
	"strconv"
)

const (
	/* Constant values for layers. Their usage is optional,
	   but (for now, at leas) recommended, because default
	   rendering functions depends on these values.
	   They are important for proper clearing characters
	   that should not be displayed, as, for example,
	   bracelet under the monster. */
	_ = iota
	UILayer
	BoardLayer
	DeadLayer
	ObjectsLayer
	CreaturesLayer
	PlayerLayer
	LookLayer
	MenuLayer
)

const CockedIcon = "Φ"

func PrintBoard(b Board, c Creatures) {
	/* Function PrintBoard is used in RenderAll function.
	   Takes level map and list of monsters as arguments
	   and iterates through Board.
	   It has to check for "]" and "[" characters, because
	   BearLibTerminal uses these symbols for config.
	   Instead of checking it here, one could just remember to
	   always pass "]]" instead of "]".
	   Prints every tile on its coords if certain conditions are met:
	   is Explored already, and:
	   - is in player's field of view (prints "normal" color) or
	   - is AlwaysVisible (prints dark color). */
	for x := 0; x < MapSizeX; x++ {
		for y := 0; y < MapSizeY; y++ {
			// Technically, "t" is new variable with own memory address...
			t := b[x][y] // Should it be *b[x][y]?
			blt.Layer(t.Layer)
			if t.Explored == true {
				ch := t.Char
				if t.Char == "[" || t.Char == "]" {
					ch = t.Char + t.Char
				}
				if IsInFOV(b, c[0].X, c[0].Y, t.X, t.Y) == true {
					glyph := "[color=" + t.Color + "]" + ch
					blt.Print(t.X, t.Y, glyph)
				} else {
					if t.AlwaysVisible == true {
						glyph := "[color=" + t.ColorDark + "]" + ch
						blt.Print(t.X, t.Y, glyph)
					}
				}
			}
		}
	}
}

func PrintObjects(b Board, o Objects, c Creatures) {
	/* Function PrintObjects is used in RenderAll function.
	   Takes map of level, slice of objects, and all monsters
	   as arguments.
	   Iterates through Objects.
	   It has to check for "]" and "[" characters, because
	   BearLibTerminal uses these symbols for config.
	   Instead of checking it here, one could just remember to
	   always pass "]]" instead of "]".
	   Prints every object on its coords if certain conditions are met:
	   AlwaysVisible bool is set to true, or is in player fov. */
	for _, v := range o {
		if (IsInFOV(b, c[0].X, c[0].Y, v.X, v.Y) == true) ||
			((v.AlwaysVisible == true) && (b[v.X][v.Y].Explored == true)) {
			blt.Layer(v.Layer)
			ch := v.Char
			if v.Char == "]" || v.Char == "[" {
				ch = v.Char + v.Char
			}
			for l := 0; l < v.Layer; l++ {
				blt.Layer(l)
				blt.ClearArea(v.X, v.Y, 1, 1)
			}
			glyph := "[color=" + v.Color + "]" + ch
			blt.Print(v.X, v.Y, glyph)
		}
	}
}

func PrintAliveCreatures(b Board, c Creatures) {
	/* Function PrintCreatures is used in RenderAll function.
	   Takes map of level and slice of Creatures as arguments.
	   Iterates through Creatures.
	   It has to check for "]" and "[" characters, because
	   BearLibTerminal uses these symbols for config.
	   Instead of checking it here, one could just remember to
	   always pass "]]" instead of "]".
	   Checks for every creature on its coords if certain conditions are met:
	   AlwaysVisible bool is set to true, or is in player fov. */
	for _, v := range c {
		if ((IsInFOV(b, c[0].X, c[0].Y, v.X, v.Y) == true) ||
			(v.AlwaysVisible == true)) && v.HPCurrent > 0 {
			blt.Layer(v.Layer)
			ch := v.Char
			if v.Char == "]" || v.Char == "[" {
				ch = v.Char + v.Char
			}
			for l := 0; l < v.Layer; l++ {
				blt.Layer(l)
				blt.ClearArea(v.X, v.Y, 1, 1)
			}
			glyph := "[color=" + v.Color + "]" + ch
			blt.Print(v.X, v.Y, glyph)
		}
	}
}

func PrintDeadCreatures(b Board, c Creatures) {
	/* Function PrintCreatures is used in RenderAll function.
	   Takes map of level and slice of Creatures as arguments.
	   Iterates through Creatures.
	   It has to check for "]" and "[" characters, because
	   BearLibTerminal uses these symbols for config.
	   Instead of checking it here, one could just remember to
	   always pass "]]" instead of "]".
	   Checks for every creature on its coords if certain conditions are met:
	   AlwaysVisible bool is set to true, or is in player fov. */
	for _, v := range c {
		if ((IsInFOV(b, c[0].X, c[0].Y, v.X, v.Y) == true) ||
			(v.AlwaysVisible == true)) && v.HPCurrent <= 0 {
			blt.Layer(v.Layer)
			ch := v.Char
			if v.Char == "]" || v.Char == "[" {
				ch = v.Char + v.Char
			}
			for l := 0; l < v.Layer; l++ {
				blt.Layer(l)
				blt.ClearArea(v.X, v.Y, 1, 1)
			}
			glyph := "[color=" + v.Color + "]" + ch
			blt.Print(v.X, v.Y, glyph)
		}
	}
}

func PrintUI(c *Creature) {
	/* Function PrintUI takes *Creature (it's supposed to be player) as argument.
	   It prints UI infos on the right side of screen.
	   For now its functionality is very modest, but it will expand when
	   new elements of game mechanics will be introduced. So, for now, it
	   provides only one basic, yet essential information: player's HP. */
	blt.Layer(UILayer)
	name := "Player"
	blt.Print(UIPosX, UIPosY, name)
	hpBar := "[color=red]"
	for hp := 1; hp <= c.HPMax; hp++ {
		character := "♥"
		if hp > c.HPCurrent {
			character = "♡"
		}
		hpBar = hpBar + character
	}
	blt.Print(UIPosX, UIPosY+1, hpBar)
	pos := 3
	for i := 0; i < len(c.Equipment); i++ {
		color := "[color=gray]"
		if i == c.ActiveWeapon {
			color = "[color=white]"
		}
		blt.Print(UIPosX, UIPosY+pos,
			color+strconv.Itoa(i+1)+". "+c.Equipment[i].Name)
		for j := 0; j < c.Equipment[i].AmmoCurrent; j++ {
			blt.Print(UIPosX+j+4+3+1, UIPosY+pos+1, "[color=dark yellow]|")
		}
		for k := c.Equipment[i].AmmoCurrent; k < c.Equipment[i].AmmoMax; k++ {
			blt.Print(UIPosX+k+4+3+1, UIPosY+pos+1, "[color=darkest yellow]|")
		}
		if c.Equipment[i].Cock == true {
			cockedPosX := UIPosX + c.Equipment[i].AmmoMax + 1 + 4 + 3 + 1
			cockedIcon := CockedIcon
			if c.Equipment[i].Cocked == true {
				cockedIcon = "[color=dark green]" + CockedIcon
			} else {
				cockedIcon = "[color=dark red]" + CockedIcon
			}
			blt.Print(cockedPosX, UIPosY+pos+1, cockedIcon)
		}
		var ranges = []string{"▇", "▇", "▇"}
		rangesStr := ""
		for j, _ := range ranges {
			if c.Equipment[i].Slot == SlotWeaponMelee {
				break
			}
			val := c.Equipment[i].Ranges[j]
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
		blt.Print(UIPosX+3, UIPosY+pos+1, rangesStr)
		pos += 2
	}
	instructions := "[color=#EDEAE0]────────────────────\n1 - grab rifle |\n2 - grab revolver -\n3 - grab melee }\nf - target/fire/cock\nr - reload\ni - inspect\ng - pick up\np - pull lever\n────────────────────"
	blt.Print(UIPosX, UIPosY+pos, instructions)
}

func PrintLog() {
	/* Function PrintLog prints game messages at the bottom of screen. */
	blt.Layer(UILayer)
	PrintMessages(LogPosX, LogPosY, "")
}

func RenderAll(b Board, o Objects, c Creatures) {
	/* Function RenderAll prints every tile and character on game screen.
	   Takes board slice (ie level map), slice of objects, and slice of creatures
	   as arguments.
	   At first, it clears whole terminal window, then uses arguments:
	   CastRays (for raycasting FOV) of first object (assuming that it is player),
	   then calls functions for printing map, objects and creatures.
	   Calls PrintLog that writes message log.
	   At the end, RenderAll calls blt.Refresh() that makes
	   changes to the game window visible. */
	blt.Clear()
	CastRays(b, c[0].X, c[0].Y)
	PrintBoard(b, c)
	PrintDeadCreatures(b, c)
	PrintObjects(b, o, c)
	PrintAliveCreatures(b, c)
	PrintUI((c)[0])
	PrintLog()
	blt.Refresh()
}
