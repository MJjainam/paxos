package main

import "fmt"

const TOTAL_ACTORS int = 5

func main() {
	fmt.Println("Good to have you back, Go!")

	network := newNetwork(5)

	/* Initialize actors */
	var Actors [TOTAL_ACTORS + 1]Actor
	for i := 1; i <= TOTAL_ACTORS; i++ {
		Actors[i] = Actor{
			N:    i,
			Network: network,
	}

	for i := 1; i <= TOTAL_ACTORS; i++ {
		for j := 1; j <= TOTAL_ACTORS; j++ {
			ch := make(chan Message)
			Actors[i].To[j] = ch
			Actors[j].From[i] = ch
		}
	}

}
