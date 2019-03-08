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
	"errors"
	"fmt"
	"sort"
	"strconv"
	"unicode/utf8"
)

const (
	MaxMessageBuffer = WindowSizeY - MapSizeY
)

func PrintMenu(x, y int, header string, options []string) {
	/* Function PrintMenu takes four arguments: two ints that are
	   top-left corner of menu, header, and slice of options.
	   If header is empty, text is moved one tile higher to
	   avoid wasting space.
	   During execution, it joins header and all of options in
	   one text, with additional formatting.
	   For example, header "MyMenu" and options ["first", "two"]
	   would produce that kind of output:
	       MyMenu
	       a) first
	       b) two
	    It refreshed terminal and waits for player input at the end. */
	blt.ClearArea(UIPosX, UIPosY, UISizeX, UISizeY)
	if header == "" {
		y--
	}
	txt := header
	for i, v := range options {
		txt = txt + "\n" + OrderToCharacter(i) + ") " + v
	}
	txt = txt + "\n[[ESC]] back"
	blt.Print(x, y, txt)
	blt.Refresh()
}

func PrintInventoryMenu(x, y int, header string, options Objects) {
	/* PrintInventoryMenu is helper function that takes Objects
	   as its main argument, and adds their names (currently
	   their symbol representation, due to some strange decisions
	   made by dev, objects doesn't have names yet) to the opts
	   slice of strings, then calls PrintMenu using that list.
	   Unfortunately that kind of "hack" is necessary, because
	   Go doesn't support generics and optional arguments,
	   and still doesn't provide sensible alternatives.
	   I'd like to just pass Objects to the PrintMenu func. */
	var opts = []string{}
	for _, v := range options {
		opts = append(opts, v.Name)
	}
	PrintMenu(x, y, header, opts)
}

func PrintEquipmentMenu(x, y int, header string, options Objects) {
	/* Similar to PrintInventoryMenu, but it sorts options
	   by their Slots initially, and slot in showed before
	   item name.
	   Note that it shows Creature's slots,
	   not all equippable objects in inventory.
	   Because of this, it is necessary to find "true" length
	   of options, skipping all nil pointers.
	   Unfortunately, it may crash in future, with
	   more slots involved. */
	var opts = []string{}
	for i := 0; i < len(options); i++ {
		txt := ""
		if options[i] != nil {
			txt = "[[" + SlotStrings[i] + "]] " + options[i].Name
		} else {
			txt = "[[" + SlotStrings[i] + "]] empty"
		}
		opts = append(opts, txt)
	}
	PrintMenu(x, y, header, opts)
}

func PrintEquippables(x, y int, header string, options Objects) {
	/* PrintEquippables is function that prints list of equippables. */
	var opts = []string{}
	for _, v := range options {
		opts = append(opts, v.Name)
	}
	PrintMenu(x, y, header, opts)
}

func PrintMessages(x, y int, header string) {
	/* PrintMessages works as PrintMenu, but it
	   will not format text in special way. */
	if header == "" {
		y--
	}
	txt := header
	for _, v := range MsgBuf {
		txt = txt + "\n" + v
	}
	blt.Print(x, y, txt)
}

func AddMessage(message string) {
	/* AddMessage is function that adds message
	   to the MessageBuffer. It removes the oldest
	   line to keep size set in MaxMessageBuffer.
	   But first, it checks if passed message is
	   not longer than whole message log.
	   This is mostly harmless, so AddMessage
	   does not returns error, but prints it
	   at its own. */
	var err error
	messageLen := utf8.RuneCountInString(message)
	if messageLen > LogSizeX {
		txt := MessageLengthError(message, messageLen, LogSizeX)
		err = errors.New("Message is too long to fit message log. " + txt)
		fmt.Println(err)
	}
	if len(MsgBuf) < MaxMessageBuffer {
		MsgBuf = append(MsgBuf, message)
	} else {
		MsgBuf = append(MsgBuf[1:], message)
	}
	PrintLog()
	blt.Refresh()
}

func RemoveLastMessage() {
	/* Function RemoveLastMessage is called when it is necessary to remove
	   last message from buffer, even if said buffer is not full.
	   It removes last message, clears its area, and reprints log. */
	MsgBuf = MsgBuf[:len(MsgBuf)-1]
	blt.Layer(UILayer)
	blt.ClearArea(LogPosX, LogPosY, LogPosX+LogSizeX, LogPosY+LogSizeY)
	PrintLog()
	blt.Refresh()
}

func HandleHighScores() {
	Scores.Scores = append(Scores.Scores, Config.Score)
	sort.Sort(sort.Reverse(sort.IntSlice(Scores.Scores)))
	size := len(Scores.Scores)
	if size > 10 {
		size = 10
		Scores.Scores = Scores.Scores[:10]
	}
	blt.Clear()
	for i := 0; i < size; i++ {
		blt.Color(blt.ColorFromName("white"))
		if Scores.Scores[i] == Config.Score {
			blt.Color(blt.ColorFromName("yellow"))
		}
		txt := strconv.Itoa(Scores.Scores[i])
		blt.Print(((WindowSizeX/2) - (utf8.RuneCountInString(txt))), 5+i, strconv.Itoa(i+1) + ". "+ txt)
	}
	for {
		blt.Refresh()
		key := blt.Read()
		if key == blt.TK_ENTER || key == blt.TK_SPACE || key == blt.TK_ESCAPE {
			break
		}
	}
	SaveScores()
}

func PrintVictoryScreen() {
	for {
		yourScore := Stats.Killed - Stats.Lost
		if yourScore < 0 {
			yourScore = 0
		}
		yourScore = (yourScore * Config.Score) / 100
		if yourScore < 10 {
			yourScore = 10
		}
		Config.Score = yourScore
		blt.Clear()
		blt.Layer(UILayer)
		line1 := "You did it!"
		line2 := "Finally did it!"
		line3 := "You killed everyone in the train,"
		line4 := "reached engine, killed driver as well,"
		line5 := "and pulled breake lever!"
		line6 := "Now, you can just unload these gold-filled chests"
		line7 := "from this train and live rich and well."
		line8 := "Your score is: " + strconv.Itoa(Config.Score)
		line1len := utf8.RuneCountInString(line1)
		line2len := utf8.RuneCountInString(line2)
		line3len := utf8.RuneCountInString(line3)
		line4len := utf8.RuneCountInString(line4)
		line5len := utf8.RuneCountInString(line5)
		line6len := utf8.RuneCountInString(line6)
		line7len := utf8.RuneCountInString(line7)
		line8len := utf8.RuneCountInString(line8)
		posy := (WindowSizeY / 2) - 5
		blt.Print(((WindowSizeX / 2) - (line1len / 2)), posy, line1)
		blt.Print(((WindowSizeX / 2) - (line2len / 2)), posy+1, line2)
		blt.Print(((WindowSizeX / 2) - (line3len / 2)), posy+2, line3)
		blt.Print(((WindowSizeX / 2) - (line4len / 2)), posy+3, line4)
		blt.Print(((WindowSizeX / 2) - (line5len / 2)), posy+4, line5)
		blt.Print(((WindowSizeX / 2) - (line6len / 2)), posy+5, line6)
		blt.Print(((WindowSizeX / 2) - (line7len / 2)), posy+6, line7)
		blt.Print(((WindowSizeX / 2) - (line8len / 2)), posy+8, line8)
		blt.Refresh()
		key := blt.Read()
		if key == blt.TK_ESCAPE || key == blt.TK_ENTER || key == blt.TK_SPACE {
			break
		}
	}
	HandleHighScores()
}

func DeadScreen() {
	for {
		yourScore := Stats.Killed - Stats.Lost
		if yourScore < 0 {
			yourScore = 0
		}
		yourScore = (yourScore * Config.Score) / 100
		if yourScore < 10 {
			yourScore = 10
		}
		Config.Score = yourScore
		blt.Clear()
		blt.Layer(UILayer)
		line1 := "Maybe robbing this train was bad idea?"
		line2 := "Maybe it was good idea, but you made serie of mistakes?"
		line3 := "You'll never know."
		line4 := "They got you. Ten times."
		line5 := "No chances to leave this damn train alive..."
		line1len := utf8.RuneCountInString(line1)
		line2len := utf8.RuneCountInString(line2)
		line3len := utf8.RuneCountInString(line3)
		line4len := utf8.RuneCountInString(line4)
		line5len := utf8.RuneCountInString(line5)
		line8 := "Your score is: " + strconv.Itoa(Config.Score)
		line8len := utf8.RuneCountInString(line8)
		posy := (WindowSizeY / 2) - 4
		blt.Print(((WindowSizeX / 2) - (line1len / 2)), posy, line1)
		blt.Print(((WindowSizeX / 2) - (line2len / 2)), posy+1, line2)
		blt.Print(((WindowSizeX / 2) - (line3len / 2)), posy+2, line3)
		blt.Print(((WindowSizeX / 2) - (line4len / 2)), posy+3, line4)
		blt.Print(((WindowSizeX / 2) - (line5len / 2)), posy+4, line5)
		blt.Print(((WindowSizeX / 2) - (line8len / 2)), posy+6, line8)
		blt.Refresh()
		key := blt.Read()
		if key == blt.TK_ESCAPE || key == blt.TK_ENTER || key == blt.TK_SPACE {
			break
		}
	}
	HandleHighScores()
}

func MainMenu(cfg *Cfg) {
	lives := livesNormal
	monsters := MonstersNormal
	reloading := AmmoUnlimited
	animations := AnimationsFalse
	score := 100
	if CfgIsHere == true {
		lives = Config.Lives
		monsters = Config.Monsters
		reloading = Config.Reloading
		animations = Config.Animations
		score = Config.Score
	}
	for {
		blt.Clear()
		livesString := ""
		if lives == livesNormal {
			livesString = "normal"
		} else if lives == livesEasy {
			livesString = "easy"
		} else if lives == livesHard {
			livesString = "hard"
		}
		line1 := "<a> ← Lives: " + livesString + " → <A>"
		monstersString := ""
		if monsters == MonstersNormal {
			monstersString = "normal"
		} else if monsters == MonstersEasy {
			monstersString = "fewer"
		} else if monsters == MonstersHard {
			monstersString = "more"
		}
		line2 := "<b> ← Enemies: " + monstersString + " → <B>"
		reloadingString := ""
		if reloading == AmmoUnlimited {
			reloadingString = "unlimited"
		} else {
			reloadingString = "limited"
		}
		line3 := "<c> ← Reloading: " + reloadingString + " → <C>"
		animationsString := ""
		if animations == AnimationsFalse {
			animationsString = "off"
		} else {
			animationsString = "on"
		}
		line4 := "<d> ← Animations: " + animationsString + " → <D>"
		line5 := "Score multiplier: " + strconv.Itoa(score) + "%"
		line6 := "Press <ENTER> to proceed."
		line1len := utf8.RuneCountInString(line1)
		line2len := utf8.RuneCountInString(line2)
		line3len := utf8.RuneCountInString(line3)
		line4len := utf8.RuneCountInString(line4)
		line5len := utf8.RuneCountInString(line5)
		line6len := utf8.RuneCountInString(line6)
		posy := (WindowSizeY / 2) - 3
		blt.Print(((WindowSizeX / 2) - (line1len / 2)), posy, line1)
		blt.Print(((WindowSizeX / 2) - (line2len / 2)), posy+1, line2)
		blt.Print(((WindowSizeX / 2) - (line3len / 2)), posy+2, line3)
		blt.Print(((WindowSizeX / 2) - (line4len / 2)), posy+3, line4)
		blt.Print(((WindowSizeX / 2) - (line5len / 2)), posy+5, line5)
		blt.Print(((WindowSizeX / 2) - (line6len / 2)), posy+7, line6)
		blt.Refresh()
		key := blt.Read()
		if key == blt.TK_ENTER {
			cfg.Lives = lives
			cfg.Monsters = monsters
			cfg.Reloading = reloading
			cfg.Score = score
			cfg.Animations = animations
			break
		}
		if key == blt.TK_A && blt.Check(blt.TK_SHIFT) != 0 {
			if lives == livesEasy {
				lives = livesNormal
				score += 25
			} else if lives == livesNormal {
				lives = livesHard
				score += 25
			} else {
				continue
			}
		} else if key == blt.TK_A {
			if lives == livesNormal {
				lives = livesEasy
				score -= 25
			} else if lives == livesHard {
				lives = livesNormal
				score -= 25
			} else {
				continue
			}
		} else if key == blt.TK_B && blt.Check(blt.TK_SHIFT) != 0 {
			if monsters == MonstersEasy {
				monsters = MonstersNormal
				score += 25
			} else if monsters == MonstersNormal {
				monsters = MonstersHard
				score += 25
			} else {
				continue
			}
		} else if key == blt.TK_B {
			if monsters == MonstersNormal {
				monsters = MonstersEasy
				score -= 25
			} else if monsters == MonstersHard {
				monsters = MonstersNormal
				score -= 25
			} else {
				continue
			}
		} else if key == blt.TK_C && blt.Check(blt.TK_SHIFT) != 0 {
			if reloading == AmmoUnlimited {
				reloading = AmmoLimited
				score += 25
			} else {
				reloading = AmmoUnlimited
				score -= 25
			}
		} else if key == blt.TK_C {
			if reloading == AmmoUnlimited {
				reloading = AmmoLimited
				score += 25
			} else {
				reloading = AmmoUnlimited
				score -= 25
			}
		} else if key == blt.TK_D && blt.Check(blt.TK_SHIFT) != 0 {
			if animations == AnimationsFalse {
				animations = AnimationsTrue
			} else {
				animations = AnimationsFalse
			}
		} else if key == blt.TK_D {
			if animations == AnimationsFalse {
				animations = AnimationsTrue
			} else {
				animations = AnimationsFalse
			}
		}
	}
	err := SaveConfig()
	if err != nil {
		fmt.Println("Error during saving config file!")
		fmt.Println(err)
	}
}
