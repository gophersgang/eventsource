package eventsource

import (
	"context"
	"sort"
	"sync"
)

// Record provides the serialized representation of the event
type Record struct {
	// Version contains the version associated with the serialized event
	Version int

	// Data contains the event in serialized form
	Data []byte
}

// History represents
type History []Record

// Len implements sort.Interface
func (h History) Len() int {
	return len(h)
}

// Swap implements sort.Interface
func (h History) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Less implements sort.Interface
func (h History) Less(i, j int) bool {
	return h[i].Version < h[j].Version
}

// Store provides an abstraction for the Repository to save data
type Store interface {
	// Save the provided serialized records to the store
	Save(ctx context.Context, aggregateID string, records ...Record) error

	// Load the history of events up to the version specified; when version is
	// 0, all events will be loaded
	Load(ctx context.Context, aggregateID string, version int) (History, error)
}

// memoryStore provides an in-memory implementation of Store
type memoryStore struct {
	mux        *sync.Mutex
	eventsByID map[string]History
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		mux:        &sync.Mutex{},
		eventsByID: map[string]History{},
	}
}

func (m *memoryStore) Save(ctx context.Context, aggregateID string, records ...Record) error {
	if _, ok := m.eventsByID[aggregateID]; !ok {
		m.eventsByID[aggregateID] = History{}
	}

	history := append(m.eventsByID[aggregateID], records...)
	sort.Sort(history)
	m.eventsByID[aggregateID] = history

	return nil
}

func (m *memoryStore) Load(ctx context.Context, aggregateID string, version int) (History, error) {
	history, ok := m.eventsByID[aggregateID]
	if !ok {
		return nil, NewError(nil, AggregateNotFound, "no aggregate found with id, %v", aggregateID)
	}

	return history, nil
}
