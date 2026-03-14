package main

import (
	"fmt"
	"slices"
	"strconv"
	"sync"
	"time"
)

type ParkingLot struct {
	// multiple parking floors
	slotStore map[string]*Slot
	// availableSlots to store available
	availableSlots map[VehicleType][]*Slot // tbd
	// pricing strategy
	pricingStrategy Pricing
	// ticketStore
	ticketStore map[string]*ParkingTicket

	markSlotAvailableLock sync.RWMutex
	checkInLock sync.Map
}

// will return > 0 if a > b, < 0 otherwise
func compareSlot(a *Slot, b *Slot) int {
	if a.FloorId != b.FloorId {
        return a.FloorId - b.FloorId
    }
    return a.Id - b.Id
	
}

func (pl *ParkingLot) markSlotAvailable(slot *Slot) {
	// lock here 
	pl.markSlotAvailableLock.Lock()
	defer pl.markSlotAvailableLock.Unlock()
	vtype := slot.GetVehicleType()
	slots := pl.availableSlots[vtype]
	e := len(slots)
	// change it with comp slots method
	if compareSlot(slot, slots[e-1]) == 1{
		slots = append(slots, slot)
		pl.availableSlots[vtype] = slots
		return
	} 
	// lowerBound logic : find the largest element smaller that slot.id
	slot.IsOccupied = false
	idx := findInsertIndex(slots, slot)

	slots = append(slots, nil)
	copy(slots[idx+1:], slots[idx:])
	slots[idx] = slot

	pl.availableSlots[vtype] = slots
	
}
func findInsertIndex(slots []*Slot, slot *Slot) int {
    l := 0
    r := len(slots)

    for l < r {
        mid := l + (r-l)/2

        if compareSlot(slots[mid], slot) < 0 {
            l = mid + 1
        } else {
            r = mid
        }
    }

    return l
}
func(pl *ParkingLot) GetAvailableSlots(vtype VehicleType) []*Slot {
	return pl.availableSlots[vtype]
} 

// NewParkingLot Init
// floors -> including ground floor (floor = 3 will mean -> 0,1,2)
func NewParkingLot(floors int, slots int, ps Pricing) *ParkingLot {
	var slotStore = make(map[string]*Slot)
	availableSlots := make(map[VehicleType][]*Slot)
	for i := range(floors) {
		for j := range(slots) {
			slot := &Slot{}
			slotId := strconv.Itoa(i) + "-" + strconv.Itoa(j) 
			slot.Id = j
			slot.FloorId = i
			slotStore[slotId] = slot

			vtype := slot.GetVehicleType()

			if _, exists := availableSlots[vtype]; !exists {
				availableSlots[vtype] = make([]*Slot, 0)
			}

			list := availableSlots[vtype]
			list = append(list, slot)
			availableSlots[vtype] = list
		}
		
	}

	// sort all the lists once
	for i := range(availableSlots) {
		list := availableSlots[i]
		slices.SortFunc(list, func (a *Slot, b *Slot) int {
			return compareSlot(a, b)
		})
		availableSlots[i] = list
	}

	return &ParkingLot{
		slotStore: slotStore,
		availableSlots: availableSlots,
		pricingStrategy: ps,
		ticketStore: make(map[string]*ParkingTicket),
		markSlotAvailableLock: sync.RWMutex{},
	}
}

func (pl *ParkingLot) getLock(vtype VehicleType) *sync.RWMutex {
	val, _ := pl.checkInLock.LoadOrStore(vtype, &sync.RWMutex{})
    return val.(*sync.RWMutex)
}

// will return error if unable to park
func (pl *ParkingLot) CheckIn(vehicle Vehicle) (ParkingTicket, error) {
	// take lock based on vtype
	vtype := vehicle.Type

	lock := pl.getLock(vtype)

	lock.Lock()
	defer lock.Unlock()

	availableSlots := pl.GetAvailableSlots(vtype)

	if len(availableSlots) == 0 {
		return ParkingTicket{}, fmt.Errorf("No Parking Slot available for %v", vtype.ToString()) 
	}

	// get the first slot available

	slot := availableSlots[0]

	// create a parking ticket

	parkingTicket := ParkingTicket{
		Id: vtype.ToString() + "_" + vehicle.RegistrationNumber,
		VehicleParked: &vehicle,
		CheckinTime: time.Now().UnixNano(),
		SlotDetails: slot,
	}

	pl.ticketStore[parkingTicket.Id] = &parkingTicket
	slot.IsOccupied = true

	// remove slot from availability list
	var newList = make([]*Slot, 0)
	if len(availableSlots) > 1 {
		newList = availableSlots[1:]
	}
	
	pl.availableSlots[vtype] = newList

	return parkingTicket, nil
}

func (pl *ParkingLot) Checkout(parkingTicket ParkingTicket) int {

	vehicle := parkingTicket.VehicleParked
	// take lock on vehicleType
	vtype := vehicle.Type

	lock := pl.getLock(vtype)
	lock.Lock()
	defer lock.Unlock()


	price := pl.pricingStrategy.CalculatePrice(parkingTicket)
	// mark slot available
	pl.markSlotAvailable(parkingTicket.SlotDetails)
	return int(price)
}