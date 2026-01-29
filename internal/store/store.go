package store

import (
	"time"

	"tg-schedule-bot/internal/domain"
)

type Store interface {
	Load() error                                                  // Загружает сохранённые слоты в память при запуске приложения.
	Save() error                                                  // Сохраняет текущее состояние слотов из памяти во внешний файл
	CreateSlot(startAt, endAt time.Time) (domain.Slot, error)     // Создаёт новый временной слот и делает его доступным для записи.
	CancelSlot(slotID string) (domain.Slot, error)                // Помечает слот как отменённый — в него больше нельзя записаться.
	GetSlot(slotID string) (domain.Slot, bool)                    // Возвращает слот по идентификатору, если он существует.
	ListSlots() []domain.Slot                                     // Возвращает список всех слотов в системе.
	ListFreeSlots() []domain.Slot                                 // Возвращает только те слоты, которые сейчас свободны.
	BookSlot(slotID string, userID int) (domain.Slot, error)      // Записывает пользователя в свободный слот.
	CancelBooking(slotID string, userID int) (domain.Slot, error) // Отменяет запись пользователя и освобождает слот.
	ListSlotsByUser(userID int) []domain.Slot                     // Возвращает все слоты, на которые записан конкретный пользователь.
}
