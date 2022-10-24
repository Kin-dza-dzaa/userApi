package apierror

import (
	"encoding/json"
)

type ErrorStruct struct {
	Message          string `json:"message,omitempty"`
	Result           string `json:"result,omitempty"`
}

func (errorStruct *ErrorStruct) Marshal() []byte {
	responseBytes, err := json.Marshal(errorStruct)
	if err != nil {
		return nil
	}
	return responseBytes
}

func NewErrorStruct(Message, Result string) *ErrorStruct {
	return &ErrorStruct{
		Result:             Result,
		Message:          	Message,
	}
}