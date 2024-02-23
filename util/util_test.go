package util

import "testing"

func TestCheckValidPassword(t *testing.T) {
	var passwords = map[string]bool{
		"_}24I:9t58Tu?m@e":  false,
		"|YlzEc|1":          false,
		"m_4xF%t\"Bu5jeb$":  true,
		"Password1!":        true,
		"12345678aA":        false,
		"^@},^@},^@},":      false,
		"12345678aA@":       false,
		"12345678 aA@":      false,
		"AbCdEfGhIi":        false,
		"abcdefgh123":       false,
		"ABCDEFGH1234":      false,
		"A1234a.":           false,
		"a1234A.":           false,
		"A1234567890abcde.": false,
		"a1234567890Abcde.": false,
	}

	for k, v := range passwords {
		if res, err := CheckValidPassword(k); res != v {
			if err != nil {
				t.Error(err)
			}
			t.Errorf("CheckValid: %v, result: %v, want: %v", k, res, v)
		}
	}
}
