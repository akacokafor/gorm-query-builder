package querybuilder_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/akacokafor/gorm-query-builder/pkg/querybuilder"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGormAdapter_ExecuteOnUrl(t *testing.T) {
	dialect := "sqlite"
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{DryRun: true})

	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		instance *querybuilder.GormAdapter
		db                *gorm.DB
		fieldsWhiteList   []interface{}
		includesWhitelist []interface{}
		filtersWhitelist  []interface{}
		sortWhitelist     []interface{}
		defaultSort       querybuilder.Sortable
	}
	type args struct {
		url string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		validator func(t *testing.T, f *fields, db *gorm.DB, err error)
	}{
		{
			name: "Should throw error when filter supplied is not in allowed filters",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: nil,
				filtersWhitelist: []interface{}{
					"status",
					"is_completed",
				},
				sortWhitelist: nil,
			},
			args: args{
				url: "https://example.com?filter[active]=1",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, querybuilder.ErrInvalidFilterQuery))
			},
		},
		{
			name: "Should throw error when one of the filters supplied is not in allowed filters",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: nil,
				filtersWhitelist: []interface{}{
					"status",
					"is_completed",
				},
				sortWhitelist: nil,
			},
			args: args{
				url: "https://example.com?filter[active]=1&filter[status]=completed",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, querybuilder.ErrInvalidFilterQuery))
			},
		},
		{
			name: "Should successfully append filter to query",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: nil,
				filtersWhitelist: []interface{}{
					"status",
					querybuilder.NewGormAllowedFilterExact("is_completed"),
				},
				sortWhitelist: nil,
			},
			args: args{
				url: "https://example.com?filter[is_completed]=1",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				stmt := db.Scan(&map[string]interface{}{}).Statement
				sqlString := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
				assert.Nil(t, err)
				assert.NotNil(t, stmt)
				assert.Contains(t, sqlString, "`is_completed` = 1")
			},
		},
		{
			name: "Should successfully accept GormAllowedFilterExact",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: nil,
				filtersWhitelist: []interface{}{
					"status",
					querybuilder.NewGormAllowedFilterExact("is_completed"),
				},
				sortWhitelist: nil,
				
			},
			args: args{
				url: "https://example.com?filter[is_completed]=1&filter[status]=completed",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				stmt := db.Scan(&map[string]interface{}{}).Statement
				sqlString := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
				assert.Nil(t, err)
				assert.NotNil(t, stmt)
				assert.Contains(t, sqlString, "`is_completed` = 1")
				if dialect == "mysql" {
					assert.Contains(t, sqlString, "`status` LIKE '%completed%'")
				}
				if dialect == "sqlite" {
					assert.Contains(t, sqlString, "`status` LIKE \"%completed%\"")
				}
			},
		},
		{
			name: "Should successfully accept GormAllowedFilterExact and pagination params",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: nil,
				filtersWhitelist: []interface{}{
					"status",
					querybuilder.NewGormAllowedFilterExact("is_completed"),
				},
				sortWhitelist: nil,

			},
			args: args{
				url: "https://example.com?filter[is_completed]=1&filter[status]=completed&page=2&size=10",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				stmt := db.Scan(&map[string]interface{}{}).Statement
				sqlString := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
				fmt.Printf(sqlString)
				assert.Nil(t, err)
				assert.NotNil(t, stmt)
				assert.Contains(t, sqlString, "`is_completed` = 1")
				if dialect == "mysql" {
					assert.Contains(t, sqlString, "`status` LIKE '%completed%'")
				}
				if dialect == "sqlite" {
					assert.Contains(t, sqlString, "`status` LIKE \"%completed%\"")
				}
					assert.Contains(t, sqlString, "OFFSET 10")
					assert.Contains(t, sqlString, "LIMIT 10")
			},
		},
		{
			name: "Should throw error when sort supplied is not in allowed sorts",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: nil,
				sortWhitelist: []interface{}{
					"created_at",
					"id",
				},
				
			},
			args: args{
				url: "https://example.com?sort=created_at,-age",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, querybuilder.ErrInvalidSortQuery))
			},
		},
		{
			name: "Should successfully append ascending sort to query",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: nil,
				sortWhitelist: []interface{}{
					"created_at",
					"id",
				},
				
			},
			args: args{
				url: "https://example.com?sort=created_at",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				stmt := db.Scan(&map[string]interface{}{}).Statement
				sqlString := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
				assert.Nil(t, err)
				assert.NotNil(t, stmt)
				assert.Contains(t, sqlString, "ORDER BY `created_at` ASC")
			},
		},
		{
			name: "Should successfully apply default sort to query",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: nil,
				sortWhitelist: []interface{}{
					"created_at",
					"id",
				},
				defaultSort: querybuilder.Sort{
					Name:      "created_at",
					Ascending: true,
				},
			},
			args: args{
				url: "https://example.com",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				stmt := db.Scan(&map[string]interface{}{}).Statement
				sqlString := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
				assert.Nil(t, err)
				assert.NotNil(t, stmt)
				assert.Contains(t, sqlString, "ORDER BY `created_at` ASC")
			},
		},
		{
			name: "Should successfully append ascending sort to query for custom sorter",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: nil,
				sortWhitelist: []interface{}{
					"created_at",
					querybuilder.NewGormAllowedSortCustom(
						"name_length",
						func(db *gorm.DB, ascending bool, propertyName string) error {
							direction := "ASC"
							if !ascending {
								direction = "DESC"
							}
							if propertyName == "name_length" {
								db.Order(fmt.Sprintf("LENGTH(%s) %s", "name", direction))
							}
							return nil
						}),
				},
				
			},
			args: args{
				url: "https://example.com?sort=-name_length",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				stmt := db.Scan(&map[string]interface{}{}).Statement
				sqlString := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
				assert.Nil(t, err)
				assert.NotNil(t, stmt)
				assert.Contains(t, sqlString, "ORDER BY LENGTH(name) DESC")
			},
		},
		{
			name: "Should successfully append ascending sort to query for custom sorter and normal sorter as well",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: nil,
				sortWhitelist: []interface{}{
					"created_at",
					querybuilder.NewGormAllowedSortCustom(
						"name_length",
						func(db *gorm.DB, ascending bool, propertyName string) error {
							direction := "ASC"
							if !ascending {
								direction = "DESC"
							}
							if propertyName == "name_length" {
								db.Order(fmt.Sprintf("LENGTH(%s) %s", "name", direction))
							}
							return nil
						}),
				},
				
			},
			args: args{
				url: "https://example.com?sort=-name_length,created_at",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				stmt := db.Scan(&map[string]interface{}{}).Statement
				sqlString := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
				assert.Nil(t, err)
				assert.NotNil(t, stmt)
				assert.Contains(t, sqlString, "ORDER BY LENGTH(name) DESC, `created_at` ASC")
			},
		},

		{
			name: "Should successfully append includes as preload options",
			fields: fields{
				db:                db,
				fieldsWhiteList:   nil,
				includesWhitelist: []interface{}{
					"wallet",
					"wallet.bank_account",
				},
				
			},
			args: args{
				url: "https://example.com?include=wallet,wallet.bank_account",
			},
			validator: func(t *testing.T, f *fields, db *gorm.DB, err error) {
				assert.Nil(t, err)
				assert.Contains(t, f.instance.GetRelationships(), "Wallet")
				assert.Contains(t, f.instance.GetRelationships(), "Wallet.BankAccount")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.db = tt.fields.db.Session(&gorm.Session{DryRun: true})
			g := querybuilder.NewGormAdapter(tt.fields.db.Table("users")).
				AllowedFilters(tt.fields.filtersWhitelist).
				AllowedSorts(tt.fields.sortWhitelist).
				AllowedIncludes(tt.fields.includesWhitelist).
				AllowedFields(tt.fields.fieldsWhiteList).
				DefaultSort(tt.fields.defaultSort)

			tt.fields.instance = g
			got, err := g.ExecuteOnUrl(tt.args.url)
			tt.validator(t, &tt.fields, got, err)
		})
	}
}
