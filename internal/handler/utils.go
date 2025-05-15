package handler

import (
	// "regexp"
	"strconv"
)

type loginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func ValidateLoginPassword(lp loginPassword) bool {
	if lp.Login == "" {
		return false
	}
	if lp.Password == "" {
		return false
	}
	return true
}

func isOnlyDigits(s string) bool {
	// var re = regexp.MustCompile(`^[0-9]+s`)
	// if re.MatchString(s) {
	// 	return true
	// }
	// return false // todo set false + fix re
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}
