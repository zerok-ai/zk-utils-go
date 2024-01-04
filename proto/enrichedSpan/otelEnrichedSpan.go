package enrichedSpan

import (
	logger "github.com/zerok-ai/zk-utils-go/logs"
	protoSpan "github.com/zerok-ai/zk-utils-go/proto/opentelemetry"
	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

var LogTag = "enrichedSpan"

type OtelEnrichedRawSpan struct {
	Span *v1.Span `json:"span"`

	// Span Attributes
	SpanAttributes         GenericMap   `json:"span_attributes,omitempty"`
	SpanEvents             []GenericMap `json:"span_events,omitempty"`
	ResourceAttributesHash string       `json:"resource_attributes_hash,omitempty"`
	ScopeAttributesHash    string       `json:"scope_attributes_hash,omitempty"`

	// ZeroK Properties
	WorkloadIdList []string   `json:"workload_id_list,omitempty"`
	GroupBy        GroupByMap `json:"group_by,omitempty"`
}

type GenericMap map[string]interface{}

type GroupByMap map[ScenarioId]GroupByValues

type GroupByValueItem struct {
	WorkloadId string `json:"workload_id"`
	Title      string `json:"title"`
	Hash       string `json:"hash"`
}
type GroupByValues []*GroupByValueItem
type ScenarioId string

func (x *OtelEnrichedRawSpan) GetProtoEnrichedSpan() *protoSpan.OtelEnrichedRawSpanForProto {
	span := protoSpan.OtelEnrichedRawSpanForProto{
		Span:                   x.Span,
		SpanAttributes:         ConvertMapToKVList(x.SpanAttributes),
		SpanEvents:             ConvertListOfMapToKVList(x.SpanEvents),
		ResourceAttributesHash: x.ResourceAttributesHash,
		ScopeAttributesHash:    x.ScopeAttributesHash,
		WorkloadIdList:         x.WorkloadIdList,
		GroupBy:                ConvertGroupByMapToKVList(x.GroupBy),
	}
	return &span
}

func ConvertGroupByMapToKVList(groupByMap GroupByMap) *protoSpan.KeyValueList {
	var groupBy []*v11.KeyValue
	for k, v := range groupByMap {
		groupBy = append(groupBy, &v11.KeyValue{
			Key:   string(k),
			Value: ConvertToAnyValue(v),
		})
	}
	return &protoSpan.KeyValueList{KeyValueList: groupBy}
}

func ConvertListOfMapToKVList(attrMap []GenericMap) []*protoSpan.KeyValueList {
	var attr []*protoSpan.KeyValueList
	for _, item := range attrMap {
		attr = append(attr, ConvertMapToKVList(item))
	}
	return attr
}

func ConvertMapToKVList(attrMap GenericMap) *protoSpan.KeyValueList {
	var attr []*v11.KeyValue
	for k, v := range attrMap {
		attr = append(attr, &v11.KeyValue{
			Key:   k,
			Value: ConvertToAnyValue(v),
		})
	}

	return &protoSpan.KeyValueList{KeyValueList: attr}
}

func ConvertToAnyValue(value interface{}) *v11.AnyValue {
	anyValue := &v11.AnyValue{}
	switch v := value.(type) {
	case string:
		anyValue.Value = &v11.AnyValue_StringValue{StringValue: v}
	case []interface{}:
		var arr []*v11.AnyValue
		for _, item := range v {
			arr = append(arr, ConvertToAnyValue(item))
		}
		anyValue.Value = &v11.AnyValue_ArrayValue{ArrayValue: &v11.ArrayValue{Values: arr}}
	case bool:
		anyValue.Value = &v11.AnyValue_BoolValue{BoolValue: v}
	case float64:
		anyValue.Value = &v11.AnyValue_DoubleValue{DoubleValue: v}
	case []byte:
		anyValue.Value = &v11.AnyValue_BytesValue{BytesValue: v}
	case int64:
		anyValue.Value = &v11.AnyValue_IntValue{IntValue: v}
	default:
		logger.Debug(LogTag, "Unknown type ", v)
	}
	return anyValue
}
