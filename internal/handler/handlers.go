package handler

import (
	"net/http"

	"github.com/VicShved/loyalty/internal/common"
	"github.com/VicShved/loyalty/internal/logger"
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

func (h Handler) InitRouter(mdwr []func(http.Handler) http.Handler) *chi.Mux {
	router := chi.NewRouter()
	for _, mw := range mdwr {
		router.Use(mw)
	}
	router.Post("/api/user/register", h.PostRegister)
	router.Post("/api/user/login", h.PostLogin)
	router.Post("/api/user/orders", h.PostOrders)
	router.Get("/api/user/orders", h.GetOrders)
	router.Get("/api/user/balance", h.GetBalance)
	router.Post("/api/user/balance/withdraw", h.PostWithDraw)
	router.Get("/api/user/withdrawals", h.GetWithDrawals)
	router.Get("/api/ping", h.PingDB)
	return router
}

type loginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h Handler) PostRegister(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	var indata loginPassword
}

func (h Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	var indata loginPassword
}

func (h Handler) PingDB(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

	err := h.serv.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h Handler) PostOrders(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

}

func (h Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

}

func (h Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

}

func (h Handler) PostWithDraw(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

}

func (h Handler) GetWithDrawals(w http.ResponseWriter, r *http.Request) {
	// Вытаскиваю userID из контекста
	userID := r.Context().Value(common.ContextUser).(string)
	logger.Log.Debug("Context User ", zap.Any("ID", userID))

}
