package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	repo "github.com/ilindan-dev/shortener/internal/domain/repository"
	"github.com/ilindan-dev/shortener/internal/service"
	"github.com/rs/zerolog"
	"net/http"
	"net/url"
)

// Handlers encapsulates all the HTTP handlers for the shortener service.
type Handlers struct {
	urlService       *service.URLService
	analyticsService *service.AnalyticsService
	logger           zerolog.Logger
	baseURL          string // Base URL for constructing short links, e.g., "http://localhost:8080"
}

// NewHandlers creates a new instance of Handlers.
func NewHandlers(urlService *service.URLService, analyticsService *service.AnalyticsService, logger *zerolog.Logger, baseURL string) *Handlers {
	return &Handlers{
		urlService:       urlService,
		analyticsService: analyticsService,
		logger:           logger.With().Str("layer", "http_handler").Logger(),
		baseURL:          baseURL,
	}
}

// RegisterRoutes sets up the routing for the application.
func (h *Handlers) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/shorten", h.CreateShortURL)
		api.GET("/analytics/:short_code", h.GetAnalytics)
	}

	router.GET("/s/:short_code", h.Redirect)
}

// CreateShortURL handles the request to create a new short URL.
func (h *Handlers) CreateShortURL(c *gin.Context) {
	var req CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	createdURL, err := h.urlService.CreateShortURL(c.Request.Context(), req.URL)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create short URL")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create short URL"})
		return
	}

	shortURL, _ := url.JoinPath(h.baseURL, "s", createdURL.ShortCode)
	c.JSON(http.StatusCreated, URLResponse{
		OriginalURL: createdURL.OriginalURL,
		ShortURL:    shortURL,
	})
}

// Redirect handles the redirection from a short URL to the original URL.
func (h *Handlers) Redirect(c *gin.Context) {
	shortCode := c.Param("short_code")
	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()

	gotURL, err := h.urlService.ProcessRedirect(c.Request.Context(), shortCode, userAgent, ipAddress)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Short URL not found"})
			return
		}
		h.logger.Error().Err(err).Str("short_code", shortCode).Msg("Failed to process redirect")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, gotURL.OriginalURL)
}

// GetAnalytics handles the request to fetch analytics for a short URL.
func (h *Handlers) GetAnalytics(c *gin.Context) {
	shortCode := c.Param("short_code")

	report, err := h.analyticsService.GetFullAnalyticsReport(c.Request.Context(), shortCode)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Analytics not found for this URL"})
			return
		}
		h.logger.Error().Err(err).Str("short_code", shortCode).Msg("Failed to get analytics")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve analytics"})
		return
	}

	// Map domain model to response DTO
	clicksByDay := make([]StatItem, len(report.ClicksByDay))
	for i, stat := range report.ClicksByDay {
		clicksByDay[i] = StatItem{Key: stat.Key, Value: stat.Value}
	}
	clicksByUserAgent := make([]StatItem, len(report.ClicksByUserAgent))
	for i, stat := range report.ClicksByUserAgent {
		clicksByUserAgent[i] = StatItem{Key: stat.Key, Value: stat.Value}
	}
	recentClicks := make([]ClickDTO, len(report.RecentClicks))
	for i, click := range report.RecentClicks {
		recentClicks[i] = ClickDTO{Timestamp: click.CreatedAt, UserAgent: click.UserAgent}
	}
	shortURL, _ := url.JoinPath(h.baseURL, "s", report.URL.ShortCode)

	c.JSON(http.StatusOK, AnalyticsResponse{
		OriginalURL:       report.URL.OriginalURL,
		ShortURL:          shortURL,
		TotalClicks:       report.TotalClicks,
		ClicksByDay:       clicksByDay,
		ClicksByUserAgent: clicksByUserAgent,
		RecentClicks:      recentClicks,
	})
}
