package models

import (
	"errors"
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	EmailVerification	bool	`gorm:"default:false" json:"email_verification"`
	DateVerification time.Time `gorm:"default:null" json:"date_verification"`
	TokenVerification string `gorm:"size:255;not null" json:"token_verification"`
	CodeVerification string `gorm:"size:255;not null" json:"code_verification"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ResultUser struct {
	ID	string `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
}


func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Username == "" {
			return errors.New("Username tidak boleh kosong!")
		}
		if u.Password == "" {
			return errors.New("Password tidak boleh kosong!")
		}
		if u.Email == "" {
			return errors.New("Email tidak boleh kosong!")
		}

		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Password tidak boleh kosong!")
		}
		if u.Email == "" {
			return errors.New("Email tidak boleh kosong!")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Username == "" {
			return errors.New("Username tidak boleh kosong!")
		}
		if u.Password == "" {
			return errors.New("Password tidak boleh kosong!")
		}
		if u.Email == "" {
			return errors.New("Email tidak boleh kosong!")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]ResultUser, error) {
	var err error
	users := []ResultUser{}
	err = db.Debug().Table("users").Select("id, email, username").Limit(100).Find(&users).Error
	if err != nil {
		return &[]ResultUser{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User tidak ditemukan")
	}
	return u, err
}

func (u *User) TokenMe(db *gorm.DB, uid uint32) (*ResultUser, error) {
	var err error
	users := ResultUser{}
	err = db.Debug().Table("users").Where("id = ?", uid).Take(&users).Error
	if err != nil {
		return &ResultUser{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &users, errors.New("User tidak ditemukan")
	}
	return &users, err
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"username":  u.Username,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) UpdateToken(db *gorm.DB, token string, uid uint32, code string) (*User, error)  {
	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"date_verification": time.Now(),
			"token_verification": token,
			"code_verification": code,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateVerifyUser(db *gorm.DB, token string, code string) (*User, error)  {
	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	user := User{}

	db.Where("token_verification = ?", token).First(&user)
	code_verify := user.CodeVerification
	data_now := time.Now()
	date := user.DateVerification

	hs := data_now.Sub(date).Hours()

	if hs > 1 {
		err := errors.New("gagal verifikasi")
		fmt.Println(err)
		return &User{}, err
	}
	if user.EmailVerification == true {
		err := errors.New("sudah verifikasi")
		fmt.Println(err)
		return &User{}, err
	}
	if code_verify != code{
		err := errors.New("code verifikasi")
		fmt.Println(err)
		return &User{}, err
	}

	db = db.Debug().Model(&User{}).Where("token_verification = ?", token).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"email_verification": true,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("token_verification = ?", token).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}
