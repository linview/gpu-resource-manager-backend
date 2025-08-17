package event

import (
	"context"
	"log"
	"time"
)

// MockEventBus 模拟事件总线
type MockEventBus struct{}

// NewMockEventBus 创建模拟事件总线
func NewMockEventBus() EventBus {
	return &MockEventBus{}
}

// Publish 模拟发布事件
func (m *MockEventBus) Publish(ctx context.Context, topic string, event interface{}) error {
	// 检查Context是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 模拟发布延迟
	select {
	case <-time.After(10 * time.Millisecond):
		log.Printf("Mock: Published event to topic %s", topic)
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Subscribe 模拟订阅事件
func (m *MockEventBus) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	// 检查Context是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 模拟订阅
	log.Printf("Mock: Subscribed to topic %s", topic)
	return nil
}

// Close 模拟关闭连接
func (m *MockEventBus) Close() {
	log.Println("Mock: EventBus closed")
}
