// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.3
// source: google/ai/generativelanguage/v1beta/prediction_service.proto

package generativelanguagepb

import (
	context "context"
	reflect "reflect"
	sync "sync"

	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request message for
// [PredictionService.Predict][google.ai.generativelanguage.v1beta.PredictionService.Predict].
type PredictRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required. The name of the model for prediction.
	// Format: `name=models/{model}`.
	Model string `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
	// Required. The instances that are the input to the prediction call.
	Instances []*structpb.Value `protobuf:"bytes,2,rep,name=instances,proto3" json:"instances,omitempty"`
	// Optional. The parameters that govern the prediction call.
	Parameters *structpb.Value `protobuf:"bytes,3,opt,name=parameters,proto3" json:"parameters,omitempty"`
}

func (x *PredictRequest) Reset() {
	*x = PredictRequest{}
	mi := &file_google_ai_generativelanguage_v1beta_prediction_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PredictRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PredictRequest) ProtoMessage() {}

func (x *PredictRequest) ProtoReflect() protoreflect.Message {
	mi := &file_google_ai_generativelanguage_v1beta_prediction_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PredictRequest.ProtoReflect.Descriptor instead.
func (*PredictRequest) Descriptor() ([]byte, []int) {
	return file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDescGZIP(), []int{0}
}

func (x *PredictRequest) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *PredictRequest) GetInstances() []*structpb.Value {
	if x != nil {
		return x.Instances
	}
	return nil
}

func (x *PredictRequest) GetParameters() *structpb.Value {
	if x != nil {
		return x.Parameters
	}
	return nil
}

// Response message for [PredictionService.Predict].
type PredictResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The outputs of the prediction call.
	Predictions []*structpb.Value `protobuf:"bytes,1,rep,name=predictions,proto3" json:"predictions,omitempty"`
}

func (x *PredictResponse) Reset() {
	*x = PredictResponse{}
	mi := &file_google_ai_generativelanguage_v1beta_prediction_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PredictResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PredictResponse) ProtoMessage() {}

func (x *PredictResponse) ProtoReflect() protoreflect.Message {
	mi := &file_google_ai_generativelanguage_v1beta_prediction_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PredictResponse.ProtoReflect.Descriptor instead.
func (*PredictResponse) Descriptor() ([]byte, []int) {
	return file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDescGZIP(), []int{1}
}

func (x *PredictResponse) GetPredictions() []*structpb.Value {
	if x != nil {
		return x.Predictions
	}
	return nil
}

var File_google_ai_generativelanguage_v1beta_prediction_service_proto protoreflect.FileDescriptor

var file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDesc = []byte{
	0x0a, 0x3c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x76, 0x65, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x2f, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x23,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x69, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x76, 0x65, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x17, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68,
	0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcf, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x45, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x2f, 0xe0, 0x41, 0x02, 0xfa, 0x41, 0x29, 0x0a, 0x27,
	0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x76, 0x65, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61,
	0x67, 0x65, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x39,
	0x0a, 0x09, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x09,
	0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x3b, 0x0a, 0x0a, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x03, 0xe0, 0x41, 0x01, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x61,
	0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x22, 0x4b, 0x0a, 0x0f, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0b, 0x70, 0x72, 0x65,
	0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0b, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x32, 0xef, 0x01, 0x0a, 0x11, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0xb3, 0x01, 0x0a, 0x07, 0x50, 0x72,
	0x65, 0x64, 0x69, 0x63, 0x74, 0x12, 0x33, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61,
	0x69, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x76, 0x65, 0x6c, 0x61, 0x6e, 0x67,
	0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x2e, 0x50, 0x72, 0x65, 0x64,
	0x69, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x61, 0x69, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x76,
	0x65, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x2e, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x3d, 0xda, 0x41, 0x0f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2c, 0x69, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x73, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x25, 0x3a, 0x01, 0x2a, 0x22, 0x20, 0x2f,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x2f, 0x7b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x3d, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2f, 0x2a, 0x7d, 0x3a, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x1a,
	0x24, 0xca, 0x41, 0x21, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x76, 0x65, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x61, 0x70, 0x69,
	0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x42, 0xa2, 0x01, 0x0a, 0x27, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x61, 0x69, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x76, 0x65, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x42, 0x16, 0x50, 0x72, 0x65, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x5d, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f,
	0x2f, 0x61, 0x69, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x76, 0x65, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x76, 0x65, 0x6c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x70, 0x62, 0x3b, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x69, 0x76, 0x65,
	0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDescOnce sync.Once
	file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDescData = file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDesc
)

func file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDescGZIP() []byte {
	file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDescOnce.Do(func() {
		file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDescData)
	})
	return file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDescData
}

var file_google_ai_generativelanguage_v1beta_prediction_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_google_ai_generativelanguage_v1beta_prediction_service_proto_goTypes = []any{
	(*PredictRequest)(nil),  // 0: google.ai.generativelanguage.v1beta.PredictRequest
	(*PredictResponse)(nil), // 1: google.ai.generativelanguage.v1beta.PredictResponse
	(*structpb.Value)(nil),  // 2: google.protobuf.Value
}
var file_google_ai_generativelanguage_v1beta_prediction_service_proto_depIdxs = []int32{
	2, // 0: google.ai.generativelanguage.v1beta.PredictRequest.instances:type_name -> google.protobuf.Value
	2, // 1: google.ai.generativelanguage.v1beta.PredictRequest.parameters:type_name -> google.protobuf.Value
	2, // 2: google.ai.generativelanguage.v1beta.PredictResponse.predictions:type_name -> google.protobuf.Value
	0, // 3: google.ai.generativelanguage.v1beta.PredictionService.Predict:input_type -> google.ai.generativelanguage.v1beta.PredictRequest
	1, // 4: google.ai.generativelanguage.v1beta.PredictionService.Predict:output_type -> google.ai.generativelanguage.v1beta.PredictResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_google_ai_generativelanguage_v1beta_prediction_service_proto_init() }
func file_google_ai_generativelanguage_v1beta_prediction_service_proto_init() {
	if File_google_ai_generativelanguage_v1beta_prediction_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_google_ai_generativelanguage_v1beta_prediction_service_proto_goTypes,
		DependencyIndexes: file_google_ai_generativelanguage_v1beta_prediction_service_proto_depIdxs,
		MessageInfos:      file_google_ai_generativelanguage_v1beta_prediction_service_proto_msgTypes,
	}.Build()
	File_google_ai_generativelanguage_v1beta_prediction_service_proto = out.File
	file_google_ai_generativelanguage_v1beta_prediction_service_proto_rawDesc = nil
	file_google_ai_generativelanguage_v1beta_prediction_service_proto_goTypes = nil
	file_google_ai_generativelanguage_v1beta_prediction_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PredictionServiceClient is the client API for PredictionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PredictionServiceClient interface {
	// Performs a prediction request.
	Predict(ctx context.Context, in *PredictRequest, opts ...grpc.CallOption) (*PredictResponse, error)
}

type predictionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPredictionServiceClient(cc grpc.ClientConnInterface) PredictionServiceClient {
	return &predictionServiceClient{cc}
}

func (c *predictionServiceClient) Predict(ctx context.Context, in *PredictRequest, opts ...grpc.CallOption) (*PredictResponse, error) {
	out := new(PredictResponse)
	err := c.cc.Invoke(ctx, "/google.ai.generativelanguage.v1beta.PredictionService/Predict", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PredictionServiceServer is the server API for PredictionService service.
type PredictionServiceServer interface {
	// Performs a prediction request.
	Predict(context.Context, *PredictRequest) (*PredictResponse, error)
}

// UnimplementedPredictionServiceServer can be embedded to have forward compatible implementations.
type UnimplementedPredictionServiceServer struct {
}

func (*UnimplementedPredictionServiceServer) Predict(context.Context, *PredictRequest) (*PredictResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Predict not implemented")
}

func RegisterPredictionServiceServer(s *grpc.Server, srv PredictionServiceServer) {
	s.RegisterService(&_PredictionService_serviceDesc, srv)
}

func _PredictionService_Predict_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PredictRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PredictionServiceServer).Predict(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ai.generativelanguage.v1beta.PredictionService/Predict",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PredictionServiceServer).Predict(ctx, req.(*PredictRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _PredictionService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ai.generativelanguage.v1beta.PredictionService",
	HandlerType: (*PredictionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Predict",
			Handler:    _PredictionService_Predict_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ai/generativelanguage/v1beta/prediction_service.proto",
}
