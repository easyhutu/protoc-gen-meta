package generate

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"strings"
)

func withDefaultValue(field *descriptor.FieldDescriptorProto, dl interface{}, nested map[string]interface{}) interface{} {
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		if field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			if dl != nil {
				return []interface{}{dl}
			}
			spname := strings.Split(field.GetTypeName(), ".")
			return nested[spname[len(spname)-1]]

		}
		return dl

	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return false
	case descriptor.FieldDescriptorProto_TYPE_INT32, descriptor.FieldDescriptorProto_TYPE_INT64:
		return 0
	case descriptor.FieldDescriptorProto_TYPE_UINT32, descriptor.FieldDescriptorProto_TYPE_UINT64:
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
