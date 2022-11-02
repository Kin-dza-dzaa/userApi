package apierror

import (
	"encoding/json"
)

type ErrorStruct struct {
	Message string `json:"message,omitempty"`
	Result  string `json:"result,omitempty"`
	Code    int    `json:"code"`
}

func (errorStruct ErrorStruct) Error() string {
	return errorStruct.Message
}

func (errorStruct *ErrorStruct) Marshal() []byte {
	responseBytes, err := json.Marshal(errorStruct)
	if err != nil {
		return nil
	}
	return responseBytes
}

func NewErrorStruct(Message string, Result string, code int) *ErrorStruct {
	return &ErrorStruct{
		Result:  Result,
		Message: Message,
		Code:    code,
	}
}
