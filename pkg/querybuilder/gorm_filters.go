package querybuilder

import (
	"fmt"
	"gorm.io/gorm"
)

type GormAllowedFilterExact struct {
	propName string

}

func (g *GormAllowedFilterExact) Keys() []string {
	return []string{g.propName}
}

func (g *GormAllowedFilterExact) Execute(db *gorm.DB, options *Options) error {
	val := options.Filters[g.propName]
	if val == nil {
		return nil
	}
	db.Where(g.propName,val)
	return nil
}

func NewGormAllowedFilterExact(propName string) *GormAllowedFilterExact {
	return &GormAllowedFilterExact{propName: propName}
}


type GormAllowedFilterSearch struct {
	propName string

}

func (g *GormAllowedFilterSearch) Keys() []string {
	return []string{g.propName}
}

func (g *GormAllowedFilterSearch) Execute(db *gorm.DB, options *Options) error {
	val := options.Filters[g.propName]
	if val == nil {
		return nil
	}
	sqlQuery := fmt.Sprintf("`%s` LIKE ?",g.propName)
	db.Where(sqlQuery,fmt.Sprintf("%%%s%%",val))
	return nil
}

func NewGormAllowedFilterSearch(propName string) *GormAllowedFilterSearch {
	return &GormAllowedFilterSearch{propName: propName}
}
