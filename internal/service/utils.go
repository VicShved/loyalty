package service

import (
	"regexp"
	"strconv"
)

type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func ValidateLoginPassword(lp LoginPassword) bool {
	if lp.Login == "" {
		return false
	}
	if lp.Password == "" {
		return false
	}
	return true
}

func IsOnlyDigits(s string) bool {
	result, err := regexp.MatchString("^[0-9]*$", s)
	if err != nil {
		return false
	}
	return result
}

func CheckLuhn(s string) bool {
	sum := 0
	nDigits := len(s)
	parity := nDigits % 2
	for i, ch := range s {
		digit, err := strconv.Atoi(string(ch))
		if err != nil {
			return false
		}
		if i%2 == parity {
			digit = digit * 2
			if digit > 9 {
				digit = digit - 9
			}
		}
		sum = sum + digit
	}
	return sum%10 == 0
}
