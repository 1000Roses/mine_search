package delivery

import (
	"mine/internal"
	"mine/internal/repositories"
	"mine/internal/services"
	"mine/internal/settings"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AppHandlers interface {
	// MiddleWare
	RateLimit(*fiber.Ctx) error
	CustomRateLimit(max int, expiration time.Duration, typeResponse string) fiber.Handler
	RequireTokenWeb(*fiber.Ctx) error
	RequireTokenLocal(*fiber.Ctx) error
	Health(*fiber.Ctx) error
	FreeOSMemory(*fiber.Ctx) error
	ReadMemStat(*fiber.Ctx) error
	EventSampleHandlers
}
type appHandlers struct {
	stt *settings.AppSettings
	svc *services.AppServices
	rdb *redis.Client
	EventSampleHandlers
}

func NewAppHandlers(
	appSettings *settings.AppSettings,
	appService *services.AppServices,
	rdb *redis.Client,
	repo *repositories.Repositories,
) AppHandlers {
	return &appHandlers{
		appSettings,
		appService,
		rdb,
		NewEventSampleHandlers(appSettings, appService, repo),
	}
}

func (bk *appHandlers) Health(ctx *fiber.Ctx) error {
	statusCode, result := 1, "Ok"
	return ctx.Status(statusCode).JSON(result)
}

func (bk *appHandlers) FreeOSMemory(ctx *fiber.Ctx) error {
	debug.FreeOSMemory()
	numGoroutines := runtime.NumGoroutine()
	internal.Log.Info("NumGoroutine", zap.Any("NumGoroutine", numGoroutines))
	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"NumGoroutine": numGoroutines,
	})
}

func (bk *appHandlers) ReadMemStat(ctx *fiber.Ctx) error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	type Temp struct {
		Alloc     float64
		Sys       float64
		HeapAlloc float64
		HeapSys   float64
		HeapIdle  float64
		HeapInuse float64
	}
	temp := Temp{
		Alloc:     float64(memStats.Alloc / 1048576),
		Sys:       float64(memStats.Sys / 1048576),
		HeapAlloc: float64(memStats.HeapAlloc / 1048576),
		HeapSys:   float64(memStats.HeapSys / 1048576),
		HeapIdle:  float64(memStats.HeapIdle / 1048576),
		HeapInuse: float64(memStats.HeapInuse / 1048576),
	}
	internal.Log.Info("ReadMemStat", zap.Any("stat", memStats))
	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"1memstats": temp,
		"2detail":   memStats,
	})
}
