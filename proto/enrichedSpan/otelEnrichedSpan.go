package enrichedSpan

import (
	"github.com/zerok-ai/zk-utils-go/common"
	logger "github.com/zerok-ai/zk-utils-go/logs"
	protoSpan "github.com/zerok-ai/zk-utils-go/proto/opentelemetry"
	otlpCommon "go.opentelemetry.io/proto/otlp/common/v1"
	otlpTrace "go.opentelemetry.io/proto/otlp/trace/v1"
)

var LogTag = "enrichedSpan"

type OtelEnrichedRawSpan struct {
	Span *otlpTrace.Span `json:"span"`

	// Span Attributes
	SpanAttributes         common.GenericMap   `json:"span_attributes,omitempty"`
	SpanEvents             []common.GenericMap `json:"span_events,omitempty"`
	ResourceAttributesHash string              `json:"resource_attributes_hash,omitempty"`
	ScopeAttributesHash    string              `json:"scope_attributes_hash,omitempty"`

	// ZeroK Properties
	WorkloadIdList []string          `json:"workload_id_list,omitempty"`
	GroupBy        common.GroupByMap `json:"group_by,omitempty"`
}

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

func ConvertGroupByMapToKVList(groupByMap common.GroupByMap) *protoSpan.KeyValueList {
	var groupBy []*otlpCommon.KeyValue
	for k, v := range groupByMap {
		groupBy = append(groupBy, &otlpCommon.KeyValue{
			Key:   string(k),
			Value: ConvertToAnyValue(v),
		})
	}
	return &protoSpan.KeyValueList{KeyValueList: groupBy}
}

func ConvertListOfMapToKVList(attrMap []common.GenericMap) []*protoSpan.KeyValueList {
	var attr []*protoSpan.KeyValueList
	for _, item := range attrMap {
		attr = append(attr, ConvertMapToKVList(item))
	}
	return attr
}

func ConvertMapToKVList(attrMap common.GenericMap) *protoSpan.KeyValueList {
	var attr []*otlpCommon.KeyValue
	for k, v := range attrMap {
		attr = append(attr, &otlpCommon.KeyValue{
			Key:   k,
			Value: ConvertToAnyValue(v),
		})
	}

	return &protoSpan.KeyValueList{KeyValueList: attr}
}

func ConvertToAnyValue(value interface{}) *otlpCommon.AnyValue {
	anyValue := &otlpCommon.AnyValue{}
	switch v := value.(type) {
	case string:
		anyValue.Value = &otlpCommon.AnyValue_StringValue{StringValue: v}
	case []interface{}:
		var arr []*otlpCommon.AnyValue
		for _, item := range v {
			arr = append(arr, ConvertToAnyValue(item))
		}
		anyValue.Value = &otlpCommon.AnyValue_ArrayValue{ArrayValue: &otlpCommon.ArrayValue{Values: arr}}
	case bool:
		anyValue.Value = &otlpCommon.AnyValue_BoolValue{BoolValue: v}
	case float64:
		anyValue.Value = &otlpCommon.AnyValue_DoubleValue{DoubleValue: v}
	case []byte:
		anyValue.Value = &otlpCommon.AnyValue_BytesValue{BytesValue: v}
	case int64:
		anyValue.Value = &otlpCommon.AnyValue_IntValue{IntValue: v}
	default:
		logger.Debug(LogTag, "Unknown type ", v)
	}
	return anyValue
}

func GetEnrichedSpan(x *protoSpan.OtelEnrichedRawSpanForProto) *OtelEnrichedRawSpan {
	span := OtelEnrichedRawSpan{
		Span:                   x.Span,
		SpanAttributes:         ConvertKVListToMap(x.SpanAttributes),
		SpanEvents:             ConvertListOfKVListToMap(x.SpanEvents),
		ResourceAttributesHash: x.ResourceAttributesHash,
		ScopeAttributesHash:    x.ScopeAttributesHash,
		WorkloadIdList:         x.WorkloadIdList,
		GroupBy:                ConvertKVListToGroupByMap(x.GroupBy),
	}
	return &span

}

func ConvertListOfKVListToMap(attrMap []*protoSpan.KeyValueList) []common.GenericMap {
	var attr []common.GenericMap
	for _, item := range attrMap {
		attr = append(attr, ConvertKVListToMap(item))
	}

	return attr
}

func ConvertKVListToMap(attr *protoSpan.KeyValueList) map[string]interface{} {
	attrMap := map[string]interface{}{}
	for _, kv := range attr.KeyValueList {
		value := GetAnyValue(kv.Value)
		if value != nil {
			attrMap[kv.Key] = value
		}
	}

	return attrMap
}

func GetAnyValue(value *otlpCommon.AnyValue) interface{} {
	switch v := value.Value.(type) {
	case *otlpCommon.AnyValue_StringValue:
		return v.StringValue
	case *otlpCommon.AnyValue_ArrayValue:
		var arr []interface{}
		for _, item := range v.ArrayValue.Values {
			arr = append(arr, GetAnyValue(item))
		}
		return arr
	case *otlpCommon.AnyValue_BoolValue:
		return v.BoolValue
	case *otlpCommon.AnyValue_DoubleValue:
		return v.DoubleValue
	case *otlpCommon.AnyValue_BytesValue:
		return v.BytesValue
	case *otlpCommon.AnyValue_IntValue:
		return v.IntValue
	default:
		logger.Debug(LogTag, "Unknown type ", v)
	}
	return nil
}

func ConvertKVListToGroupByMap(attr *protoSpan.KeyValueList) common.GroupByMap {
	attrMap := common.GroupByMap{}
	for _, kv := range attr.KeyValueList {
		value := GetAnyValue(kv.Value)
		if value != nil {
			attrMap[common.ScenarioId(kv.Key)] = ConvertToGroupByValues(value)
		}
	}
	return attrMap
}

func ConvertToGroupByValues(value interface{}) common.GroupByValues {
	var arr common.GroupByValues
	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			arr = append(arr, ConvertToGroupByValueItem(item))
		}
	default:
		logger.Debug(LogTag, "Unknown type ", v)
	}
	return arr
}

func ConvertToGroupByValueItem(value interface{}) *common.GroupByValueItem {
	switch v := value.(type) {
	case map[string]interface{}:
		return &common.GroupByValueItem{
			WorkloadId: v["workload_id"].(string),
			Title:      v["title"].(string),
			Hash:       v["hash"].(string),
		}
	default:
		logger.Debug(LogTag, "Unknown type ", v)
	}
	return nil

}
