package middware

import (
	"context"
	"net/http"

	"github.com/VicShved/loyalty/internal/common"
	"github.com/VicShved/loyalty/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func parseTokenUserID(tokenStr string) (*jwt.Token, uint, error) {
	claims := &common.CustClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(common.ServerConfig.SecretKey), nil
	})
	return token, (*claims).UserID, err
}

// auth middleware
func AuthMiddlewareCookie(next http.Handler) http.Handler {
	authFunc := func(w http.ResponseWriter, r *http.Request) {
		var userID uint
		var token *jwt.Token
		cook, err := r.Cookie(common.AuthorizationName)
		if err == http.ErrNoCookie {
			logger.Log.Debug("ErrNoCookie")
		} else {
			token, userID, _ = parseTokenUserID(cook.Value)
			logger.Log.Debug("AuthMiddleware", zap.Any("token.Claims", token.Claims), zap.String("cookie", cook.Value))
			// Если токен не валидный,  то создаю нoвый userID
			if !token.Valid {
				logger.Log.Debug("Not valid token")
				userID = 0
			}
		}
		logger.Log.Debug("User ", zap.Uint("userID", userID))
		// добавляю userID в контекст
		ctx := context.WithValue(r.Context(), common.ContextUser, userID)
		// Вызываю след.обработчик
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(authFunc)
}

func AuthMiddlewareHeader(next http.Handler) http.Handler {
	authFunc := func(w http.ResponseWriter, r *http.Request) {
		var userID uint
		var token *jwt.Token
		tokenString := r.Header.Get(common.AuthorizationName)
		if tokenString == "" {
			logger.Log.Debug("ErrNoAuthHeader")
		} else {
			token, userID, _ = parseTokenUserID(tokenString)
			logger.Log.Debug("AuthMiddleware", zap.Any("token.Claims", token.Claims), zap.String("cookie", tokenString))
			if !token.Valid {
				logger.Log.Debug("Not valid token")
				userID = 0
			}
		}
		logger.Log.Debug("User ", zap.Uint("userID", userID))
		// добавляю userID в контекст
		ctx := context.WithValue(r.Context(), common.ContextUser, userID)
		// Вызываю след.обработчик
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(authFunc)
}
