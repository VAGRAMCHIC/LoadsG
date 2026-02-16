package scheduler

import (
	"errors"
	"fmt"
	"sync"
)

var ErrEventTypeNotFound = errors.New("error: event type not found")

type Event interface {
	Run(data interface{}) (interface{}, error)
	GetName() string
	GetVersion() string
	GetCapabilities() []string
}

type EventRegistry struct {
	events      map[string]Event
	mu          sync.RWMutex
	initialized bool
}

// Global registry instance
var globalRegistry = &EventRegistry{
	events: make(map[string]Event),
}

// GetRegistry возвращает глобальный реестр
func GetRegistry() *EventRegistry {
	return globalRegistry
}

func (r *EventRegistry) Register(e Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := e.GetName()
	if _, exists := r.events[name]; exists {
		return fmt.Errorf("processor '%s' already registered", name)
	}

	r.events[name] = e
	fmt.Printf("✅ Processor registered: %s v%s\n", name, e.GetVersion())
	return nil
}

func (r *EventRegistry) GetEvent(name string) (Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	event, exists := r.events[name]
	if !exists {
		return nil, fmt.Errorf("processor '%s' not found", name)
	}

	return event, nil
}

func (r *EventRegistry) List() []Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Event, 0, len(r.events))
	for _, p := range r.events {
		list = append(list, p)
	}

	return list
}

func (r *EventRegistry) FindByCapability(capability string) []Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Event
	for _, p := range r.events {
		for _, c := range p.GetCapabilities() {
			if c == capability {
				result = append(result, p)
				break
			}
		}
	}

	return result
}

func (r *EventRegistry) Initialize() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.initialized {
		return nil
	}

	// Здесь может быть загрузка конфигурации, плагинов и т.д.
	r.initialized = true
	fmt.Println("📋 Registry initialized")

	return nil
}

// Size возвращает количество зарегистрированных процессоров
func (r *EventRegistry) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.events)
}

// Remove удаляет процессор из реестра
func (r *EventRegistry) Remove(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[name]; exists {
		delete(r.events, name)
		fmt.Printf("🗑️ Processor removed: %s\n", name)
		return true
	}

	return false
}
