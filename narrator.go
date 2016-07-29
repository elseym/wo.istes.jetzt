package main

import "fmt"

// Narrator prints a line to stdout with a given prefix
type Narrator struct {
	Prefix string
}

// NewNarrator returns a new Narrator with the given prefix
func NewNarrator(prefix string) Narrator {
	return Narrator{Prefix: prefix}
}

// Say prints a line to stdout
func (n Narrator) Say(a ...interface{}) {
	fmt.Print("[" + n.Prefix + "] ")
	fmt.Println(a...)
}

// Sayf prints a line to stdout honouring the given format
func (n Narrator) Sayf(format string, a ...interface{}) {
	format = fmt.Sprintf("[%s] %s", n.Prefix, format)
	fmt.Printf(format+"\n", a...)
}
