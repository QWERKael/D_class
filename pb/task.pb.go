// Code generated by protoc-gen-go. DO NOT EDIT.
// source: task.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type CommonCmdRequest_CmdType int32

const (
	CommonCmdRequest_COMMON_CMD    CommonCmdRequest_CmdType = 0
	CommonCmdRequest_FILE_TRANSFER CommonCmdRequest_CmdType = 1
	CommonCmdRequest_ASYNC_TASK    CommonCmdRequest_CmdType = 2
)

var CommonCmdRequest_CmdType_name = map[int32]string{
	0: "COMMON_CMD",
	1: "FILE_TRANSFER",
	2: "ASYNC_TASK",
}

var CommonCmdRequest_CmdType_value = map[string]int32{
	"COMMON_CMD":    0,
	"FILE_TRANSFER": 1,
	"ASYNC_TASK":    2,
}

func (x CommonCmdRequest_CmdType) String() string {
	return proto.EnumName(CommonCmdRequest_CmdType_name, int32(x))
}

func (CommonCmdRequest_CmdType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{2, 0}
}

type CommonCmdReply_ExeState int32

const (
	CommonCmdReply_Ok  CommonCmdReply_ExeState = 0
	CommonCmdReply_Err CommonCmdReply_ExeState = 1
)

var CommonCmdReply_ExeState_name = map[int32]string{
	0: "Ok",
	1: "Err",
}

var CommonCmdReply_ExeState_value = map[string]int32{
	"Ok":  0,
	"Err": 1,
}

func (x CommonCmdReply_ExeState) String() string {
	return proto.EnumName(CommonCmdReply_ExeState_name, int32(x))
}

func (CommonCmdReply_ExeState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{3, 0}
}

type TransferInfo_TransferType int32

const (
	TransferInfo_Upload   TransferInfo_TransferType = 0
	TransferInfo_Download TransferInfo_TransferType = 1
)

var TransferInfo_TransferType_name = map[int32]string{
	0: "Upload",
	1: "Download",
}

var TransferInfo_TransferType_value = map[string]int32{
	"Upload":   0,
	"Download": 1,
}

func (x TransferInfo_TransferType) String() string {
	return proto.EnumName(TransferInfo_TransferType_name, int32(x))
}

func (TransferInfo_TransferType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{5, 0}
}

type TransferInfo_TransferState int32

const (
	TransferInfo_Apply    TransferInfo_TransferState = 0
	TransferInfo_Begin    TransferInfo_TransferState = 1
	TransferInfo_Complete TransferInfo_TransferState = 2
	TransferInfo_Error    TransferInfo_TransferState = 3
)

var TransferInfo_TransferState_name = map[int32]string{
	0: "Apply",
	1: "Begin",
	2: "Complete",
	3: "Error",
}

var TransferInfo_TransferState_value = map[string]int32{
	"Apply":    0,
	"Begin":    1,
	"Complete": 2,
	"Error":    3,
}

func (x TransferInfo_TransferState) String() string {
	return proto.EnumName(TransferInfo_TransferState_name, int32(x))
}

func (TransferInfo_TransferState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{5, 1}
}

type Row struct {
	Row                  []string `protobuf:"bytes,1,rep,name=row,proto3" json:"row,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Row) Reset()         { *m = Row{} }
func (m *Row) String() string { return proto.CompactTextString(m) }
func (*Row) ProtoMessage()    {}
func (*Row) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{0}
}

func (m *Row) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Row.Unmarshal(m, b)
}
func (m *Row) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Row.Marshal(b, m, deterministic)
}
func (m *Row) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Row.Merge(m, src)
}
func (m *Row) XXX_Size() int {
	return xxx_messageInfo_Row.Size(m)
}
func (m *Row) XXX_DiscardUnknown() {
	xxx_messageInfo_Row.DiscardUnknown(m)
}

var xxx_messageInfo_Row proto.InternalMessageInfo

func (m *Row) GetRow() []string {
	if m != nil {
		return m.Row
	}
	return nil
}

type Table struct {
	Header               *Row     `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	Body                 []*Row   `protobuf:"bytes,2,rep,name=body,proto3" json:"body,omitempty"`
	Footer               *Row     `protobuf:"bytes,3,opt,name=footer,proto3" json:"footer,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Table) Reset()         { *m = Table{} }
func (m *Table) String() string { return proto.CompactTextString(m) }
func (*Table) ProtoMessage()    {}
func (*Table) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{1}
}

func (m *Table) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Table.Unmarshal(m, b)
}
func (m *Table) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Table.Marshal(b, m, deterministic)
}
func (m *Table) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Table.Merge(m, src)
}
func (m *Table) XXX_Size() int {
	return xxx_messageInfo_Table.Size(m)
}
func (m *Table) XXX_DiscardUnknown() {
	xxx_messageInfo_Table.DiscardUnknown(m)
}

var xxx_messageInfo_Table proto.InternalMessageInfo

func (m *Table) GetHeader() *Row {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Table) GetBody() []*Row {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Table) GetFooter() *Row {
	if m != nil {
		return m.Footer
	}
	return nil
}

type CommonCmdRequest struct {
	Type                 CommonCmdRequest_CmdType `protobuf:"varint,1,opt,name=type,proto3,enum=pb.CommonCmdRequest_CmdType" json:"type,omitempty"`
	Plugin               string                   `protobuf:"bytes,2,opt,name=plugin,proto3" json:"plugin,omitempty"`
	Cmd                  string                   `protobuf:"bytes,3,opt,name=cmd,proto3" json:"cmd,omitempty"`
	SubCmd               []string                 `protobuf:"bytes,4,rep,name=sub_cmd,json=subCmd,proto3" json:"sub_cmd,omitempty"`
	Flags                []string                 `protobuf:"bytes,5,rep,name=flags,proto3" json:"flags,omitempty"`
	Args                 map[string]string        `protobuf:"bytes,6,rep,name=args,proto3" json:"args,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *CommonCmdRequest) Reset()         { *m = CommonCmdRequest{} }
func (m *CommonCmdRequest) String() string { return proto.CompactTextString(m) }
func (*CommonCmdRequest) ProtoMessage()    {}
func (*CommonCmdRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{2}
}

func (m *CommonCmdRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonCmdRequest.Unmarshal(m, b)
}
func (m *CommonCmdRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonCmdRequest.Marshal(b, m, deterministic)
}
func (m *CommonCmdRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonCmdRequest.Merge(m, src)
}
func (m *CommonCmdRequest) XXX_Size() int {
	return xxx_messageInfo_CommonCmdRequest.Size(m)
}
func (m *CommonCmdRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonCmdRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CommonCmdRequest proto.InternalMessageInfo

func (m *CommonCmdRequest) GetType() CommonCmdRequest_CmdType {
	if m != nil {
		return m.Type
	}
	return CommonCmdRequest_COMMON_CMD
}

func (m *CommonCmdRequest) GetPlugin() string {
	if m != nil {
		return m.Plugin
	}
	return ""
}

func (m *CommonCmdRequest) GetCmd() string {
	if m != nil {
		return m.Cmd
	}
	return ""
}

func (m *CommonCmdRequest) GetSubCmd() []string {
	if m != nil {
		return m.SubCmd
	}
	return nil
}

func (m *CommonCmdRequest) GetFlags() []string {
	if m != nil {
		return m.Flags
	}
	return nil
}

func (m *CommonCmdRequest) GetArgs() map[string]string {
	if m != nil {
		return m.Args
	}
	return nil
}

type CommonCmdReply struct {
	Status               CommonCmdReply_ExeState `protobuf:"varint,1,opt,name=status,proto3,enum=pb.CommonCmdReply_ExeState" json:"status,omitempty"`
	ResultMsg            string                  `protobuf:"bytes,2,opt,name=result_msg,json=resultMsg,proto3" json:"result_msg,omitempty"`
	ResultTable          *Table                  `protobuf:"bytes,3,opt,name=result_table,json=resultTable,proto3" json:"result_table,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *CommonCmdReply) Reset()         { *m = CommonCmdReply{} }
func (m *CommonCmdReply) String() string { return proto.CompactTextString(m) }
func (*CommonCmdReply) ProtoMessage()    {}
func (*CommonCmdReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{3}
}

func (m *CommonCmdReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CommonCmdReply.Unmarshal(m, b)
}
func (m *CommonCmdReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CommonCmdReply.Marshal(b, m, deterministic)
}
func (m *CommonCmdReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CommonCmdReply.Merge(m, src)
}
func (m *CommonCmdReply) XXX_Size() int {
	return xxx_messageInfo_CommonCmdReply.Size(m)
}
func (m *CommonCmdReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CommonCmdReply.DiscardUnknown(m)
}

var xxx_messageInfo_CommonCmdReply proto.InternalMessageInfo

func (m *CommonCmdReply) GetStatus() CommonCmdReply_ExeState {
	if m != nil {
		return m.Status
	}
	return CommonCmdReply_Ok
}

func (m *CommonCmdReply) GetResultMsg() string {
	if m != nil {
		return m.ResultMsg
	}
	return ""
}

func (m *CommonCmdReply) GetResultTable() *Table {
	if m != nil {
		return m.ResultTable
	}
	return nil
}

type Chunks struct {
	TransferId           int32    `protobuf:"varint,1,opt,name=transfer_id,json=transferId,proto3" json:"transfer_id,omitempty"`
	Size                 int64    `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Content              []byte   `protobuf:"bytes,3,opt,name=Content,proto3" json:"Content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Chunks) Reset()         { *m = Chunks{} }
func (m *Chunks) String() string { return proto.CompactTextString(m) }
func (*Chunks) ProtoMessage()    {}
func (*Chunks) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{4}
}

func (m *Chunks) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Chunks.Unmarshal(m, b)
}
func (m *Chunks) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Chunks.Marshal(b, m, deterministic)
}
func (m *Chunks) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Chunks.Merge(m, src)
}
func (m *Chunks) XXX_Size() int {
	return xxx_messageInfo_Chunks.Size(m)
}
func (m *Chunks) XXX_DiscardUnknown() {
	xxx_messageInfo_Chunks.DiscardUnknown(m)
}

var xxx_messageInfo_Chunks proto.InternalMessageInfo

func (m *Chunks) GetTransferId() int32 {
	if m != nil {
		return m.TransferId
	}
	return 0
}

func (m *Chunks) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *Chunks) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type TransferInfo struct {
	Type                 TransferInfo_TransferType  `protobuf:"varint,1,opt,name=type,proto3,enum=pb.TransferInfo_TransferType" json:"type,omitempty"`
	State                TransferInfo_TransferState `protobuf:"varint,2,opt,name=state,proto3,enum=pb.TransferInfo_TransferState" json:"state,omitempty"`
	FileName             string                     `protobuf:"bytes,3,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	FilePath             string                     `protobuf:"bytes,4,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"`
	TransferId           int32                      `protobuf:"varint,5,opt,name=transfer_id,json=transferId,proto3" json:"transfer_id,omitempty"`
	ErrorMsg             string                     `protobuf:"bytes,6,opt,name=error_msg,json=errorMsg,proto3" json:"error_msg,omitempty"`
	Md5                  string                     `protobuf:"bytes,7,opt,name=md5,proto3" json:"md5,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *TransferInfo) Reset()         { *m = TransferInfo{} }
func (m *TransferInfo) String() string { return proto.CompactTextString(m) }
func (*TransferInfo) ProtoMessage()    {}
func (*TransferInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce5d8dd45b4a91ff, []int{5}
}

func (m *TransferInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransferInfo.Unmarshal(m, b)
}
func (m *TransferInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransferInfo.Marshal(b, m, deterministic)
}
func (m *TransferInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransferInfo.Merge(m, src)
}
func (m *TransferInfo) XXX_Size() int {
	return xxx_messageInfo_TransferInfo.Size(m)
}
func (m *TransferInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_TransferInfo.DiscardUnknown(m)
}

var xxx_messageInfo_TransferInfo proto.InternalMessageInfo

func (m *TransferInfo) GetType() TransferInfo_TransferType {
	if m != nil {
		return m.Type
	}
	return TransferInfo_Upload
}

func (m *TransferInfo) GetState() TransferInfo_TransferState {
	if m != nil {
		return m.State
	}
	return TransferInfo_Apply
}

func (m *TransferInfo) GetFileName() string {
	if m != nil {
		return m.FileName
	}
	return ""
}

func (m *TransferInfo) GetFilePath() string {
	if m != nil {
		return m.FilePath
	}
	return ""
}

func (m *TransferInfo) GetTransferId() int32 {
	if m != nil {
		return m.TransferId
	}
	return 0
}

func (m *TransferInfo) GetErrorMsg() string {
	if m != nil {
		return m.ErrorMsg
	}
	return ""
}

func (m *TransferInfo) GetMd5() string {
	if m != nil {
		return m.Md5
	}
	return ""
}

func init() {
	proto.RegisterEnum("pb.CommonCmdRequest_CmdType", CommonCmdRequest_CmdType_name, CommonCmdRequest_CmdType_value)
	proto.RegisterEnum("pb.CommonCmdReply_ExeState", CommonCmdReply_ExeState_name, CommonCmdReply_ExeState_value)
	proto.RegisterEnum("pb.TransferInfo_TransferType", TransferInfo_TransferType_name, TransferInfo_TransferType_value)
	proto.RegisterEnum("pb.TransferInfo_TransferState", TransferInfo_TransferState_name, TransferInfo_TransferState_value)
	proto.RegisterType((*Row)(nil), "pb.Row")
	proto.RegisterType((*Table)(nil), "pb.Table")
	proto.RegisterType((*CommonCmdRequest)(nil), "pb.CommonCmdRequest")
	proto.RegisterMapType((map[string]string)(nil), "pb.CommonCmdRequest.ArgsEntry")
	proto.RegisterType((*CommonCmdReply)(nil), "pb.CommonCmdReply")
	proto.RegisterType((*Chunks)(nil), "pb.Chunks")
	proto.RegisterType((*TransferInfo)(nil), "pb.TransferInfo")
}

func init() { proto.RegisterFile("task.proto", fileDescriptor_ce5d8dd45b4a91ff) }

var fileDescriptor_ce5d8dd45b4a91ff = []byte{
	// 701 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x54, 0xdd, 0x6e, 0xd3, 0x50,
	0x0c, 0x6e, 0x92, 0x36, 0x5d, 0xbc, 0xae, 0x0a, 0xd6, 0xc4, 0xa2, 0x95, 0xb1, 0x2a, 0x57, 0x15,
	0x42, 0x15, 0x74, 0x4c, 0x43, 0x08, 0x21, 0x95, 0xac, 0x93, 0x26, 0x68, 0x87, 0x4e, 0x8b, 0x10,
	0x57, 0x55, 0x4a, 0x4e, 0x7f, 0xd4, 0x24, 0x27, 0x24, 0x27, 0x94, 0xf0, 0x42, 0x5c, 0xf2, 0x1a,
	0xbc, 0x0b, 0x2f, 0x81, 0xce, 0x49, 0xba, 0x75, 0x5d, 0xef, 0xec, 0xcf, 0x76, 0xfc, 0xd9, 0xfe,
	0x72, 0x00, 0xb8, 0x9b, 0x2c, 0xdb, 0x51, 0xcc, 0x38, 0x43, 0x35, 0x9a, 0xd8, 0x47, 0xa0, 0x11,
	0xb6, 0x42, 0x13, 0xb4, 0x98, 0xad, 0x2c, 0xa5, 0xa9, 0xb5, 0x0c, 0x22, 0x4c, 0x7b, 0x0a, 0x95,
	0x91, 0x3b, 0xf1, 0x29, 0x9e, 0x82, 0x3e, 0xa7, 0xae, 0x47, 0x63, 0x4b, 0x69, 0x2a, 0xad, 0xfd,
	0x4e, 0xb5, 0x1d, 0x4d, 0xda, 0x84, 0xad, 0x48, 0x01, 0x63, 0x03, 0xca, 0x13, 0xe6, 0x65, 0x96,
	0xda, 0xd4, 0x36, 0xc3, 0x12, 0x14, 0xd5, 0x53, 0xc6, 0x38, 0x8d, 0x2d, 0x6d, 0xab, 0x3a, 0x87,
	0xed, 0xbf, 0x2a, 0x98, 0x0e, 0x0b, 0x02, 0x16, 0x3a, 0x81, 0x47, 0xe8, 0xf7, 0x94, 0x26, 0x1c,
	0x5f, 0x40, 0x99, 0x67, 0x11, 0x95, 0x1d, 0xeb, 0x9d, 0x27, 0xa2, 0x66, 0x3b, 0xa7, 0xed, 0x04,
	0xde, 0x28, 0x8b, 0x28, 0x91, 0x99, 0xf8, 0x18, 0xf4, 0xc8, 0x4f, 0x67, 0x8b, 0xd0, 0x52, 0x9b,
	0x4a, 0xcb, 0x20, 0x85, 0x27, 0x06, 0xfb, 0x16, 0x78, 0xb2, 0xb9, 0x41, 0x84, 0x89, 0x47, 0x50,
	0x4d, 0xd2, 0xc9, 0x58, 0xa0, 0x65, 0x39, 0xae, 0x9e, 0xa4, 0x13, 0x27, 0xf0, 0xf0, 0x10, 0x2a,
	0x53, 0xdf, 0x9d, 0x25, 0x56, 0x45, 0xc2, 0xb9, 0x83, 0x1d, 0x28, 0xbb, 0xf1, 0x2c, 0xb1, 0x74,
	0x39, 0xdd, 0xd3, 0x9d, 0x54, 0xba, 0xf1, 0x2c, 0xe9, 0x85, 0x3c, 0xce, 0x88, 0xcc, 0x3d, 0xbe,
	0x00, 0xe3, 0x16, 0x12, 0x0c, 0x96, 0x34, 0x93, 0xa3, 0x18, 0x44, 0x98, 0xa2, 0xd1, 0x0f, 0xd7,
	0x4f, 0x69, 0x41, 0x35, 0x77, 0xde, 0xa8, 0xaf, 0x15, 0xfb, 0x2d, 0x54, 0x8b, 0xb1, 0xb0, 0x0e,
	0xe0, 0xdc, 0xf4, 0xfb, 0x37, 0x83, 0xb1, 0xd3, 0xbf, 0x34, 0x4b, 0xf8, 0x08, 0x0e, 0xae, 0xae,
	0x3f, 0xf6, 0xc6, 0x23, 0xd2, 0x1d, 0x0c, 0xaf, 0x7a, 0xc4, 0x54, 0x44, 0x4a, 0x77, 0xf8, 0x75,
	0xe0, 0x8c, 0x47, 0xdd, 0xe1, 0x07, 0x53, 0xb5, 0xff, 0x28, 0x50, 0xdf, 0xe0, 0x16, 0xf9, 0x19,
	0x9e, 0x81, 0x9e, 0x70, 0x97, 0xa7, 0x49, 0xb1, 0xca, 0xc6, 0x16, 0xff, 0xc8, 0xcf, 0xda, 0xbd,
	0x9f, 0x74, 0xc8, 0x5d, 0x4e, 0x49, 0x91, 0x8a, 0x27, 0x00, 0x31, 0x4d, 0x52, 0x9f, 0x8f, 0x83,
	0x64, 0x56, 0x90, 0x34, 0x72, 0xa4, 0x9f, 0xcc, 0xf0, 0x39, 0xd4, 0x8a, 0x30, 0x17, 0x02, 0x29,
	0x0e, 0x6b, 0x88, 0x2f, 0x4b, 0xc5, 0x90, 0xfd, 0x3c, 0x2c, 0x1d, 0xbb, 0x01, 0x7b, 0xeb, 0x06,
	0xa8, 0x83, 0x7a, 0xb3, 0x34, 0x4b, 0x58, 0x05, 0xad, 0x17, 0xc7, 0xa6, 0x62, 0x7f, 0x01, 0xdd,
	0x99, 0xa7, 0xe1, 0x32, 0xc1, 0x53, 0xd8, 0xe7, 0xb1, 0x1b, 0x26, 0x53, 0x1a, 0x8f, 0x17, 0x9e,
	0x64, 0x5b, 0x21, 0xb0, 0x86, 0xae, 0x3d, 0x44, 0x28, 0x27, 0x8b, 0x5f, 0xf9, 0xce, 0x34, 0x22,
	0x6d, 0xb4, 0xa0, 0xea, 0xb0, 0x90, 0xd3, 0x90, 0x4b, 0x12, 0x35, 0xb2, 0x76, 0xed, 0x7f, 0x2a,
	0xd4, 0x46, 0xeb, 0xe2, 0x70, 0xca, 0xf0, 0xe5, 0x3d, 0x45, 0x9d, 0x48, 0xb2, 0x1b, 0xf1, 0x5b,
	0x67, 0x43, 0x52, 0xaf, 0xa0, 0x22, 0x16, 0x92, 0xb7, 0xac, 0xe7, 0xa7, 0xdf, 0x59, 0x93, 0x6f,
	0x2f, 0x4f, 0xc6, 0x06, 0x18, 0xd3, 0x85, 0x4f, 0xc7, 0xa1, 0x1b, 0xd0, 0x42, 0x76, 0x7b, 0x02,
	0x18, 0xb8, 0xc1, 0x5d, 0x30, 0x72, 0xf9, 0xdc, 0x2a, 0xdf, 0x05, 0x3f, 0xb9, 0x7c, 0xbe, 0xbd,
	0x82, 0xca, 0x83, 0x15, 0x34, 0xc0, 0xa0, 0x71, 0xcc, 0x62, 0x79, 0x16, 0x3d, 0xaf, 0x96, 0x80,
	0xb8, 0x8a, 0x09, 0x5a, 0xe0, 0x9d, 0x5b, 0xd5, 0x5c, 0x66, 0x81, 0x77, 0x6e, 0xb7, 0xee, 0x56,
	0x20, 0x15, 0x05, 0xa0, 0x7f, 0x8e, 0x7c, 0xe6, 0x7a, 0x66, 0x09, 0x6b, 0xb0, 0x77, 0xc9, 0x56,
	0xa1, 0xf4, 0x14, 0xfb, 0x1d, 0x1c, 0xdc, 0x9b, 0x05, 0x0d, 0xa8, 0x74, 0xa3, 0xc8, 0xcf, 0xcc,
	0x92, 0x30, 0xdf, 0xd3, 0xd9, 0x22, 0x34, 0x15, 0x51, 0xe4, 0xb0, 0x20, 0xf2, 0x29, 0xa7, 0xa6,
	0x2a, 0x02, 0x3d, 0xd1, 0xdc, 0xd4, 0x3a, 0xbf, 0x15, 0x30, 0x84, 0xa8, 0xdc, 0x50, 0xbc, 0x07,
	0x17, 0xb9, 0x23, 0x15, 0x86, 0x87, 0xbb, 0x7e, 0x98, 0x63, 0x7c, 0x28, 0x43, 0xbb, 0x84, 0xe7,
	0x70, 0x20, 0xbb, 0xae, 0xb9, 0xa0, 0xb9, 0xbd, 0xf2, 0xe3, 0x07, 0x88, 0x5d, 0xc2, 0x67, 0xeb,
	0xb9, 0x10, 0xe4, 0x67, 0xa5, 0xa0, 0x76, 0x65, 0xb6, 0x94, 0x89, 0x2e, 0x5f, 0xbe, 0xb3, 0xff,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x99, 0x11, 0xe0, 0x60, 0x07, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CommanderClient is the client API for Commander service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CommanderClient interface {
	CommonCmd(ctx context.Context, in *CommonCmdRequest, opts ...grpc.CallOption) (*CommonCmdReply, error)
	ApplyTransfer(ctx context.Context, in *TransferInfo, opts ...grpc.CallOption) (*TransferInfo, error)
	Upload(ctx context.Context, opts ...grpc.CallOption) (Commander_UploadClient, error)
}

type commanderClient struct {
	cc *grpc.ClientConn
}

func NewCommanderClient(cc *grpc.ClientConn) CommanderClient {
	return &commanderClient{cc}
}

func (c *commanderClient) CommonCmd(ctx context.Context, in *CommonCmdRequest, opts ...grpc.CallOption) (*CommonCmdReply, error) {
	out := new(CommonCmdReply)
	err := c.cc.Invoke(ctx, "/pb.Commander/CommonCmd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commanderClient) ApplyTransfer(ctx context.Context, in *TransferInfo, opts ...grpc.CallOption) (*TransferInfo, error) {
	out := new(TransferInfo)
	err := c.cc.Invoke(ctx, "/pb.Commander/ApplyTransfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commanderClient) Upload(ctx context.Context, opts ...grpc.CallOption) (Commander_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Commander_serviceDesc.Streams[0], "/pb.Commander/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &commanderUploadClient{stream}
	return x, nil
}

type Commander_UploadClient interface {
	Send(*Chunks) error
	CloseAndRecv() (*TransferInfo, error)
	grpc.ClientStream
}

type commanderUploadClient struct {
	grpc.ClientStream
}

func (x *commanderUploadClient) Send(m *Chunks) error {
	return x.ClientStream.SendMsg(m)
}

func (x *commanderUploadClient) CloseAndRecv() (*TransferInfo, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(TransferInfo)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CommanderServer is the server API for Commander service.
type CommanderServer interface {
	CommonCmd(context.Context, *CommonCmdRequest) (*CommonCmdReply, error)
	ApplyTransfer(context.Context, *TransferInfo) (*TransferInfo, error)
	Upload(Commander_UploadServer) error
}

// UnimplementedCommanderServer can be embedded to have forward compatible implementations.
type UnimplementedCommanderServer struct {
}

func (*UnimplementedCommanderServer) CommonCmd(ctx context.Context, req *CommonCmdRequest) (*CommonCmdReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommonCmd not implemented")
}
func (*UnimplementedCommanderServer) ApplyTransfer(ctx context.Context, req *TransferInfo) (*TransferInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApplyTransfer not implemented")
}
func (*UnimplementedCommanderServer) Upload(srv Commander_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}

func RegisterCommanderServer(s *grpc.Server, srv CommanderServer) {
	s.RegisterService(&_Commander_serviceDesc, srv)
}

func _Commander_CommonCmd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommonCmdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommanderServer).CommonCmd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Commander/CommonCmd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommanderServer).CommonCmd(ctx, req.(*CommonCmdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Commander_ApplyTransfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommanderServer).ApplyTransfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Commander/ApplyTransfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommanderServer).ApplyTransfer(ctx, req.(*TransferInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Commander_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CommanderServer).Upload(&commanderUploadServer{stream})
}

type Commander_UploadServer interface {
	SendAndClose(*TransferInfo) error
	Recv() (*Chunks, error)
	grpc.ServerStream
}

type commanderUploadServer struct {
	grpc.ServerStream
}

func (x *commanderUploadServer) SendAndClose(m *TransferInfo) error {
	return x.ServerStream.SendMsg(m)
}

func (x *commanderUploadServer) Recv() (*Chunks, error) {
	m := new(Chunks)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Commander_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Commander",
	HandlerType: (*CommanderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CommonCmd",
			Handler:    _Commander_CommonCmd_Handler,
		},
		{
			MethodName: "ApplyTransfer",
			Handler:    _Commander_ApplyTransfer_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _Commander_Upload_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "task.proto",
}
