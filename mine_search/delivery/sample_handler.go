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

type EventSampleHandlers interface {
	SampleGroupAPIs(*fiber.Ctx) error
}

type eventSampleHandlers struct {
	stt  *settings.AppSettings
	svc  *services.AppServices
	repo *repositories.Repositories
}

func NewEventSampleHandlers(
	appSettings *settings.AppSettings,
	appService *services.AppServices,
	repo *repositories.Repositories,
) EventSampleHandlers {
	return &eventSampleHandlers{
		stt:  appSettings,
		svc:  appService,
		repo: repo,
	}
}

func (tk *eventSampleHandlers) SampleGroupAPIs(ctx *fiber.Ctx) error {
	startTime := time.Now()
	status, msg := tk.stt.ErrMsgs.Params.Code, tk.stt.ErrMsgs.Params.Msg

	var detail interface{}
	type reqBody struct {
		TypeRequest string      `json:"type_request"`
		Data        interface{} `json:"data"`
	}
	bodyData := &reqBody{}
	err := ctx.BodyParser(&bodyData)
	if err != nil {
		tk.stt.Log.Debug("Cannot parse", zap.Error(err))
		return ctx.Status(fiber.StatusOK).JSON(models.Resp{
			Status: tk.stt.ErrMsgs.Params.Code,
			Msg:    tk.stt.ErrMsgs.Params.Msg,
		})
	}
	app_version, ok := ctx.Locals("app_version").(string)
	if !ok {
		tk.stt.Log.Error("app_version error", zap.Any("data", ctx.Locals("app_version")))
		app_version = ""
	}
	infoUser := &models.InfoUser{
		CustomerId:    ctx.Locals("customer_id").(string),
		CustomerPhone: ctx.Locals("customer_phone").(string),
	}
	funcRefs := map[string]func(tk *eventSampleHandlers, reqData interface{}, infoUser *models.InfoUser) (status int, msg string, detail interface{}){
		"do_func_sample": DoFuncSample,
	}
	if _, ok := funcRefs[bodyData.TypeRequest]; ok {
		tk.stt.Log.Info("Call ", zap.String("type_request:", bodyData.TypeRequest),
			zap.Any("data", bodyData.Data))
		status, msg, detail = funcRefs[bodyData.TypeRequest](tk, bodyData.Data, infoUser)
		tk.stt.Log.Info("Resp ", zap.String("type_request: ", bodyData.TypeRequest), zap.Any("input", bodyData.Data),
			zap.Int("status", status), zap.String("msg", msg),
			zap.Any("detail", detail))

	}
	// TODO: Send to response to kibana server
	utils.Pool.Execute(&utils.SendLogToKibanaTask{
		BootstrapServer: tk.stt.Brokers,
		TopicName:       tk.stt.KafkaTopicName,
		Message: utils.KibanaMessage{
			ServiceName: internal.ServiceName,
			UserAgent:   string(ctx.Context().UserAgent()),
			FuncName:    bodyData.TypeRequest,
			Token:       string(ctx.Request().Header.Peek("TOKEN")),
			Input:       bodyData,
			Output: fmt.Sprintf("(customerId %v, customerPhone %v) --> status %v, msg %v, detail %v", infoUser.CustomerId, infoUser.CustomerPhone,
				status, msg, detail),
			ExecutedTime: time.Since(startTime).Seconds(),
			Version:      utils.GetTimeUTC7().String(),
			Url:          ctx.Request().URI().String(),
		},
	})
	// TODO: Send to response to kibana all server
	errPoolAll := utils.Pool.Execute(&utils.SendLogToKibanaAllTask{
		BootstrapServer: tk.stt.Brokers,
		TopicName:       tk.stt.KafkaTopicNameAll,
		Message: utils.KibanaMessageAll{
			Phone:        infoUser.CustomerPhone,
			CustomerId:   infoUser.CustomerId,
			AppVersion:   app_version,
			Status:       status,
			FunctionName: "mine/v1/",
			ActionName:   bodyData.TypeRequest,
			DateAction:   utils.GetStringTimeUTC7("Y-M-D H:M:S"),
			Url:          ctx.Context().URI().String(),
			Note:         fmt.Sprintf("status:%v,msg:%v", status, msg),
			TypeLog:      "Webkit",
			ProcessTime:  fmt.Sprintf("%f", time.Since(startTime).Seconds()),
			Topic_name:   tk.stt.KafkaTopicNameAll,
			ScreenId:     "",
			ServiceName:  internal.ServiceName,
		},
	})
	if errPoolAll != nil {
		tk.stt.Log.Error("Error Execute Pool Kibana All", zap.Error(errPoolAll))
	}
	return ctx.Status(fiber.StatusOK).JSON(models.Resp{
		Status: status,
		Msg:    msg,
		Detail: detail,
	})
}

func DoFuncSample(tk *eventSampleHandlers, reqData interface{}, customer *models.InfoUser) (status int, msg string, detail interface{}) {
	status, msg = tk.stt.ErrMsgs.Params.Code, tk.stt.ErrMsgs.Params.Msg
	_, ok := reqData.(map[string]interface{})
	if ok {
		jsonbody, err := json.Marshal(reqData)
		if err != nil {
			tk.stt.Log.Error("Error DoFuncSample", zap.Any("error", err), zap.Any("reqData", reqData))
			return status, msg, detail
		}
		type DataInput struct {
			Phone string `json:"phone" validate:"required"`
			Mail  string `json:"mail" validate:"required"`
		}
		dataInput := &DataInput{}
		if err := json.Unmarshal(jsonbody, &dataInput); err != nil {
			tk.stt.Log.Error("Error DoFuncSample", zap.Any("error", err))
			return status, "Không parse được dữ liệu", nil
		}
		validateError := models.Validate.Struct(dataInput)
		if validateError != nil {
			tk.stt.Log.Error("validateError", zap.Any("input", dataInput), zap.Error(validateError))
			status, msg = tk.stt.ErrMsgs.Params.Code, tk.stt.ErrMsgs.Params.Msg
			return
		}
		checkValidPhone := utils.CheckRegexFrType(dataInput.Phone, tk.stt.Regexs.RegexVNPhoneNumber)
		if !checkValidPhone {
			return 0, "số điện thoại không đúng định dạng", nil
		}
		result, err := tk.svc.EventSampleService.DoFuncSample(dataInput.Phone, dataInput.Mail)
		if err != nil {
			tk.stt.Log.Error("Error DoFuncSample", zap.Any("error", err))
			return 0, err.Error(), nil
		}
		return 1, "Ok", result
	}
	return status, msg, detail
}
