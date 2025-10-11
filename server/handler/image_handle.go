package handler

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ServeImage xử lý yêu cầu trả về một file ảnh.
// Nó lấy tên file từ URL parameter và phục vụ file từ một thư mục cố định (ví dụ: "uploads").
func ServeImage(c *gin.Context) {
	// Lấy tên file từ URL. Ví dụ: /images/my-avatar.png -> filename = "my-avatar.png"
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image filename is required"})
		return
	}

	// Tạo đường dẫn an toàn đến file ảnh.
	// Điều này giúp ngăn chặn các cuộc tấn công "Directory Traversal".
	// Ví dụ: "uploads/my-avatar.png"
	// Bạn cần tạo thư mục "uploads" ở thư mục gốc của dự án.
	imagePath := filepath.Join("uploads", filename)

	// Phục vụ file. Gin sẽ tự động thiết lập Content-Type header phù hợp.
	// Nếu file không tồn tại, c.File() sẽ tự động trả về lỗi 404 Not Found.
	c.File(imagePath)
}
