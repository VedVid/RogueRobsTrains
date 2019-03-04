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
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Game struct {
	LevelInt  int
	LevelStr string
	Levels []string
}

var MsgBuf = []string{}
var LastTarget *Creature

var G = new(Game)

func main() {
	var cells = new(Board)
	var objs = new(Objects)
	var actors = new(Creatures)
	StartGame(cells, actors, objs, G)
	for {
		if G.LevelStr != G.Levels[G.LevelInt] {
			var err error
			player := (*actors)[0]
			*cells, *actors, err = LoadJsonMap(G.Levels[G.LevelInt])
			if err != nil {
				fmt.Println(err)
			}
			player.X, player.Y = (*actors)[0].X, (*actors)[0].Y
			(*actors)[0] = player
			G.LevelStr = G.Levels[G.LevelInt]
		}
		RenderAll(*cells, *objs, *actors)
		key := blt.Read()
		if key == blt.TK_S && blt.Check(blt.TK_SHIFT) != 0 {
			err := SaveGame(*cells, *actors, *objs, *G)
			if err != nil {
				fmt.Println(err)
			}
			break
		} else if key == blt.TK_Q && blt.Check(blt.TK_SHIFT) != 0 ||
			(*actors)[0].HPCurrent <= 0 {
			DeleteSaves()
			break
		} else {
			turnSpent := Controls(key, (*actors)[0], cells, actors, objs)
			if turnSpent == true {
				CreaturesTakeTurn(*cells, *actors, *objs)
			}
		}
	}
	blt.Close()
}

func NewGame(b *Board, c *Creatures, o *Objects, g *Game) {
	/* Function NewGame initializes game state - creates player, monsters, and game map.
	   This implementation is generic-placeholder, for testing purposes. */
	playerMelee, err := NewObject(0, 0, "BowieKnife.json")
	if err != nil {
		fmt.Println(err)
	}
	playerSecondary, err := NewObject(0, 0, "Remington1875.json")
	if err != nil {
		fmt.Println(err)
	}
	playerPrimary, err := NewObject(0, 0, "SpencerRepeater.json")
	if err != nil {
		fmt.Println(err)
	}
	*o = Objects{}
	*b, *c, err = LoadJsonMap("trainStart.json")
	if err != nil {
		fmt.Println(err)
	}
	(*c)[0].Equipment = Objects{playerPrimary, playerSecondary, playerMelee}
	for i := 0; i < len(*c); i++ {
		monster := (*c)[i]
		weapon := monster.ActiveWeapon
		if monster.Equipment[weapon] == nil {
			if weapon == SlotWeaponMelee {
				monster.Equipment[weapon], _ = NewObject(0, 0, "BowieKnife.json")
			} else if weapon == SlotWeaponSecondary {
				monster.Equipment[weapon], _ = NewObject(0, 0, "Remington1875.json")
			} else if weapon == SlotWeaponPrimary {
				monster.Equipment[weapon], _ = NewObject(0, 0, "SpencerRepeater.json")
			}
		}
	}
	g.LevelInt = 0
	g.Levels = []string{"trainStart.json", "train1.json", "train2.json", "train3.json", "train4.json",
						"trainFinal1.json", "trainFinal2.json"}
	g.LevelStr = g.Levels[g.LevelInt]
}

func StartGame(b *Board, c *Creatures, o *Objects, g *Game) {
	/* Function StartGame determines if game save is present (and valid), then
	   loads data, or initializes new game.
	   Panics if some-but-not-all save files are missing. */
	_, errBoard := os.Stat(MapPathGob)
	_, errCreatures := os.Stat(CreaturesPathGob)
	_, errObjects := os.Stat(ObjectsPathGob)
	_, errGame := os.Stat(GamePathGob)
	if errBoard == nil && errCreatures == nil && errObjects == nil && errGame == nil {
		LoadGame(b, c, o, g)
	} else if errBoard != nil && errCreatures != nil && errObjects != nil && errGame != nil {
		NewGame(b, c, o, g)
	} else {
		txt := CorruptedSaveError(errBoard, errCreatures, errObjects)
		fmt.Println("Error: save files are corrupted: " + txt)
		panic(-1)
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	InitializeFOVTables()
	InitializeBLT()
}
