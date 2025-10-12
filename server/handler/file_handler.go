package handler

import (
	"net/http"
	"project/pkg/storage"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	Spaces *storage.SpacesClient
}

func NewFileHandler(spaces *storage.SpacesClient) *FileHandler {
	return &FileHandler{Spaces: spaces}
}

func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open file"})
		return
	}
	defer src.Close()

	url, err := h.Spaces.UploadFile(src, file.Filename, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func (h *FileHandler) GetFile(c *gin.Context) {
	filename := c.Param("filename")
	data, contentType, err := h.Spaces.GetFile(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, contentType, data)
}
