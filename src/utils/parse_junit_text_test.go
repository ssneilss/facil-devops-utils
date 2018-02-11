package utils

import (
	"testing"
)

func TestParseJunitText(t *testing.T) {
	t.Run("Should parse xml correctly", func(t *testing.T) {
		text, _ := ParseJunitText("junit.xml")
		if text == "" {
			t.Fail()
		}
	})
}
