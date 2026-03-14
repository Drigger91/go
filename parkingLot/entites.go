package main

type VehicleType int
const (
	Bike VehicleType = iota 
	Car 
	Truck 
)
func (v VehicleType) BasePrice() int {
    switch v {
    case Bike:
        return 10
    case Car:
        return 20
    case Truck:
        return 30
    default:
        return 0
    }
}
func (v VehicleType) ToString() string {
    switch v {
    case Bike:
        return "Bike"
    case Car:
        return "Car"
    case Truck:
        return "Truck"
    default:
        return ""
    }
}


type Vehicle struct {
	RegistrationNumber string
	Type VehicleType
}

type Slot struct {
	IsOccupied bool
	Id int
	FloorId int
}
func (s *Slot) GetVehicleType() VehicleType {
    if s.Id == 1 {
		return Truck
	} 
	if s.Id > 1 && s.Id <= 4 {
		return Bike
	}
	return Car
}

type ParkingFloor struct {
	Slots []Slot
	FloorNum int
}

type ParkingTicket struct {
	Id string
	VehicleParked *Vehicle
	CheckinTime int64
	SlotDetails *Slot
}

