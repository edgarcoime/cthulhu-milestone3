package service

import (
	"github.com/cthulhu-platform/gateway/internal/service/auth"
	"github.com/cthulhu-platform/gateway/internal/service/diagnose"
)

type ServiceContainer struct {
	AuthService     auth.AuthService
	DiagnoseService diagnose.DiagnoseService
}