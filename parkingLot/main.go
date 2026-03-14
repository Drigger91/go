package main

import (
	"fmt"
	"time"
)


func main() {
	fmt.Println("Parking lot")
	var pricing NormalPricing
	pl := NewParkingLot(2, 8, &pricing)

	ptck, err := pl.CheckIn(Vehicle{
		RegistrationNumber: "1",
		Type: Car,
	})

	fmt.Println(ptck, err)

	time.Sleep(3 * time.Second)

	price := pl.Checkout(ptck)

	fmt.Println("Price", price)
}