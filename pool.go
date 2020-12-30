package zip

import (
	"context"
	"io"
	"sync"

	"github.com/hashicorp/go-multierror"
)

var (
	pool = sync.Pool{
		New: func() interface{} {
			buffer := make([]byte, 32*1024)
			return &buffer
		},
	}
)

type Request struct {
	Name   string
	Reader io.ReadCloser
}

type CompressorPool struct {
	ctx    context.Context
	cancel func()
	c      chan *Request
	wg     sync.WaitGroup
	zip    *Writer

	err error
}

func (self *CompressorPool) check_error(err error) {
	if err == nil {
		return
	}
	self.err = multierror.Append(self.err, err)
}

func (self *CompressorPool) Compress(request *Request) {
	self.c <- request
}

func (self *CompressorPool) Close() error {
	self.cancel()
	self.wg.Wait()

	return self.err
}

func NewCompressorPool(ctx context.Context, zip *Writer, size int) *CompressorPool {
	sub_ctx, cancel := context.WithCancel(ctx)

	self := &CompressorPool{
		ctx:    sub_ctx,
		cancel: cancel,
		c:      make(chan *Request, size),
		zip:    zip,
	}

	for i := 0; i < size; i++ {
		self.wg.Add(1)

		go func() {
			defer self.wg.Done()

			for {
				select {
				case <-self.ctx.Done():
					return

				case request := <-self.c:
					out_fd, err := self.zip.Create(request.Name)
					self.check_error(err)
					_, err = Copy(self.ctx, out_fd, request.Reader)
					self.check_error(err)
					self.check_error(out_fd.Close())
					self.check_error(request.Reader.Close())
				}
			}
		}()
	}

	return self
}

// An io.Copy() that respects context cancellations.
func Copy(ctx context.Context, dst io.Writer, src io.Reader) (n int, err error) {
	offset := 0
	buff := pool.Get().(*[]byte)
	defer pool.Put(buff)

	for {
		select {
		case <-ctx.Done():
			return n, nil

		default:
			n, err = src.Read(*buff)
			if err != nil && err != io.EOF {
				return offset, err
			}

			if n == 0 {
				return offset, nil
			}

			_, err = dst.Write((*buff)[:n])
			if err != nil {
				return offset, err
			}
			offset += n
		}
	}
}
