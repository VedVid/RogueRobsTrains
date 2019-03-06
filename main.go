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
	LevelInt int
	LevelStr string
	Levels   []string
	Alive    int
}

var MsgBuf = []string{}
var LastTarget *Creature
var RailsMod = false
var TimerMod = 10
var G = new(Game)

func main() {
	var cells = new(Board)
	var objs = new(Objects)
	var actors = new(Creatures)
	StartGame(cells, actors, objs)
	timer := 0
	for {
		if timer > 100 {
			timer = 0
		}
		if (*actors)[0].HPCurrent <= 0 {
			DeleteSaves()
			break
		}
		if G.Alive == 0 {
			G.Alive = -1
			AddMessage("All enemies are down. You may proceed.")
			for x := len(*cells) - 8; x < len(*cells); x++ {
				for y := 0; y < len((*cells)[0]); y++ {
					if (*cells)[x][y].Name == "doors to next carriage" {
						(*cells)[x][y].Color = "#FFCC00"
						(*cells)[x][y].ColorDark = "#FFCC00"
					}
				}
			}
		}
		if G.LevelStr != G.Levels[G.LevelInt] {
			var err error
			player := (*actors)[0]
			*cells, *actors, err = LoadJsonMap(G.Levels[G.LevelInt])
			if err != nil {
				fmt.Println(err)
			}
			player.X, player.Y = (*actors)[0].X, (*actors)[0].Y
			(*actors)[0] = player
			player.HPCurrent = player.HPMax
			G.LevelStr = G.Levels[G.LevelInt]
			G.Alive = len(*actors) - 1
			for i := 0; i < len(*objs); i++ {
				(*objs)[i] = nil
			}
			*objs = (*objs)[:0]
		}
		if timer%TimerMod == 0 {
			cells.MoveMap()
		}
		RenderAll(*cells, *objs, *actors)
		if blt.HasInput() == true {
			key := blt.Read()
			if key == blt.TK_S && blt.Check(blt.TK_SHIFT) != 0 {
				err := SaveGame(*cells, *actors, *objs)
				if err != nil {
					fmt.Println(err)
				}
				break
			} else if key == blt.TK_Q && blt.Check(blt.TK_SHIFT) != 0 {
				AddMessage("Do you want to quit the game?")
				AddMessage("It will delete the saves as well. [[Y/N]]")
				RenderAll(*cells, *objs, *actors)
				confirm := false
				for {
					keyConfirm := blt.Read()
					if keyConfirm == blt.TK_Y {
						confirm = true
						break
					} else if keyConfirm == blt.TK_N {
						break
					} else {
						continue
					}
				}
				if confirm == true {
					DeleteSaves()
					break
				} else {
					AddMessage("OK, then...")
				}
			} else {
				turnSpent := Controls(key, (*actors)[0], cells, actors, objs)
				if turnSpent == true {
					CreaturesTakeTurn(*cells, *actors, objs)
				}
			}
		}
		timer++
	}
	blt.Close()
}

func NewGame(b *Board, c *Creatures, o *Objects) {
	/* Function NewGame initializes game state - creates player, monsters, and game map.
	   This implementation is generic-placeholder, for testing purposes. */
	playerMelee, err := NewObject(0, 0, MeleeWeapons[RandInt(len(MeleeWeapons)-1)])
	if err != nil {
		fmt.Println(err)
	}
	playerSecondary, err := NewObject(0, 0, SecondaryWeapons[RandInt(len(SecondaryWeapons)-1)])
	if err != nil {
		fmt.Println(err)
	}
	playerPrimary, err := NewObject(0, 0, PrimaryWeapons[RandInt(len(PrimaryWeapons)-1)])
	if err != nil {
		fmt.Println(err)
	}
	*o = Objects{}
	*b, *c, err = LoadJsonMap("trainStart.json")
	if err != nil {
		fmt.Println(err)
	}
	(*c)[0].Equipment = Objects{playerPrimary, playerSecondary, playerMelee}
	G.LevelInt = 0
	G.Levels = []string{"trainStart.json"}
	var middleLevels = []string{"train1.json", "train2.json", "train3.json", "train4.json"}
	rand.Shuffle(len(middleLevels), func(i, j int) {
		middleLevels[i], middleLevels[j] = middleLevels[j], middleLevels[i]
	})
	G.Levels = append(G.Levels, middleLevels...)
	G.Levels = append(G.Levels, "trainFinal1.json", "trainFinal2.json")
	G.LevelStr = G.Levels[G.LevelInt]
	G.Alive = len(*c) - 1
}

func StartGame(b *Board, c *Creatures, o *Objects) {
	/* Function StartGame determines if game save is present (and valid), then
	   loads data, or initializes new game.
	   Panics if some-but-not-all save files are missing. */
	_, errBoard := os.Stat(MapPathGob)
	_, errCreatures := os.Stat(CreaturesPathGob)
	_, errObjects := os.Stat(ObjectsPathGob)
	_, errGame := os.Stat(GamePathGob)
	_, errTimer := os.Stat(TimerPathGob)
	_, errRails := os.Stat(RailsPathGob)
	if errBoard == nil && errCreatures == nil && errObjects == nil &&
		errGame == nil && errTimer == nil && errRails == nil {
		LoadGame(b, c, o)
	} else if errBoard != nil && errCreatures != nil && errObjects != nil &&
		errGame != nil && errTimer != nil && errRails != nil {
		NewGame(b, c, o)
	} else {
		txt := CorruptedSaveError(errBoard, errCreatures, errObjects, errGame, errTimer, errRails)
		fmt.Println("Error: save files are corrupted: " + txt)
		panic(-1)
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	InitializeFOVTables()
	InitializeBLT()
}
