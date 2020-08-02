package view

import (
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/terminal/termbox"
)

type Dashboard interface {
	Build(terminal *termbox.Terminal) *container.Container
}
