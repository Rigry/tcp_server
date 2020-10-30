package main

import "fmt"
// import "time"
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
		// time.Sleep(1 * time.Second)
			s.HandleRequest()
	}

}