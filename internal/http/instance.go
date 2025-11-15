package http

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"io"
	"time"
)

func (c *Client) CreateInstances(instances []proto.Message) error {
	fmt.Printf("\nupdate instances\n")

	url := fmt.Sprintf("http://%v/naming/%v/instances", c.Address, c.Version)

	body, err := JSONFromProtoMessages(instances)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	response, err := c.SendRequest("POST", url, body)
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	ret, err := GetBatchWriteResponse(response)
	if err != nil {
		if err == io.EOF {
			time.Sleep(time.Second)
			return nil
		}

		fmt.Printf("%v\n", err)
		return err
	}
	fmt.Printf("update instance resp:%+v\n", ret)
	if ret != nil {
		fmt.Printf("update instance resp should be nil, %v\n", ret)
		return fmt.Errorf("update instance resp should be nil, %v", ret)
	}
	time.Sleep(time.Second)
	return nil
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
