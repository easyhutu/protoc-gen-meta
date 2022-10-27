package generate

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

func withDefaultValue(protoType descriptor.FieldDescriptorProto_Type) interface{} {
	protoType.Type()
	switch protoType {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return false
	case descriptor.FieldDescriptorProto_TYPE_INT32, descriptor.FieldDescriptorProto_TYPE_INT64:
		return 0
	case  descriptor.FieldDescriptorProto_TYPE_UINT32, descriptor.FieldDescriptorProto_TYPE_UINT64:
		return 0
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE, descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return 0.0
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return ""
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return 0
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return []byte("")

	default:
		return ""
	}
}
