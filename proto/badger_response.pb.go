// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.24.3
// source: badger_response.proto

package __

import (
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

type BadgerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string                       `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value *OtelEnrichedRawSpanForProto `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *BadgerResponse) Reset() {
	*x = BadgerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_badger_response_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BadgerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BadgerResponse) ProtoMessage() {}

func (x *BadgerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_badger_response_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BadgerResponse.ProtoReflect.Descriptor instead.
func (*BadgerResponse) Descriptor() ([]byte, []int) {
	return file_badger_response_proto_rawDescGZIP(), []int{0}
}

func (x *BadgerResponse) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *BadgerResponse) GetValue() *OtelEnrichedRawSpanForProto {
	if x != nil {
		return x.Value
	}
	return nil
}

type BadgerEbpfResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string                `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value *EbpfEntryDataForSpan `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *BadgerEbpfResponse) Reset() {
	*x = BadgerEbpfResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_badger_response_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BadgerEbpfResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BadgerEbpfResponse) ProtoMessage() {}

func (x *BadgerEbpfResponse) ProtoReflect() protoreflect.Message {
	mi := &file_badger_response_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BadgerEbpfResponse.ProtoReflect.Descriptor instead.
func (*BadgerEbpfResponse) Descriptor() ([]byte, []int) {
	return file_badger_response_proto_rawDescGZIP(), []int{1}
}

func (x *BadgerEbpfResponse) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *BadgerEbpfResponse) GetValue() *EbpfEntryDataForSpan {
	if x != nil {
		return x.Value
	}
	return nil
}

type BadgerResponseList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ResponseList     []*BadgerResponse     `protobuf:"bytes,1,rep,name=response_list,json=responseList,proto3" json:"response_list,omitempty"`
	EbpfResponseList []*BadgerEbpfResponse `protobuf:"bytes,2,rep,name=ebpf_response_list,json=ebpfResponseList,proto3" json:"ebpf_response_list,omitempty"`
}

func (x *BadgerResponseList) Reset() {
	*x = BadgerResponseList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_badger_response_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BadgerResponseList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BadgerResponseList) ProtoMessage() {}

func (x *BadgerResponseList) ProtoReflect() protoreflect.Message {
	mi := &file_badger_response_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BadgerResponseList.ProtoReflect.Descriptor instead.
func (*BadgerResponseList) Descriptor() ([]byte, []int) {
	return file_badger_response_proto_rawDescGZIP(), []int{2}
}

func (x *BadgerResponseList) GetResponseList() []*BadgerResponse {
	if x != nil {
		return x.ResponseList
	}
	return nil
}

func (x *BadgerResponseList) GetEbpfResponseList() []*BadgerEbpfResponse {
	if x != nil {
		return x.EbpfResponseList
	}
	return nil
}

var File_badger_response_proto protoreflect.FileDescriptor

var file_badger_response_proto_rawDesc = []byte{
	0x0a, 0x15, 0x62, 0x61, 0x64, 0x67, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13,
	0x6f, 0x70, 0x65, 0x6e, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x65, 0x62, 0x70, 0x66, 0x5f, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x0e, 0x42, 0x61, 0x64, 0x67, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x38, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4f, 0x74, 0x65, 0x6c, 0x45, 0x6e, 0x72, 0x69, 0x63, 0x68, 0x65, 0x64, 0x52, 0x61, 0x77,
	0x53, 0x70, 0x61, 0x6e, 0x46, 0x6f, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x59, 0x0a, 0x12, 0x42, 0x61, 0x64, 0x67, 0x65, 0x72, 0x45, 0x62, 0x70,
	0x66, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x31, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x45, 0x62, 0x70, 0x66, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x44, 0x61, 0x74, 0x61,
	0x46, 0x6f, 0x72, 0x53, 0x70, 0x61, 0x6e, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x99,
	0x01, 0x0a, 0x12, 0x42, 0x61, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x3a, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x61, 0x64, 0x67, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x47, 0x0a, 0x12, 0x65, 0x62, 0x70, 0x66, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x61, 0x64, 0x67, 0x65, 0x72, 0x45, 0x62, 0x70, 0x66,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x10, 0x65, 0x62, 0x70, 0x66, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_badger_response_proto_rawDescOnce sync.Once
	file_badger_response_proto_rawDescData = file_badger_response_proto_rawDesc
)

func file_badger_response_proto_rawDescGZIP() []byte {
	file_badger_response_proto_rawDescOnce.Do(func() {
		file_badger_response_proto_rawDescData = protoimpl.X.CompressGZIP(file_badger_response_proto_rawDescData)
	})
	return file_badger_response_proto_rawDescData
}

var file_badger_response_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_badger_response_proto_goTypes = []interface{}{
	(*BadgerResponse)(nil),              // 0: proto.BadgerResponse
	(*BadgerEbpfResponse)(nil),          // 1: proto.BadgerEbpfResponse
	(*BadgerResponseList)(nil),          // 2: proto.BadgerResponseList
	(*OtelEnrichedRawSpanForProto)(nil), // 3: proto.OtelEnrichedRawSpanForProto
	(*EbpfEntryDataForSpan)(nil),        // 4: proto.EbpfEntryDataForSpan
}
var file_badger_response_proto_depIdxs = []int32{
	3, // 0: proto.BadgerResponse.value:type_name -> proto.OtelEnrichedRawSpanForProto
	4, // 1: proto.BadgerEbpfResponse.value:type_name -> proto.EbpfEntryDataForSpan
	0, // 2: proto.BadgerResponseList.response_list:type_name -> proto.BadgerResponse
	1, // 3: proto.BadgerResponseList.ebpf_response_list:type_name -> proto.BadgerEbpfResponse
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_badger_response_proto_init() }
func file_badger_response_proto_init() {
	if File_badger_response_proto != nil {
		return
	}
	file_opentelemetry_proto_init()
	file_ebpf_entry_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_badger_response_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BadgerResponse); i {
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
		file_badger_response_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BadgerEbpfResponse); i {
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
		file_badger_response_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BadgerResponseList); i {
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
			RawDescriptor: file_badger_response_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_badger_response_proto_goTypes,
		DependencyIndexes: file_badger_response_proto_depIdxs,
		MessageInfos:      file_badger_response_proto_msgTypes,
	}.Build()
	File_badger_response_proto = out.File
	file_badger_response_proto_rawDesc = nil
	file_badger_response_proto_goTypes = nil
	file_badger_response_proto_depIdxs = nil
}
