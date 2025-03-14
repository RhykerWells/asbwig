package bot

import "testing"

func TestToInt64(t *testing.T) {
	// String test
	i := ToInt(1)
	t.Logf("integer test: %d || SHOULD WORK", i)
	i = ToInt("2")
	t.Logf("String test: %d || SHOULD WORK", i)
	i = ToInt(3.0)
	t.Logf("flat float test: %d || SHOULD WORK", i)
	i = ToInt("4.0")
	t.Logf("string flat float test: %d || SHOULD WORK", i)
	i = ToInt(5.1)
	t.Logf("non flat float test: %d", i)
	i = ToInt("6.1")
	t.Logf("string non flat float test: %d", i)
}