package postgres

import (
	"fmt"

	"github.com/trennepohl/istio-auth-poc/authorization/internal"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type postgres struct {
	db *gorm.DB
}

func New(settings *internal.DatabaseConfig) (internal.Database, error) {
	connection, err := openConnection(settings)
	if err != nil {
		return &postgres{}, err
	}

	return &postgres{db: connection}, nil
}

func openConnection(conf *internal.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", conf.DatabaseHost,
		conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseName, conf.DatabasePort)
	db, err := gorm.Open(driver.Open(dsn), &gorm.Config{})
	return db, err
}

func (p *postgres) RemoveRole(roleID uint64, userID uint64) error {
	return p.db.Model(&internal.User{ID: userID}).Association("Roles").Delete(&internal.Role{ID: roleID})
}

func (p *postgres) AssignRole(roleID uint64, userID uint64) error {
	return p.db.Model(&internal.User{ID: userID}).Association("Roles").Append(&internal.Role{ID: roleID})
}

func (p *postgres) GetRoles() ([]internal.Role, error) {
	roles := make([]internal.Role, 0)
	err := p.db.Find(&roles).Error
	return roles, err
}

func (p *postgres) GetUsers() ([]internal.User, error) {
	result := make([]internal.User, 0)
	err := p.db.Preload("Roles").Find(&result).Error
	return result, err
}

func (p *postgres) CreateRole(role internal.Role) error {
	return p.db.Create(&role).Error
}

func (p *postgres) UpdateUser(user internal.User) error {
	return p.db.Save(&user).Error
}

func (p *postgres) CreateUser(user *internal.User, options ...internal.UserOption) error {
	for _, opt := range options {
		opt(user)
	}

	return p.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error
}

func (p *postgres) Migrate(i ...interface{}) error {
	return p.db.AutoMigrate(i...)
}

func (p *postgres) GetUser(email string) (internal.User, error) {
	user := internal.User{}
	err := p.db.Preload("Roles").Where("email=?", email).First(&user).Error
	if err == nil {
		return user, err
	}

	if err.Error() == gorm.ErrRecordNotFound.Error() {
		return user, internal.ErrUserNotFound
	}

	return user, err
}
