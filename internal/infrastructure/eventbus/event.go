package eventbus

import "time"

// Event 领域事件接口
type Event interface {
	// EventName 返回事件名称
	EventName() string
	// EventTime 事件发生时间
	EventTime() time.Time
}

// BaseEvent 基础事件实现
type BaseEvent struct {
	Name string    `json:"name"`
	Time time.Time `json:"time"`
}

func NewBaseEvent(name string) BaseEvent {
	return BaseEvent{
		Name: name,
		Time: time.Now(),
	}
}

func (e BaseEvent) EventName() string {
	return e.Name
}

func (e BaseEvent) EventTime() time.Time {
	return e.Time
}
