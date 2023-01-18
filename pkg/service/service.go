package service

import (
	"errors"
	"fmt"
	"metalflow/pkg/global"
	"metalflow/pkg/response"
	"metalflow/pkg/utils"
	"reflect"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MysqlService struct {
	TX *gorm.DB // 事务对象实例
	DB *gorm.DB // 无事务对象实例
}

// New 初始化服务
func New(c *gin.Context) MysqlService {
	// 获取事务对象
	tx := global.GetTx(c)
	return MysqlService{
		TX: tx,
		DB: global.Mysql,
	}
}

var findCountCache = cache.New(5*time.Minute, 48*time.Hour) //nolint:gomnd

// Create 创建, model需使用指针, 否则可能无法插入数据
func (s *MysqlService) Create(req, model any) (err error) {
	i := model
	rv := reflect.ValueOf(req)

	if rv.Kind() == reflect.Struct && reflect.ValueOf(i).Kind() != reflect.Ptr {
		return errors.New("model must be a pointer")
	}

	if rv.Kind() == reflect.Slice {
		mv := reflect.Indirect(reflect.ValueOf(model))
		if mv.Kind() == reflect.Struct {
			slice := reflect.MakeSlice(reflect.SliceOf(mv.Type()), 0, 0)
			arr := reflect.New(slice.Type())
			i = arr.Interface()
		} else if mv.Kind() == reflect.Slice {
			slice := reflect.MakeSlice(mv.Type(), 0, 0)
			arr := reflect.New(slice.Type())
			i = arr.Interface()
		}
	}
	utils.Struct2StructByJson(req, i)
	// 创建数据
	err = s.TX.Create(i).Error
	return
}

// UpdateById 根据编号更新
func (s *MysqlService) UpdateById(id uint, req, model interface{}) error {
	// 获取model值
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || (rv.IsNil() || rv.Elem().Kind() != reflect.Struct) {
		return errors.New("model must be a pointer")
	}

	q := s.DB.Model(rv.Interface()).Where("id = ?", id).First(rv.Interface())
	if errors.Is(q.Error, gorm.ErrRecordNotFound) {
		return q.Error
	}

	m := make(map[string]interface{}, 0)
	utils.CompareDifferenceStruct2SnakeKeyByJson(rv.Elem().Interface(), req, &m)

	return q.Updates(&m).Error
}

// DeleteByIds 批量删除, model需使用指针, 否则可能无法插入数据
func (s *MysqlService) DeleteByIds(ids []uint, model any) (err error) {
	return s.TX.Where("id IN (?)", ids).Delete(model).Error
}

// Find 查询, model需使用指针, 否则可能无法绑定数据
//
//nolint:gocyclo
func (s *MysqlService) Find(query *gorm.DB, page *response.PageInfo, model any) (err error) { //nolint:gocyclo
	// 获取model值
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || (rv.IsNil() || rv.Elem().Kind() != reflect.Slice) {
		return fmt.Errorf("model must be a pointer")
	}

	countCache := false
	if page.CountCache != nil {
		countCache = *page.CountCache
	}
	if !page.NoPagination {
		if !page.SkipCount {
			// 查询条数
			fromCache := false
			// 以sql语句作为缓存键
			stmt := query.Session(&gorm.Session{DryRun: true}).Count(&page.Total).Statement
			cacheKey := s.TX.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
			if countCache {
				countCache, ok := findCountCache.Get(cacheKey)
				if ok {
					total, _ := countCache.(int64)
					page.Total = total
					fromCache = true
				}
			}
			if !fromCache {
				err = query.Count(&page.Total).Error
				if err == nil {
					findCountCache.Set(cacheKey, page.Total, cache.DefaultExpiration)
				}
			} else {
				global.Log.Debugf("条数缓存命中: %s, total: %d", cacheKey, page.Total)
			}
		}
		if page.Total > 0 || page.SkipCount {
			// 获取分页参数
			limit, offset := page.GetLimit()
			if page.LimitPrimary == "" {
				err = query.Limit(limit).Offset(offset).Find(model).Error
			} else {
				// 解析model
				if query.Statement.Model != nil {
					err = query.Statement.Parse(query.Statement.Model)
					if err != nil {
						return
					}
				}
				err = query.Joins(
					// 通过索引先分页再获取join其他字段, 以提高查询效率
					fmt.Sprintf(
						"JOIN (?) AS `OFFSET_T` ON `%s`.`id` = `OFFSET_T`.`%s`",
						query.Statement.Table,
						page.LimitPrimary,
					),
					query.
						Session(&gorm.Session{}).
						Select(
							fmt.Sprintf("`%s`.`%s`", query.Statement.Table, page.LimitPrimary),
						).
						Limit(limit).
						Offset(offset),
				).Find(model).Error
			}
		}
	} else {
		// 不使用分页
		err = query.Find(model).Error
		if err == nil {
			page.Total = int64(rv.Elem().Len())
			// 获取分页参数
			page.GetLimit()
		}
	}
	page.CountCache = &countCache
	return err
}
