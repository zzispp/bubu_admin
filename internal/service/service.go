package service

import (
	"bubu_admin/internal/repository"
	"bubu_admin/pkg/jwt"
	"bubu_admin/pkg/log"
	"bubu_admin/pkg/sid"
	"github.com/spf13/viper"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
	conf   *viper.Viper
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
	conf *viper.Viper,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		tm:     tm,
		conf:   conf,
	}
}
