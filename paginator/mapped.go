package paginator

import "gorm.io/gorm"

type Builder[In any, Out any] func(In) Out

func MappedPaginate[T any, Out any](pageInfo PageInfo, db *gorm.DB, mapper Builder[T, Out]) (Paginated[Out], error) {
	offset := (pageInfo.Page - 1) * pageInfo.Limit

	var empty []T

	var count int64
	res := db.Count(&count)

	if res.Error != nil {
		return Paginated[Out]{}, res.Error
	}
	res = db.Offset(offset).Limit(pageInfo.Limit).Find(&empty)

	if res.Error != nil {
		return Paginated[Out]{}, res.Error
	}

	payload := buildResponse[Out](count, pageInfo)
	payload.Data = mapData(empty, mapper)

	return payload, nil
}

func mapData[In any, Out any](data []In, mapper Builder[In, Out]) []Out {
	var output []Out
	for _, value := range data {
		mappedValue := mapper(value)
		output = append(output, mappedValue)
	}

	return output
}
