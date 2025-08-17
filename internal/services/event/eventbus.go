package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// EventBus 事件总线接口
type EventBus interface {
	Publish(ctx context.Context, topic string, event interface{}) error
	Subscribe(ctx context.Context, topic string, handler EventHandler) error
	Close()
}

// EventHandler 事件处理器接口
type EventHandler func(ctx context.Context, event []byte) error

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
func (e *NATSEventBus) Publish(ctx context.Context, topic string, event interface{}) error {
	// 检查Context是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 使用Context控制发布超时
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return e.conn.Publish(topic, data)
	}
}

// Subscribe 订阅事件
func (e *NATSEventBus) Subscribe(ctx context.Context, topic string, handler EventHandler) error {
	_, err := e.conn.Subscribe(topic, func(msg *nats.Msg) {
		// 为每个消息创建Context
		msgCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := handler(msgCtx, msg.Data); err != nil {
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
	AlertID  string `json:"alert_id"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// EventType 事件类型常量
const (
	EventTypeHardwareDiscovered  = "hardware.discovered"
	EventTypeHardwareProvisioned = "hardware.provisioned"
	EventTypeHardwareFailed      = "hardware.failed"
	EventTypeAlertRaised         = "alert.raised"
	EventTypeAlertResolved       = "alert.resolved"
)
