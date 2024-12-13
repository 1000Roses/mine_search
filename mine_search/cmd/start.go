package cmd

import (
	"fmt"
	"mine/delivery"
	"mine/internal"
	"mine/internal/repositories"
	"mine/internal/services"
	"mine/internal/settings"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func FreeOSMemory() {
	for {
		debug.FreeOSMemory()
		internal.Log.Info("FreeOSMemory", zap.Any("NumGoroutine", runtime.NumGoroutine()), zap.Any("NumCPU", runtime.NumCPU()))
		time.Sleep(10 * time.Second)
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long:  `Start cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		// Todo: Load application settings
		// init timezone
		os.Setenv("TZ", "Asia/Ho_Chi_Minh")
		// app setting
		appSettings := settings.NewAppSettings()
		mysqlDB := settings.NewSQLDB(appSettings.Cfgs)
		repositories := repositories.NewRepositories(mysqlDB, settings.NewLogger())
		rdbCache := settings.NewConnectRedis("", "", 0)
		if appSettings != nil {
			fmt.Println("INIT SERVICE...")
			// Todo: Init service
			// Todo: Compose into one app service
			appSvcs := services.NewAppServices(appSettings, repositories, rdbCache)
			// Todo: Init handler
			appHandler := delivery.NewAppHandlers(appSettings, appSvcs, rdbCache, repositories)
			// ticketHandler := delivery.NewTicketHandlers(appSettings, appSvcs, repositories)
			AppServer := fiber.New()
			AppServer.Use(cors.New(cors.Config{
				AllowOrigins: "*",
			}))
			AppServer.Use(logger.New(logger.Config{
				Format:     "${method} - ${path} - header:${reqHeaders} - body:${body} - resp-status:${status} - resp_body:${resBody}\n\n",
				TimeFormat: "2006-01-02 15:04:05",
				TimeZone:   "Asia/Ho_Chi_Minh",
			}))
			fmt.Println("Start Schedule...")
			Schedule(appSettings, appSvcs)
			fmt.Println("INIT ROUTE")
			AppServer.Get("/mine/health", appHandler.Health)
			AppServer.Post("/mine/v1/public/sample", appHandler.RequireTokenWeb, appHandler.RateLimit, appHandler.SampleGroupAPIs)
			// search text
			AppServer.Post("/mine/v1/public/typesense/text-search", appHandler.RequireTokenWeb, appHandler.RateLimit, appHandler.SampleGroupAPIs)
			AppServer.Post("/mine/v1/public/redis/text-search", appHandler.RequireTokenWeb, appHandler.RateLimit, appHandler.SampleGroupAPIs)
			// search location
			AppServer.Post("/mine/v1/public/typesense/location-search", appHandler.RequireTokenWeb, appHandler.RateLimit, appHandler.SampleGroupAPIs)
			AppServer.Post("/mine/v1/public/redis/location-search", appHandler.RequireTokenWeb, appHandler.RateLimit, appHandler.SampleGroupAPIs)
			// blooming filter

			// vector + advanced RAG

			fmt.Println("INIT ROUTE SUCCESS")
			if err := AppServer.Listen(":8386"); err != nil {
				fmt.Println("Fiber server got error ", err)
			}
			fmt.Println("App settings: ", appSettings)
		} else {
			fmt.Println("Error config!!!!!!! ")
		}
	},
}

func Schedule(appSettings *settings.AppSettings, appSvcs *services.AppServices) {
	// Hẹn lịch lại các deal sau khi chạy lại service
	go FreeOSMemory()
	// Cron
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	s := gocron.NewScheduler(loc)
	// s.Every(30).Seconds().Do(appSvcCallJobOrderHomeService, 30)
	s.StartAsync()
	fmt.Println("Start Schedule SUCCESS")
}

func GetIpKeyLimit(ctx *fiber.Ctx) error {
	ip := string(ctx.Request().Header.Peek("X-Real-Ip"))
	ctx.Locals("keyLimit", ip)
	return ctx.Next()
}

func init() {
	rootCmd.AddCommand(startCmd)
}
