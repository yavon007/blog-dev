package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage = 1
	DefaultSize = 10
	MaxSize     = 50
)

type Params struct {
	Page int
	Size int
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.Size
}

func FromQuery(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = DefaultPage
	}
	if size < 1 || size > MaxSize {
		size = DefaultSize
	}

	return Params{Page: page, Size: size}
}
