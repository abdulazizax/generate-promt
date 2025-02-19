package handlers

import (
	"generate-promt-v1/api/models"
	"generate-promt-v1/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
)

func (h Handler) HandleDbError(c *gin.Context, err error, message string) bool {
	if err == nil {
		return false
	}

	var errorResponse models.ErrorResponse
	statusCode := http.StatusInternalServerError

	if err == pgx.ErrNoRows {
		errorResponse = models.ErrorResponse{
			Message: "The requested resource was not found.",
			Code:    config.ErrorNotFound,
		}
		c.JSON(http.StatusNotFound, errorResponse)
		return true
	}

	switch e := err.(type) {
	case *pgconn.PgError:
		// Handle PostgreSQL-specific errors
		switch e.Code {
		case "23505":
			// Unique constraint violation
			errorResponse = models.ErrorResponse{
				Message: "Duplicate key error (unique constraint violation).",
				Code:    config.ErrorDuplicateKey,
			}
			statusCode = http.StatusBadRequest
		case "23503":
			// Foreign key violation
			errorResponse = models.ErrorResponse{
				Message: "The record could not be deleted because it is used in other records.",
				Code:    config.ErrorConflict,
			}
			statusCode = http.StatusBadRequest
		case "22001":
			// Value too long for column
			errorResponse = models.ErrorResponse{
				Message: "Value too long for column.",
				Code:    config.ErrorInvalidRequest,
			}
			statusCode = http.StatusBadRequest
		default:
			// General PostgreSQL error
			errorResponse = models.ErrorResponse{
				Message: "Ooops! Something went wrong.",
				Code:    config.ErrorInternalServer,
			}
		}
	default:
		if strings.Contains(err.Error(), "BAD_REQUEST") {
			errorResponse = models.ErrorResponse{
				Message: strings.TrimPrefix(err.Error(), "BAD_REQUEST"),
				Code:    config.ErrorBadRequest,
			}
		} else {
			// General PostgreSQL error
			errorResponse = models.ErrorResponse{
				Message: "Ooops! Something went wrong.",
				Code:    config.ErrorInternalServer,
			}
		}
	}

	c.JSON(statusCode, errorResponse)
	return true
}

func (h Handler) ReturnError(c *gin.Context, code string, message string, statusCode int) {
	errorResponse := models.ErrorResponse{
		Message: message,
		Code:    code,
	}
	c.JSON(statusCode, errorResponse)
}
