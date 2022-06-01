package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var score int
var enemy []*sprite
var puntos int
var vidas = 1
var base Config
var pacman sprite

type Color int

var level []string

const reset = "\x1b[0m"

const (
	BLACK Color = iota
	RED
	GREEN
	BROWN
	BLUE
	MAGENTA
	CYAN
	GREY
)

func fondo(text string) string {
	return "\x1b[44m" + text + reset
}

type sprite struct {
	row int
	col int
}

var clrs = map[Color]string{
	BROWN:   "\x1b[1;33;43m",
	BLUE:    "\x1b[1;34;44m",
	MAGENTA: "\x1b[1;35;45m",
	CYAN:    "\x1b[1;36;46m",
	GREY:    "\x1b[1;37;47m",
	BLACK:   "\x1b[1;30;40m",
	RED:     "\x1b[1;31;41m",
	GREEN:   "\x1b[1;32;42m",
}

type Config struct {
	Player   string `json:"pacman"`
	Enemy    string `json:"enemy"`
	Wall     string `json:"wall"`
	Dot      string `json:"dot"`
	Lose     string `json:"lose"`
	Space    string `json:"space"`
	UseEmoji bool   `json:"use_emoji"`
}

func fondoN(text string, color Color) string {
	if c, ok := clrs[color]; ok {
		return c + text + reset
	}
	return fondo(text)
}

func dibujar() {
	clean()
	for _, line := range level {
		for _, chr := range line {
			switch chr {
			case '#':
				fmt.Print(fondo(base.Wall))
			case '.':
				fmt.Print(base.Dot)
			default:
				fmt.Print(base.Space)
			}
		}
		fmt.Println()

	}
	cursorMov(pacman.row, pacman.col)
	fmt.Print(base.Player)

	for _, g := range enemy {
		cursorMov(g.row, g.col)
		fmt.Print(base.Enemy)
	}

	cursorMov(len(level)+1, 0)

	fmt.Println("SCORE: ", score, "\tvidas: ", vidas)
}

func iniciar() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin

	badLevel := cbTerm.Run()
	if badLevel != nil {
		log.Fatalln("Unable to activate cbreak mode: ", badLevel)
	}
}
func cursorMov(row, col int) {
	if base.UseEmoji {
		cursor(row, col*2)
	} else {
		cursor(row, col)
	}
}
func clean() {
	fmt.Print("\x1b[2J")
	cursorMov(0, 0)
}
func cargarLevel(file string) error {
	f, badLevel := os.Open(file)
	if badLevel != nil {
		return badLevel
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	badLevel = decoder.Decode(&base)
	if badLevel != nil {
		return badLevel
	}
	return nil
}
func cursor(row, col int) {
	fmt.Printf("\x1b[%d;%df", row+1, col+1)
}
func limpieza() {
	cookedTerm := exec.Command("stty", "-cbreak", "echo")
	cookedTerm.Stdin = os.Stdin

	badLevel := cookedTerm.Run()
	if badLevel != nil {
		log.Fatalln("Unable to restore cooked mode: ", badLevel)
	}
}
func pacmanM(dir string) {
	pacman.row, pacman.col = moving(pacman.row, pacman.col, dir)

	removeDot := func(row, col int) {
		level[row] = level[row][0:col] + " " + level[row][col+1:]
	}

	switch level[pacman.row][pacman.col] {
	case '.':
		puntos--
		score++
		removeDot(pacman.row, pacman.col)
	case 'X':
		score += 10
		removeDot(pacman.row, pacman.col)
	}
}
func moving(oldRow, oldCol int, dir string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol

	switch dir {
	case "UP":
		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(level) - 1
		}
	case "DOWN":
		newRow = newRow + 1
		if newRow == len(level) {
			newRow = 0
		}
	case "RIGHT":
		newCol = newCol + 1
		if newCol == len(level[0]) {
			newCol = 0
		}
	case "LEFT":
		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(level[0]) - 1
		}
	}
	if level[newRow][newCol] == '#' {
		newRow = oldRow
		newCol = oldCol
	}
	return
}
func levelL(file string) error {
	f, badLevel := os.Open(file)
	if badLevel != nil {
		return badLevel
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		level = append(level, line)
	}

	for row, line := range level {
		for col, char := range line {
			switch char {
			case 'P':
				pacman = sprite{row, col}
			case 'G':
				enemy = append(enemy, &sprite{row, col})
			case '.':
				puntos++
			}
		}
	}

	return nil
}
func leerConsola() (string, error) {
	buffer := make([]byte, 100)

	_, badLevel := os.Stdin.Read(buffer)
	if badLevel != nil {
		return "", badLevel
	}

	return string(buffer[0]), nil

}

func fantasmaM() string {
	dir := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}
	return move[dir]
}
func controles() (string, error) {
	buffer := make([]byte, 100)

	cnt, badLevel := os.Stdin.Read(buffer)
	if badLevel != nil {
		return "", badLevel
	}

	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	} else if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			switch buffer[2] {
			case 'A':
				return "UP", nil
			case 'B':
				return "DOWN", nil
			case 'C':
				return "RIGHT", nil
			case 'D':
				return "LEFT", nil
			}

		}
	}

	return "", nil
}

func fantasmaC() {
	for _, g := range enemy {
		dir := fantasmaM()
		g.row, g.col = moving(g.row, g.col, dir)
		time.Sleep(800 * time.Millisecond)
	}

}

var gameover string = "\n" +
	"▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄       ▄▄  ▄▄▄▄▄▄▄▄▄▄▄       ▄▄▄▄▄▄▄▄▄▄▄  ▄               ▄  ▄▄▄▄▄▄▄▄▄▄▄  ▄▄▄▄▄▄▄▄▄▄▄\n" +
	"▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌▐░░▌     ▐░░▌▐░░░░░░░░░░░▌     ▐░░░░░░░░░░░▌▐░▌             ▐░▌▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌\n" +
	"▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌▐░▌░▌   ▐░▐░▌▐░█▀▀▀▀▀▀▀▀▀      ▐░█▀▀▀▀▀▀▀█░▌ ▐░▌           ▐░▌ ▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀▀▀▀█░▌\n" +
	"▐░▌          ▐░▌       ▐░▌▐░▌▐░▌ ▐░▌▐░▌▐░▌               ▐░▌       ▐░▌  ▐░▌         ▐░▌  ▐░▌          ▐░▌       ▐░▌\n" +
	"▐░▌ ▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄█░▌▐░▌ ▐░▐░▌ ▐░▌▐░█▄▄▄▄▄▄▄▄▄      ▐░▌       ▐░▌   ▐░▌       ▐░▌   ▐░█▄▄▄▄▄▄▄▄▄ ▐░█▄▄▄▄▄▄▄█░▌\n" +
	"▐░▌▐░░░░░░░░▌▐░░░░░░░░░░░▌▐░▌  ▐░▌  ▐░▌▐░░░░░░░░░░░▌     ▐░▌       ▐░▌    ▐░▌     ▐░▌    ▐░░░░░░░░░░░▌▐░░░░░░░░░░░▌\n" +
	"▐░▌ ▀▀▀▀▀▀█░▌▐░█▀▀▀▀▀▀▀█░▌▐░▌   ▀   ▐░▌▐░█▀▀▀▀▀▀▀▀▀      ▐░▌       ▐░▌     ▐░▌   ▐░▌     ▐░█▀▀▀▀▀▀▀▀▀ ▐░█▀▀▀▀█░█▀▀ \n" +
	"▐░▌       ▐░▌▐░▌       ▐░▌▐░▌       ▐░▌▐░▌               ▐░▌       ▐░▌      ▐░▌ ▐░▌      ▐░▌          ▐░▌     ▐░▌  \n" +
	"▐░█▄▄▄▄▄▄▄█░▌▐░▌       ▐░▌▐░▌       ▐░▌▐░█▄▄▄▄▄▄▄▄▄      ▐░█▄▄▄▄▄▄▄█░▌       ▐░▐░▌       ▐░█▄▄▄▄▄▄▄▄▄ ▐░▌      ▐░▌ \n" +
	"▐░░░░░░░░░░░▌▐░▌       ▐░▌▐░▌       ▐░▌▐░░░░░░░░░░░▌     ▐░░░░░░░░░░░▌        ▐░▌        ▐░░░░░░░░░░░▌▐░▌       ▐░▌\n" +
	" ▀▀▀▀▀▀▀▀▀▀▀  ▀         ▀  ▀         ▀  ▀▀▀▀▀▀▀▀▀▀▀       ▀▀▀▀▀▀▀▀▀▀▀          ▀          ▀▀▀▀▀▀▀▀▀▀▀  ▀         ▀ \n"

func main() {

	fmt.Println("**WELCOME TO PACMAN**")
	fmt.Println("Please enter the numer of lives ")
	liv, error := leerConsola()
	if error != nil {
		fmt.Println(error)
	} else {
		fmt.Println("You are going to play: " + liv + " vidas")
		vidas, _ = strconv.Atoi(liv)
	}
	iniciar()
	defer limpieza()

	badLevel := levelL("level1.txt")
	if badLevel != nil {
		log.Println("We canot find the level ", badLevel)
		return
	}

	badLevel = cargarLevel("config.json")
	if badLevel != nil {
		log.Println("Failed to load configuration", badLevel)
		return
	}

	entrada := make(chan string)
	go func(ch chan<- string) {
		for {
			entrada, badLevel := controles()
			if badLevel != nil {
				log.Println("error reading entrada:", badLevel)
				ch <- "ESC"
			}
			ch <- entrada
		}
	}(entrada)

	for {

		dibujar()

		select {
		case inp := <-entrada:
			if inp == "ESC" {
				vidas = 0
			}
			pacmanM(inp)
		default:

		}

		go fantasmaC()

		for _, g := range enemy {
			if pacman == *g {
				vidas--
			}
		}

		if puntos == 0 || vidas <= 0 {
			if vidas == 0 {
				cursorMov(pacman.row, pacman.col)
				fmt.Print(base.Lose)
				cursorMov(len(level)+2, 0)
				fmt.Println(gameover)
			}
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
}
