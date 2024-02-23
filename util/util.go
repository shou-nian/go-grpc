package util

import (
	"regexp"
	"strings"
)

// CheckValidPassword
//
//	password rule:
//	1. Length is 8 to 16
//	2. First character is a-zA-Z
//	3. At least one capital letter and one lowercase letter
//	4. At least one 0-9 character
//	5. At least one special character: ,./?<>!#=-+[](){}|@%"$
//	6. Cannot contain Spaces
func CheckValidPassword(password string) (bool, error) {
	// check rule 1
	if len(password) < 8 || len(password) > 16 {
		return false, nil
	}

	// check rule 6
	if strings.Contains(password, " ") {
		return false, nil
	}

	// check rule 2
	if ok, err := regexp.MatchString(`^[a-zA-Z]`, string(password[0])); !ok {
		return false, err
	} else {
		// check rule 3
		if ok, err = regexp.MatchString(`[a-z]`, password); !ok {
			return false, err
		} else {
			if ok, err = regexp.MatchString(`[A-Z]`, password); !ok {
				return false, err
			}
		}
	}

	// check rule 4
	if ok, err := regexp.MatchString(`[0-9]`, password); !ok {
		return false, err
	}

	// check rule 5
	if ok, err := regexp.MatchString(`[,./?<>!#=\-+\[\](){}|@%"$]`, password); !ok {
		return false, err
	}

	return true, nil
}
