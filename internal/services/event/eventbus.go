package event

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

// EventBus 事件总线接口
type EventBus interface {
	Publish(topic string, event interface{}) error
	Subscribe(topic string, handler EventHandler) error
	Close()
}

// EventHandler 事件处理器接口
type EventHandler func(event []byte) error

// NATSEventBus NATS事件总线实现
type NATSEventBus struct {
	conn *nats.Conn
}

// NewEventBus 创建新的事件总线
func NewEventBus(natsURL string) EventBus {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	return &NATSEventBus{
		conn: conn,
	}
}

// Publish 发布事件
func (e *NATSEventBus) Publish(topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return e.conn.Publish(topic, data)
}

// Subscribe 订阅事件
func (e *NATSEventBus) Subscribe(topic string, handler EventHandler) error {
	_, err := e.conn.Subscribe(topic, func(msg *nats.Msg) {
		if err := handler(msg.Data); err != nil {
			log.Printf("Error handling event on topic %s: %v", topic, err)
		}
	})

	return err
}

// Close 关闭连接
func (e *NATSEventBus) Close() {
	if e.conn != nil {
		e.conn.Close()
	}
}

// Event 基础事件结构
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
}

// HardwareEvent 硬件事件
type HardwareEvent struct {
	Event
	HardwareID string `json:"hardware_id"`
	Status     string `json:"status"`
}

// AlertEvent 告警事件
type AlertEvent struct {
	Event
	AlertID   string `json:"alert_id"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
}

// EventType 事件类型常量
const (
	EventTypeHardwareDiscovered = "hardware.discovered"
	EventTypeHardwareProvisioned = "hardware.provisioned"
	EventTypeHardwareFailed = "hardware.failed"
	EventTypeAlertRaised = "alert.raised"
	EventTypeAlertResolved = "alert.resolved"
)
