package zip

import "github.com/davecgh/go-spew/spew"

func Debug(v interface{}) {
	spew.Dump(v)
}
