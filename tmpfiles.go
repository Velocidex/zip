package zip

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
)

var TmpfileFactory TmpfileProvider = DefaultTmpfileProvider(0)

type TempFile interface {
	io.WriteCloser

	// Open the tempfile for reading.
	Open() (io.ReadCloser, error)

	// Remove the tmpfile and clean up.
	Remove()

	// For information purposes only.
	Name() string
}

// A provider for temp files. Can be overridden by callers to
// customized tmpfile management.
type TmpfileProvider interface {
	TempFile() (TempFile, error)
}

type DefaultTempFile struct {
	*bufio.Writer
	fd *os.File
}

func (self *DefaultTempFile) Name() string {
	return self.fd.Name()
}

func (self *DefaultTempFile) Close() error {
	self.Writer.Flush()
	return self.fd.Close()
}

func (self *DefaultTempFile) Open() (io.ReadCloser, error) {
	self.Writer.Flush()
	self.fd.Close()

	fd, err := os.Open(self.fd.Name())
	return fd, err
}

func (self *DefaultTempFile) Remove() {
	_ = os.Remove(self.fd.Name())
}

type DefaultTmpfileProvider int

func (self DefaultTmpfileProvider) TempFile() (TempFile, error) {
	fd, err := ioutil.TempFile("", "tmp")
	if err != nil {
		return nil, err
	}

	return &DefaultTempFile{
		Writer: bufio.NewWriterSize(fd, 1024*1204),
		fd:     fd,
	}, nil
}

func SetTmpfileProvider(provider TmpfileProvider) {
	TmpfileFactory = provider
}
