# 🚗 Parking Lot Management System

## 📌 Problem Statement

Design and implement a **Parking Lot Management System**.

The parking lot can have multiple floors, and each floor contains multiple parking slots for different vehicle types. The system should allow vehicles to **park**, **unpark**, and **query available slots** efficiently.

Your solution should be **extensible, thread-safe, and well-structured** to meet SDE2 expectations.

---

## 🛠️ Functional Requirements

### 1. Create Parking Lot
Initialize a parking lot with:
- `number_of_floors`
- `slots_per_floor`

Each slot supports a specific **vehicle type** (e.g., Bike, Car, Truck). You may assume slots are pre-classified by type.

**Example Configuration (Floor 1):**
- Slot 1 → Bike
- Slot 2 → Bike
- Slot 3 → Car
- Slot 4 → Car
- Slot 5 → Truck

### 2. Park Vehicle
**Method:** `Park(vehicle)`

The system should:
- Assign the **nearest available slot** on the lowest possible floor.
- Generate and return a **parking ticket**.

**Ticket Details:**
The ticket should contain `ticketId`, `floorNumber`, `slotNumber`, `vehicleNumber`, `vehicleType`, and `entryTime`.
- Format: `<parkingLotId>_<floor>_<slot>` (e.g., `PR123_2_5`)

### 3. Unpark Vehicle
**Method:** `Unpark(ticketId)`

The system should:
- Free the occupied slot.
- Calculate the parking duration.
- Return the parking fee based on the following pricing:
  - **Bike:** ₹10/hour
  - **Car:** ₹20/hour
  - **Truck:** ₹30/hour

### 4. Display Free Slots
Support queries like:
- `DisplayFreeSlots(vehicleType)`
- `DisplayOccupiedSlots(vehicleType)`
- `DisplayFreeSlotCount(vehicleType)`

**Example Output:**
```text
Free slots for CAR:
Floor 1 : 2, 5
Floor 2 : 3