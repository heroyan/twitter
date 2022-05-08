package api

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_GenSessionId(t *testing.T) {
	convey.Convey("SecessionId test", t, func() {
		t.Logf(genSessionId())
	})
}
