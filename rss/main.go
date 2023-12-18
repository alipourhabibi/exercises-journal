package main

import (
	"github.com/alipourhabibi/exercises-journal/rss/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
