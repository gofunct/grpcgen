package main

import (
	"github.com/gofunct/gen/cmd"
	"github.com/gofunct/viperize"
	"sync"
)

func main() {
	x := sync.WaitGroup{}
	x.Add(1)
	go viperize.ViperizeCli()
	x.Done()
	cmd.Execute()
}
