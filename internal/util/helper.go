package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FindAll[M any](c *fiber.Ctx, db *gorm.DB, preloads ...string) (int64, int64, int64, []M, error) {
	var models []M
	var total int64

	query := db.Model(new(M))

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	query = queryParams(query, c)

	if err := query.Count(&total).Error; err != nil {
		return 0, 0, 0, nil, err
	}

	sortBy := c.Query("sort_by", "id")
	sortOrder := c.Query("sort_order", "asc")

	if strings.ToLower(sortOrder) == "desc" {
		query = query.Order(fmt.Sprintf("%s desc", sortBy))
	} else {
		query = query.Order(fmt.Sprintf("%s asc", sortBy))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return 0, 0, 0, nil, err
	}

	return total, int64(page), int64(limit), models, nil
}

func FindAllByCondition[M any](c *fiber.Ctx, db *gorm.DB, column string, value interface{}, preloads ...string) (int64, int64, int64, []M, error) {
	var models []M
	var total int64

	query := db.Model(new(M))

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Where(fmt.Sprintf("%s = ?", column), value).Find(&models).Error; err != nil {
		return 0, 0, 0, nil, err
	}

	query = queryParams(query, c)

	if err := query.Count(&total).Error; err != nil {
		return 0, 0, 0, nil, err
	}

	sortBy := c.Query("sort_by", "id")
	sortOrder := c.Query("sort_order", "asc")

	if strings.ToLower(sortOrder) == "desc" {
		query = query.Order(fmt.Sprintf("%s desc", sortBy))
	} else {
		query = query.Order(fmt.Sprintf("%s asc", sortBy))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&models).Error; err != nil {
		return 0, 0, 0, nil, err
	}

	return total, int64(page), int64(limit), models, nil
}

func FindOne[M any](c *fiber.Ctx, db *gorm.DB, id int64, preloads ...string) (*M, error) {
	var models M
	query := db.Model(new(M))

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.First(&models, id).Error; err != nil {
		return nil, err
	}

	return &models, nil
}

func FindOneByCondition[M any](c *fiber.Ctx, db *gorm.DB, column string, value interface{}, preloads ...string) (*M, error) {
	var models M
	query := db.Model(new(M))

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Where(fmt.Sprintf("%s = ?", column), value).First(&models).Error; err != nil {
		return nil, err
	}

	return &models, nil
}

func UpdateOne[M any](c *fiber.Ctx, db *gorm.DB, id int64, data interface{}, preloads ...string) (*M, error) {
	var Models M

	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	err := db.First(&Models, id).Error
	if err != nil {
		return nil, err
	}

	err = db.Model(&Models).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, err
	}
	return &Models, nil
}

func queryParams(query *gorm.DB, c *fiber.Ctx) *gorm.DB {
	queryParams := c.Queries()

	for key, value := range queryParams {
		values := strings.Split(value, ",")
		switch {
		case strings.HasPrefix(key, "search[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "search["), "]")
			query = query.Where(fmt.Sprintf("%s ILIKE ?", colName), "%"+value+"%")

		case strings.HasPrefix(key, "search_cols[") && strings.HasSuffix(key, "]"):
			colNames := strings.Split(strings.TrimSuffix(strings.TrimPrefix(key, "search_cols["), "]"), "|")
			var orConditions []string
			var orArgs []interface{}
			for _, col := range colNames {
				orConditions = append(orConditions, fmt.Sprintf("%s ILIKE ?", col))
				orArgs = append(orArgs, "%"+value+"%")
			}
			query = query.Where(strings.Join(orConditions, " OR "), orArgs...)

		case strings.HasPrefix(key, "filter_not[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "filter_not["), "]")
			query = query.Where(fmt.Sprintf("%s NOT IN (?)", colName), values)

		case strings.HasPrefix(key, "filterrange[") && strings.HasSuffix(key, "]"):
			colName := strings.TrimSuffix(strings.TrimPrefix(key, "filterrange["), "]")
			rangeValues := strings.Split(value, "|")
			if len(rangeValues) == 2 {
				if rangeValues[0] != "-" {
					query = query.Where(fmt.Sprintf("%s >= ?", colName), rangeValues[0])
				}
				if rangeValues[1] != "-" {
					query = query.Where(fmt.Sprintf("%s <= ?", colName), rangeValues[1])
				}
			}

		case key != "page" && key != "limit" && key != "sort_by" && key != "sort_order":
			if len(values) == 1 {
				if strings.ToLower(values[0]) == "null" {
					query = query.Where(fmt.Sprintf("%s IS NULL", key))
				} else {
					query = query.Where(fmt.Sprintf("%s = ?", key), values[0])
				}
			} else {
				query = query.Where(fmt.Sprintf("%s IN (?)", key), values)
			}
		}
	}

	return query
}
