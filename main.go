package main

import (
	"google.golang.org/appengine"
)

func main() {
	s := Server{}
	s.routes()
	appengine.Main()
}
