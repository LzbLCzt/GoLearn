package test

import (
	"bytes"
	"flag"
	apiV1Model "git.woa.com/polaris/polaris-server-api/api/v1/model"
	"github.com/golang/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	"testing"
)

func TestProtoToJson(t *testing.T) {
	defer glog.Flush()

	// 设置glog标志，强制输出到stderr
	flag.Set("logtostderr", "true")
	flag.Set("stderrthreshold", "INFO")
	flag.Parse()

	instances := CreateInstance()
	glog.Infof("instances: %v", instances)

	// 将[]*apiV1Model.Instance转换为[]proto.Message
	var protoMessages []proto.Message
	for _, instance := range instances {
		protoMessages = append(protoMessages, instance)
	}

	buffer, err := JSONFromProtoMessages(protoMessages)
	if err != nil {
		glog.Errorf("JSONFromProtoMessages error: %v", err)
		return
	}
	glog.Infof("buffer: %s", buffer.String())
}

func CreateInstance() []*apiV1Model.Instance {
	var instances []*apiV1Model.Instance
	instances = append(instances, &apiV1Model.Instance{
		Service:   &wrappers.StringValue{Value: "polaris.test.grpc"},
		Namespace: &wrappers.StringValue{Value: "Test"},
		Host:      &wrappers.StringValue{Value: "127.0.0.1"},
		Port:      &wrappers.UInt32Value{Value: 8080},
		Protocol:  &wrappers.StringValue{Value: "grpc"},
	})

	instances = append(instances, &apiV1Model.Instance{
		Service:   &wrappers.StringValue{Value: "polaris.test.grpc"},
		Namespace: &wrappers.StringValue{Value: "Test2"},
		Host:      &wrappers.StringValue{Value: "127.0.0.2"},
		Port:      &wrappers.UInt32Value{Value: 8080},
		Protocol:  &wrappers.StringValue{Value: "grpc"},
		Metadata: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	})
	return instances
}

// JSONFromProtoMessages 通用protobuf消息数组转JSON
func JSONFromProtoMessages(messages []proto.Message) (*bytes.Buffer, error) {
	marshaler := jsonpb.Marshaler{Indent: " "}

	buffer := bytes.NewBuffer([]byte{})
	buffer.Write([]byte("["))

	for i, msg := range messages {
		if i > 0 {
			buffer.Write([]byte(",\n"))
		}

		msgBuffer := bytes.NewBuffer([]byte{})
		err := marshaler.Marshal(msgBuffer, msg)
		if err != nil {
			return nil, err
		}
		buffer.Write(msgBuffer.Bytes())
	}

	buffer.Write([]byte("]"))
	return buffer, nil
}
