package dao

//
//func (page *Page[T]) SelectPage(query interface{}, args []interface{}, order string) (e error) {
//	e = nil
//	var model T
//	if len(args) == 0 {
//		DB.Model(&model).Count(&page.Total)
//		if page.Total == 0 {
//			page.Data = []T{}
//			return
//		}
//		e = DB.Model(&model).Scopes(Paginate(page)).Order(order).Find(&page.Data).Error
//	} else {
//		DB.Model(&model).Where(query, args...).Count(&page.Total)
//		if page.Total == 0 {
//			page.Data = []T{}
//			return
//		}
//		e = DB.Model(&model).Where(query, args).Scopes(Paginate(page)).Order(order).Find(&page.Data).Error
//	}
//
//	return
//}

//func Paginate[T any](page *model.Page[T]) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if page.CurrentPage <= 0 {
//			page.CurrentPage = 0
//		}
//		switch {
//		case page.PageSize > 100:
//			page.PageSize = 100
//		case page.PageSize <= 0:
//			page.PageSize = 10
//		}
//		page.Pages = page.Total / page.PageSize
//		if page.Total%page.PageSize != 0 {
//			page.Pages++
//		}
//		p := page.CurrentPage
//		if page.CurrentPage > page.Pages {
//			p = page.Pages
//		}
//		size := page.PageSize
//		offset := int((p - 1) * size)
//		return db.Offset(offset).Limit(int(size))
//	}
//}
