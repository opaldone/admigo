package users

import (
	"time"

	"admigo/models/mcom"
)

// SessionModel model for a session
type SessionModel struct {
	UID       string
	UserID    int
	CreatedAt time.Time
}

func sessionByUUID(uuid string) *SessionModel {
	que := `
		select uid, user_id, created_at
		from sessions
		where uid = $1
	`
	s := SessionModel{}
	err := mcom.Dbc.QueryRow(que, uuid).Scan(&s.UID, &s.UserID, &s.CreatedAt)
	if err != nil {
		return nil
	}

	return &s
}

func (s *SessionModel) user() *UserModel {
	u, err := getUser("where us.id = $1 and us.confirmed = 1", s.UserID, true)
	if err != nil {
		return nil
	}

	return u
}

// SessionUser users session
func SessionUser(uuid string) *UserModel {
	session := sessionByUUID(uuid)
	if session == nil {
		return nil
	}

	return session.user()
}
