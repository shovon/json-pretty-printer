package json

import "testing"

func TestGetEscapedRune(t *testing.T) {
	assert(getEscapedRune('<') == "&lt;", "Should have escaped '<' to '&lt;'")
	assert(getEscapedRune('>') == "&gt;", "Should have escaped '>' to '&gt;'")
	assert(getEscapedRune('&') == "&amp;", "Should have escaped '&' to '&amp;")
	assert(getEscapedRune('"') == "&quot;", "Should have escaped '\"' to '&quot'")
	assert(getEscapedRune('\'') == "&apos;", "Should have escaped ''' to '&apos'")
	assert(getEscapedRune('f') == "f", "Should not have escaped 'f'")
}
