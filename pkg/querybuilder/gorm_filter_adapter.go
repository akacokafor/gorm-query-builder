package querybuilder

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"reflect"
)

type GormAllowedFilter interface {
	Keys() []string
	Execute(db *gorm.DB, options OptionsInterface) error
}

func (g *GormAdapter) isValidFilterKey(key string) bool {
	filterKeys := g.getFilterKeys(g.filtersWhitelist)
	for _, validKey := range filterKeys {
		if key == validKey {
			return true
		}
	}
	return false
}

func (g *GormAdapter) getFilterKeys(whitelist []interface{}) []string {
	var keys []string
	for _, entry := range whitelist {
		if val, ok := entry.(string); ok {
			keys = append(keys, val)
		}

		if val, ok := entry.(GormAllowedFilter); ok {
			keys = append(keys, val.Keys()...)
		}
	}
	return keys
}

func (g *GormAdapter) validateFilters(instance OptionsInterface) error {
	if len(g.filtersWhitelist) == 0 {
		return nil
	}

	for _, entry := range g.filtersWhitelist {
		log.Printf("entry = %s", reflect.TypeOf(entry).String())
		_, isString := entry.(string)
		_, isAllowedFilter := entry.(GormAllowedFilter)

		log.Printf("entry = %s isString = %v, isAllowedFilter = %v", reflect.TypeOf(entry).String(), isString, isAllowedFilter)
		if !isAllowedFilter && !isString {
			return errors.New("all filters must be string or objects that implement GormAllowedFilter")
		}
	}

	for key, _ := range instance.GetFilters() {
		if !g.isValidFilterKey(key) {
			return fmt.Errorf("invalid filter key %s, %w", key, ErrInvalidFilterQuery)
		}
	}

	return nil
}

func (g *GormAdapter) applyFilters(instance OptionsInterface) error {
	if len(g.filtersWhitelist) == 0 {
		for key, _ := range instance.GetFilters() {
			if err := NewGormAllowedFilterSearch(key).Execute(g.db, instance); err != nil {
				return err
			}
		}
		return nil
	}

	for suppliedFilterKey, _ := range instance.GetFilters() {
		for _, whiteListFilterEntry := range g.filtersWhitelist {
			if _k, ok := whiteListFilterEntry.(string); ok {
				if _k == suppliedFilterKey {
					if err := NewGormAllowedFilterSearch(_k).Execute(g.db, instance); err != nil {
						return err
					}
				}
			}

			if op, ok := whiteListFilterEntry.(GormAllowedFilter); ok {
				for _, _k := range op.Keys() {
					if _k == suppliedFilterKey {
						if err :=  op.Execute(g.db,instance); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func (g *GormAdapter) applyQuery(instance OptionsInterface) error {
	if len(g.filtersWhitelist) == 0 || instance.GetQuery() == nil {
		return nil
	}

	var q *gorm.DB
	for _, whiteListFilterEntry := range g.filtersWhitelist {
		if _k, ok := whiteListFilterEntry.(string); ok {
			if q == nil {
				q = g.db.Or(fmt.Sprintf("%s LIKE ?",_k),fmt.Sprintf("%%%s%%",*instance.GetQuery()))
			}else {
				q = q.Or(fmt.Sprintf("%s LIKE ?",_k),fmt.Sprintf("%%%s%%",*instance.GetQuery()))
			}
		}
	}

	g.db.Where(q)

	return nil
}
