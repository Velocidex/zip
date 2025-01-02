package zip

import (
	"io/ioutil"
	"os"
)

var TmpfileFactory TmpfileProvider = DefaultTmpfileProvider(0)

// A provider for temp files. Can be overridden by callers to
// customized tmpfile management.
type TmpfileProvider interface {
	TempFile() (*os.File, error)
	RemoveTempFile(filename string)
}

type DefaultTmpfileProvider int

func (self DefaultTmpfileProvider) TempFile() (*os.File, error) {
	return ioutil.TempFile("", "tmp")
}

func (self DefaultTmpfileProvider) RemoveTempFile(filename string) {
	_ = os.Remove(filename)
}

func SetTmpfileProvider(provider TmpfileProvider) {
	TmpfileFactory = provider
}
