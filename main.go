package main

import (
	"fmt"
	"time"
)

/*
- periodically list all windows, remember all different titles
- save to file
- at start, restore from file
- if a new window (unknown id) is found, check if the title is known, if yes, move to remembered desktop
*/

func main() {
	ws := NewWindows()
	ws.ReadFromFile(true)
	//fmt.Printf("restored %#v\n", ws)
	// TODO maybe use a "stabilization delay" before saving a window title?
	for {
		ws.Refresh()
		err := ws.SaveToFile()
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(1 * time.Second)
	}
}
