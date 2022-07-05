package paginator

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PageInfo struct {
	Page  int
	Limit int
}

func ExtractPageData(c *gin.Context) PageInfo {
	pageStr := c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "20"
	}
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)

	return PageInfo{
		Limit: limit,
		Page:  page,
	}
}

type Paginated[T any] struct {
	TotalPages   int
	NextPage     int
	PreviousPage int
	CurrentPage  int
	Data         []T
}

func (p Paginated[T]) Render() map[string]interface{} {
	return map[string]interface{}{
		"nextPage":      p.NextPage,
		"previousPage":  p.PreviousPage,
		"currentPage":   p.CurrentPage,
		"numberOfPages": p.TotalPages,
		"data":          p.Data,
	}
}

func Paginate[T any](pageInfo PageInfo, db *gorm.DB) (Paginated[T], error) {
	offset := (pageInfo.Page - 1) * pageInfo.Limit

	var empty []T

	var count int64
	res := db.Count(&count)

	if res.Error != nil {
		return Paginated[T]{}, res.Error
	}
	res = db.Offset(offset).Limit(pageInfo.Limit).Find(&empty)

	if res.Error != nil {
		return Paginated[T]{}, res.Error
	}

	payload := buildResponse[T](count, pageInfo)
	payload.Data = empty

	return payload, nil
}

func buildResponse[T any](count int64, pageInfo PageInfo) Paginated[T] {
	payload := Paginated[T]{}

	total := float64(count) / float64(pageInfo.Limit)

	payload.TotalPages = int(math.Ceil(total))
	payload.CurrentPage = pageInfo.Page

	if payload.CurrentPage < payload.TotalPages {
		payload.NextPage = payload.CurrentPage + 1
	}
	if payload.CurrentPage > 1 {
		payload.PreviousPage = payload.CurrentPage - 1
	}

	return payload
}
