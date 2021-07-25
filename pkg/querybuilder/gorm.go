package querybuilder

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrInvalidFilterQuery = errors.New("filters contains an invalid filter")
	ErrInvalidSortQuery = errors.New("sorts contains an invalid sort")
)



type GormAdapter struct {
	db *gorm.DB
	fieldsWhiteList []interface{}
	includesWhitelist []interface{}
	filtersWhitelist []interface{}
	sortWhitelist []interface{}
	defaultSort string
}

//AllowedFilters white lists only the acceptable filters that can be applied from the query parameters
func (g *GormAdapter) AllowedFilters(filtersWhitelist []interface{}) *GormAdapter {
	g.filtersWhitelist = filtersWhitelist
	return g
}


//AllowedIncludes white lists only the acceptable includes that can be applied from the query parameters
func (g *GormAdapter) AllowedIncludes(includesWhitelist []interface{}) *GormAdapter {
	g.includesWhitelist = includesWhitelist
	return g
}

//AllowedIncludes white lists only the acceptable includes that can be applied from the query parameters
func (g *GormAdapter) AllowedFields(fieldsWhiteList []interface{}) *GormAdapter {
	g.fieldsWhiteList = fieldsWhiteList
	return g
}

//DefaultSort sets the default sort to apply on query
func (g *GormAdapter) DefaultSort(defaultSort string) *GormAdapter {
	g.defaultSort = defaultSort
	return g
}

//AllowedSorts white lists only the acceptable sort columns that can be applied from the query parameters
func (g *GormAdapter) AllowedSorts(sortWhitelist []interface{}) *GormAdapter {
	g.sortWhitelist = sortWhitelist
	return g
}

func (g *GormAdapter) ExecuteOnUrl(url string) (*gorm.DB, error) {
	optionsInstance, err := ParseUrl(url)
	if err != nil {
		return g.db, err
	}

	if err := g.validate(optionsInstance); err != nil {
		return g.db, err
	}

	if err := g.applyOptions(optionsInstance); err != nil {
		return g.db, err
	}

	return g.db, nil
}

func (g *GormAdapter) validate(instance *Options) error {
	if err := g.validateFilters(instance); err != nil {
		return err
	}

	if err := g.validateSorts(instance); err != nil {
		return err
	}

	return nil
}

func (g *GormAdapter) applyOptions(instance *Options) error {
	if err := g.applyFilters(instance); err != nil {
		return err
	}

	if err := g.applySorts(instance); err != nil {
		return err
	}

	return nil
}


func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{db: db}
}


