package utils

import (
	"path/filepath"
	"testing"
)

func TestValidFile(t *testing.T) {
	valid := validName("something_1.toml", "something")
	if valid != true {
		t.Error("invalid parsing for validName")
	}
}
