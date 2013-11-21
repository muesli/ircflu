package app

import (
	"flag"
)

type CliFlag struct {
	V     interface{}
	Name  string
	Value string
	Desc  string
}

var (
	appflags []CliFlag
)

func AddFlags(flags []CliFlag) {
	for _, flag := range flags {
		appflags = append(appflags, flag)
	}
}

func Run() {
	for _, f := range appflags {
		flag.StringVar((f.V).(*string), f.Name, f.Value, f.Desc)
	}

	flag.Parse()
}
