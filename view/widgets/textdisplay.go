package widgets

import (
	"cjavellana.me/launchpad/agent/errors"
	"github.com/mum4k/termdash/widgets/text"
)

type render func(*text.Text)

type RollContentDisplay struct {

	Display *text.Text
}

func NewRollContentDisplay() *RollContentDisplay {
	display, err := text.New(text.RollContent(), text.WrapAtWords())
	errors.CheckFatal(err)
	return &RollContentDisplay{Display: display}
}

func (d *RollContentDisplay) Update(text string) {
	err := d.Display.Write(text)
	errors.CheckFatal(err)
}
