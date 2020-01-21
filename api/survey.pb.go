// Code generated by protoc-gen-go. DO NOT EDIT.
// source: survey.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SurveyItem struct {
	Id          string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Key         string      `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Follows     []string    `protobuf:"bytes,3,rep,name=follows,proto3" json:"follows,omitempty"`
	Condition   *Expression `protobuf:"bytes,4,opt,name=condition,proto3" json:"condition,omitempty"`
	Priority    float32     `protobuf:"fixed32,5,opt,name=priority,proto3" json:"priority,omitempty"`
	Version     int32       `protobuf:"varint,6,opt,name=version,proto3" json:"version,omitempty"`
	VersionTags []string    `protobuf:"bytes,7,rep,name=versionTags,proto3" json:"versionTags,omitempty"`
	// Question group attributes ->
	Items           []*SurveyItem `protobuf:"bytes,8,rep,name=items,proto3" json:"items,omitempty"`
	SelectionMethod *Expression   `protobuf:"bytes,9,opt,name=selectionMethod,proto3" json:"selectionMethod,omitempty"`
	// Question attributes ->
	Type                 string   `protobuf:"bytes,10,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SurveyItem) Reset()         { *m = SurveyItem{} }
func (m *SurveyItem) String() string { return proto.CompactTextString(m) }
func (*SurveyItem) ProtoMessage()    {}
func (*SurveyItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_a40f94eaa8e6ca46, []int{0}
}

func (m *SurveyItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SurveyItem.Unmarshal(m, b)
}
func (m *SurveyItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SurveyItem.Marshal(b, m, deterministic)
}
func (m *SurveyItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SurveyItem.Merge(m, src)
}
func (m *SurveyItem) XXX_Size() int {
	return xxx_messageInfo_SurveyItem.Size(m)
}
func (m *SurveyItem) XXX_DiscardUnknown() {
	xxx_messageInfo_SurveyItem.DiscardUnknown(m)
}

var xxx_messageInfo_SurveyItem proto.InternalMessageInfo

func (m *SurveyItem) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SurveyItem) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *SurveyItem) GetFollows() []string {
	if m != nil {
		return m.Follows
	}
	return nil
}

func (m *SurveyItem) GetCondition() *Expression {
	if m != nil {
		return m.Condition
	}
	return nil
}

func (m *SurveyItem) GetPriority() float32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *SurveyItem) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *SurveyItem) GetVersionTags() []string {
	if m != nil {
		return m.VersionTags
	}
	return nil
}

func (m *SurveyItem) GetItems() []*SurveyItem {
	if m != nil {
		return m.Items
	}
	return nil
}

func (m *SurveyItem) GetSelectionMethod() *Expression {
	if m != nil {
		return m.SelectionMethod
	}
	return nil
}

func (m *SurveyItem) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type Expression struct {
	Name                 string           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Dtype                string           `protobuf:"bytes,2,opt,name=dtype,proto3" json:"dtype,omitempty"`
	Data                 []*ExpressionArg `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Expression) Reset()         { *m = Expression{} }
func (m *Expression) String() string { return proto.CompactTextString(m) }
func (*Expression) ProtoMessage()    {}
func (*Expression) Descriptor() ([]byte, []int) {
	return fileDescriptor_a40f94eaa8e6ca46, []int{1}
}

func (m *Expression) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Expression.Unmarshal(m, b)
}
func (m *Expression) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Expression.Marshal(b, m, deterministic)
}
func (m *Expression) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Expression.Merge(m, src)
}
func (m *Expression) XXX_Size() int {
	return xxx_messageInfo_Expression.Size(m)
}
func (m *Expression) XXX_DiscardUnknown() {
	xxx_messageInfo_Expression.DiscardUnknown(m)
}

var xxx_messageInfo_Expression proto.InternalMessageInfo

func (m *Expression) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Expression) GetDtype() string {
	if m != nil {
		return m.Dtype
	}
	return ""
}

func (m *Expression) GetData() []*ExpressionArg {
	if m != nil {
		return m.Data
	}
	return nil
}

type ExpressionArg struct {
	// Types that are valid to be assigned to Data:
	//	*ExpressionArg_Exp
	//	*ExpressionArg_Str
	//	*ExpressionArg_Num
	Data                 isExpressionArg_Data `protobuf_oneof:"data"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ExpressionArg) Reset()         { *m = ExpressionArg{} }
func (m *ExpressionArg) String() string { return proto.CompactTextString(m) }
func (*ExpressionArg) ProtoMessage()    {}
func (*ExpressionArg) Descriptor() ([]byte, []int) {
	return fileDescriptor_a40f94eaa8e6ca46, []int{2}
}

func (m *ExpressionArg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExpressionArg.Unmarshal(m, b)
}
func (m *ExpressionArg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExpressionArg.Marshal(b, m, deterministic)
}
func (m *ExpressionArg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExpressionArg.Merge(m, src)
}
func (m *ExpressionArg) XXX_Size() int {
	return xxx_messageInfo_ExpressionArg.Size(m)
}
func (m *ExpressionArg) XXX_DiscardUnknown() {
	xxx_messageInfo_ExpressionArg.DiscardUnknown(m)
}

var xxx_messageInfo_ExpressionArg proto.InternalMessageInfo

type isExpressionArg_Data interface {
	isExpressionArg_Data()
}

type ExpressionArg_Exp struct {
	Exp *Expression `protobuf:"bytes,1,opt,name=exp,proto3,oneof"`
}

type ExpressionArg_Str struct {
	Str string `protobuf:"bytes,2,opt,name=str,proto3,oneof"`
}

type ExpressionArg_Num struct {
	Num float64 `protobuf:"fixed64,3,opt,name=num,proto3,oneof"`
}

func (*ExpressionArg_Exp) isExpressionArg_Data() {}

func (*ExpressionArg_Str) isExpressionArg_Data() {}

func (*ExpressionArg_Num) isExpressionArg_Data() {}

func (m *ExpressionArg) GetData() isExpressionArg_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *ExpressionArg) GetExp() *Expression {
	if x, ok := m.GetData().(*ExpressionArg_Exp); ok {
		return x.Exp
	}
	return nil
}

func (m *ExpressionArg) GetStr() string {
	if x, ok := m.GetData().(*ExpressionArg_Str); ok {
		return x.Str
	}
	return ""
}

func (m *ExpressionArg) GetNum() float64 {
	if x, ok := m.GetData().(*ExpressionArg_Num); ok {
		return x.Num
	}
	return 0
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*ExpressionArg) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ExpressionArg_Exp)(nil),
		(*ExpressionArg_Str)(nil),
		(*ExpressionArg_Num)(nil),
	}
}

type ResponseItem struct {
	Key                  string          `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string          `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Items                []*ResponseItem `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
	Dtype                string          `protobuf:"bytes,4,opt,name=dtype,proto3" json:"dtype,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ResponseItem) Reset()         { *m = ResponseItem{} }
func (m *ResponseItem) String() string { return proto.CompactTextString(m) }
func (*ResponseItem) ProtoMessage()    {}
func (*ResponseItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_a40f94eaa8e6ca46, []int{3}
}

func (m *ResponseItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseItem.Unmarshal(m, b)
}
func (m *ResponseItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseItem.Marshal(b, m, deterministic)
}
func (m *ResponseItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseItem.Merge(m, src)
}
func (m *ResponseItem) XXX_Size() int {
	return xxx_messageInfo_ResponseItem.Size(m)
}
func (m *ResponseItem) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseItem.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseItem proto.InternalMessageInfo

func (m *ResponseItem) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *ResponseItem) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *ResponseItem) GetItems() []*ResponseItem {
	if m != nil {
		return m.Items
	}
	return nil
}

func (m *ResponseItem) GetDtype() string {
	if m != nil {
		return m.Dtype
	}
	return ""
}

func init() {
	proto.RegisterType((*SurveyItem)(nil), "influenzanet.survey.SurveyItem")
	proto.RegisterType((*Expression)(nil), "influenzanet.survey.Expression")
	proto.RegisterType((*ExpressionArg)(nil), "influenzanet.survey.ExpressionArg")
	proto.RegisterType((*ResponseItem)(nil), "influenzanet.survey.ResponseItem")
}

func init() { proto.RegisterFile("survey.proto", fileDescriptor_a40f94eaa8e6ca46) }

var fileDescriptor_a40f94eaa8e6ca46 = []byte{
	// 401 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x3d, 0x73, 0xd4, 0x30,
	0x10, 0x45, 0xd6, 0xf9, 0x92, 0xdb, 0x0b, 0x1f, 0xb3, 0xa4, 0xd0, 0xd0, 0x20, 0x5c, 0xb9, 0xba,
	0x22, 0x19, 0xa0, 0xa2, 0x20, 0x33, 0xcc, 0x90, 0x82, 0x46, 0x50, 0xd1, 0x99, 0x78, 0x13, 0x34,
	0xd8, 0x92, 0x47, 0x92, 0x8f, 0x98, 0x9e, 0x9f, 0x4c, 0xcf, 0x48, 0x3e, 0x9f, 0x0f, 0x26, 0x4c,
	0xe8, 0xf6, 0xbd, 0xdb, 0x77, 0xfb, 0xf6, 0x69, 0x0d, 0x27, 0xbe, 0x77, 0x5b, 0x1a, 0x36, 0x9d,
	0xb3, 0xc1, 0xe2, 0x53, 0x6d, 0xae, 0x9b, 0x9e, 0xcc, 0x8f, 0xca, 0x50, 0xd8, 0x8c, 0x3f, 0x15,
	0xbf, 0x32, 0x80, 0x8f, 0xa9, 0xbc, 0x0c, 0xd4, 0xe2, 0x23, 0xc8, 0x74, 0x2d, 0x98, 0x64, 0xe5,
	0x4a, 0x65, 0xba, 0xc6, 0x27, 0xc0, 0xbf, 0xd1, 0x20, 0xb2, 0x44, 0xc4, 0x12, 0x05, 0x1c, 0x5d,
	0xdb, 0xa6, 0xb1, 0xdf, 0xbd, 0xe0, 0x92, 0x97, 0x2b, 0x35, 0x41, 0x7c, 0x03, 0xab, 0x2b, 0x6b,
	0x6a, 0x1d, 0xb4, 0x35, 0x62, 0x21, 0x59, 0xb9, 0x3e, 0x7b, 0xbe, 0xb9, 0x63, 0xe6, 0xe6, 0xdd,
	0x6d, 0xe7, 0xc8, 0x7b, 0x6d, 0x8d, 0x9a, 0x15, 0xf8, 0x0c, 0x8e, 0x3b, 0xa7, 0xad, 0xd3, 0x61,
	0x10, 0xb9, 0x64, 0x65, 0xa6, 0xf6, 0x38, 0x0e, 0xdd, 0x92, 0x8b, 0x0a, 0xb1, 0x94, 0xac, 0xcc,
	0xd5, 0x04, 0x51, 0xc2, 0x7a, 0x57, 0x7e, 0xaa, 0x6e, 0xbc, 0x38, 0x4a, 0x96, 0x0e, 0x29, 0x7c,
	0x09, 0xb9, 0x0e, 0xd4, 0x7a, 0x71, 0x2c, 0xf9, 0x3f, 0x2d, 0xcd, 0x11, 0xa8, 0xb1, 0x1b, 0x2f,
	0xe1, 0xb1, 0xa7, 0x86, 0xae, 0xa2, 0xb7, 0x0f, 0x14, 0xbe, 0xda, 0x5a, 0xac, 0xfe, 0x6f, 0xa7,
	0xbf, 0x75, 0x88, 0xb0, 0x08, 0x43, 0x47, 0x02, 0x52, 0x8a, 0xa9, 0x2e, 0x0c, 0xc0, 0x2c, 0x89,
	0x1d, 0xa6, 0x6a, 0x69, 0x17, 0x7c, 0xaa, 0xf1, 0x14, 0xf2, 0x3a, 0xc9, 0xc6, 0xf0, 0x47, 0x80,
	0xaf, 0x60, 0x51, 0x57, 0xa1, 0x4a, 0xd9, 0xaf, 0xcf, 0x8a, 0x7b, 0xbc, 0xbc, 0x75, 0x37, 0x2a,
	0xf5, 0x17, 0x01, 0x1e, 0xfe, 0x41, 0xe3, 0x39, 0x70, 0xba, 0xed, 0xd2, 0xc4, 0xfb, 0x77, 0x7a,
	0xff, 0x40, 0xc5, 0x6e, 0x44, 0xe0, 0x3e, 0xb8, 0xd1, 0x51, 0xe4, 0x7c, 0x70, 0x91, 0x33, 0x7d,
	0x2b, 0xb8, 0x64, 0x25, 0x8b, 0x9c, 0xe9, 0xdb, 0x8b, 0xe5, 0xe8, 0xb2, 0xf8, 0xc9, 0xe0, 0x44,
	0x91, 0xef, 0xac, 0xf1, 0x94, 0xee, 0x6b, 0x77, 0x4f, 0x6c, 0xbe, 0xa7, 0x53, 0xc8, 0xb7, 0x55,
	0xd3, 0xef, 0xd7, 0x4c, 0x00, 0x5f, 0x4f, 0x8f, 0x36, 0xee, 0xf9, 0xe2, 0x4e, 0x7f, 0x87, 0xff,
	0x3c, 0x3d, 0xdb, 0x3e, 0xb5, 0xc5, 0x41, 0x6a, 0x17, 0xf9, 0x67, 0x5e, 0x75, 0xfa, 0xcb, 0x32,
	0x7d, 0x08, 0xe7, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x9f, 0x4b, 0x11, 0x1e, 0x18, 0x03, 0x00,
	0x00,
}
