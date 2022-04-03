package main

import (
	"fmt"
	"image"
)

type DeckDrawer struct {
	cards        []*Card
	backSideName string
	columns      int
	rows         int
	deckFileName string

	images        []image.Image
	imageBackSide image.Image
}

func NewDeckDrawer(deck *Deck) *DeckDrawer {
	return &DeckDrawer{
		cards:        deck.Cards,
		backSideName: deck.GetBackSideName(),
		columns:      deck.Columns,
		rows:         deck.Rows,
		deckFileName: deck.FileName,
	}
}

func (d *DeckDrawer) log(logType string, progress, total int) {
	if total > 0 {
		fmt.Printf("\r[ %s ] %s %d / %d", logType, d.deckFileName, progress, total)
	} else {
		fmt.Printf("\r[ %s ] %s          ", logType, d.deckFileName)
	}
	if logType == "DONE" {
		fmt.Println()
	}
}
func (d *DeckDrawer) loadCards() {
	for i, card := range d.cards {
		d.images = append(d.images, OpenImage(card.GetFilePath()))
		d.log("LOAD", i+1, len(d.cards))
	}
	d.imageBackSide = OpenImage(GetConfig().CachePath + d.backSideName)
	return
}
func (d *DeckDrawer) collectDeckImage() *Image {
	bound := d.images[0].Bounds().Max
	deckImage := CreateImage(bound.X, bound.Y, d.columns, d.rows)
	// Draw front side images
	for row := 0; row < d.rows; row++ {
		for col := 0; col < d.columns; col++ {
			if len(d.images) <= (row*d.columns + col) {
				continue
			}
			img := d.images[row*d.columns+col]
			deckImage.Draw(col, row, img)
			d.log("DRAW", row*d.columns+col+1, len(d.images))
		}
	}
	// On bottom right place draw back side image
	deckImage.Draw(d.columns-1, d.rows-1, d.imageBackSide)
	return deckImage
}
func (d *DeckDrawer) Draw() {
	d.loadCards()
	deckImage := d.collectDeckImage()
	d.log("SAVE", 0, 0)
	deckImage.SaveImage(GetConfig().ResultDir + d.deckFileName)
	d.log("DONE", 0, 0)
}
