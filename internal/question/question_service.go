package question

import (
	"course/internal/domain"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QuestionService struct {
	db *gorm.DB
}

func NewQuestionService(database *gorm.DB) *QuestionService {
	return &QuestionService{
		db: database,
	}
}

func (ex QuestionService) GetQuestion(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid question id",
		})
		return
	}
	var question domain.Question
	err = ex.db.Where("id = ?", id).Take(&question).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}
	ctx.JSON(200, question)
}

// func (ex QuestionService) GetUserScore(ctx *gin.Context) {
// 	paramQuestionID := ctx.Param("id")
// 	questionID, err := strconv.Atoi(paramQuestionID)
// 	if err != nil {
// 		ctx.JSON(400, gin.H{
// 			"message": "invalid question id",
// 		})
// 		return
// 	}
// 	var question domain.Question
// 	err = ex.db.Where("id = ?", questionID).Preload("Questions").Take(&question).Error
// 	if err != nil {
// 		ctx.JSON(404, gin.H{
// 			"message": "not found",
// 		})
// 		return
// 	}

// 	userID := int(ctx.Request.Context().Value("user_id").(float64))
// 	var questions []domain.Question
// 	err = ex.db.Where("question_id = ? AND user_id = ?", questionID, userID).Find(&questions).Error

// 	if err != nil {
// 		ctx.JSON(200, gin.H{
// 			"score": 0,
// 		})
// 		return
// 	}

// 	mapQA := make(map[int]domain.Question)
// 	for _, question := range questions {
// 		mapQA[question.QuestionID] = question
// 	}

// 	var score int
// 	for _, question := range question.Questions {
// 		if strings.EqualFold(question.CorrectQuestion, mapQA[question.ID].Question) {
// 			score += question.Score
// 		}
// 	}
// 	ctx.JSON(200, gin.H{
// 		"score": score,
// 	})
// }
