package exercise

import (
	"course/internal/domain"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExerciseService struct {
	db *gorm.DB
}

func NewExerciseService(database *gorm.DB) *ExerciseService {
	return &ExerciseService{
		db: database,
	}
}

func (ex ExerciseService) GetExercise(ctx *gin.Context) {
	paramID := ctx.Param("id")
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise domain.Exercise
	err = ex.db.Where("id = ?", id).Preload("Questions").Take(&exercise).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}
	ctx.JSON(200, exercise)
}

func (ex ExerciseService) GetUserScore(ctx *gin.Context) {
	paramExerciseID := ctx.Param("id")
	exerciseID, err := strconv.Atoi(paramExerciseID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}
	var exercise domain.Exercise
	err = ex.db.Where("id = ?", exerciseID).Preload("Questions").Take(&exercise).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}

	userID := int(ctx.Request.Context().Value("user_id").(float64))
	var answers []domain.Answer
	err = ex.db.Where("exercise_id = ? AND user_id = ?", exerciseID, userID).Find(&answers).Error

	if err != nil {
		ctx.JSON(200, gin.H{
			"score": 0,
		})
		return
	}

	mapQA := make(map[int]domain.Answer)
	for _, answer := range answers {
		mapQA[answer.QuestionID] = answer
	}

	var score int
	for _, question := range exercise.Questions {
		if strings.EqualFold(question.CorrectAnswer, mapQA[question.ID].Answer) {
			score += question.Score
		}
	}
	ctx.JSON(200, gin.H{
		"score": score,
	})
}

func (ex ExerciseService) CreateQuestions(ctx *gin.Context) {
	var questionRequest domain.Question
	paramExerciseID := ctx.Param("id")
	exerciseID, err := strconv.Atoi(paramExerciseID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}

	err = ctx.ShouldBind(&questionRequest)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid Input",
		})
		return
	}

	var exercise domain.Exercise
	err = ex.db.Where("id = ?", exerciseID).Take(&exercise).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "exercise id not found",
		})
		return
	}

	userID := int(ctx.Request.Context().Value("user_id").(float64))

	question := domain.Question{
		ExerciseID:    exerciseID,
		Body:          questionRequest.Body,
		OptionA:       questionRequest.OptionA,
		OptionB:       questionRequest.OptionB,
		OptionC:       questionRequest.OptionC,
		OptionD:       questionRequest.OptionD,
		CorrectAnswer: questionRequest.CorrectAnswer,
		CreatorID:     userID,
		Score:         10,
	}

	if err := ex.db.Create(&question).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed when create user",
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "success",
	})
}

func (ex ExerciseService) CreateAnswer(ctx *gin.Context) {
	var answerRequest domain.Answer
	paramExerciseID := ctx.Param("id")
	fmt.Println(paramExerciseID)
	exerciseID, err := strconv.Atoi(paramExerciseID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid exercise id",
		})
		return
	}

	paramQuestionID := ctx.Param("questionID")
	questionID, err := strconv.Atoi(paramQuestionID)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid question id",
		})
		return
	}

	err = ctx.ShouldBind(&answerRequest)

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid Input",
		})
		return
	}

	var question domain.Question
	err = ex.db.Where("id = ? AND exercise_id = ?", questionID, exerciseID).Take(&question).Error
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "question not found",
		})
		return
	}

	userID := int(ctx.Request.Context().Value("user_id").(float64))

	answer := domain.Answer{
		ExerciseID: exerciseID,
		QuestionID: questionID,
		UserID:     userID,
		Answer:     answerRequest.Answer,
	}

	fmt.Println(answer)

	if err := ex.db.Create(&answer).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed when submit answer of the question",
			"res":     err,
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "success",
	})
}
