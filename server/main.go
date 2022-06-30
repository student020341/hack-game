package main

import "server/pkg/server"

func main() {
	s := server.MakeServer(nil)
	s.Start()
}
