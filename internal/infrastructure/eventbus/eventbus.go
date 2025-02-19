package eventbus

import (
	"context"
	"sync"
)

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event Event) error

// EventBus 事件总线接口
type EventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event Event) error
	// Subscribe 订阅事件
	Subscribe(eventName string, handler EventHandler)
}

// LocalEventBus 本地事件总线实现
type LocalEventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewLocalEventBus 创建本地事件总线
func NewLocalEventBus() *LocalEventBus {
	return &LocalEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Publish 发布事件
func (bus *LocalEventBus) Publish(ctx context.Context, event Event) error {
	bus.mu.RLock()
	defer bus.mu.RUnlock()

	handlers, exists := bus.handlers[event.EventName()]
	if !exists {
		return nil
	}

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

// Subscribe 订阅事件
func (bus *LocalEventBus) Subscribe(eventName string, handler EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.handlers[eventName] = append(bus.handlers[eventName], handler)
}
