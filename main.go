package main

import (
	"fmt"

	"github.com/k0kubun/pp"
        "github.com/sytallax/pizzatime/dominos"
)

func main() {
	a := dominos.Address{Street: "1600 PENNSYLVANIA AVE NW", City: "WASHINGTON", Region: "DC", PostalCode: 20500}
	s, err := a.GetNearestStore()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pp.Print(s)

	m, err := s.GetMenu()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pp.Print(m)
}
