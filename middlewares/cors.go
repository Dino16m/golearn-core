package middlewares

import (
	"github.com/dino16m/golearn-core/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CORSMiddleware gin.HandlerFunc

func New(cfg config.CORSConfig) CORSMiddleware {
	defaultCfg := cors.DefaultConfig()
	defaultCfg.AllowHeaders = append(defaultCfg.AllowHeaders, "X-CSRF-TOKEN")
	validCfg := updateCORSCfg(cfg, defaultCfg)
	return CORSMiddleware(cors.New(validCfg))
}

func updateCORSCfg(update config.CORSConfig, original cors.Config) cors.Config {
	original.AllowAllOrigins = update.AllowAllOrigins
	original.AllowCredentials = update.AllowCredentials
	original.AllowWildcard = update.AllowWildcard
	if len(update.AllowHeaders) > 0 {
		original.AllowHeaders = append(original.AllowHeaders, update.AllowHeaders...)
	}
	if len(update.AllowMethods) > 0 {
		original.AllowMethods = update.AllowMethods
	}
	if len(update.AllowOrigins) > 0 {
		original.AllowOrigins = update.AllowOrigins
	}
	if len(update.ExposeHeaders) > 0 {
		original.ExposeHeaders = update.ExposeHeaders
	}
	return original
}
