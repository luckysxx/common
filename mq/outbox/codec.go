package outbox

import "encoding/json"

// EncodePayload 将任意业务事件编码为 JSON payload。
// 这样 service 层不需要到处手写 json.Marshal。
func EncodePayload(v any) ([]byte, error) {
	return json.Marshal(v)
}

// EncodeHeaders 将事件头编码为 JSON。
// headers 为空时返回 nil，方便调用方按需落库。
// 当前可以先不用 headers，但预留出来后，后续 trace_id、source 等元数据更好接入。
func EncodeHeaders(headers map[string]string) ([]byte, error) {
	if len(headers) == 0 {
		return nil, nil
	}
	return json.Marshal(headers)
}

// NewJSONRecord 构造一条以 JSON 为载荷的 outbox 记录。
// 这是业务层最推荐直接调用的入口：
// 给它业务 payload，它返回一条标准化的 Record。
func NewJSONRecord(id, aggregateType, aggregateID, eventType string, payload any, headers map[string]string) (*Record, error) {
	payloadBytes, err := EncodePayload(payload)
	if err != nil {
		return nil, err
	}

	headerBytes, err := EncodeHeaders(headers)
	if err != nil {
		return nil, err
	}

	return NewRecord(id, aggregateType, aggregateID, eventType, payloadBytes, headerBytes), nil
}
