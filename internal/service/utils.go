package service

import (
	"regexp"
	"strconv"

	"github.com/VicShved/loyalty/internal/logger"
	"go.uber.org/zap"
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
	// var re = regexp.MustCompile(`^[0-9]*$`)
	// return re.MatchString(s)
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
		logger.Log.Debug("CheckLuhn", zap.Int("In digit", digit))
		if i%2 == parity {
			digit = digit * 2
			if digit > 9 {
				digit = digit - 9
			}
		}
		logger.Log.Debug("CheckLuhn", zap.Int("Out digit", digit))
		sum = sum + digit
		logger.Log.Debug("CheckLuhn", zap.Int("sum", sum))
	}
	return sum%10 == 0
}
