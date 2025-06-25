package truoraHttp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/truora/microservice/internal/dto"
	"github.com/truora/microservice/internal/usecase"
)

type Handler struct {
	stockRatingSvc    usecase.StockRatingService
	stockAlgorithmSvc usecase.StockAlgorithmService
}

func NewHandler(stockRatingSvc usecase.StockRatingService, stockAlgorithmSvc usecase.StockAlgorithmService) *Handler {
	return &Handler{
		stockRatingSvc:    stockRatingSvc,
		stockAlgorithmSvc: stockAlgorithmSvc,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/api/hello", h.HelloWorld)
	r.Get("/api/external/hello", h.GetExternalHello)

	r.Route("/api/stock-ratings", func(r chi.Router) {
		r.Get("/", h.GetPaginatedStockRatings)
		r.Post("/", h.CreateStockRating)
		r.Post("/batch", h.CreateStockRatingBatch)
		r.Get("/{id}", h.GetStockRatingByID)
		r.Get("/ticker/{ticker}", h.GetStockRatingsByTicker)
		r.Get("/ticker/{ticker}/latest", h.GetLatestStockRatingByTicker)
	})

	r.Route("/api/algorithms", func(r chi.Router) {
		r.Get("/best-time-to-buy-sell/{ticker}", h.GetBestTimeToBuyAndSell)
		r.Post("/best-time-to-buy-sell/multiple", h.GetBestTimeToBuyAndSellMultiple)
		r.Get("/best-time-to-buy-sell/global", h.GetBestTimeToBuyAndSellGlobal)
	})

	r.Route("/api/jobs", func(r chi.Router) {
		r.Get("/{jobId}", h.GetJobByID)
	})
}

func (h *Handler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Hello World!",
	})
}

func (h *Handler) GetExternalHello(w http.ResponseWriter, r *http.Request) {
	job, err := h.stockRatingSvc.GetHello(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusAccepted, map[string]interface{}{
		"job_id":  job.ID,
		"status":  job.Status,
		"message": "Job created successfully. Use /api/jobs/{job_id} to check status.",
	})
}

func (h *Handler) GetJobByID(w http.ResponseWriter, r *http.Request) {
	jobIDStr := chi.URLParam(r, "jobId")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID format")
		return
	}

	job, err := h.stockRatingSvc.GetJobByID(r.Context(), jobID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if job == nil {
		respondWithError(w, http.StatusNotFound, "Job not found")
		return
	}

	respondWithJSON(w, http.StatusOK, job)
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

func (h *Handler) GetPaginatedStockRatings(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	// Set defaults
	page := 1
	pageSize := 20

	// Parse page parameter
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid page parameter")
			return
		}
	}

	// Parse page_size parameter
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid page_size parameter (must be between 1 and 100)")
			return
		}
	}

	// Get paginated ratings
	response, err := h.stockRatingSvc.GetPaginatedStockRatings(r.Context(), page, pageSize)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *Handler) GetBestTimeToBuyAndSell(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "Ticker is required")
		return
	}

	// Parse optional date parameters
	var startDate, endDate *time.Time

	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
			return
		}
	}

	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
			return
		}
	}

	// Validate date range
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		respondWithError(w, http.StatusBadRequest, "start_date cannot be after end_date")
		return
	}

	recommendation, err := h.stockAlgorithmSvc.BestTimeToBuyAndSell(r.Context(), ticker, startDate, endDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, recommendation)
}

func (h *Handler) GetBestTimeToBuyAndSellMultiple(w http.ResponseWriter, r *http.Request) {
	var request dto.TradingAnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if len(request.Tickers) == 0 {
		respondWithError(w, http.StatusBadRequest, "At least one ticker is required")
		return
	}

	// Validate date range if provided
	if !request.StartDate.IsZero() && !request.EndDate.IsZero() && request.StartDate.After(request.EndDate) {
		respondWithError(w, http.StatusBadRequest, "start_date cannot be after end_date")
		return
	}

	var startDate, endDate *time.Time
	if !request.StartDate.IsZero() {
		startDate = &request.StartDate
	}
	if !request.EndDate.IsZero() {
		endDate = &request.EndDate
	}

	recommendations, err := h.stockAlgorithmSvc.BestTimeToBuyAndSellMultiple(r.Context(), request.Tickers, startDate, endDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, recommendations)
}

func (h *Handler) GetBestTimeToBuyAndSellGlobal(w http.ResponseWriter, r *http.Request) {
	// Parse optional date parameters
	var startDate, endDate *time.Time

	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid start_date format. Use YYYY-MM-DD")
			return
		}
	}

	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		} else {
			respondWithError(w, http.StatusBadRequest, "Invalid end_date format. Use YYYY-MM-DD")
			return
		}
	}

	// Validate date range
	if startDate != nil && endDate != nil && startDate.After(*endDate) {
		respondWithError(w, http.StatusBadRequest, "start_date cannot be after end_date")
		return
	}

	recommendation, err := h.stockAlgorithmSvc.BestTimeToBuyAndSellGlobal(r.Context(), startDate, endDate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, recommendation)
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
