package validator

import "regexp"

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)

func IsValidUsername(username string) bool {

	return usernameRegex.MatchString(username)

}
