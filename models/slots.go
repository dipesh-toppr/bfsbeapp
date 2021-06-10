package models

type Slot struct {
	ID            uint
	TeacherId     uint
	Date          string
	AvailableSlot uint
	IsBooked      bool //changed from int
}
