package main

import (
	"Betelgeuse/shogi"
	"fmt"
)

func main() {
	shogi.Init()
	shogi.IniRand(5489)
	shogi.IniRandomTable()
	fmt.Printf("Hello, World!")
}
