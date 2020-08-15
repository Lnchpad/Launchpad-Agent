package widgets

import (
	"github.com/mum4k/termdash/widgets/text"
	"log"
)

type render func(*text.Text)

type RollContentDisplay struct {
	Display *text.Text
}

func NewRollContentDisplay() *RollContentDisplay {
	display, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		log.Fatal(err)
	}

	return &RollContentDisplay{Display: display}
}

func (d *RollContentDisplay) Update(text string) {
	err := d.Display.Write(text)
	if err != nil {
		log.Fatal(err)
	}
}
