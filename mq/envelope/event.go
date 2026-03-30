package envelope

import (
	"encoding/json"
	"time"
)

// Event 是跨服务共享的统一领域事件外层结构。
// Payload 采用 json.RawMessage，便于与具体业务事件 DTO 解耦。
//
// 这一层更偏“消息语义模型”：
// 业务事件先组织成统一 Envelope，再决定是直接发送、写 Outbox，
// 还是交给 CDC 搬运到 Kafka。
type Event struct {
	Version       string          `json:"version"`
	EventType     string          `json:"event_type"`
	AggregateType string          `json:"aggregate_type,omitempty"`
	AggregateID   string          `json:"aggregate_id,omitempty"`
	Timestamp     int64           `json:"timestamp"`
	Payload       json.RawMessage `json:"payload"`
}

// New 构造统一事件外层结构。
// Timestamp 在这里统一生成，避免各服务各自拼装时字段风格不一致。
func New(version, eventType, aggregateType, aggregateID string, payload json.RawMessage) Event {
	return Event{
		Version:       version,
		EventType:     eventType,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Timestamp:     time.Now().Unix(),
		Payload:       payload,
	}
}
