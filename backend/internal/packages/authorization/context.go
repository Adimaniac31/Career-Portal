package authorization

import (
	"errors"
)

type AuthContext struct {
	UserID    uint
	Role      string
	CollegeID *uint
}

func BuildAuthContext(claims map[string]interface{}) (*AuthContext, error) {
	// user_id (JWT numbers decode as float64)
	uid, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id claim")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, errors.New("invalid role claim")
	}

	var collegeID *uint
	if cid, ok := claims["college_id"].(float64); ok {
		c := uint(cid)
		collegeID = &c
	}

	return &AuthContext{
		UserID:    uint(uid),
		Role:      role,
		CollegeID: collegeID,
	}, nil
}
