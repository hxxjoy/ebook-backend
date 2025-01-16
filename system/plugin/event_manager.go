// system/core/plugin/event_manager.go
package plugin

import (
	"sync"
)

// EventHandler 事件处理函数类型
type EventHandler func(interface{}) error

// EventManager 事件管理器
type EventManager struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewEventManager 创建事件管理器
func NewEventManager() *EventManager {
	return &EventManager{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (em *EventManager) Subscribe(event string, handler EventHandler) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if _, exists := em.handlers[event]; !exists {
		em.handlers[event] = make([]EventHandler, 0)
	}
	em.handlers[event] = append(em.handlers[event], handler)
}

// Unsubscribe 取消订阅事件
func (em *EventManager) Unsubscribe(event string, handler EventHandler) {
	em.mu.Lock()
	defer em.mu.Unlock()

	if handlers, exists := em.handlers[event]; exists {
		for i, h := range handlers {
			if &h == &handler {
				em.handlers[event] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Emit 触发事件
func (em *EventManager) Emit(event string, data interface{}) error {
	em.mu.RLock()
	handlers := em.handlers[event]
	em.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(data); err != nil {
			return err
		}
	}
	return nil
}
