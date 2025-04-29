package middware

import (
	"context"
	"net/http"

	"github.com/VicShved/loyalty/internal/common"
	"github.com/VicShved/loyalty/internal/logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func setAuthCook(w http.ResponseWriter, userID *string) {

	token, _ := common.GetJWTTokenString(userID)
	http.SetCookie(w, &http.Cookie{
		Name:  common.AuthorizationCookName,
		Value: token,
	})
}

func parseTokenUserID(tokenStr string) (*jwt.Token, string, error) {
	claims := &common.CustClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(common.ServerConfig.SecretKey), nil
	})
	return token, (*claims).UserID, err
}

// auth middleware
func AuthMiddleware(next http.Handler) http.Handler {
	authFunc := func(w http.ResponseWriter, r *http.Request) {
		var userID string
		var token *jwt.Token
		cook, err := r.Cookie(common.AuthorizationCookName)
		//  если нет куки, то создаю новую
		if err == http.ErrNoCookie {
			logger.Log.Debug("ErrNoCookie")
			userID, _ = common.GetNewUUID()
			setAuthCook(w, &userID)
		} else {
			token, userID, _ = parseTokenUserID(cook.Value)
			// Если токен не валидный,  то создаю нвый userID
			if !token.Valid {
				logger.Log.Debug("Not valid token")
				userID, _ = common.GetNewUUID()
				setAuthCook(w, &userID)
			}
		}
		// Если кука не содержит ид пользователя, то возвращаю 401
		if userID == "" {
			logger.Log.Debug("Empty userID")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Log.Debug("User ", zap.String("ID", string(userID)))
		// добавляю userID в контекст
		ctx := context.WithValue(r.Context(), common.ContextUser, userID)
		// Вызываю след.обработчик
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(authFunc)
}
