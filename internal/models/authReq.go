package models

import "html"

// easyjson:json
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// easyjson:json
type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// easyjson:json
type DummyLoginReq struct {
	Role string `json:"role"`
}

func (s *LoginReq) Sanitize() {
	s.Email = html.EscapeString(s.Email)
	s.Password = html.EscapeString(s.Password)
}

func (s *RegisterReq) Sanitize() {
	s.Email = html.EscapeString(s.Email)
	s.Password = html.EscapeString(s.Password)
}

func (s *DummyLoginReq) Sanitize() {
	s.Role = html.EscapeString(s.Role)
}
