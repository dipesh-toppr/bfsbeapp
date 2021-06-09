package models

type Slot struct {
	ID            uint
	TeacherId     uint
	AvailableSlot uint
	IsBooked      bool //changed from int
}
