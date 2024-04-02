package service

import (
	"go_todo/src/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type SessionService struct {
	DB *gorm.DB
}

func (h *SessionService) SessionUser(token string) *model.User {
	session := model.Session{}
	h.DB.Where("id = ?", token).Find(&session)
	if session.ID == "" || session.ExpiresAt.Before(time.Now()) {
		return nil
	}

	sessionUser := model.User{}
	h.DB.Where("id = ?", session.UserID).Find(&sessionUser)

	return &sessionUser
}

func (h *SessionService) NewSession(
	user *model.User,
) string {
	session := model.Session{}
	session.ID = uuid.New().String()
	session.UserID = user.ID
	session.ExpiresAt = time.Now().Add(time.Hour * 24)
	h.DB.Create(&session)
	return session.ID
}

func (h *SessionService) DeleteSession(
	user *model.User,
	token string,
) bool {
	session := h.DB.Where("id = ?", token).Where("user_id = ?", user.ID).Find(
		&model.Session{},
	)
	if session.RowsAffected == 0 {
		return false
	}

	h.DB.Where("id = ?", token).Delete(&model.Session{})
	return true
}
