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
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	KB_QWERTY = iota
	KB_QWERTZ
	KB_AZERTY
	KB_Dvorak
)

var KeyMap map[rune]int

var HardcodedKeys = []int{
	blt.TK_RETURN,
	blt.TK_ENTER,
	blt.TK_ESCAPE,
	blt.TK_BACKSPACE,
	blt.TK_TAB,
	blt.TK_SPACE,
	blt.TK_PAUSE,
	blt.TK_INSERT,
	blt.TK_HOME,
	blt.TK_PAGEUP,
	blt.TK_DELETE,
	blt.TK_END,
	blt.TK_PAGEDOWN,
	blt.TK_RIGHT,
	blt.TK_LEFT,
	blt.TK_DOWN,
	blt.TK_UP,
	blt.TK_KP_DIVIDE,
	blt.TK_KP_MULTIPLY,
	blt.TK_KP_MINUS,
	blt.TK_KP_PLUS,
	blt.TK_KP_ENTER,
	blt.TK_KP_1,
	blt.TK_KP_2,
	blt.TK_KP_3,
	blt.TK_KP_4,
	blt.TK_KP_5,
	blt.TK_KP_6,
	blt.TK_KP_7,
	blt.TK_KP_8,
	blt.TK_KP_9,
	blt.TK_KP_0,
	blt.TK_KP_PERIOD,
}

var QWERTYLayoutRunesToCodes = map[rune]int{
	'q':  blt.TK_Q,
	'Q':  blt.TK_Q,
	'w':  blt.TK_W,
	'W':  blt.TK_W,
	'e':  blt.TK_E,
	'E':  blt.TK_E,
	'r':  blt.TK_R,
	'R':  blt.TK_R,
	't':  blt.TK_T,
	'T':  blt.TK_T,
	'y':  blt.TK_Y,
	'Y':  blt.TK_Y,
	'u':  blt.TK_U,
	'U':  blt.TK_U,
	'i':  blt.TK_I,
	'I':  blt.TK_I,
	'o':  blt.TK_O,
	'O':  blt.TK_O,
	'p':  blt.TK_P,
	'P':  blt.TK_P,
	'a':  blt.TK_A,
	'A':  blt.TK_A,
	's':  blt.TK_S,
	'S':  blt.TK_S,
	'd':  blt.TK_D,
	'D':  blt.TK_D,
	'f':  blt.TK_F,
	'F':  blt.TK_F,
	'g':  blt.TK_G,
	'G':  blt.TK_G,
	'h':  blt.TK_H,
	'H':  blt.TK_H,
	'j':  blt.TK_J,
	'J':  blt.TK_J,
	'k':  blt.TK_K,
	'K':  blt.TK_K,
	'l':  blt.TK_L,
	'L':  blt.TK_L,
	'z':  blt.TK_Z,
	'Z':  blt.TK_Z,
	'x':  blt.TK_X,
	'X':  blt.TK_X,
	'c':  blt.TK_C,
	'C':  blt.TK_C,
	'v':  blt.TK_V,
	'V':  blt.TK_V,
	'b':  blt.TK_B,
	'B':  blt.TK_B,
	'n':  blt.TK_N,
	'N':  blt.TK_N,
	'm':  blt.TK_M,
	'M':  blt.TK_M,
	',':  blt.TK_COMMA,
	'<':  blt.TK_COMMA,
	'.':  blt.TK_PERIOD,
	'>':  blt.TK_PERIOD,
	'/':  blt.TK_SLASH,
	'?':  blt.TK_SLASH,
	';':  blt.TK_SEMICOLON,
	':':  blt.TK_SEMICOLON,
	'\'': blt.TK_APOSTROPHE,
	'"':  blt.TK_APOSTROPHE,
	'[':  blt.TK_LBRACKET,
	'{':  blt.TK_LBRACKET,
	']':  blt.TK_RBRACKET,
	'}':  blt.TK_RBRACKET,
	'1':  blt.TK_1,
	'!':  blt.TK_1,
	'2':  blt.TK_2,
	'@':  blt.TK_2,
	'3':  blt.TK_3,
	'#':  blt.TK_3,
	'4':  blt.TK_4,
	'$':  blt.TK_4,
	'5':  blt.TK_5,
	'%':  blt.TK_5,
	'6':  blt.TK_6,
	'^':  blt.TK_6,
	'7':  blt.TK_7,
	'&':  blt.TK_7,
	'8':  blt.TK_8,
	'*':  blt.TK_8,
	'9':  blt.TK_9,
	'(':  blt.TK_9,
	'0':  blt.TK_0,
	')':  blt.TK_0,
	'-':  blt.TK_MINUS,
	'_':  blt.TK_MINUS,
	'=':  blt.TK_EQUALS,
	'+':  blt.TK_EQUALS,
}

var DvorakLayoutRunesToCodes map[rune]int

var QWERTZLayoutRunesToCodes = map[rune]int{}

var AZERTYLayoutRunesToCodes = map[rune]int{}

func InitializeKeyboardLayouts() {
	InitializeQWERTZ()
	InitializeAZERTY()
	InitializeDvorak()
}

func ChooseKeyboardLayout() {
	switch KeyboardLayout {
	case KB_QWERTY:
		KeyMap = QWERTYLayoutRunesToCodes
	case KB_QWERTZ:
		KeyMap = QWERTZLayoutRunesToCodes
	case KB_AZERTY:
		KeyMap = AZERTYLayoutRunesToCodes
	case KB_Dvorak:
		KeyMap = DvorakLayoutRunesToCodes
	}
}

func InitializeQWERTZ() {
	for k, v := range QWERTYLayoutRunesToCodes {
		QWERTZLayoutRunesToCodes[k] = v
	}
	QWERTZLayoutRunesToCodes['z'] = blt.TK_Y
	QWERTZLayoutRunesToCodes['Z'] = blt.TK_Y
	QWERTZLayoutRunesToCodes['y'] = blt.TK_Z
	QWERTZLayoutRunesToCodes['Y'] = blt.TK_Z
	QWERTZLayoutRunesToCodes[';'] = blt.TK_COMMA
	QWERTZLayoutRunesToCodes[':'] = blt.TK_PERIOD
	QWERTZLayoutRunesToCodes['-'] = blt.TK_SLASH
	QWERTZLayoutRunesToCodes['_'] = blt.TK_SLASH
	QWERTZLayoutRunesToCodes['ö'] = blt.TK_SEMICOLON
	QWERTZLayoutRunesToCodes['Ö'] = blt.TK_SEMICOLON
	QWERTZLayoutRunesToCodes['ä'] = blt.TK_APOSTROPHE
	QWERTZLayoutRunesToCodes['Ä'] = blt.TK_APOSTROPHE
	QWERTZLayoutRunesToCodes['ü'] = blt.TK_LBRACKET
	QWERTZLayoutRunesToCodes['Ü'] = blt.TK_LBRACKET
	QWERTZLayoutRunesToCodes['+'] = blt.TK_RBRACKET
	QWERTZLayoutRunesToCodes['*'] = blt.TK_RBRACKET
	QWERTZLayoutRunesToCodes['"'] = blt.TK_2
	QWERTZLayoutRunesToCodes['§'] = blt.TK_3
	QWERTZLayoutRunesToCodes['&'] = blt.TK_6
	QWERTZLayoutRunesToCodes['/'] = blt.TK_7
	QWERTZLayoutRunesToCodes['('] = blt.TK_8
	QWERTZLayoutRunesToCodes[')'] = blt.TK_9
	QWERTZLayoutRunesToCodes['='] = blt.TK_0
	QWERTZLayoutRunesToCodes['ß'] = blt.TK_MINUS
	QWERTZLayoutRunesToCodes['?'] = blt.TK_MINUS
	QWERTZLayoutRunesToCodes['´'] = blt.TK_EQUALS
	QWERTZLayoutRunesToCodes['`'] = blt.TK_EQUALS
}

func InitializeAZERTY() {
	for k, v := range QWERTYLayoutRunesToCodes {
		AZERTYLayoutRunesToCodes[k] = v
	}
	AZERTYLayoutRunesToCodes['a'] = blt.TK_Q
	AZERTYLayoutRunesToCodes['A'] = blt.TK_Q
	AZERTYLayoutRunesToCodes['z'] = blt.TK_W
	AZERTYLayoutRunesToCodes['Z'] = blt.TK_W
	AZERTYLayoutRunesToCodes['q'] = blt.TK_A
	AZERTYLayoutRunesToCodes['Q'] = blt.TK_A
	AZERTYLayoutRunesToCodes['w'] = blt.TK_Z
	AZERTYLayoutRunesToCodes['W'] = blt.TK_Z
	AZERTYLayoutRunesToCodes[','] = blt.TK_M
	AZERTYLayoutRunesToCodes['?'] = blt.TK_M
	AZERTYLayoutRunesToCodes[';'] = blt.TK_COMMA
	AZERTYLayoutRunesToCodes['.'] = blt.TK_COMMA
	AZERTYLayoutRunesToCodes[':'] = blt.TK_PERIOD
	AZERTYLayoutRunesToCodes['/'] = blt.TK_PERIOD
	AZERTYLayoutRunesToCodes['!'] = blt.TK_SLASH
	AZERTYLayoutRunesToCodes['§'] = blt.TK_SLASH
	AZERTYLayoutRunesToCodes['m'] = blt.TK_SEMICOLON
	AZERTYLayoutRunesToCodes['M'] = blt.TK_SEMICOLON
	AZERTYLayoutRunesToCodes['ù'] = blt.TK_APOSTROPHE
	AZERTYLayoutRunesToCodes['%'] = blt.TK_APOSTROPHE
	AZERTYLayoutRunesToCodes['^'] = blt.TK_LBRACKET
	AZERTYLayoutRunesToCodes['¨'] = blt.TK_LBRACKET
	AZERTYLayoutRunesToCodes['$'] = blt.TK_RBRACKET
	AZERTYLayoutRunesToCodes['£'] = blt.TK_RBRACKET
	AZERTYLayoutRunesToCodes['&'] = blt.TK_1
	AZERTYLayoutRunesToCodes['é'] = blt.TK_2
	AZERTYLayoutRunesToCodes['"'] = blt.TK_3
	AZERTYLayoutRunesToCodes['\''] = blt.TK_4
	AZERTYLayoutRunesToCodes['('] = blt.TK_5
	AZERTYLayoutRunesToCodes['-'] = blt.TK_6
	AZERTYLayoutRunesToCodes['è'] = blt.TK_7
	AZERTYLayoutRunesToCodes['_'] = blt.TK_8
	AZERTYLayoutRunesToCodes['ç'] = blt.TK_9
	AZERTYLayoutRunesToCodes['à'] = blt.TK_0
	AZERTYLayoutRunesToCodes[')'] = blt.TK_MINUS
	AZERTYLayoutRunesToCodes['°'] = blt.TK_MINUS
}

func InitializeDvorak() {
	DvorakLayoutRunesToCodes = map[rune]int{
		'\'': blt.TK_Q,
		'"':  blt.TK_Q,
		',':  blt.TK_W,
		'<':  blt.TK_W,
		'.':  blt.TK_E,
		'>':  blt.TK_E,
		'p':  blt.TK_R,
		'P':  blt.TK_R,
		'y':  blt.TK_T,
		'Y':  blt.TK_T,
		'f':  blt.TK_Y,
		'F':  blt.TK_Y,
		'g':  blt.TK_U,
		'G':  blt.TK_U,
		'c':  blt.TK_I,
		'C':  blt.TK_I,
		'r':  blt.TK_O,
		'R':  blt.TK_O,
		'l':  blt.TK_P,
		'L':  blt.TK_P,
		'a':  blt.TK_A,
		'A':  blt.TK_A,
		'o':  blt.TK_S,
		'O':  blt.TK_S,
		'e':  blt.TK_D,
		'E':  blt.TK_D,
		'u':  blt.TK_F,
		'U':  blt.TK_F,
		'i':  blt.TK_G,
		'I':  blt.TK_G,
		'd':  blt.TK_H,
		'D':  blt.TK_H,
		'h':  blt.TK_J,
		'H':  blt.TK_J,
		't':  blt.TK_K,
		'T':  blt.TK_K,
		'n':  blt.TK_L,
		'N':  blt.TK_L,
		';':  blt.TK_Z,
		':':  blt.TK_Z,
		'q':  blt.TK_X,
		'Q':  blt.TK_X,
		'j':  blt.TK_C,
		'J':  blt.TK_C,
		'k':  blt.TK_V,
		'K':  blt.TK_V,
		'x':  blt.TK_B,
		'X':  blt.TK_B,
		'b':  blt.TK_N,
		'B':  blt.TK_N,
		'm':  blt.TK_M,
		'M':  blt.TK_M,
		'w':  blt.TK_COMMA,
		'W':  blt.TK_COMMA,
		'v':  blt.TK_PERIOD,
		'V':  blt.TK_PERIOD,
		'z':  blt.TK_SLASH,
		'Z':  blt.TK_SLASH,
		's':  blt.TK_SEMICOLON,
		'S':  blt.TK_SEMICOLON,
		'-':  blt.TK_APOSTROPHE,
		'_':  blt.TK_APOSTROPHE,
		'/':  blt.TK_LBRACKET,
		'?':  blt.TK_LBRACKET,
		'=':  blt.TK_RBRACKET,
		'+':  blt.TK_RBRACKET,
		'1':  blt.TK_1,
		'!':  blt.TK_1,
		'2':  blt.TK_2,
		'@':  blt.TK_2,
		'3':  blt.TK_3,
		'#':  blt.TK_3,
		'4':  blt.TK_4,
		'$':  blt.TK_4,
		'5':  blt.TK_5,
		'%':  blt.TK_5,
		'6':  blt.TK_6,
		'^':  blt.TK_6,
		'7':  blt.TK_7,
		'&':  blt.TK_7,
		'8':  blt.TK_8,
		'*':  blt.TK_8,
		'9':  blt.TK_9,
		'(':  blt.TK_9,
		'0':  blt.TK_0,
		')':  blt.TK_0,
		'[':  blt.TK_MINUS,
		'{':  blt.TK_MINUS,
		']':  blt.TK_EQUALS,
		'}':  blt.TK_EQUALS,
	}
}

func ReadOptionsControls() {
	f, err := os.Open("options_controls.cfg")
	if err != nil {
		panic("Can't find options_controls.cfg file!")
	}
	defer f.Close()
	var opts= []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var lines= []string{}
		line := scanner.Text()
		if utf8.RuneCountInString(line) > 0 && []rune(line)[0] != '#' {
			line = strings.Replace(line, "\r", "\n", -1)
			lines = strings.Split(line, "\n")
			for i := 0; i < len(lines); i++ {
				opts = append(opts, strings.ToUpper(strings.TrimSpace(lines[i])))
			}
		}
	}
	for _, v := range opts {
		var results= strings.Split(v, "=")
		resKey := strings.TrimSpace(results[0])
		if resKey == "KB_LAYOUT" {
			val := strings.TrimSpace(results[1])
			switch val {
			case "QWERTY":
				KeyboardLayout = KB_QWERTY
			case "QWERTZ":
				KeyboardLayout = KB_QWERTZ
			case "AZERTY":
				KeyboardLayout = KB_AZERTY
			case "DVORAK":
				KeyboardLayout = KB_Dvorak
			default:
				fmt.Println("Wrong value in KB_LAYOUT; using QWERTY.")
				KeyboardLayout = KB_QWERTY
			}
		} else if resKey == "CUSTOM_CONTROLS" {
			val := strings.TrimSpace(results[1])
			if val == "TRUE" {
				CustomControls = true
			} else if val == "FALSE" {
				CustomControls = false
			} else {
				fmt.Println("Wrong value is CUSTOM_CONTROLS; using FALSE.")
				CustomControls = false
			}
		}
	}
	for _, v := range opts {
		var results= strings.Split(v, "=")
		resKey := strings.TrimSpace(results[0])
		resValue := strings.TrimSpace(results[1])
		if utf8.RuneCountInString(resKey) > 0 && []rune(resKey)[0] != '#' &&
			resKey != "KB_LAYOUT" && resKey != "CUSTOM_CONTROLS" {
			addKeyToCustomLayout(resKey, resValue)
		}
	}
}

func addKeyToCustomLayout(resKey string, resValue string) {
	var tempMap map[rune]int
	switch KeyboardLayout {
	case KB_QWERTY: tempMap = QWERTYLayoutRunesToCodes
	case KB_QWERTZ: tempMap = QWERTZLayoutRunesToCodes
	case KB_AZERTY: tempMap = AZERTYLayoutRunesToCodes
	case KB_Dvorak: tempMap = DvorakLayoutRunesToCodes
	}
	var s string
	valid := false
	for _, v := range Actions {
		if resKey == v {
			valid = true
		}
	}
	if valid == true {
		s = resKey
	} else {
		panic("Wrong value: " + resKey)
	}
	var i int
	switch resValue {
	case "RETURN": i = blt.TK_RETURN
	case "ENTER": i = blt.TK_ENTER
	case "TAB": i = blt.TK_TAB
	case "SPACE": i = blt.TK_SPACE
	case "PAUSE": i = blt.TK_PAUSE
	case "INSERT": i = blt.TK_INSERT
	case "HOME": i = blt.TK_HOME
	case "PAGEUP": i = blt.TK_PAGEUP
	case "DELETE": i = blt.TK_DELETE
	case "END": i = blt.TK_END
	case "PAGEDOWN": i = blt.TK_PAGEDOWN
	case "RIGHT": i = blt.TK_RIGHT
	case "LEFT": i = blt.TK_LEFT
	case "DOWN": i = blt.TK_DOWN
	case "UP": i = blt.TK_UP
	case "KP_DIVIDE": i = blt.TK_KP_DIVIDE
	case "KP_MULTIPLY": i = blt.TK_KP_MULTIPLY
	case "KP_MINUS": i = blt.TK_KP_MINUS
	case "KP_PLUS": i = blt.TK_KP_PLUS
	case "KP_ENTER": i = blt.TK_KP_ENTER
	case "KP_1": i = blt.TK_KP_1
	case "KP_2": i = blt.TK_KP_2
	case "KP_3": i = blt.TK_KP_3
	case "KP_4": i = blt.TK_KP_4
	case "KP_5": i = blt.TK_KP_5
	case "KP_6": i = blt.TK_KP_6
	case "KP_7": i = blt.TK_KP_7
	case "KP_8": i = blt.TK_KP_8
	case "KP_9": i = blt.TK_KP_9
	case "KP_0": i = blt.TK_KP_0
	case "KP_PERIOD": i = blt.TK_KP_PERIOD
	default:
		if utf8.RuneCountInString(resValue) == 1 {
			i = tempMap[[]rune(resValue)[0]]
		} else {
			panic("Wrong value: " + resValue)
		}
	}
	CustomCommandKeys[i] = s
}
