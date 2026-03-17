package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"store-service/internal/point"

	"github.com/gin-gonic/gin"
)

type PointAPI struct {
	PointService point.PointInterface
}

// @Summary Deduct points from user
// @Description Deduct points from user's point balance
// @Tags point
// @Accept json
// @Produce json
// @Param request body point.SubmitedPoint true "Point deduction request"
// @Success 200 {object} point.Point
// @Failure 400 {string} string "Bad request error"
// @Failure 500
// @Router /api/v1/point [post]
func (api PointAPI) DeductPointHandler(context *gin.Context) {
	ctx := context.Request.Context()

	var request point.SubmitedPoint
	if err := context.BindJSON(&request); err != nil {
		slog.ErrorContext(ctx, "Point deduct bad request",
			"log_type", "error",
			"error_code", "INVALID_REQUEST",
			"error_message", err.Error(),
			"user_id", 0,
		)
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	uid, uidErr := strconv.Atoi(context.GetHeader("uid"))
	if uidErr != nil {
		uid = 1
	}

	res, err := api.PointService.DeductPoint(ctx, uid, request)
	if err != nil {
		slog.ErrorContext(ctx, "PointService.DeductPoint failed",
			"log_type", "error",
			"error_code", "POINT_DEDUCTION_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
			slog.Any("request", map[string]any{"amount": request.Amount}),
		)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.InfoContext(ctx, "Points deducted",
		"log_type", "business",
		"event", "points_deducted",
		"entity_type", "point",
		"entity_id", uid,
		"actor_id", uid,
		slog.Any("metadata", map[string]any{
			"amount":          request.Amount,
			"remaining_point": res.Point,
		}),
	)

	context.JSON(http.StatusOK, res)
}

// @Summary Get total points
// @Description Get user's total point balance
// @Tags point
// @Accept json
// @Produce json
// @Success 200 {object} point.Point
// @Failure 500
// @Router /api/v1/point [get]
func (api PointAPI) TotalPointHandler(context *gin.Context) {
	uid, uidErr := strconv.Atoi(context.GetHeader("uid"))
	if uidErr != nil {
		uid = 1
	}

	ctx := context.Request.Context()
	res, err := api.PointService.TotalPoint(ctx, uid)

	if err != nil {
		slog.ErrorContext(ctx, "PointService.TotalPoint failed",
			"log_type", "error",
			"error_code", "POINT_QUERY_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
		)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, res)
}
