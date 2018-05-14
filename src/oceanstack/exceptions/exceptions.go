package exceptions

import (
	"fmt"
)

type Exception struct {
	Msg string
}

func (err Exception) Error() string {
	return fmt.Sprintf("Exception:%s\n", err.Msg)
}

func NewException(msg string) Exception {
	return Exception{
		Msg: msg,
	}
}

/*****************************************************************************/
type JsonMarshallException struct {
	Exception
}

func (err JsonMarshallException) Error() string {
	return fmt.Sprintf("Cannot marshall the object to json:%s\n", err.Msg)
}

func NewJsonMarshallException(msg string) JsonMarshallException {
	err := JsonMarshallException{}
	err.Msg = msg
	return err
}

/*****************************************************************************/

type InvalidJsonException struct {
	Exception
}

func (err InvalidJsonException) Error() string {
	return fmt.Sprintf("Invalid Json Body:%s\n", err.Msg)
}

func NewInvalidJsonException(msg string) InvalidJsonException {
	err := InvalidJsonException{}
	err.Msg = msg
	return err
}

/*****************************************************************************/

type NotFoundException struct {
	Exception
}

func (err NotFoundException) Error() string {
	return fmt.Sprintf("Object cannot be found:%s\n", err.Msg)
}

func NewNotFoundException(msg string) NotFoundException {
	err := NotFoundException{}
	err.Msg = msg
	return err
}

/*****************************************************************************/

type SQLException struct {
	Exception
}

func (err SQLException) Error() string {
	return fmt.Sprintf("Database SQL ERROR:%s\n", err.Msg)
}

func NewSQLException(msg string) SQLException {
	err := SQLException{}
	err.Msg = msg
	return err
}

/*****************************************************************************/

type ConnectionException struct {
	Exception
}

func (err ConnectionException) Error() string {
	return fmt.Sprintf("Cannot Connect to server:%s\n", err.Msg)
}

func NewConnectionException(msg string) ConnectionException {
	err := ConnectionException{}
	err.Msg = msg
	return err
}

/*****************************************************************************/
