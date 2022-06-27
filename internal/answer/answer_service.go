package answer

import (
	"course/internal/domain"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnswerService struct {
	db *gorm.DB
}

func NewAnswerService(database *gorm.DB) *AnswerService {
	return &AnswerService{
		db: database,
	}
}

func (ex AnswerService) GetAnswer(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid answer id",
		})
		return
	}
	var answer domain.Answer
	err = ex.db.Where("id = ?", id).Take(&answer).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}
	ctx.JSON(200, answer)
}
