// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.4
// source: like.proto

package like

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

type GeteTotalFavoritedNumReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId []int64 `protobuf:"varint,1,rep,packed,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GeteTotalFavoritedNumReq) Reset() {
	*x = GeteTotalFavoritedNumReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_like_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeteTotalFavoritedNumReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeteTotalFavoritedNumReq) ProtoMessage() {}

func (x *GeteTotalFavoritedNumReq) ProtoReflect() protoreflect.Message {
	mi := &file_like_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeteTotalFavoritedNumReq.ProtoReflect.Descriptor instead.
func (*GeteTotalFavoritedNumReq) Descriptor() ([]byte, []int) {
	return file_like_proto_rawDescGZIP(), []int{0}
}

func (x *GeteTotalFavoritedNumReq) GetUserId() []int64 {
	if x != nil {
		return x.UserId
	}
	return nil
}

type GetFavoriteCountByUserIdReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId []int64 `protobuf:"varint,1,rep,packed,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GetFavoriteCountByUserIdReq) Reset() {
	*x = GetFavoriteCountByUserIdReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_like_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFavoriteCountByUserIdReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFavoriteCountByUserIdReq) ProtoMessage() {}

func (x *GetFavoriteCountByUserIdReq) ProtoReflect() protoreflect.Message {
	mi := &file_like_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFavoriteCountByUserIdReq.ProtoReflect.Descriptor instead.
func (*GetFavoriteCountByUserIdReq) Descriptor() ([]byte, []int) {
	return file_like_proto_rawDescGZIP(), []int{1}
}

func (x *GetFavoriteCountByUserIdReq) GetUserId() []int64 {
	if x != nil {
		return x.UserId
	}
	return nil
}

type IsFavoriteReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoId []int64 `protobuf:"varint,1,rep,packed,name=videoId,proto3" json:"videoId,omitempty"`
	UserId  []int64 `protobuf:"varint,2,rep,packed,name=userId,proto3" json:"userId,omitempty"`
}

func (x *IsFavoriteReq) Reset() {
	*x = IsFavoriteReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_like_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsFavoriteReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsFavoriteReq) ProtoMessage() {}

func (x *IsFavoriteReq) ProtoReflect() protoreflect.Message {
	mi := &file_like_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsFavoriteReq.ProtoReflect.Descriptor instead.
func (*IsFavoriteReq) Descriptor() ([]byte, []int) {
	return file_like_proto_rawDescGZIP(), []int{2}
}

func (x *IsFavoriteReq) GetVideoId() []int64 {
	if x != nil {
		return x.VideoId
	}
	return nil
}

func (x *IsFavoriteReq) GetUserId() []int64 {
	if x != nil {
		return x.UserId
	}
	return nil
}

type GetFavoriteCountByVideoIdReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  []int64 `protobuf:"varint,1,rep,packed,name=userId,proto3" json:"userId,omitempty"`
	VideoId []int64 `protobuf:"varint,2,rep,packed,name=videoId,proto3" json:"videoId,omitempty"`
}

func (x *GetFavoriteCountByVideoIdReq) Reset() {
	*x = GetFavoriteCountByVideoIdReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_like_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFavoriteCountByVideoIdReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFavoriteCountByVideoIdReq) ProtoMessage() {}

func (x *GetFavoriteCountByVideoIdReq) ProtoReflect() protoreflect.Message {
	mi := &file_like_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFavoriteCountByVideoIdReq.ProtoReflect.Descriptor instead.
func (*GetFavoriteCountByVideoIdReq) Descriptor() ([]byte, []int) {
	return file_like_proto_rawDescGZIP(), []int{3}
}

func (x *GetFavoriteCountByVideoIdReq) GetUserId() []int64 {
	if x != nil {
		return x.UserId
	}
	return nil
}

func (x *GetFavoriteCountByVideoIdReq) GetVideoId() []int64 {
	if x != nil {
		return x.VideoId
	}
	return nil
}

// 响应
type GeteTotalFavoritedNumReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count []int64 `protobuf:"varint,1,rep,packed,name=count,proto3" json:"count,omitempty"`
}

func (x *GeteTotalFavoritedNumReply) Reset() {
	*x = GeteTotalFavoritedNumReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_like_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeteTotalFavoritedNumReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeteTotalFavoritedNumReply) ProtoMessage() {}

func (x *GeteTotalFavoritedNumReply) ProtoReflect() protoreflect.Message {
	mi := &file_like_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeteTotalFavoritedNumReply.ProtoReflect.Descriptor instead.
func (*GeteTotalFavoritedNumReply) Descriptor() ([]byte, []int) {
	return file_like_proto_rawDescGZIP(), []int{4}
}

func (x *GeteTotalFavoritedNumReply) GetCount() []int64 {
	if x != nil {
		return x.Count
	}
	return nil
}

type GetFavoriteCountByUserIdReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count []int64 `protobuf:"varint,1,rep,packed,name=count,proto3" json:"count,omitempty"`
}

func (x *GetFavoriteCountByUserIdReply) Reset() {
	*x = GetFavoriteCountByUserIdReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_like_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFavoriteCountByUserIdReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFavoriteCountByUserIdReply) ProtoMessage() {}

func (x *GetFavoriteCountByUserIdReply) ProtoReflect() protoreflect.Message {
	mi := &file_like_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFavoriteCountByUserIdReply.ProtoReflect.Descriptor instead.
func (*GetFavoriteCountByUserIdReply) Descriptor() ([]byte, []int) {
	return file_like_proto_rawDescGZIP(), []int{5}
}

func (x *GetFavoriteCountByUserIdReply) GetCount() []int64 {
	if x != nil {
		return x.Count
	}
	return nil
}

type IsFavoriteReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsFavorite []bool `protobuf:"varint,1,rep,packed,name=isFavorite,proto3" json:"isFavorite,omitempty"`
}

func (x *IsFavoriteReply) Reset() {
	*x = IsFavoriteReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_like_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IsFavoriteReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IsFavoriteReply) ProtoMessage() {}

func (x *IsFavoriteReply) ProtoReflect() protoreflect.Message {
	mi := &file_like_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IsFavoriteReply.ProtoReflect.Descriptor instead.
func (*IsFavoriteReply) Descriptor() ([]byte, []int) {
	return file_like_proto_rawDescGZIP(), []int{6}
}

func (x *IsFavoriteReply) GetIsFavorite() []bool {
	if x != nil {
		return x.IsFavorite
	}
	return nil
}

type EtFavoriteCountByUserIdReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count []int64 `protobuf:"varint,1,rep,packed,name=count,proto3" json:"count,omitempty"`
}

func (x *EtFavoriteCountByUserIdReply) Reset() {
	*x = EtFavoriteCountByUserIdReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_like_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EtFavoriteCountByUserIdReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EtFavoriteCountByUserIdReply) ProtoMessage() {}

func (x *EtFavoriteCountByUserIdReply) ProtoReflect() protoreflect.Message {
	mi := &file_like_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EtFavoriteCountByUserIdReply.ProtoReflect.Descriptor instead.
func (*EtFavoriteCountByUserIdReply) Descriptor() ([]byte, []int) {
	return file_like_proto_rawDescGZIP(), []int{7}
}

func (x *EtFavoriteCountByUserIdReply) GetCount() []int64 {
	if x != nil {
		return x.Count
	}
	return nil
}

var File_like_proto protoreflect.FileDescriptor

var file_like_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6c, 0x69, 0x6b, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x22, 0x32, 0x0a, 0x18, 0x67, 0x65, 0x74, 0x65, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x46,
	0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x64, 0x4e, 0x75, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x16,
	0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x35, 0x0a, 0x1b, 0x67, 0x65, 0x74, 0x46, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x41, 0x0a,
	0x0d, 0x69, 0x73, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x12, 0x18,
	0x0a, 0x07, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52,
	0x07, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x22, 0x50, 0x0a, 0x1c, 0x67, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x52, 0x65, 0x71,
	0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x69, 0x64, 0x65,
	0x6f, 0x49, 0x64, 0x18, 0x02, 0x20, 0x03, 0x28, 0x03, 0x52, 0x07, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x49, 0x64, 0x22, 0x32, 0x0a, 0x1a, 0x67, 0x65, 0x74, 0x65, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x46,
	0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x64, 0x4e, 0x75, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x35, 0x0a, 0x1d, 0x67, 0x65, 0x74, 0x46, 0x61, 0x76,
	0x6f, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x31, 0x0a,
	0x0f, 0x69, 0x73, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x73, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65,
	0x22, 0x34, 0x0a, 0x1c, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0xe7, 0x02, 0x0a, 0x04, 0x6c, 0x69, 0x6b, 0x65, 0x12,
	0x5c, 0x0a, 0x15, 0x67, 0x65, 0x74, 0x65, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x46, 0x61, 0x76, 0x6f,
	0x72, 0x69, 0x74, 0x65, 0x64, 0x4e, 0x75, 0x6d, 0x12, 0x21, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e,
	0x67, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x71, 0x1a, 0x20, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x2e, 0x67, 0x65, 0x74, 0x65, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x46, 0x61, 0x76, 0x6f,
	0x72, 0x69, 0x74, 0x65, 0x64, 0x4e, 0x75, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x62, 0x0a,
	0x18, 0x67, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x21, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x2e, 0x67, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x71, 0x1a, 0x23, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x67, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x63, 0x0a, 0x19, 0x67, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x12, 0x22,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x67, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x52,
	0x65, 0x71, 0x1a, 0x22, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x65, 0x74, 0x46, 0x61, 0x76, 0x6f,
	0x72, 0x69, 0x74, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x38, 0x0a, 0x0a, 0x69, 0x73, 0x46, 0x61, 0x76, 0x6f,
	0x72, 0x69, 0x74, 0x65, 0x12, 0x13, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x69, 0x73, 0x46, 0x61,
	0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x15, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x2e, 0x69, 0x73, 0x46, 0x61, 0x76, 0x6f, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x6c, 0x69, 0x6b, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_like_proto_rawDescOnce sync.Once
	file_like_proto_rawDescData = file_like_proto_rawDesc
)

func file_like_proto_rawDescGZIP() []byte {
	file_like_proto_rawDescOnce.Do(func() {
		file_like_proto_rawDescData = protoimpl.X.CompressGZIP(file_like_proto_rawDescData)
	})
	return file_like_proto_rawDescData
}

var file_like_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_like_proto_goTypes = []interface{}{
	(*GeteTotalFavoritedNumReq)(nil),      // 0: user.geteTotalFavoritedNumReq
	(*GetFavoriteCountByUserIdReq)(nil),   // 1: user.getFavoriteCountByUserIdReq
	(*IsFavoriteReq)(nil),                 // 2: user.isFavoriteReq
	(*GetFavoriteCountByVideoIdReq)(nil),  // 3: user.getFavoriteCountByVideoIdReq
	(*GeteTotalFavoritedNumReply)(nil),    // 4: user.geteTotalFavoritedNumReply
	(*GetFavoriteCountByUserIdReply)(nil), // 5: user.getFavoriteCountByUserIdReply
	(*IsFavoriteReply)(nil),               // 6: user.isFavoriteReply
	(*EtFavoriteCountByUserIdReply)(nil),  // 7: user.etFavoriteCountByUserIdReply
}
var file_like_proto_depIdxs = []int32{
	1, // 0: user.like.geteTotalFavoritedNum:input_type -> user.getFavoriteCountByUserIdReq
	1, // 1: user.like.getFavoriteCountByUserId:input_type -> user.getFavoriteCountByUserIdReq
	3, // 2: user.like.getFavoriteCountByVideoId:input_type -> user.getFavoriteCountByVideoIdReq
	2, // 3: user.like.isFavorite:input_type -> user.isFavoriteReq
	4, // 4: user.like.geteTotalFavoritedNum:output_type -> user.geteTotalFavoritedNumReply
	5, // 5: user.like.getFavoriteCountByUserId:output_type -> user.getFavoriteCountByUserIdReply
	7, // 6: user.like.getFavoriteCountByVideoId:output_type -> user.etFavoriteCountByUserIdReply
	6, // 7: user.like.isFavorite:output_type -> user.isFavoriteReply
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_like_proto_init() }
func file_like_proto_init() {
	if File_like_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_like_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GeteTotalFavoritedNumReq); i {
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
		file_like_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFavoriteCountByUserIdReq); i {
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
		file_like_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsFavoriteReq); i {
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
		file_like_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFavoriteCountByVideoIdReq); i {
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
		file_like_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GeteTotalFavoritedNumReply); i {
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
		file_like_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFavoriteCountByUserIdReply); i {
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
		file_like_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IsFavoriteReply); i {
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
		file_like_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EtFavoriteCountByUserIdReply); i {
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
			RawDescriptor: file_like_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_like_proto_goTypes,
		DependencyIndexes: file_like_proto_depIdxs,
		MessageInfos:      file_like_proto_msgTypes,
	}.Build()
	File_like_proto = out.File
	file_like_proto_rawDesc = nil
	file_like_proto_goTypes = nil
	file_like_proto_depIdxs = nil
}
