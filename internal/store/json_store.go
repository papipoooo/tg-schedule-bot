package store

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"tg-schedule-bot/internal/domain"
	"time"
)

type JsonStore struct {
	slots    map[string]domain.Slot
	mu       sync.Mutex
	dataPath string
}

func NewStore() Store {
	return &JsonStore{
		slots: make(map[string]domain.Slot),
	}
}

func (j *JsonStore) Load() error {
	if j.slots == nil {
		j.slots = make(map[string]domain.Slot)
	}

	if j.dataPath == "" {
		return nil
	}

	file, err := os.ReadFile(j.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read file: %w", err)
	}

	var slotsList []domain.Slot
	errParsing := json.Unmarshal(file, &slotsList)
	if errParsing != nil {
		return fmt.Errorf("failed to parse json: %w", err)
	}

	for _, slot := range slotsList {
		if slot.ID == "" {
			return fmt.Errorf("invalid data: slot has empty ID")
		}
		j.slots[slot.ID] = slot
	}

	return nil
}

func (j *JsonStore) Save() error {
	var allSlots = []domain.Slot{}

	if j.dataPath == "" {
		return nil
	}

	for _, item := range j.slots {
		allSlots = append(allSlots, item)

	}

	jsonText, err := json.Marshal(allSlots)
	if err != nil {
		return fmt.Errorf("could not be translated into json text %w", err)
	}

	writeErr := os.WriteFile(j.dataPath, jsonText, 0644)
	if writeErr != nil {
		return fmt.Errorf("could not write the file %w", writeErr)
	}

	return nil
}

func (j *JsonStore) CreateSlot(startAt, endAt time.Time) (domain.Slot, error) {
	j.mu.Lock()
	defer j.mu.Unlock()

	if startAt.Equal(endAt) || startAt.After(endAt) {
		return domain.Slot{}, fmt.Errorf("invalid time range: startAt must be before endAt")
	}

	for _, existing := range j.slots {
		if existing.Status == domain.StatusCanceled {
			continue
		}

		if startAt.Before(existing.EndAt) && endAt.After(existing.StartAt) {
			return domain.Slot{}, fmt.Errorf("time slot conflicts with existing slot")
		}

	}

	var id string

	for {
		id = strconv.Itoa(rand.Int())
		if _, exists := j.slots[id]; !exists {
			break
		}
	}

	newSlot := domain.Slot{
		ID:        id,
		StartAt:   startAt,
		EndAt:     endAt,
		Status:    domain.StatusFree,
		BookedBy:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	j.slots[newSlot.ID] = newSlot

	return newSlot, nil
}

func (j *JsonStore) CancelSlot(slotID string) (domain.Slot, error) {

	j.mu.Lock()
	defer j.mu.Unlock()

	if _, ok := j.slots[slotID]; !ok {
		return domain.Slot{}, fmt.Errorf("the specified slot does not exist")
	}

	slot := j.slots[slotID]

	switch slot.Status {
	case domain.StatusCanceled:
		return domain.Slot{}, fmt.Errorf("slot already canceled")
	case domain.StatusBusy, domain.StatusFree:
		slot.Status = domain.StatusCanceled
		slot.UpdatedAt = time.Now()
		j.slots[slotID] = slot
		return slot, nil
	default:
		return domain.Slot{}, fmt.Errorf("invalid status")
	}
}

func (j *JsonStore) GetSlot(slotID string) (domain.Slot, bool) {

	j.mu.Lock()
	defer j.mu.Unlock()

	if _, ok := j.slots[slotID]; !ok {
		return domain.Slot{}, false
	}
	return j.slots[slotID], true
}

func (j *JsonStore) ListSlots() []domain.Slot {

	j.mu.Lock()
	defer j.mu.Unlock()

	var listOfSlots = []domain.Slot{}

	for _, item := range j.slots {
		listOfSlots = append(listOfSlots, item)
	}
	return listOfSlots
}

func (j *JsonStore) ListFreeSlots() []domain.Slot {

	j.mu.Lock()
	defer j.mu.Unlock()

	var listOfFreeSlots = []domain.Slot{}

	for _, item := range j.slots {
		if item.Status == domain.StatusFree {
			listOfFreeSlots = append(listOfFreeSlots, item)
		}
	}

	return listOfFreeSlots
}

func (j *JsonStore) BookSlot(slotID string, userID int) (domain.Slot, error) {

	j.mu.Lock()
	defer j.mu.Unlock()

	if _, ok := j.slots[slotID]; !ok {
		return domain.Slot{}, fmt.Errorf("slot not found")
	}

	slot := j.slots[slotID]

	switch slot.Status {
	case domain.StatusBusy:
		return domain.Slot{}, fmt.Errorf("slot already booked")
	case domain.StatusCanceled:
		return domain.Slot{}, fmt.Errorf("slot canceled")
	case domain.StatusFree:
		if slot.BookedBy != 0 {
			return domain.Slot{}, fmt.Errorf("invalid slot state")
		}
		slot.Status = domain.StatusBusy
		slot.BookedBy = userID
		slot.UpdatedAt = time.Now()
	default:
		return domain.Slot{}, fmt.Errorf("status error")
	}

	j.slots[slotID] = slot

	return slot, nil
}

func (j *JsonStore) CancelBooking(slotID string, userID int) (domain.Slot, error) {

	j.mu.Lock()
	j.mu.Unlock()

	if _, ok := j.slots[slotID]; !ok {
		return domain.Slot{}, fmt.Errorf("slot not found")
	}

	slot := j.slots[slotID]

	switch slot.Status {
	case domain.StatusBusy:
		return domain.Slot{}, fmt.Errorf("slot already booked")
	case domain.StatusCanceled:
		return domain.Slot{}, fmt.Errorf("slot canceled")
	case domain.StatusFree:
		return domain.Slot{}, fmt.Errorf("slot is not booked")

	default:
		return domain.Slot{}, fmt.Errorf("status error")
	}
}

func (j *JsonStore) ListSlotsByUser(userID int) []domain.Slot {
	return
}
