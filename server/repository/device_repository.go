package repository

import (
	"project/database"
	"project/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateDevice tạo mới device
func CreateDevice(device *models.Device) (*models.Device, error) {
	if device.ID == uuid.Nil {
		device.ID = uuid.New()
	}

	if err := database.DB.Create(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

// GetDeviceByInfo tìm device theo UserID + IP + UserAgent
func GetDeviceByInfo(userID uuid.UUID, ip, userAgent string) (*models.Device, error) {
	var device models.Device
	err := database.DB.Table("devices").Joins("JOIN tokens on tokens.device_id = devices.id").Where("devices.user_id = ? AND devices.ip = ? AND devices.user_agent = ?", userID, ip, userAgent).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // chưa tồn tại
		}
		return nil, err
	}
	return &device, nil
}

// Optional: Get all devices of a user
func GetDevicesByUser(userID uuid.UUID) ([]*models.Device, error) {
	var devices []*models.Device
	if err := database.DB.Where("user_id = ?", userID).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
