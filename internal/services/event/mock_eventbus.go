package event

// MockEventBus 模拟事件总线
type MockEventBus struct{}

// Publish 模拟发布事件
func (m *MockEventBus) Publish(topic string, event interface{}) error {
	// 模拟发布事件，实际不做任何操作
	return nil
}

// Subscribe 模拟订阅事件
func (m *MockEventBus) Subscribe(topic string, handler EventHandler) error {
	// 模拟订阅事件，实际不做任何操作
	return nil
}

// Close 模拟关闭连接
func (m *MockEventBus) Close() {
	// 模拟关闭连接，实际不做任何操作
}
