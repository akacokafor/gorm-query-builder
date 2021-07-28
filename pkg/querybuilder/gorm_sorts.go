package querybuilder

import "gorm.io/gorm"

type (
	GormAllowedSorter func(db *gorm.DB, ascending bool, propertyName string) error
	GormAllowedSortCustom struct {
		propName string
		sorter   GormAllowedSorter
	}
)

func (g *GormAllowedSortCustom) Names() []string {
	return []string{
		g.propName,
	}
}

func (g *GormAllowedSortCustom) Execute(db *gorm.DB, options OptionsInterface) error {
	for _, sort := range options.GetSort() {
		if sort.GetName() == g.propName {
			return g.sorter(db, sort.IsAscending(), g.propName)
		}
	}
	return nil
}

func NewGormAllowedSortCustom(propName string, sorter GormAllowedSorter) *GormAllowedSortCustom {
	return &GormAllowedSortCustom{
		propName: propName,
		sorter: sorter,
	}
}
