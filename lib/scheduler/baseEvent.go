package scheduler

import (
	"reflect"
	"runtime"
)

// BaseEvent - базовая структура для всех процессоров
type BaseEvent struct {
	name         string
	version      string
	description  string
	capabilities []string
}

func (bp *BaseEvent) GetName() string {
	return bp.name
}

func (bp *BaseEvent) GetVersion() string {
	return bp.version
}

func (bp *BaseEvent) GetCapabilities() []string {
	return bp.capabilities
}

func (bp *BaseEvent) GetDescription() string {
	return bp.description
}

// Register регистрирует процессор в глобальном реестре
func Register(p Event) error {
	return GetRegistry().Register(p)
}

// GetFunctionName возвращает имя функции для отладки
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
