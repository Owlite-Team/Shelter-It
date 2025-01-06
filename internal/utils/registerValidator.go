package utils

import (
	"errors"
	"regexp"
	"shelter-it-be/internal/model/request"
)

type RegisterReq struct {
	*request.RegisterReq
}

func (req RegisterReq) Validate() error {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}
	return nil
}
