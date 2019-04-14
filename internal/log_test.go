package internal

import (
	"strconv"
	"testing"
)

func Test_JustLogErr(t *testing.T) {
	JustLogErr(strconv.Atoi("a"))
}
