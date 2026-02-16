package http

import (
    "context"
    "database/sql"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type HealthHandler struct {
    db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
    return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    if h.db != nil {
        if err := h.db.PingContext(ctx); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "status": "unavailable",
                "detail": err.Error(),
            })
            return
        }
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
