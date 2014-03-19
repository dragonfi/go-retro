package main

import (
	"flag"
)

func main() {
	var player_number int
	flag.IntVar(&player_number, "p", 1, "The number of players. (1-4)")
	flag.Parse()
	x, y := Init()
	defer Close()
	offsetx, offsety := 2, 2
	aw := NewArenaWidget(offsetx, offsety, x-2*offsetx, y-2*offsety, player_number)

	aw.Run()
}
