package dto

import (
	"encoding/json"
	"fmt"
)

type FactorialMessage struct {
	Number int64 `json:"number"`
}

func (m *FactorialMessage) Bytes() []byte {
	return fmt.Appendf(nil, "{\"number\": %d}", m.Number)
}

func (m *FactorialMessage) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}
