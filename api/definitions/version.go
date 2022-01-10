package definitions

import "strings"

// MatchVersion checks if a version matches a spec,
// a spec is a version with zero or more of the numbers replaced with an 'x'.
// version can also be a spec
func MatchVersion(spec, version string) bool {
	if !IsValidVersionSpec(spec) || !IsValidVersionSpec(version) {
		return false
	}

	snum, vnum := "", ""
	snext, vnext := false, false
	for {
		snum, spec, snext = nextNum(spec)
		vnum, version, vnext = nextNum(version)
		if isNumber(snum) {
			if snum != vnum {
				return false
			}
		} else if snum == "x" {
			if vnum != "x" && !isNumber(vnum) {
				return false
			}
		} // else
		// {
		//    // no need for an else here, the spec is already checked to be valid.
		// }

		if !snext {
			return !vnext
		}
	}
}

// IsValidVersion returns true for a string of the form "0.0.1" or "1.1"
func IsValidVersion(version string) bool {
	num, next, n := "", false, 0
	for {
		n++
		num, version, next = nextNum(version)
		if next {
			if !isNumber(num) {
				return false
			}
		} else {
			return isNumber(num) && n <= 4
		}
	}
}

// IsValidVersionSpec returns true for a string of the form "0.0.1" or "1.x.x"
func IsValidVersionSpec(version string) bool {
	num, next, xOnly := "", false, false
	for {
		num, version, next = nextNum(version)

		if next {
			if num == "x" {
				xOnly = true
				continue
			}

			if isNumber(num) && !xOnly {
				continue
			}

			return false
		} else {
			if num == "x" {
				return true
			}

			if isNumber(num) {
				return !xOnly
			}

			return false
		}
	}
}

func nextNum(s string) (num string, remaining string, next bool) {
	dot := strings.Index(s, ".")
	if dot < 0 {
		return s, "", false
	} else {
		return s[:dot], s[dot+1:], true
	}
}

func isNumber(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if !isDigit(r) {
			return false
		}
	}

	return true
}

func isDigit(r rune) bool { return '0' <= r && r <= '9' }
