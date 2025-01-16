package main

import (
	"strconv"
	"strings"
)

/*
This file stores functions and struct related to proposal request number (PRN)
*/

type PRN string

func (base PRN) GreaterThan(in PRN) bool {
	baseCycle, baseActorNumber := base.Parse()
	inCycle, inActorNumber := in.Parse()

	if baseCycle > inCycle {
		return true
	} else if baseCycle < inCycle {
		return false
	} else {
		if baseActorNumber > inActorNumber {
			return true
		} else {
			return false
		}
	}

}

func (p PRN) Parse() (int, int) {

	parts := strings.Split(string(p), ".")
	firstPart, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0
	}

	// Convert the second part to an integer
	secondPart, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0
	}
	// Return the two integers
	return firstPart, secondPart

}

func (p PRN) getCycle() int {
	c, _ := p.Parse()
	return c

}

func (p PRN) getActorNumber() int {
	_, n := p.Parse()
	return n
}

func newPRN(c, n int) PRN {
	cString := strconv.Itoa(c)
	nString := strconv.Itoa(n)
	return PRN(cString + "." + nString)
}
