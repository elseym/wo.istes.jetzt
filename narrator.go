package main

import "fmt"

// Narrator prints a line to stdout with a given prefix
type Narrator struct {
	Prefix string
	Quiet  bool
}

// NewNarrator returns a new Narrator with the given prefix
func NewNarrator(prefix string, quiet bool) Narrator {
	return Narrator{Prefix: prefix, Quiet: quiet}
}

// Say prints a line to stdout
func (n Narrator) Say(a ...interface{}) {
	n.printIt(fmt.Sprint(a...))
}

// Sayf prints a line to stdout honouring the given format
func (n Narrator) Sayf(format string, a ...interface{}) {
	n.printIt(fmt.Sprintf(format, a...))
}

func (n Narrator) printIt(s string) {
	if !n.Quiet {
		fmt.Println("["+n.Prefix+"]", s)
	}
}
