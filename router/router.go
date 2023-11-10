package router

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type Router struct {
	r gin.IRouter
}

func New(r gin.IRouter, args ...interface{}) *Router {
	return &Router{r: r}
}

func (r *Router) Group(relativePath string, handlers ...gin.HandlerFunc) *Router {
	return &Router{r.r.Group(relativePath, handlers...)}
}

func (r *Router) GET(relativePath string, handler interface{}, handlers ...gin.HandlerFunc) *Router {
	r.r.GET(relativePath, append(handlers, r.handle(handler))...)
	return r
}

func (r *Router) POST(relativePath string, handler interface{}, handlers ...gin.HandlerFunc) *Router {
	r.r.POST(relativePath, append(handlers, r.handle(handler))...)
	return r
}

func (r *Router) handle(fun interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var buff [4096]byte
				n := runtime.Stack(buff[:], false)
				fmt.Printf("==> %s\n", string(buff[:n]))
				ctx.String(http.StatusInternalServerError, fmt.Sprint("err:", err))
			}
		}()

		var rst interface{}
		switch fn := fun.(type) {
		case func(c *gin.Context) interface{}:
			rst = fn(ctx)

		case func(c *Context) interface{}:
			rst = fn(&Context{Context: ctx})

		}

		switch v := rst.(type) {
		case error:
			ctx.JSON(http.StatusOK, gin.H{"code": 500, "msg": v.Error()})

		case func():
			v()

		default:
			result := gin.H{"code": 200}
			if v != nil {
				result["data"] = v
			}
			ctx.JSON(http.StatusOK, result)
		}
	}
}
