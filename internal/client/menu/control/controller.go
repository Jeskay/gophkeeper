package control

import (
	"context"
	authControl "gophkeeper/internal/client/authentication/control"
	abstraction "gophkeeper/internal/client/menu/abstraction"
)

type menuController struct {
	authController authControl.Controller
	repository     abstraction.Repository
}

func NewController(authController authControl.Controller) *menuController {
	return &menuController{authController: authController, repository: abstraction.NewRepository()}
}

func (c *menuController) IsAuthenticated() bool {
	return c.repository.GetToken() != ""
}

func (c *menuController) SaveToken(value string) {
	c.repository.SaveToken(value)
}

func (c *menuController) Authorize(ctx context.Context) context.Context {
	return c.authController.Authenticate(ctx, c.repository.GetToken())
}
