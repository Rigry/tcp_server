package main

import "fmt"
import "time"
import "github.com/tcp_server/server"

func main() {

	fmt.Println("Launching server...")

	s := server.Make()
	s.SetLifeTime(10) // 10 seconds

	err := s.Listen()
	if err != nil {
		fmt.Println("failed listeninng", err)
	}
	defer s.Close()

	for {
		time.Sleep(1 * time.Second)
	}

}