package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func CaretLocation(te *walk.TextEdit) int {
	idx := int(te.SendMessage(win.EM_LINEFROMCHAR, ^uintptr(0), 0))
	// idx := int(te.SendMessage(win.EM_LINEFROMCHAR, 0xffffffff, 0))
	return idx
}

func main() {
	var inTE *walk.TextEdit
	var outLabel *walk.Label

	MainWindow{
		Title:   "Linijka",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			VSplitter{
				Children: []Widget{
					Label{AssignTo: &outLabel, Text: "Hello"},
					TextEdit{AssignTo: &inTE, OnKeyPress: func(key walk.Key) { log.Printf("Key %v, caret:%d", key, CaretLocation(inTE)) }},
				},
			},
		},
	}.Run()
}
