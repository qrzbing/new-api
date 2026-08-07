package controller

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/system_setting"
	"github.com/stretchr/testify/assert"
)

func TestPaymentReturnPathUsesDefaultDashboardRoutes(t *testing.T) {
	previousAddress := system_setting.ServerAddress
	system_setting.ServerAddress = "https://dashboard.example.com/"
	t.Cleanup(func() { system_setting.ServerAddress = previousAddress })

	assert.Equal(
		t,
		"https://dashboard.example.com/wallet?pay=success",
		paymentReturnPath("/wallet?pay=success"),
	)
	assert.Equal(
		t,
		"https://dashboard.example.com/usage-logs",
		paymentReturnPath("/usage-logs"),
	)
}

func TestPaymentReturnPathUsesAppBasePathWithoutServerAddress(t *testing.T) {
	previousAddress := system_setting.ServerAddress
	previousBasePath := common.AppBasePath
	system_setting.ServerAddress = ""
	common.AppBasePath = "/new-api"
	t.Cleanup(func() {
		system_setting.ServerAddress = previousAddress
		common.AppBasePath = previousBasePath
	})

	assert.Equal(t, "/new-api/wallet?pay=success", paymentReturnPath("/wallet?pay=success"))
}
