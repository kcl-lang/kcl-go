package service

import (
	"errors"

	"github.com/golang/protobuf/proto"
)

// CallRestMethod call an restful method.
func CallRestMethod(host, method string, input, output proto.Message) error {
	var result RestfulResult
	result.Result = output
	if err := httpPost(host+"/api:protorpc/"+method, input, &result); err != nil {
		return err
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	return nil
}
