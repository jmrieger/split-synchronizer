package proxy

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/splitio/go-agent/splitio/proxy/middleware"
)

func Run(port string, adminPort string, apiKeys []string) {
	//gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	//CORS - Allows all origins
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"SplitSDKMachineName",
		"SplitSDKMachineIP",
		"SplitSDKVersion",
		"Authorization"}
	router.Use(cors.New(corsConfig))

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	//TODO add custom logger as middleware (?)
	router.Use(gin.Logger())
	router.Use(middleware.ValidateAPIKeys(apiKeys))

	go func() {
		adminRouter := gin.Default()
		// Admin routes
		admin := adminRouter.Group("/admin")
		{
			admin.GET("/ping", ping)
			admin.GET("/version", version)
			admin.GET("/uptime", uptime)
			admin.GET("/stats", showStats)
			admin.GET("/dashboard", showDashboard)
		}

		adminRouter.Run(adminPort)
	}()

	// API routes
	api := router.Group("/api")
	{
		api.GET("/splitChanges", splitChanges)
		api.GET("/segmentChanges/:name", segmentChanges)
		api.GET("/mySegments/:key", mySegments)
		api.POST("/testImpressions/bulk", postBulkImpressions)
		api.POST("/metrics/times", postMetricsTimes)
		api.POST("/metrics/counters", postMetricsCounters)
		api.POST("/metrics/gauge", postMetricsGauge)
		api.POST("/metrics/time", postMetricsTime)
		api.POST("/metrics/counter", postMetricsCounter)
	}
	router.Run(port)
}
