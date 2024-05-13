package utilities

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/orestiskosko/hours-app/templates"
)

func IsHtmxRequest(c echo.Context) bool {
	return c.Request().Header.Get("hx-request") == "true"
}

// This custom Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	isHxRequest := IsHtmxRequest(ctx)

	var template templ.Component
	if isHxRequest {
		template = t
	} else {
		template = templates.Layout(t, ctx.Path())
	}

	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	renderCtx := context.WithValue(ctx.Request().Context(), templates.HxRequestContextKey, isHxRequest)
	return template.Render(renderCtx, ctx.Response().Writer)
}
