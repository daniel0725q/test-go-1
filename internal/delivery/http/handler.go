package truoraHttp

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/truora/microservice/internal/dto"
	"github.com/truora/microservice/internal/usecase"
)

type Handler struct {
	stockRatingSvc usecase.StockRatingService
}

func NewHandler(stockRatingSvc usecase.StockRatingService) *Handler {
	return &Handler{
		stockRatingSvc: stockRatingSvc,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/hello", h.HelloWorld)
	r.Get("/api/external/hello", h.GetExternalHello)

	r.Route("/api/stock-ratings", func(r chi.Router) {
		r.Post("/", h.CreateStockRating)
		r.Post("/batch", h.CreateStockRatingBatch)
		r.Get("/{id}", h.GetStockRatingByID)
		r.Get("/ticker/{ticker}", h.GetStockRatingsByTicker)
		r.Get("/ticker/{ticker}/latest", h.GetLatestStockRatingByTicker)
	})
}

func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Hello World!",
	})
}

func (h *Handler) GetExternalHello(w http.ResponseWriter, r *http.Request) {
	response, err := h.stockRatingSvc.GetHello(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "Successfully retrieved hello message",
		"data":    response,
	})
}

func (h *Handler) CreateStockRating(w http.ResponseWriter, r *http.Request) {
	var rating dto.StockRatingResponse
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.stockRatingSvc.CreateStockRating(r.Context(), &rating); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, rating)
}

func (h *Handler) CreateStockRatingBatch(w http.ResponseWriter, r *http.Request) {
	var ratings []*dto.StockRatingResponse
	if err := json.NewDecoder(r.Body).Decode(&ratings); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.stockRatingSvc.CreateStockRatingBatch(r.Context(), ratings); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, ratings)
}

func (h *Handler) GetStockRatingByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID format")
		return
	}

	rating, err := h.stockRatingSvc.GetStockRatingByID(r.Context(), uint(id))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if rating == nil {
		respondWithError(w, http.StatusNotFound, "Stock rating not found")
		return
	}

	respondWithJSON(w, http.StatusOK, rating)
}

func (h *Handler) GetStockRatingsByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "Ticker is required")
		return
	}

	ratings, err := h.stockRatingSvc.GetStockRatingsByTicker(r.Context(), ticker)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ratings)
}

func (h *Handler) GetLatestStockRatingByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "Ticker is required")
		return
	}

	rating, err := h.stockRatingSvc.GetLatestStockRatingByTicker(r.Context(), ticker)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if rating == nil {
		respondWithError(w, http.StatusNotFound, "No stock rating found for ticker")
		return
	}

	respondWithJSON(w, http.StatusOK, rating)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
