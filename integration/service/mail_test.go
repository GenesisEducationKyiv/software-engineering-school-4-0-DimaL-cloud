package service

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/configs"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-DimaL-cloud/internal/service"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendEmails_Success(t *testing.T) {
	conf, err := configs.NewConfig("../../configs/config.yml")
	if err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	mailService := service.NewMailService(conf.Mail)
	err = mailService.SendEmails("Integration Test", "test sending emails", []string{conf.Mail.Username})
	assert.NoError(t, err)
}
