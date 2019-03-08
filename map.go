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
	"math/rand"
	"unicode/utf8"
)

type Tile struct {
	// Tiles are map cells - floors, walls, doors.
	BasicProperties
	VisibilityProperties
	Explored bool
	CollisionProperties
}

type MapJson struct {
	// For unmarshalling json data.
	Cells          []string
	Data           [][]int
	Layouts        [][][]string
	Char           map[string]string
	Name           map[string]string
	Color          map[string]string
	ColorDark      map[string]string
	Layer          map[string]int
	AlwaysVisible  map[string]bool
	Explored       map[string]bool
	Blocked        map[string]bool
	BlocksSight    map[string]bool
	MonstersCoords [][]int
	MonstersTypes  []string
}

/* Board is map representation, that uses 2d slice
   to hold data of its every cell. */
type Board [][]*Tile

var grassColors = []string{
	"#8F9779", "#8F9779", "#4F7942", "#4F7942", "#6c7c59", "#6c7c59", "#A9BA9D", "#A9BA9D",
	"#8A9A5B", "#8A9A5B", "#6C7C59", "#6C7C59", "#4B5320", "#4B5320", "#355E3B", "#355E3B",
	"#444C38", "#444C38", "#679267", "#679267",
	"#C3B091", "#826644", "#D2B48C", "#5C5248", "#C19A6B",
	"#5E716A", "#98817B",
}

var stoneColors = []string{
	"#989898", "#555555", "#B2BEB5", "#727472", "#928E85", "#708090", "#AA98A9", "#98817B",
}

func NewTile(layer, x, y int, character, name, color, colorDark string,
	alwaysVisible, explored, blocked, blocksSight bool) (*Tile, error) {
	/* Function NewTile takes all values necessary by its struct,
	   and creates then returns Tile. */
	var err error
	if layer < 0 {
		txt := LayerError(layer)
		err = errors.New("Tile layer is smaller than 0." + txt)
	}
	if x < 0 || x >= MapSizeX || y < 0 || y >= MapSizeY {
		txt := CoordsError(x, y)
		err = errors.New("Tile coords is out of window range." + txt)
	}
	if utf8.RuneCountInString(character) != 1 {
		txt := CharacterLengthError(character)
		err = errors.New("Tile character string length is not equal to 1." + txt)
	}
	tileBasicProperties := BasicProperties{x, y, character, name, color,
		colorDark}
	tileVisibilityProperties := VisibilityProperties{layer, alwaysVisible}
	tileCollisionProperties := CollisionProperties{blocked, blocksSight}
	tileNew := &Tile{tileBasicProperties, tileVisibilityProperties,
		explored, tileCollisionProperties}
	return tileNew, err
}

func InitializeEmptyMap() Board {
	/* Function InitializeEmptyMap returns new Board, filled with
	   generic (ie "empty") tiles.
	   It starts by declaring 2d slice of *Tile - unfortunately, Go seems to
	   lack simple way to do it, therefore it's necessary to use
	   the first for loop.
	   The second, nested loop initializes specific Tiles within Board bounds. */
	b := make([][]*Tile, MapSizeX)
	for i := range b {
		b[i] = make([]*Tile, MapSizeY)
	}
	for x := 0; x < MapSizeX; x++ {
		for y := 0; y < MapSizeY; y++ {
			var err error
			b[x][y], err = NewTile(BoardLayer, x, y, ".", "floor", "light gray",
				"dark gray", true, false, false, false)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return b
}

func (b *Board) MoveMap() {
	const railY1 = 7+1
	const railY2 = 12-1
	for x := 0; x < MapSizeX; x++ {
		for y := 0; y < MapSizeY; y++ {
			if (*b)[x][y].Name == "railroad" {
				if (*b)[x][y].Char == "━" {
					continue
				} else {
					(*b)[x][y] = NewBackgroundTile(*b, x, y)
				}
			}
			if (x == MapSizeX-1) && ((*b)[x][y].Name == "grass" || (*b)[x][y].Name == "stone") {
				(*b)[x][y] = NewBackgroundTile(*b, x, y)
				continue
			}
			if (*b)[x][y].Name == "grass" || (*b)[x][y].Name == "stone" {
				if (*b)[x+1][y].Name == "grass" || (*b)[x+1][y].Name == "stone" {
					(*b)[x][y].Name = (*b)[x+1][y].Name
					(*b)[x][y].Char = (*b)[x+1][y].Char
					(*b)[x][y].Color = (*b)[x+1][y].Color
					(*b)[x][y].ColorDark = (*b)[x+1][y].ColorDark
				} else {
					(*b)[x][y] = NewBackgroundTile(*b, x, y)
				}
			}
		}
	}
	for rx := 0; rx < MapSizeX; rx++ {
		for ry := railY1; ry <= railY2; ry++ {
			if (*b)[rx][ry].Name != "grass" && (*b)[rx][ry].Name != "stone" {
				continue
			}
			if RailsMod == false && rx%2 == 0 {
				(*b)[rx][ry].Name = "railroad"
				(*b)[rx][ry].Char = "┃"
				(*b)[rx][ry].Color = "#483C32"
				(*b)[rx][ry].ColorDark = "#483C32"
			} else if RailsMod == true && rx%2 != 0 {
				(*b)[rx][ry].Name = "railroad"
				(*b)[rx][ry].Char = "┃"
				(*b)[rx][ry].Color = "#483C32"
				(*b)[rx][ry].ColorDark = "#483C32"
			}
		}
	}
	if RailsMod == false {
		RailsMod = true
	} else {
		RailsMod = false
	}
}

func NewBackgroundTile(b Board, x, y int) *Tile {
	t := b[x][y]
	val1 := RandInt(100)
	if val1 <= 85 {
		val2 := RandInt(70)
		if val2 <= 10 {
			t.Char = ";"
		} else if val2 <= 20 {
			t.Char = ":"
		} else if val2 <= 30 {
			t.Char = "'"
		} else if val2 <= 40 {
			t.Char = "\""
		} else if val2 <= 50 {
			t.Char = ","
		} else if val2 <= 60 {
			t.Char = "."
		} else if val2 <= 70 {
			t.Char = "`"
		}
		t.Name = "grass"
		t.Color = grassColors[RandInt(len(grassColors)-1)]
		t.ColorDark = t.Color
	} else {
		val2 := RandInt(20)
		if val2 <= 10 {
			t.Char = "^"
		} else if val2 <= 20 {
			t.Char = "*"
		}
		t.Name = "stone"
		t.Color = stoneColors[RandInt(len(stoneColors)-1)]
		t.ColorDark = t.Color
	}
	return t
}

func ReplaceTile(t *Tile, s string, m *MapJson) {
	/* ReplaceTile is function that takes tile, string (supposed to be
	   one-character-lenght - symbol of map tile, taken from json map) and
	   MapJson (ie unmarshalled json map).
	   It uses m's legend to overwrite old map values with data read from file. */
	t.Char = m.Char[s]
	t.Name = m.Name[s]
	t.Color = m.Color[s]
	t.ColorDark = m.ColorDark[s]
	if t.Name == "grass" {
		val1 := RandInt(100)
		if val1 <= 85 {
			val2 := RandInt(70)
			if val2 <= 10 {
				t.Char = ";"
			} else if val2 <= 20 {
				t.Char = ":"
			} else if val2 <= 30 {
				t.Char = "'"
			} else if val2 <= 40 {
				t.Char = "\""
			} else if val2 <= 50 {
				t.Char = ","
			} else if val2 <= 60 {
				t.Char = "."
			} else if val2 <= 70 {
				t.Char = "`"
			}
			t.Color = grassColors[RandInt(len(grassColors)-1)]
			t.ColorDark = t.Color
		} else {
			val2 := RandInt(20)
			if val2 <= 10 {
				t.Char = "^"
			} else if val2 <= 20 {
				t.Char = "*"
			}
			t.Name = "stone"
			t.Color = stoneColors[RandInt(len(stoneColors)-1)]
			t.ColorDark = t.Color
		}
	}
	t.Layer = m.Layer[s]
	t.AlwaysVisible = m.AlwaysVisible[s]
	t.Explored = m.Explored[s]
	t.Blocked = m.Blocked[s]
	t.BlocksSight = m.BlocksSight[s]
}

func LoadJsonMap(mapFile string) (Board, Creatures, error) {
	/* Function LoadJsonMap takes string (name of json map file) as argument,
	   and returns Board (ie map), Creatures (included in premade json maps)
	   and error.
	   It uses new type - struct MapJson - to store all values read from file.
	   Panics if unmarshalling encounters any error.
	   Other possible errors are about internal structure of json file:
	       - length of Data and Layouts has to be the same
	       - length of MonstersCoords and MonstersTypes has to be the same.
	   It is important because instead of using multi-type json lists
	   (it would be possible to store map monsters as [x: int, y: int, file: string])
	   there are independent structures. The reason is Go's limitations: bot lists
	   (slices) and dictionaries (maps) are strongly typed. (Un)Marshalling multi-type
	   lists would be cumbersome. On the other hand, it means that creating and editing
	   json maps require discipline.
	   After error checking, three major operations are queued.
	   At first, game reads json map (Cells) and modifies (previously initialized)
	   tiles regarding to json legend (Char, Name (...), BlocksSight).
	   Then it repeats this operation for every area marked as "randomly generated".
	   Some important points to make about these areas:
	       - they are not created *randomly*
	           = areas ("rooms") are specified in JsonMap.Data
	           = they are filled using prefabs (JsonMap.Layouts)
	   At the end, monsters are created and placed on map (their datas are stored
	   in json map as MonstersCoords (x, y) and MonstersTypes (their json files). */
	var jsonMap = &MapJson{}
	var err error
	err = MapFromJson(MapsPathJson+mapFile, jsonMap)
	if err != nil {
		fmt.Println(err)
		panic(-1)
	}
	cells := jsonMap.Cells
	data := jsonMap.Data
	layouts := jsonMap.Layouts
	// Number of items in data should match number of layouts.
	if len(data) != len(layouts) {
		txt := MapDataLayoutsError((len(data)), len(layouts), mapFile)
		err = errors.New("Length of data and layouts does not match. " + txt)
	}
	thisMap := InitializeEmptyMap()
	for x := 0; x < len(cells[0]); x++ {
		for y := 0; y < len(cells); y++ {
			// y,x because - due to 2darray nature - there is height first, width later...
			ReplaceTile(thisMap[x][y], string(cells[y][x]), jsonMap)
		}
	}
	for i, room := range data {
		layoutsToChoose := layouts[i]
		layout := layoutsToChoose[rand.Intn(len(layoutsToChoose))]
		for x := 0; x < len(layout[0]); x++ {
			for y := 0; y < len(layout); y++ {
				ReplaceTile(thisMap[room[0]+x][room[1]+y], string(layout[y][x]), jsonMap)
			}
		}
	}
	coords := jsonMap.MonstersCoords
	aiTypes := jsonMap.MonstersTypes
	if len(coords) != len(aiTypes) {
		txt := MapMonstersCoordsAiError(len(coords), len(aiTypes), mapFile)
		err = errors.New("Length of MonstersCoords and MonstersTypes does not match. " + txt)
	}
	var enemies = []string{"dumbMelee.json", "patherMelee.json", "patherMelee.json", "dumbRanged.json",
		"dumbRanged.json", "patherRanged.json", "patherRanged.json", "patherRanged.json"}
	var creatures = Creatures{}
	for k := 0; k < len(coords); k++ {
		if aiTypes[k] == "player" {
			player, err := NewPlayer(coords[k][0], coords[k][1])
			if err != nil {
				fmt.Println(err)
			}
			creatures = append(creatures, player)
			continue
		}
		aitype := aiTypes[k] + ".json"
		if aiTypes[k] == "any" {
			aitype = enemies[RandRange(0, len(enemies)-1)]
		}
		if aiTypes[k] == "maybe" {
			if RandInt(100) <= 33 {
				aitype = enemies[RandRange(0, len(enemies)-1)]
			} else {
				continue
			}
		}
		monster, err := NewCreature(coords[k][0], coords[k][1], aitype)
		if err != nil {
			fmt.Println(err)
		}
		creatures = append(creatures, monster)
	}
	for _, area := range data {
		n := 0
		val := RandInt(100)
		if val < 50 {
			val2 := RandInt(75)
			if val2 < 50 {
				n = 1
			} else {
				n = 2
			}
		}
		for x := area[0]; x < area[0]+area[2]; x++ {
			for y := area[1]; y < area[1]+area[3]; y++ {
				if n == 0 {
					goto Areas
				}
				if thisMap[x][y].Blocked == false && thisMap[x][y].BlocksSight == false {
					chances := 50
					if Config.Monsters == MonstersEasy {
						chances = 25
					} else if Config.Monsters == MonstersHard {
						chances = 75
					}
					if RandInt(100) <= chances {
						aitype := enemies[RandRange(0, len(enemies)-1)]
						monster, err := NewCreature(x, y, aitype)
						if err != nil {
							fmt.Println(err)
						}
						creatures = append(creatures, monster)
						n--
					}
				}
			}
		}
	Areas:
		continue
	}
	for i := 0; i < len(creatures); i++ {
		monster := creatures[i]
		if monster.AIType == MeleeDumbAI || monster.AIType == MeleePatherAI {
			monster.ActiveWeapon = SlotWeaponMelee
		} else {
			monster.ActiveWeapon = RandRange(0, 1)
		}
		weapon := monster.ActiveWeapon
		if monster.Equipment[weapon] == nil {
			if weapon == SlotWeaponMelee {
				monster.Equipment[weapon], _ = NewObject(0, 0, MeleeWeapons[RandRange(0, len(MeleeWeapons)-1)])
			} else if weapon == SlotWeaponSecondary {
				monster.Equipment[weapon], _ = NewObject(0, 0, SecondaryWeapons[RandRange(0, len(SecondaryWeapons)-1)])
			} else if weapon == SlotWeaponPrimary {
				monster.Equipment[weapon], _ = NewObject(0, 0, PrimaryWeapons[RandRange(0, len(PrimaryWeapons)-1)])
			}
		}
	}
	if mapFile == "trainFinal2.json" {
		RailsMod = true
	}
	return thisMap, creatures, err
}
