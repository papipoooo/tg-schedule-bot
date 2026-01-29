package domain

import (
	"time"
)

type Status string

const (
	StatusFree     Status = "free"
	StatusBusy     Status = "busy"
	StatusCanceled Status = "canceled"
)

type User struct {
	ID          int64  // Telegram user id
	Username    string // @ из Telegram
	DisplayName string // имя из Telegram
}

type Slot struct { // временной слот для записи
	ID        string
	StartAt   time.Time // начало
	EndAt     time.Time // конец
	Status    Status    // свободен/занят/отменен
	BookedBy  int       // кто бронировал
	CreatedAt time.Time
	UpdatedAt time.Time
}
