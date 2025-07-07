package dao

import (
	"github.com/mjcode-max/TurboGin/internal/model"
	"gorm.io/gorm"
)

// IUserDAO 用户数据操作接口
type IUserDAO interface {
	GetByID(id uint) (*model.User, error)
	Create(user *model.User) error
	// FindByName 扩展自定义查询方法
	FindByName(name string) ([]model.User, error) // 自定义查询
}

// UserDAO 实现 IUserDAO
type UserDAO struct {
	IBaseDAO[model.User] // 嵌入泛型 DAO
}

func NewUserDAO(db *gorm.DB) IUserDAO {
	return &UserDAO{
		IBaseDAO: NewBaseDAO[model.User](db), // 初始化泛型 DAO
	}
}

// FindByName 如果不借用IBaseDAO实现访问数据库，可以通过DB()获取db
func (u UserDAO) FindByName(name string) ([]model.User, error) {
	var entities []model.User
	err := u.DB().Where("name = ?", name).Find(&entities).Error
	if err != nil {
		return nil, err
	}
	return entities, nil
}
