package main

import (
	"time"
)

type Pricing interface {
	CalculatePrice(ticket ParkingTicket) int
}


type NormalPricing struct {}

func (ps *NormalPricing) CalculatePrice(ticket ParkingTicket) int {
	totalTimeInNano := time.Now().UnixNano() - ticket.CheckinTime

	timeDur := time.Duration(totalTimeInNano)
	seconds := timeDur.Seconds()

	return int(seconds) * (ticket.VehicleParked.Type).BasePrice()
}