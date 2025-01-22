// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.0
// 	protoc        v5.29.2
// source: core/sec/insecure/pb/plaintext.proto

package pb

import (
	pb "github.com/dep2p/core/crypto/pb"
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

type Exchange struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            []byte                 `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Pubkey        *pb.PublicKey          `protobuf:"bytes,2,opt,name=pubkey" json:"pubkey,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Exchange) Reset() {
	*x = Exchange{}
	mi := &file_core_sec_insecure_pb_plaintext_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Exchange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Exchange) ProtoMessage() {}

func (x *Exchange) ProtoReflect() protoreflect.Message {
	mi := &file_core_sec_insecure_pb_plaintext_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Exchange.ProtoReflect.Descriptor instead.
func (*Exchange) Descriptor() ([]byte, []int) {
	return file_core_sec_insecure_pb_plaintext_proto_rawDescGZIP(), []int{0}
}

func (x *Exchange) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Exchange) GetPubkey() *pb.PublicKey {
	if x != nil {
		return x.Pubkey
	}
	return nil
}

var File_core_sec_insecure_pb_plaintext_proto protoreflect.FileDescriptor

var file_core_sec_insecure_pb_plaintext_proto_rawDesc = []byte{
	0x0a, 0x24, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x73, 0x65, 0x63, 0x2f, 0x69, 0x6e, 0x73, 0x65, 0x63,
	0x75, 0x72, 0x65, 0x2f, 0x70, 0x62, 0x2f, 0x70, 0x6c, 0x61, 0x69, 0x6e, 0x74, 0x65, 0x78, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x70, 0x6c, 0x61, 0x69, 0x6e, 0x74, 0x65, 0x78,
	0x74, 0x2e, 0x70, 0x62, 0x1a, 0x1b, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x63, 0x72, 0x79, 0x70, 0x74,
	0x6f, 0x2f, 0x70, 0x62, 0x2f, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x48, 0x0a, 0x08, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2c, 0x0a,
	0x06, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x4b, 0x65, 0x79, 0x52, 0x06, 0x70, 0x75, 0x62, 0x6b, 0x65, 0x79, 0x42, 0x32, 0x5a, 0x30, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69, 0x62, 0x70, 0x32, 0x70,
	0x2f, 0x67, 0x6f, 0x2d, 0x6c, 0x69, 0x62, 0x70, 0x32, 0x70, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x73, 0x65, 0x63, 0x2f, 0x69, 0x6e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x65, 0x2f, 0x70, 0x62,
}

var (
	file_core_sec_insecure_pb_plaintext_proto_rawDescOnce sync.Once
	file_core_sec_insecure_pb_plaintext_proto_rawDescData = file_core_sec_insecure_pb_plaintext_proto_rawDesc
)

func file_core_sec_insecure_pb_plaintext_proto_rawDescGZIP() []byte {
	file_core_sec_insecure_pb_plaintext_proto_rawDescOnce.Do(func() {
		file_core_sec_insecure_pb_plaintext_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_sec_insecure_pb_plaintext_proto_rawDescData)
	})
	return file_core_sec_insecure_pb_plaintext_proto_rawDescData
}

var file_core_sec_insecure_pb_plaintext_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_core_sec_insecure_pb_plaintext_proto_goTypes = []any{
	(*Exchange)(nil),     // 0: plaintext.pb.Exchange
	(*pb.PublicKey)(nil), // 1: crypto.pb.PublicKey
}
var file_core_sec_insecure_pb_plaintext_proto_depIdxs = []int32{
	1, // 0: plaintext.pb.Exchange.pubkey:type_name -> crypto.pb.PublicKey
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_core_sec_insecure_pb_plaintext_proto_init() }
func file_core_sec_insecure_pb_plaintext_proto_init() {
	if File_core_sec_insecure_pb_plaintext_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_core_sec_insecure_pb_plaintext_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_sec_insecure_pb_plaintext_proto_goTypes,
		DependencyIndexes: file_core_sec_insecure_pb_plaintext_proto_depIdxs,
		MessageInfos:      file_core_sec_insecure_pb_plaintext_proto_msgTypes,
	}.Build()
	File_core_sec_insecure_pb_plaintext_proto = out.File
	file_core_sec_insecure_pb_plaintext_proto_rawDesc = nil
	file_core_sec_insecure_pb_plaintext_proto_goTypes = nil
	file_core_sec_insecure_pb_plaintext_proto_depIdxs = nil
}
