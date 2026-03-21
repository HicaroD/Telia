package token

import "testing"

func TestErrorTypeValue(t *testing.T) {
	t.Logf("ERROR_TYPE = %d", ERROR_TYPE)
	t.Logf("UNTYPED_NULLPTR = %d", UNTYPED_NULLPTR)
	t.Logf("STRING_TYPE = %d", STRING_TYPE)
	kind, ok := KEYWORDS["error"]
	t.Logf("KEYWORDS[error] = %d, ok = %v", kind, ok)
}
