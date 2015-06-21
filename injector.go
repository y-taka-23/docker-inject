package main

import (
	"fmt"
)

type Injector struct {
}

func newInjector() *Injector {
	return &Injector{}
}

func (inj *Injector) run(args []string) {
	for i, a := range args {
		fmt.Printf("%d: %s\n", i, a)
	}
}
