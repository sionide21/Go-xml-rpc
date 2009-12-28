package xmlrpc

import (
	"time"
)

func (st StructValue) GetInt(s string) int { return int(st[s].(IntValue)) }

func (st StructValue) GetString(s string) string {
	return string(st[s].(StringValue))
}

func (st StructValue) GetDouble(s string) float {
	return float(st[s].(DoubleValue))
}

func (st StructValue) GetTime(s string) time.Time {
	return time.Time(st[s].(DateTimeValue))
}

func (st StructValue) GetBoolean(s string) bool {
	return bool(st[s].(BooleanValue))
}

func (st StructValue) GetBytes(s string) []byte {
	return []byte(st[s].(Base64Value))
}

func (st StructValue) GetArray(s string) ArrayValue {
	return st[s].(ArrayValue)
}

func (st StructValue) GetStruct(s string) StructValue {
	return st[s].(StructValue)
}
