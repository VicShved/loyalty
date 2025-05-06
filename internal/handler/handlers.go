package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/VicShved/loyalty/internal/common"
	"github.com/VicShved/loyalty/internal/logger"
	"github.com/VicShved/loyalty/internal/repository"
	"github.com/VicShved/loyalty/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type reqJSON struct {
	URL string `json:"url"`
}

type respJSON struct {
	Result string `json:"result"`
}

type Handler struct {
	serv    *service.ShortenService
	baseurl string
}

func GetHandler(serv *service.ShortenService) *Handler {
	return &Handler{
		serv: serv,
	}
}

func setAuthCook(w http.ResponseWriter, userID uint) {

	token, err := common.GetJWTTokenString(userID)
	if err != nil {
		logger.Log.Error("setAuthCook", zap.Any("Error", err))
	}
	http.SetCookie(w, &http.Cookie{
		Name:  common.AuthorizationName,
		Value: token,
	},
	)
}

func setAuthHeader(w http.ResponseWriter, userID uint) {
	token, err := common.GetJWTTokenString(userID)
	if err != nil {
		logger.Log.Error("setAuthHeader", zap.Any("Error", err))
	}
	w.Header().Add(common.AuthorizationName, token)
}

func (h Handler) InitRouter(mdwr []func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	for _, mw := range mdwr {
		router.Use(mw)
	}
	router.Post("/api/user/register", h.PostRegister)
	router.Post("/api/user/login", h.PostLogin)
	router.Post("/api/user/orders", h.PostOrders)
	router.Get("/api/user/orders", h.GetOrders)
	// router.Get("/api/user/balance", h.GetBalance)
	// router.Post("/api/user/balance/withdraw", h.PostWithDraw)
	// router.Get("/api/user/withdrawals", h.GetWithDrawals)
	router.Get("/api/ping", h.PingDB)
	return router
}

func (h Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(uint)
	logger.Log.Debug("Context User", zap.Any("userID", userID))

	err := h.serv.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h Handler) PostRegister(w http.ResponseWriter, r *http.Request) {
	var logPass loginPassword
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &logPass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ValidateLoginPassword(logPass) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Log.Debug("PostRegister", zap.Any("logPass = ", logPass))
	userID, err := h.serv.Register(logPass.Login, logPass.Password)
	if err != nil {
		if errors.Is(err, repository.ErrLoginConflict) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// setAuthCook(w, userID)
	setAuthHeader(w, userID)
}

func (h Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	var logPass loginPassword
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &logPass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !ValidateLoginPassword(logPass) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Log.Debug("PostRegister", zap.Any("logPass = ", logPass))
	userID, err := h.serv.Login(logPass.Login, logPass.Password)
	logger.Log.Debug("", zap.Uint("userID", userID), zap.Any("err", err))
	if err != nil {
		if errors.Is(err, repository.ErrLoginPassword) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// setAuthCook(w, userID)
	setAuthHeader(w, userID)
}

func (h Handler) PostOrders(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(uint)
	logger.Log.Debug("Context User ", zap.Uint("userID", userID))
	if userID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	defer r.Body.Close()
	orderBytes, _ := io.ReadAll(r.Body)
	orderNumber := string(orderBytes)
	if orderNumber == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.Log.Debug("", zap.String("orderNumber", orderNumber))
	if !isOnlyDigits(orderNumber) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	_, isNew, err := h.serv.SaveOrder(orderNumber, userID)

	if err != nil {
		if errors.Is(err, repository.ErrOrderNumberUserConflict) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isNew {
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(uint)
	logger.Log.Debug("Context User ", zap.Uint("ID", userID))

	w.Header().Set("Content-Type", "application/json")
	orders, err := h.serv.GetOrders(userID)
	if len(*orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	respData, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	lenth, err := w.Write(respData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(lenth))
	logger.Log.Debug("", zap.String("response", string(respData)))
}

// func (h Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
// 	// Вытаскиваю userID из контекста
// 	userID := r.Context().Value(common.ContextUser).(string)
// 	logger.Log.Debug("Context User ", zap.Any("ID", userID))

// }

// func (h Handler) PostWithDraw(w http.ResponseWriter, r *http.Request) {
// 	// Вытаскиваю userID из контекста
// 	userID := r.Context().Value(common.ContextUser).(string)
// 	logger.Log.Debug("Context User ", zap.Any("ID", userID))

// }

// func (h Handler) GetWithDrawals(w http.ResponseWriter, r *http.Request) {
// 	// Вытаскиваю userID из контекста
// 	userID := r.Context().Value(common.ContextUser).(string)
// 	logger.Log.Debug("Context User ", zap.Any("ID", userID))

// }
