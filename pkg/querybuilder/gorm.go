package querybuilder

import (
	"errors"
	"gorm.io/gorm"
	"strings"
)

var (
	ErrInvalidFilterQuery = errors.New("filters contains an invalid filter")
	ErrInvalidSortQuery = errors.New("sorts contains an invalid sort")
)



type GormAdapter struct {
	db                  *gorm.DB
	filtersWhitelist    []interface{}
	sortWhitelist       []interface{}
	fieldsWhiteList     []interface{}
	includesWhitelist   []interface{}
	defaultSort         string
	defaultToPagination bool
	defaultPage         int
	defaultSize         int
	relationships       []string
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

	return g.Execute(optionsInstance)
}

func (g *GormAdapter) Execute(optionsInstance *Options) (*gorm.DB, error) {
	if err := g.validate(optionsInstance); err != nil {
		return g.db, err
	}
	if err := g.applyOptions(optionsInstance); err != nil {
		return g.db, err
	}

	return g.db, nil
}

func (g *GormAdapter) Paginate(optionsInstance *Options) (*gorm.DB, error) {

	g.defaultToPagination = true
	g.defaultPage = 1
	g.defaultSize = 30

	if _, err := g.Execute(optionsInstance); err != nil {
		return g.db, err
	}

	if err := g.applyPagination(optionsInstance); err != nil {
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

	if err := g.applyQuery(instance); err != nil {
		return err
	}

	if err := g.applySorts(instance); err != nil {
		return err
	}

	if err := g.applyIncludes(instance); err != nil {
		return err
	}

	return nil
}

func (g *GormAdapter) normalizeIncludeName(name string) string {
	var sb strings.Builder
	componentParts := strings.Split(name,".")
	for index, part := range componentParts {
		if len(part) > 0 {
			sb.WriteString(strings.ToUpper(part[:1]))
			sb.WriteString(strings.ToLower(part[1:]))
			if index != len(componentParts) - 1 {
				sb.WriteString(".")
			}
		}
	}
	return sb.String()
}

func (g *GormAdapter) applyPagination(instance *Options) error {
	if g.defaultToPagination && instance.Page == nil {
		instance.Page = &g.defaultPage
	}

	if instance.Page == nil {
		return nil
	}

	if instance.Size == nil {
		instance.Size = &g.defaultSize
	}

	page := *instance.Page
	size := *instance.Size
	offset := (page - 1) * size
	g.db.Offset(offset).Limit(size)
	return nil
}

func (g *GormAdapter) addRelationship(name string) {
	g.relationships = append(g.relationships, name)
}

func (g *GormAdapter) GetRelationships() []string {
	var relationships []string
	relationships = append(relationships,g.relationships...)
	return relationships
}

func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{db: db}
}


