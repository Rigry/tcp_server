package main

import "fmt"
import "tcp_server/server"

func main() {

	fmt.Println("Launching server...")

	s := server.Make()

	err := s.Listen()
	if err != nil {
		fmt.Println("failed listeninng", err)
	}
	defer s.Close()

	for {
		go s.HandleRequest()
	}

}