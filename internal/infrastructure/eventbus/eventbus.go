package eventbus

import "context"

// EventBus 事件总线接口
type EventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event interface{}) error
	// Subscribe 订阅事件
	Subscribe(eventType string, handler func(event interface{}) error) error
}

// SimpleEventBus 简单的内存事件总线实现
type SimpleEventBus struct {
	handlers map[string][]func(event interface{}) error
}

// NewSimpleEventBus 创建一个简单的事件总线
func NewSimpleEventBus() *SimpleEventBus {
	return &SimpleEventBus{
		handlers: make(map[string][]func(event interface{}) error),
	}
}

// Publish 发布事件
func (b *SimpleEventBus) Publish(ctx context.Context, event interface{}) error {
	// 获取事件类型
	var eventType string
	switch e := event.(type) {
	case interface{ GetType() string }:
		eventType = e.GetType()
	default:
		// 如果事件没有实现 GetType 方法，使用类型名作为事件类型
		eventType = "unknown"
	}

	// 调用所有相关的处理器
	handlers := b.handlers[eventType]
	for _, handler := range handlers {
		if err := handler(event); err != nil {
			// 在实际应用中，可能需要更复杂的错误处理策略
			return err
		}
	}

	return nil
}

// Subscribe 订阅事件
func (b *SimpleEventBus) Subscribe(eventType string, handler func(event interface{}) error) error {
	if b.handlers[eventType] == nil {
		b.handlers[eventType] = make([]func(event interface{}) error, 0)
	}
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}
