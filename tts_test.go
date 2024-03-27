package main

import (
	"testing"
)

func TestEscapeXmlText(t *testing.T) {
	res := escapeXmlText("&><")
	t.Logf("voice: %s", res)
}
