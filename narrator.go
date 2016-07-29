package main

import "fmt"

type Narrator struct {
	Prefix string
}

func NewNarrator(prefix string) Narrator {
	return Narrator{Prefix: prefix}
}

func (n Narrator) Say(a ...interface{}) {
	fmt.Print("[" + n.Prefix + "] ")
	fmt.Println(a...)
}

func (n Narrator) Sayf(format string, a ...interface{}) {
	format = fmt.Sprintf("[%s] %s", n.Prefix, format)
	fmt.Printf(format+"\n", a...)
}
