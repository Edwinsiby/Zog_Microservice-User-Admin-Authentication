package db

import (
	"errors"
	"log"
	"service1/pkg/db"
	"service1/pkg/entity"

	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init() {
	DB, err = db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
}

func GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	query := "SELECT * FROM users WHERE email = ?"
	result := DB.Raw(query, email).Scan(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
func GetByPhone(phone string) (*entity.User, error) {
	var user entity.User
	query := "SELECT * FROM users WHERE phone = ?"
	result := DB.Raw(query, phone).Scan(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func Create(user *entity.User) error {
	query := "INSERT INTO users (first_name, last_name, email, phone, password) VALUES (?, ?, ?, ?, ?)"
	return DB.Exec(query, user.FirstName, user.LastName, user.Email, user.Phone, user.Password).Error
}

func CreateSignup(user *entity.Signup) error {
	return DB.Create(user).Error
}

func CreateOtpKey(otpKey *entity.OtpKey) error {

	return DB.Create(otpKey).Error
}

func GetByKey(key string) (*entity.OtpKey, error) {
	var otpKey entity.OtpKey
	result := DB.Where(&entity.OtpKey{Key: key}).First(&otpKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &otpKey, nil
}

func GetSignupByPhone(phone string) (*entity.Signup, error) {
	var user entity.Signup
	result := DB.Where(&entity.Signup{Phone: phone}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// ------------------------------------------------------------------------------------------------------------
func CheckPermission(user *entity.User) (bool, error) {
	result := DB.Where(&entity.User{Phone: user.Phone}).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	permission := user.Permission
	return permission, nil
}

func CreateAddress(address *entity.Address) error {
	return DB.Create(address).Error
}

func AdminCreate(admin *entity.Admin) error {
	return DB.Create(admin).Error
}

func AdminGetByPhone(phone string) (*entity.Admin, error) {
	var admin entity.Admin
	result := DB.Where(&entity.Admin{Phone: phone}).First(&admin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &admin, nil
}
func AdminGetByEmail(email string) (*entity.Admin, error) {
	var admin entity.Admin
	result := DB.Where(&entity.Admin{Email: email}).First(&admin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &admin, nil
}
