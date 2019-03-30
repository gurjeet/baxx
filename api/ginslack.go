package main

import (
	"fmt"

	"net/http/httputil"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/jackdoe/baxx/message"
)

func SlackRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				//debug.PrintStack()
				httprequest, _ := httputil.DumpRequest(c.Request, false)
				stack := debug.Stack()
				log.Warnf("[Recovery] panic recovered:\n%s\n%s\n%s", string(httprequest), err, stack)
				if c.Request.RequestURI != "/digitalocean" {
					m := fmt.Sprintf("%s%s ```%s```", httprequest, err, stack)
					message.SendSlackDefault("panic", m)
				}
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
