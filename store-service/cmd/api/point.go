package api

import (
	"log"
	"net/http"

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
	var request point.SubmitedPoint
	if err := context.BindJSON(&request); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		log.Printf("bad request %s", err.Error())
		return
	}

	uid := 1
	res, err := api.PointService.DeductPoint(uid, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
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
	uid := 1
	res, err := api.PointService.TotalPoint(uid)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, res)
}
