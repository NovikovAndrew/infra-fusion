// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.23.4
// source: save_filter_request.proto

package pb

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

type SaveFilterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserID uint64  `protobuf:"varint,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	Filter *Filter `protobuf:"bytes,2,opt,name=filter,proto3" json:"filter,omitempty"`
}

func (x *SaveFilterRequest) Reset() {
	*x = SaveFilterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_save_filter_request_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveFilterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveFilterRequest) ProtoMessage() {}

func (x *SaveFilterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_save_filter_request_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveFilterRequest.ProtoReflect.Descriptor instead.
func (*SaveFilterRequest) Descriptor() ([]byte, []int) {
	return file_save_filter_request_proto_rawDescGZIP(), []int{0}
}

func (x *SaveFilterRequest) GetUserID() uint64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *SaveFilterRequest) GetFilter() *Filter {
	if x != nil {
		return x.Filter
	}
	return nil
}

type SaveFilterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status FilterStatus `protobuf:"varint,1,opt,name=status,proto3,enum=proto.FilterStatus" json:"status,omitempty"`
}

func (x *SaveFilterResponse) Reset() {
	*x = SaveFilterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_save_filter_request_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveFilterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveFilterResponse) ProtoMessage() {}

func (x *SaveFilterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_save_filter_request_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveFilterResponse.ProtoReflect.Descriptor instead.
func (*SaveFilterResponse) Descriptor() ([]byte, []int) {
	return file_save_filter_request_proto_rawDescGZIP(), []int{1}
}

func (x *SaveFilterResponse) GetStatus() FilterStatus {
	if x != nil {
		return x.Status
	}
	return FilterStatus_FILTER_CREATED
}

var File_save_filter_request_proto protoreflect.FileDescriptor

var file_save_filter_request_proto_rawDesc = []byte{
	0x0a, 0x19, 0x73, 0x61, 0x76, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x5f, 0x72, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x10, 0x72, 0x70, 0x63, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x72, 0x70, 0x63, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x52, 0x0a,
	0x11, 0x53, 0x61, 0x76, 0x65, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x25, 0x0a, 0x06, 0x66, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x22, 0x41, 0x0a, 0x12, 0x53, 0x61, 0x76, 0x65, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x4e, 0x6f, 0x76, 0x69, 0x6b, 0x6f, 0x76, 0x41, 0x6e, 0x64, 0x72, 0x65, 0x77,
	0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2d, 0x66, 0x75, 0x73, 0x69, 0x6f, 0x6e, 0x2f, 0x75, 0x73,
	0x65, 0x72, 0x2d, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_save_filter_request_proto_rawDescOnce sync.Once
	file_save_filter_request_proto_rawDescData = file_save_filter_request_proto_rawDesc
)

func file_save_filter_request_proto_rawDescGZIP() []byte {
	file_save_filter_request_proto_rawDescOnce.Do(func() {
		file_save_filter_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_save_filter_request_proto_rawDescData)
	})
	return file_save_filter_request_proto_rawDescData
}

var file_save_filter_request_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_save_filter_request_proto_goTypes = []interface{}{
	(*SaveFilterRequest)(nil),  // 0: proto.SaveFilterRequest
	(*SaveFilterResponse)(nil), // 1: proto.SaveFilterResponse
	(*Filter)(nil),             // 2: proto.Filter
	(FilterStatus)(0),          // 3: proto.FilterStatus
}
var file_save_filter_request_proto_depIdxs = []int32{
	2, // 0: proto.SaveFilterRequest.filter:type_name -> proto.Filter
	3, // 1: proto.SaveFilterResponse.status:type_name -> proto.FilterStatus
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_save_filter_request_proto_init() }
func file_save_filter_request_proto_init() {
	if File_save_filter_request_proto != nil {
		return
	}
	file_rpc_filter_proto_init()
	file_rpc_filter_status_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_save_filter_request_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveFilterRequest); i {
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
		file_save_filter_request_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveFilterResponse); i {
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
			RawDescriptor: file_save_filter_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_save_filter_request_proto_goTypes,
		DependencyIndexes: file_save_filter_request_proto_depIdxs,
		MessageInfos:      file_save_filter_request_proto_msgTypes,
	}.Build()
	File_save_filter_request_proto = out.File
	file_save_filter_request_proto_rawDesc = nil
	file_save_filter_request_proto_goTypes = nil
	file_save_filter_request_proto_depIdxs = nil
}
