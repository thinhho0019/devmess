package router

import (
	"fmt"
	"os"
	"project/handler"
	"project/middleware"
	"project/pkg/storage"

	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.Engine) {

	protected := r.Group("/api", middleware.VerifyAccessToken)
	{
		spaces, err := storage.NewSpacesClient(
			os.Getenv("SPACES_KEY"),
			os.Getenv("SPACES_SECRET"),
			os.Getenv("SPACES_REGION"),
			os.Getenv("SPACES_ENDPOINT"),
			os.Getenv("SPACES_BUCKET"),
		)
		if err != nil {
			fmt.Println("Failed to init Spaces:", err)
		}

		fileHandler := handler.NewFileHandler(spaces)
		protected.GET("/images/:filename", handler.ServeImage)
		protected.POST("/upload", fileHandler.Upload)
		protected.GET("/files/:filename", fileHandler.GetFile)
	}
}
