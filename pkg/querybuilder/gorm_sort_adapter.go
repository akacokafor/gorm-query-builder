package querybuilder

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"reflect"
	"strings"
)

type GormAllowedSort interface {
	Names() []string
	Execute(db *gorm.DB, options OptionsInterface) error
}

func (g *GormAdapter) validateSorts(instance OptionsInterface) error {

	if len( g.sortWhitelist) == 0 {
		return nil
	}

	for _, entry := range g.sortWhitelist {
		log.Printf("entry = %s", reflect.TypeOf(entry).String())
		_, isString := entry.(string)
		_, isAllowedFilter := entry.(GormAllowedSort)
		if !isAllowedFilter && !isString {
			return errors.New("all sorts must be string or objects that implement GormAllowedSort")
		}
	}

	for _, name := range instance.GetSort() {
		if !g.isValidSortName(name) {
			return fmt.Errorf("invalid sort key %s, %w", name.GetName(), ErrInvalidSortQuery)
		}
	}

	return nil
}

func (g *GormAdapter) isValidSortName(name Sortable) bool {
	sortNames := g.getSortNames(g.sortWhitelist)
	for _, validKey := range sortNames {
		if name.GetName() == validKey {
			return true
		}
	}
	return false
}

func (g *GormAdapter) getSortNames(whitelist []interface{}) []string {
	var keys []string
	for _, entry := range whitelist {
		if val, ok := entry.(string); ok {
			keys = append(keys, val)
		}

		if val, ok := entry.(GormAllowedSort); ok {
			keys = append(keys, val.Names()...)
		}
	}
	return keys
}

func (g *GormAdapter) applySorts(instance OptionsInterface) error {

	sortableList := instance.GetSort()
	if g.defaultSort != nil && len(sortableList) == 0 {
		sortableList = append(sortableList, g.defaultSort)
	}

	if len(g.sortWhitelist) == 0 {
		orderStr := ""
		separator := ""
		for index, val := range sortableList {
			if index > 0 && separator != "," {
				separator = ","
			}
			direction := "ASC"
			if !val.IsAscending() {
				direction = "DESC"
			}
			orderStr = strings.TrimSpace(fmt.Sprintf("%s%s `%s` %s",orderStr, separator,val.GetName(),direction))
		}
		if orderStr != "" {
			g.db.Order(orderStr)
		}
		return nil
	}


	orderStr := ""
	separator := ""
	for index, sortEntry := range sortableList {
		if index > 0 && separator != "," {
			separator = ","
		}
		for _, sortWhiteListEntry := range g.sortWhitelist {
			if _k, ok := sortWhiteListEntry.(string); ok {
				if _k == sortEntry.GetName() {
					direction := "ASC"
					if !sortEntry.IsAscending() {
						direction = "DESC"
					}

					orderStr = strings.TrimSpace(fmt.Sprintf("%s%s `%s` %s",orderStr, separator,sortEntry.GetName(),direction))
					if orderStr[:1] == "," {
						orderStr = orderStr[1:]
					}
				}
			}

			if op, ok := sortWhiteListEntry.(GormAllowedSort); ok {
				for _, _k := range op.Names() {
					if _k == sortEntry.GetName() {
						if err := op.Execute(g.db,instance); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	if orderStr != "" {
		g.db.Order(orderStr)
	}
	return nil
}