package messages

import (
	fmt "fmt"
	"time"

	context "golang.org/x/net/context"
)

type Service struct{}

func (this *Service) Call(
	ctx context.Context, req *Request) (*Response, error) {
	msg := time.Now().Format("2006-01-02 15:04:05.999999")
	return &Response{Resp: fmt.Sprintf("Hello, Now is %s", msg)}, nil
}
