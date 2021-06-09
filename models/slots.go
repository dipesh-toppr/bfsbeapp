package models

type Slot struct {
	ID        uint
	TeacherId uint
	Date      string //added
	Time      string //added and availableSlot deleted
	IsBooked  bool   //changed from int
}
