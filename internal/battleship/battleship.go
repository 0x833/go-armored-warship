package battleship

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

type Ship interface {
	IsSunk() bool
}

type BattleShip struct {
	Position Coordinate
	Requires int
}

type Coordinate struct {
	X int
	Y int
	value string
}

func (c *Coordinate) GetValue() string {
	return c.value
}

type Board struct {
	Grid [][]Coordinate
}

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok { //if we defined a clear func for that platform:
		value()  //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}


func title() string {
	return fmt.Sprintf(`____       _______ _______ _      ______  _____ _    _ _____ _____
|  _ \   /\|__   __|__   __| |    |  ____|/ ____| |  | |_   _|  __ \
| |_) | /  \  | |     | |  | |    | |__  | (___ | |__| | | | | |__) |
|  _ < / /\ \ | |     | |  | |    |  __|  \___ \|  __  | | | |  ___/
| |_) / ____ \| |     | |  | |____| |____ ____) | |  | |_| |_| |
|____/_/    \_\_|     |_|  |______|______|_____/|_|  |_|_____|_|     `)
}

func Prompt() int {
	selection := -1
	for ; selection < 0 && selection > 3; {
		CallClear()
		fmt.Print(title())

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter a selection: ")
		text, _ := reader.ReadString('\n')

		var err error
		selection, err = strconv.Atoi(text)
		if err != nil {
			selection = -1
		}
	}

	return selection
}


