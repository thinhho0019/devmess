package repository

import (
	"errors"
	"fmt"
	"project/database"
	"project/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	CreateDevice(device *models.Device) (*models.Device, error)
	GetDeviceByInfo(userID uuid.UUID, ip, userAgent string) (*models.Device, error)
	GetDevicesByUser(userID uuid.UUID) ([]*models.Device, error)
	UpdateDevice(device *models.Device) (*models.Device, error)
}
type deviceRepo struct {
	db *gorm.DB
}

var NewDeviceRepository = func() DeviceRepository {
	return &deviceRepo{
		db: database.DB,
	}
}

// CreateDevice tạo mới device
func (r *deviceRepo) CreateDevice(device *models.Device) (*models.Device, error) {
	if device.ID == uuid.Nil {
		device.ID = uuid.New()
	}

	if err := r.db.Create(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}
func (r *deviceRepo) UpdateDevice(device *models.Device) (*models.Device, error) {
	if device == nil {
		return nil, errors.New("device is Nil")
	}

	if err := r.db.Save(device).Error; err != nil {
		return nil, fmt.Errorf("fail update device %s", err)
	}
	return device, nil
}

// GetDeviceByInfo tìm device theo UserID + IP + UserAgent
func (r *deviceRepo) GetDeviceByInfo(userID uuid.UUID, ip, userAgent string) (*models.Device, error) {
	var device models.Device
	err := r.db.Table("devices").Joins("JOIN tokens on tokens.device_id = devices.id").Where("devices.user_id = ? AND devices.ip = ? AND devices.user_agent = ?", userID, ip, userAgent).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // chưa tồn tại
		}
		return nil, err
	}
	return &device, nil
}

// Optional: Get all devices of a user
func (r *deviceRepo) GetDevicesByUser(userID uuid.UUID) ([]*models.Device, error) {
	var devices []*models.Device
	if err := r.db.Where("user_id = ?", userID).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
