package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"quote/models"
)

type TodoHandler struct {
	DB *gorm.DB
}

func (h *TodoHandler) GetAll(c *gin.Context) {
	var todos []models.Todo
	h.DB.Find(&todos)

	c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func (h *TodoHandler) Create(c *gin.Context) {
	todo := models.Todo{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Completed:   false,
	}

	h.DB.Create(&todo)

	c.JSON(http.StatusOK, gin.H{"data": todo})
}

func (h *TodoHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	if err := h.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id error, not found"})
		return
	}
	h.DB.Model(&todo).Update("completed", true)

	c.JSON(http.StatusOK, gin.H{"data": todo})

	//h.DB.Delete(&*todo)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	if err := h.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id error, not found"})
		return
	}

	h.DB.Delete(&todo)

	c.JSON(http.StatusOK, gin.H{"data": "Delete succeed"})
}
