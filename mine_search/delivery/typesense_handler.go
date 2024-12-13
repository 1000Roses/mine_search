package delivery

import (
	"encoding/json"
	"fmt"
	"mine/internal"
	"mine/internal/models"
	"mine/internal/repositories"
	"mine/internal/services"
	"mine/internal/settings"
	"mine/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type EventTypeSenseHandlers interface {
	TypeSenseTextSearchHandler(*fiber.Ctx) error
	// TypeSenseLocationSearch(*fiber.Ctx) error
}

type eventTypeSenseHandlers struct {
	stt  *settings.AppSettings
	svc  *services.AppServices
	repo *repositories.Repositories
}

func NewEventTypeSenseHandlers(
	appSettings *settings.AppSettings,
	appService *services.AppServices,
	repo *repositories.Repositories,
) EventTypeSenseHandlers {
	return &eventTypeSenseHandlers{
		stt:  appSettings,
		svc:  appService,
		repo: repo,
	}
}

func (tk *eventTypeSenseHandlers) TypeSenseTextSearchHandler(ctx *fiber.Ctx) error {
	startTime := time.Now()
	status, msg := internal.SysStatus.SystemBusy.Status, internal.SysStatus.SystemBusy.Msg

	var bodyData map[string]interface{} // Use map for parsing flexible JSON
	if err := ctx.BodyParser(&bodyData); err != nil {
		internal.Log.Error("Cannot parse request body", zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(models.Resp{
			Status: status,
			Msg:    msg,
		})
	}

	infoUser := &models.InfoUser{
		CustomerId:    ctx.Locals("customer_id").(string),
		CustomerPhone: ctx.Locals("customer_phone").(string),
	}

	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		internal.Log.Error("Failed to marshal request body", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Resp{
			Status: status,
			Msg:    msg,
		})
	}

	var dataInput models.Search
	if err := json.Unmarshal(jsonBody, &dataInput); err != nil {
		internal.Log.Error("Failed to unmarshal request data", zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Resp{
			Status: internal.SysStatus.WrongParams.Status,
			Msg:    internal.SysStatus.WrongParams.Msg,
		})
	}

	if err := models.Validate.Struct(&dataInput); err != nil {
		internal.Log.Error("Validation error", zap.Any("input", dataInput), zap.Error(err))
		return ctx.Status(fiber.StatusBadRequest).JSON(models.Resp{
			Status: internal.SysStatus.WrongParams.Status,
			Msg:    internal.SysStatus.WrongParams.Msg,
		})
	}

	result, err := tk.svc.EventTypeSenseService.SearchText(dataInput)
	if err != nil {
		tk.stt.Log.Error("Service error", zap.Error(err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.Resp{
			Status: status,
			Msg:    err.Error(),
		})
	}

	defer func() {
		// Log details to Kibana
		utils.Pool.Execute(&utils.SendLogToKibanaTask{
			BootstrapServer: tk.stt.Brokers,
			TopicName:       tk.stt.KafkaTopicName,
			Message: utils.KibanaMessage{
				ServiceName: internal.ServiceName,
				UserAgent:   string(ctx.Context().UserAgent()),
				FuncName:    fmt.Sprintf("%v", bodyData["TypeRequest"]),
				Token:       string(ctx.Request().Header.Peek("TOKEN")),
				Input:       bodyData,
				Output: fmt.Sprintf("CustomerId: %v, CustomerPhone: %v, Status: %v, Msg: %v, Detail: %v",
					infoUser.CustomerId, infoUser.CustomerPhone, status, msg, result),
				ExecutedTime: time.Since(startTime).Seconds(),
				Version:      utils.GetTimeUTC7().String(),
				Url:          ctx.Request().URI().String(),
			},
		})
		// Log to Kibana (all servers)
		if errPoolAll := utils.Pool.Execute(&utils.SendLogToKibanaAllTask{
			BootstrapServer: tk.stt.Brokers,
			TopicName:       tk.stt.KafkaTopicNameAll,
			Message: utils.KibanaMessageAll{
				Phone:        infoUser.CustomerPhone,
				CustomerId:   infoUser.CustomerId,
				Status:       status,
				FunctionName: "mine/v1/",
				ActionName:   fmt.Sprintf("%v", bodyData["TypeRequest"]),
				DateAction:   utils.GetStringTimeUTC7("Y-M-D H:M:S"),
				Url:          ctx.Request().URI().String(),
				Note:         fmt.Sprintf("Status: %v, Msg: %v", status, msg),
				TypeLog:      "Webkit",
				ProcessTime:  fmt.Sprintf("%f", time.Since(startTime).Seconds()),
				ServiceName:  internal.ServiceName,
			},
		}); errPoolAll != nil {
			tk.stt.Log.Error("Failed to log to Kibana (all)", zap.Error(errPoolAll))
		}
	}()

	return ctx.Status(fiber.StatusOK).JSON(models.Resp{
		Status: 1,
		Msg:    "Ok",
		Detail: result,
	})
}
