package repos

import (
	"pcdn-server/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

// tcRepo is the repository for the tc model
type tcRepo struct {
	DB *gorm.DB
}

// NewTcRepo creates a new tc repo
func NewTcRepo(db *gorm.DB) *tcRepo {
	return &tcRepo{DB: db}
}

// Create creates a new tc record
func (r *tcRepo) Create(tc *models.TcModel) error {
	tc.CreateTime = time.Now().UnixMilli()
	tc.UpdateTime = tc.CreateTime

	return r.DB.Create(tc).Error
}

func (r *tcRepo) Save(tc *models.TcModel) error {
	tc.CreateTime = time.Now().UnixMilli()
	tc.UpdateTime = time.Now().UnixMilli()

	// 先尝试创建记录
	err := r.DB.Create(tc).Error
	if err != nil {
		// 检查是否是唯一约束错误（SN重复）
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			// SN重复，执行更新操作
			tc.CreateTime = 0
			return r.DB.Where("sn = ?", tc.SN).Updates(tc).Error
		}
		// 其他错误直接返回
		return err
	}

	// 创建成功
	return nil
}

// GetBySN retrieves a tc record by SN
func (r *tcRepo) GetBySN(sn string) (*models.TcModel, error) {
	var tc models.TcModel
	err := r.DB.Where("sn = ?", sn).First(&tc).Error
	return &tc, err
}

// GetByTaskID retrieves a tc record by TaskID
func (r *tcRepo) GetByTaskID(taskID string) (*models.TcModel, error) {
	var tc models.TcModel
	err := r.DB.Where("task_id = ?", taskID).First(&tc).Error
	return &tc, err
}

// Update updates an existing tc record
func (r *tcRepo) Update(tc *models.TcModel) error {
	return r.DB.Save(tc).Error
}

// Delete deletes a tc record by ID
func (r *tcRepo) Delete(id uint) error {
	return r.DB.Delete(&models.TcModel{}, id).Error
}

// List retrieves a list of tc records with pagination
func (r *tcRepo) List(page, pageSize int) ([]models.TcModel, int64, error) {
	var tcs []models.TcModel
	var count int64

	db := r.DB.Model(&models.TcModel{})

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tcs).Error; err != nil {
		return nil, 0, err
	}

	return tcs, count, nil
}

// List retrieves a list of tc records with pagination
func (r *tcRepo) Find(offset, limit int) ([]models.TcModel, int64, error) {
	var tcs []models.TcModel
	var count int64

	db := r.DB.Model(&models.TcModel{})

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Offset(offset).Limit(limit).Find(&tcs).Error; err != nil {
		return nil, 0, err
	}

	return tcs, count, nil
}
