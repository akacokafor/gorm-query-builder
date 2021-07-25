package querybuilder

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type GormAllowedSort interface {
	Names() []string
	Execute(db *gorm.DB, options *Options) error
}

func (g *GormAdapter) validateSorts(instance *Options) error {

	if len( g.sortWhitelist) == 0 {
		return nil
	}

	for _, entry := range g.sortWhitelist {
		_, isString := entry.(string)
		_, isAllowedFilter := entry.(GormAllowedSort)
		if !isAllowedFilter && !isString {
			return errors.New("all sorts must be string or objects that implement GormAllowedSort")
		}
	}

	for _, name := range instance.Sort {
		if !g.isValidSortName(name) {
			return fmt.Errorf("invalid sort key %s, %w", name.Name, ErrInvalidSortQuery)
		}
	}

	return nil
}

func (g *GormAdapter) isValidSortName(name Sort) bool {
	sortNames := g.getSortNames(g.sortWhitelist)
	for _, validKey := range sortNames {
		if name.Name == validKey {
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

func (g *GormAdapter) applySorts(instance *Options) error {

	if g.defaultSort != "" && len(instance.Sort) == 0 {
		instance.addSort(g.defaultSort)
	}

	if len(g.sortWhitelist) == 0 {
		orderStr := ""
		separator := ""
		for index, val := range instance.Sort {
			if index > 0 && separator != "," {
				separator = ","
			}
			orderStr = strings.TrimSpace(fmt.Sprintf("%s%s `%s` %s",orderStr, separator,val.Name,val.Direction()))
		}
		if orderStr != "" {
			g.db.Order(orderStr)
		}
		return nil
	}


	orderStr := ""
	separator := ""
	for index, sortEntry := range instance.Sort {
		if index > 0 && separator != "," {
			separator = ","
		}
		for _, sortWhiteListEntry := range g.sortWhitelist {
			if _k, ok := sortWhiteListEntry.(string); ok {
				if _k == sortEntry.Name {
					orderStr = strings.TrimSpace(fmt.Sprintf("%s%s `%s` %s",orderStr, separator,sortEntry.Name,sortEntry.Direction()))
					if orderStr[:1] == "," {
						orderStr = orderStr[1:]
					}
				}
			}

			if op, ok := sortWhiteListEntry.(GormAllowedSort); ok {
				for _, _k := range op.Names() {
					if _k == sortEntry.Name {
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