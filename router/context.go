package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

func (c *Context) BindJSON(v any) error {
	err := c.Context.BindJSON(v)
	if err != nil {
		return err
	}
	return Validate(v)
}

func (c *Context) JSON(v any) {
	switch rst := v.(type) {
	case error:
		c.Context.JSON(http.StatusOK, gin.H{"code": 500, "msg": rst.Error()})

	default:
		result := gin.H{"code": 200}
		if rst == nil {
			result["data"] = rst
		}
		c.Context.JSON(http.StatusOK, result)
	}
}

func (c *Context) OffsetLimit() (int, int) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	return ((offset - 1) * limit), limit
}
