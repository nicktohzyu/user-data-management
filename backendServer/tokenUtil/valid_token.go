package tokenUtil

const SPECIAL string = "xxx xxx xxx xxx "
const LENGTH int = 16

func IsValidToken(token string) bool {
	if !ValidTokenFormat(token) {
		return false
	}
	// TODO: check that date has not expired
	return true
}

func ValidTokenFormat(token string) bool {
	// length is correct
	if len(token) != LENGTH {
		return false
	}
	// is not the special value of all 0s
	if token == SPECIAL {
		return false
	}
	return true
}
