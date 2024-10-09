package version

import (
	"testing"
)

func Test_IsDebug(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		Version = ""
		Revision = ""
		got := IsDebug()
		if got != true {
			t.Errorf("got = %#v, want %#v", got, true)
		}
	})
}
