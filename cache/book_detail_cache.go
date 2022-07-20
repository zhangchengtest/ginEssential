package cache

import (
	"errors"
	"ginEssential/dao"
	"time"

	"github.com/goburrow/cache"
	"github.com/sirupsen/logrus"
	"github.com/zhangchengtest/simple/sqls"

	"ginEssential/model"
)

type bookDetailCache struct {
	cache cache.LoadingCache // 标签缓存
}

var BookDetailCache = newBookDetailCache()

func newBookDetailCache() *bookDetailCache {
	return &bookDetailCache{
		cache: cache.NewLoadingCache(
			func(key cache.Key) (value cache.Value, e error) {
				value = dao.BookDetailDao.Get(sqls.DB(), key2Int64(key))
				if value == nil {
					e = errors.New("数据不存在")
				}
				return
			},
			cache.WithMaximumSize(1000),
			cache.WithExpireAfterAccess(30*time.Minute),
		),
	}
}

func (c *bookDetailCache) Get(bookDetailId int64) *model.BookDetail {
	val, err := c.cache.Get(bookDetailId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if val != nil {
		return val.(*model.BookDetail)
	}
	return nil
}

func (c *bookDetailCache) GetList(bookDetailIds []int64) (bookDetails []model.BookDetail) {
	if len(bookDetailIds) == 0 {
		return nil
	}
	for _, bookDetailId := range bookDetailIds {
		bookDetail := c.Get(bookDetailId)
		if bookDetail != nil {
			bookDetails = append(bookDetails, *bookDetail)
		}
	}
	return
}

func (c *bookDetailCache) Invalidate(bookDetailId int64) {
	c.cache.Invalidate(bookDetailId)
}
