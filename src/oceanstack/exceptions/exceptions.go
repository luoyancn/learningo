package exceptions

import (
	"fmt"
)

type JsonMarshallException struct {
	msg string
}

func (err JsonMarshallException) Error() string {
	return fmt.Sprintf("Cannot marshall the object to json:%s\n", err.msg)
}

func NewJsonMarshallException(msg string) JsonMarshallException {
	return JsonMarshallException{
		msg: msg,
	}
}

type InvalidJsonException struct {
	msg string
}

func (err InvalidJsonException) Error() string {
	return fmt.Sprintf("Invalid Json Body:%s\n", err.msg)
}

func NewInvalidJsonException(msg string) InvalidJsonException {
	return InvalidJsonException{
		msg: msg,
	}
}
