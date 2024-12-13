package delivery

import (
	"mine/internal"
	"mine/internal/models"
	"mine/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

func (s *appHandlers) RateLimit(ctx *fiber.Ctx) error {
	return ctx.Next()
}

func (s *appHandlers) CustomRateLimit(max int, expiration time.Duration, typeResponse string) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			keyLimit, _ := c.Locals("keyLimit").(string)
			return keyLimit
		},
		LimitReached: func(c *fiber.Ctx) error {
			if typeResponse == "LOCAL" {
				return c.JSON(models.RespLocal{
					StatusCode: internal.SysStatus.RateLimit.Status,
					Message:    internal.SysStatus.RateLimit.Msg,
				})
			} else {
				return c.JSON(models.Resp{
					Status: internal.SysStatus.RateLimit.Status,
					Msg:    internal.SysStatus.RateLimit.Msg,
				})
			}
		},
	})
}

func (s *appHandlers) RequireTokenWeb(ctx *fiber.Ctx) error {
	token := string(ctx.Request().Header.Peek("TOKEN"))
	ip := string(ctx.Request().Header.Peek("X-Real-Ip"))
	funcName := "RequireTokenWeb"
	clams := jwt.MapClaims{}
	resultFe := models.Resp{}
	var body interface{}
	ctx.BodyParser(&body)
	defer func() {
		internal.Log.Info(funcName, zap.Any("ip", ip), zap.Any("url", ctx.Context().URI()), zap.Any("authen", token), zap.Any("jwt", clams), zap.Any("resultFe", resultFe))
	}()
	if utils.IsEmpty(token) {
		resultFe = models.Resp{
			Status: s.stt.ErrMsgs.TokenRequired.Code,
			Msg:    s.stt.ErrMsgs.TokenRequired.Msg,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}
	decodeToken, err := jwt.ParseWithClaims(token, clams, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.stt.Keys.AuthClientKey), nil
	})
	if err != nil {
		jwtErr := err.(*jwt.ValidationError).Errors
		if jwtErr == jwt.ValidationErrorExpired {
			resultFe = models.Resp{
				Status: s.stt.ErrMsgs.TokenExpired.Code,
				Msg:    s.stt.ErrMsgs.TokenExpired.Msg,
			}
			return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
		}
	}

	if decodeToken == nil || !decodeToken.Valid {
		resultFe = models.Resp{
			Status: s.stt.ErrMsgs.InvalidToken.Code,
			Msg:    s.stt.ErrMsgs.InvalidToken.Msg,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}

	ctx.Locals("keyLimit", clams["customerId"])
	ctx.Locals("customer_id", clams["customerId"])
	ctx.Locals("customer_phone", clams["phone"])
	ctx.Locals("app_version", clams["appVersion"])
	ctx.Locals("access_token", clams["accessToken"])

	return ctx.Next()
}

func (s *appHandlers) RequireTokenLocal(ctx *fiber.Ctx) error {
	token := string(ctx.Request().Header.Peek("TOKEN"))
	resultFe := models.RespLocal{}
	funcName := "RequireTokenLocal"
	var body interface{}
	ctx.BodyParser(&body)
	defer func() {
		internal.Log.Info(funcName, zap.Any("url", ctx.Context().URI()), zap.Any("authen", token), zap.Any("result", resultFe))
	}()
	if utils.IsEmpty(token) {
		resultFe = models.RespLocal{
			StatusCode: s.stt.ErrMsgs.TokenRequired.Code,
			Message:    s.stt.ErrMsgs.TokenRequired.Msg,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}
	getToken := utils.CreateLocalToken(s.stt.Keys.AuthClientKey, s.stt.Keys.AuthClientKey)
	if token != getToken {
		detail := map[string]interface{}{}
		if !s.stt.Cfgs.UseProduction {
			detail["token"] = getToken
		}
		resultFe = models.RespLocal{
			StatusCode: s.stt.ErrMsgs.InvalidToken.Code,
			Message:    s.stt.ErrMsgs.InvalidToken.Msg,
			Data:       detail,
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(resultFe)
	}
	return ctx.Next()
}
