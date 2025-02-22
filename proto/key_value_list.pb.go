// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.24.3
// source: key_value_list.proto

package __

import (
	v1 "go.opentelemetry.io/proto/otlp/common/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type KeyValueList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KeyValueList []*v1.KeyValue `protobuf:"bytes,1,rep,name=key_value_list,json=keyValueList,proto3" json:"key_value_list,omitempty"`
}

func (x *KeyValueList) Reset() {
	*x = KeyValueList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_key_value_list_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValueList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValueList) ProtoMessage() {}

func (x *KeyValueList) ProtoReflect() protoreflect.Message {
	mi := &file_key_value_list_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValueList.ProtoReflect.Descriptor instead.
func (*KeyValueList) Descriptor() ([]byte, []int) {
	return file_key_value_list_proto_rawDescGZIP(), []int{0}
}

func (x *KeyValueList) GetKeyValueList() []*v1.KeyValue {
	if x != nil {
		return x.KeyValueList
	}
	return nil
}

var File_key_value_list_proto protoreflect.FileDescriptor

var file_key_value_list_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6b, 0x65, 0x79, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x6c, 0x69, 0x73, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x6f,
	0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x5d, 0x0a, 0x0c, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4c,
	0x69, 0x73, 0x74, 0x12, 0x4d, 0x0a, 0x0e, 0x6b, 0x65, 0x79, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x6f, 0x70,
	0x65, 0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4b, 0x65, 0x79, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x0c, 0x6b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_key_value_list_proto_rawDescOnce sync.Once
	file_key_value_list_proto_rawDescData = file_key_value_list_proto_rawDesc
)

func file_key_value_list_proto_rawDescGZIP() []byte {
	file_key_value_list_proto_rawDescOnce.Do(func() {
		file_key_value_list_proto_rawDescData = protoimpl.X.CompressGZIP(file_key_value_list_proto_rawDescData)
	})
	return file_key_value_list_proto_rawDescData
}

var file_key_value_list_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_key_value_list_proto_goTypes = []interface{}{
	(*KeyValueList)(nil), // 0: proto.KeyValueList
	(*v1.KeyValue)(nil),  // 1: opentelemetry.proto.common.v1.KeyValue
}
var file_key_value_list_proto_depIdxs = []int32{
	1, // 0: proto.KeyValueList.key_value_list:type_name -> opentelemetry.proto.common.v1.KeyValue
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_key_value_list_proto_init() }
func file_key_value_list_proto_init() {
	if File_key_value_list_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_key_value_list_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValueList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_key_value_list_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_key_value_list_proto_goTypes,
		DependencyIndexes: file_key_value_list_proto_depIdxs,
		MessageInfos:      file_key_value_list_proto_msgTypes,
	}.Build()
	File_key_value_list_proto = out.File
	file_key_value_list_proto_rawDesc = nil
	file_key_value_list_proto_goTypes = nil
	file_key_value_list_proto_depIdxs = nil
}
