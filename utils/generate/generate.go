package generate

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"io/ioutil"
	"os"
	"strings"
)

type Generate struct {
	Req      *plugin.CodeGeneratorRequest
	Resp     *plugin.CodeGeneratorResponse
	Services map[string][]*LocationService
	Msg      map[string][]*LocationMessage
}

type LocationMessage struct {
	Location        *descriptor.SourceCodeInfo_Location
	Message         *descriptor.DescriptorProto
	LeadingComments []string
}

type LocationService struct {
	ServiceName string                            `json:"service_name"`
	Method      *descriptor.MethodDescriptorProto `json:"method"`
	PackageName string                            `json:"package_name"`
	Path        string                            `json:"path"` // /package.service/method
	ReqName     string                            `json:"req_name"`
	RespName    string                            `json:"resp_name"`
	Req         *LocationMessage                  `json:"req"`
	Resp        *LocationMessage                  `json:"resp"`
	RequestMock map[string]interface{}            `json:"request_mock"`
	Filename    string                            `json:"filename"`
}

func (g *Generate) filterMessages(suffix, filename string) *LocationMessage {
	for _, message := range g.Msg[filename] {
		if strings.HasSuffix(suffix, message.Message.GetName()) {
			return message
		}
	}
	return nil
}

func (g *Generate) filterNestedMsg(suffix string, nestedInfo map[string]interface{}) interface{} {
	for key, ret := range nestedInfo {
		if strings.HasSuffix(suffix, key) {
			return ret
		}
	}
	return nil
}

func New() *Generate {
	g := &Generate{
		Req:      &plugin.CodeGeneratorRequest{},
		Resp:     &plugin.CodeGeneratorResponse{},
		Services: make(map[string][]*LocationService),
		Msg:      make(map[string][]*LocationMessage),
	}
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	if err := proto.Unmarshal(data, g.Req); err != nil {
		panic(err)
	}
	return g
}

func (g *Generate) ReqParameter() string {
	return g.Req.GetParameter()
}

func (g *Generate) Done() {
	for filename, services := range g.Services {
		outfielder := strings.Replace(filename, ".proto", ".meta.json", -1)
		var jsonFile plugin.CodeGeneratorResponse_File
		jsonFile.Name = &outfielder
		bs, _ := json.Marshal(services)
		content := string(bs)
		jsonFile.Content = &content
		g.Resp.File = append(g.Resp.File, &jsonFile)
		os.Stderr.WriteString(fmt.Sprintf("Created Proto Meta File: %s \n", outfielder))

	}
	marshalled, err := proto.Marshal(g.Resp)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(marshalled)

}

func (g *Generate) GenMeta() *Generate {
	g.getLocationServices()

	for filename, services := range g.Services {
		for _, ser := range services {
			mockInfo := make(map[string]interface{})
			nestedInfo := make(map[string]interface{})
			g.loopNested(ser.Req.Message, nestedInfo, filename)
			for _, field := range ser.Req.Message.Field {
				mockInfo[field.GetName()] = field.GetType()
			}
			g.loopField(ser.Req.Message, mockInfo, nestedInfo, filename)

			ser.RequestMock = mockInfo
		}
	}
	return g

}
func (g *Generate) loopNested(message *descriptor.DescriptorProto, nestedInfo map[string]interface{}, filename string) {
	for _, nfield := range message.GetNestedType() {
		ret := make(map[string]interface{})
		if nfield.GetOptions().GetMapEntry() {
			ret[nfield.GetField()[0].GetName()] = nfield.GetField()[1].GetName()
			nestedInfo[nfield.GetName()] = ret
			continue
		}
		g.loopField(nfield, ret, nil, filename)
		nestedInfo[nfield.GetName()] = ret
	}
}
func (g *Generate) loopField(message *descriptor.DescriptorProto, mockInfo, nestedInfo map[string]interface{}, filename string) {
	for _, field := range message.GetField() {
		// 如果是嵌套类型，递归构建mock对象，否则直接构建mock对象kv
		if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			ret := make(map[string]interface{})
			msg := g.filterMessages(field.GetTypeName(), filename)
			if msg == nil {
				mockInfo[field.GetName()] = withDefaultValue(field, nil, nestedInfo)
				continue
			}
			g.loopField(msg.Message, ret, nestedInfo, filename)
			mockInfo[field.GetName()] = withDefaultValue(field, ret, nil)
		} else {
			mockInfo[field.GetName()] = withDefaultValue(field, nil, nil)
		}
	}
}

func (g *Generate) getLocationServices() {
	for _, descriptorProto := range g.Req.ProtoFile {
		filename := descriptorProto.GetName()
		locationMessages := make([]*LocationMessage, 0)
		desc := descriptorProto.GetSourceCodeInfo()
		locations := desc.GetLocation()
		for _, location := range locations {

			if len(location.GetPath()) > 2 {
				continue
			}

			leadingComments := strings.Split(location.GetLeadingComments(), "\n")
			if len(location.GetPath()) > 1 && location.GetPath()[0] == int32(4) {
				message := descriptorProto.GetMessageType()[location.GetPath()[1]]
				locationMessages = append(locationMessages, &LocationMessage{
					Message:  message,
					Location: location,
					// Because we are only parsing messages here at the root level we will not get field comments
					LeadingComments: leadingComments[:len(leadingComments)-1],
				})
			}
		}
		g.Msg[filename] = append(g.Msg[filename], locationMessages...)
		if len(descriptorProto.Service) <= 0 {
			continue
		}
		for _, service := range descriptorProto.Service {
			for _, method := range service.Method {
				ser := &LocationService{
					ServiceName: service.GetName(),
					Method:      method,
					Filename:    filename,
					Path:        fmt.Sprintf("/%s.%s/%s", descriptorProto.GetPackage(), service.GetName(), method.GetName()),
					PackageName: descriptorProto.GetPackage(),
				}
				ser.ReqName = withTpName(method.GetInputType())
				ser.RespName = withTpName(method.GetOutputType())
				ser.Req = g.filterMessages(method.GetInputType(), filename)
				ser.Resp = g.filterMessages(method.GetOutputType(), filename)
				g.Services[filename] = append(g.Services[filename], ser)
			}
		}
	}
}
