package view

import "fmt"

type SimpleStdoutPrinter struct {
}

func (p *SimpleStdoutPrinter) Update(text string) {
	fmt.Printf("%s", text)
}
