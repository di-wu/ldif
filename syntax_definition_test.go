package ldif

import (
	"io/ioutil"
	"testing"
)

func TestExample1(t *testing.T) {
	raw, _ := ioutil.ReadFile("testdata/example1.ldif")
	s := []rune(string(raw))
	f := ldifFile(s)

	t.Run("version", func(t *testing.T) {
		if f.GetNode("version:", true) == nil {
			t.Error("version not found")
			return
		}

		if nr := f.GetNode("version-number", true); nr == nil {
			t.Error("version number not found")
		} else if v := nr.ValueString(); v != "1" {
			t.Errorf("version number was not 1, got %s", v)
		}
	})

	t.Run("records", func(t *testing.T) {
		records := f.GetAllNodes("ldif-attrval-record")
		if l := len(records); l != 2 {
			t.Errorf("did not find 2 records, got %d", l)
		}
	})
}
