package querybuilder_test

import (
	"testing"

	"github.com/akacokafor/gorm-query-builder/pkg/querybuilder"
	"github.com/stretchr/testify/assert"
)

func TestParseUrl(t *testing.T) {
	type args struct {
		originUrl string
	}
	tests := []struct {
		name     string
		args     args
		validate func(t *testing.T, p *querybuilder.Options, err error)
	}{
		{
			name: "should successfully parse url wuth query",
			args: args{
				originUrl: "https://example.com?q=search",
			},
			validate: func(t *testing.T, p *querybuilder.Options, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, p)
				assert.NotNil(t, p.Query)
				assert.Equal(t, *p.Query, "search")
			},
		},
		{
			name: "should successfully parse url with page",
			args: args{
				originUrl: "https://example.com?q=search&page=2",
			},
			validate: func(t *testing.T, p *querybuilder.Options, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, p)
				assert.NotNil(t, p.Page)
				assert.NotNil(t, p.Query)
				assert.Equal(t, *p.Query, "search")
				assert.Equal(t, *p.Page, 2)
			},
		},
		{
			name: "should successfully parse url with page and size",
			args: args{
				originUrl: "https://example.com?q=search&page=2&size=15",
			},
			validate: func(t *testing.T, p *querybuilder.Options, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, p)
				assert.NotNil(t, p.Page)
				assert.NotNil(t, p.Query)
				assert.NotNil(t, p.Size)
				assert.Equal(t, *p.Query, "search")
				assert.Equal(t, *p.Page, 2)
				assert.Equal(t, *p.Size, 15)
			},
		},
		{
			name: "should successfully parse url with page and size and descending name sort",
			args: args{
				originUrl: "https://example.com?q=search&page=2&size=15&sort=-name",
			},
			validate: func(t *testing.T, p *querybuilder.Options, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, p)
				assert.NotNil(t, p.Page)
				assert.NotNil(t, p.Query)
				assert.NotNil(t, p.Size)
				assert.NotNil(t, p.Sort)
				assert.Equal(t, *p.Query, "search")
				assert.Equal(t, *p.Page, 2)
				assert.Equal(t, *p.Size, 15)
				assert.Len(t, p.Sort, 1)
				assert.Equal(t, p.Sort[0].GetName(), "name")
				assert.False(t, p.Sort[0].IsAscending())
			},
		},
		{
			name: "should successfully parse url with page and size and ascending name sort",
			args: args{
				originUrl: "https://example.com?q=search&page=2&size=15&sort=name",
			},
			validate: func(t *testing.T, p *querybuilder.Options, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, p)
				assert.NotNil(t, p.Page)
				assert.NotNil(t, p.Query)
				assert.NotNil(t, p.Size)
				assert.NotNil(t, p.Sort)
				assert.Equal(t, *p.Query, "search")
				assert.Equal(t, *p.Page, 2)
				assert.Equal(t, *p.Size, 15)
				assert.Len(t, p.Sort, 1)
				assert.Equal(t, p.Sort[0].GetName(), "name")
				assert.True(t, p.Sort[0].IsAscending())
			},
		},
		{
			name: "should successfully parse url with page and size and ascending name sort, descending age sort",
			args: args{
				originUrl: "https://example.com?q=search&page=2&size=15&sort=name,-age",
			},
			validate: func(t *testing.T, p *querybuilder.Options, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, p)
				assert.NotNil(t, p.Page)
				assert.NotNil(t, p.Query)
				assert.NotNil(t, p.Size)
				assert.NotNil(t, p.Sort)
				assert.Equal(t, *p.Query, "search")
				assert.Equal(t, *p.Page, 2)
				assert.Equal(t, *p.Size, 15)
				assert.Len(t, p.Sort, 2)
				assert.Equal(t, p.Sort[0].GetName(), "name")
				assert.True(t, p.Sort[0].IsAscending())
				assert.Equal(t, p.Sort[1].GetName(), "age")
				assert.False(t, p.Sort[1].IsAscending())
			},
		},
		{
			name: "should successfully parse url with page and size and ascending name sort, descending age sort and 'active' filtering",
			args: args{
				originUrl: "https://example.com?q=search&page=2&size=15&sort=name,-age&filter[active]=1",
			},
			validate: func(t *testing.T, p *querybuilder.Options, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, p)
				assert.NotNil(t, p.Page)
				assert.NotNil(t, p.Query)
				assert.NotNil(t, p.Size)
				assert.NotNil(t, p.Sort)
				assert.NotNil(t, p.Filters)
				assert.Equal(t, "search", *p.Query)
				assert.Equal(t, *p.Page, 2)
				assert.Equal(t, *p.Size, 15)
				assert.Len(t, p.Sort, 2)
				assert.Len(t, p.Filters, 1)
				assert.Equal(t, "name", p.Sort[0].GetName())
				assert.True(t, p.Sort[0].IsAscending())
				assert.Equal(t, "age", p.Sort[1].GetName())
				assert.False(t, p.Sort[1].IsAscending())
				assert.Equal(t, 1, p.Filters["active"])
			},
		},
		{
			name: "should successfully parse url with page and size and ascending name sort, descending age sort and 'active' filtering, includes as well",
			args: args{
				originUrl: "https://example.com?q=search&page=2&size=15&sort=name,-age&filter[active]=1&include=user,mobile",
			},
			validate: func(t *testing.T, p *querybuilder.Options, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, p)
				assert.NotNil(t, p.Page)
				assert.NotNil(t, p.Query)
				assert.NotNil(t, p.Size)
				assert.NotNil(t, p.Sort)
				assert.NotNil(t, p.Filters)
				assert.Equal(t, "search", *p.Query)
				assert.Equal(t, *p.Page, 2)
				assert.Equal(t, *p.Size, 15)
				assert.Len(t, p.Sort, 2)
				assert.Len(t, p.Filters, 1)
				assert.Equal(t, "name", p.Sort[0].GetName())
				assert.True(t, p.Sort[0].IsAscending())
				assert.Equal(t, "age", p.Sort[1].GetName())
				assert.False(t, p.Sort[1].IsAscending())
				assert.Equal(t, 1, p.Filters["active"])
				assert.Equal(t, "user", p.Includes[0])
				assert.Equal(t, "mobile", p.Includes[1])
			},
		},
		{
			name: "should successfully parse url with page and size and ascending name sort, descending age sort and 'active' filtering, includes, fields",
			args: args{
				originUrl: "https://example.com?q=search&page=2&size=15&sort=name,-age&filter[active]=1&include=user,mobile&fields[user]=id,name",
			},
			validate: func(t *testing.T, p *querybuilder.Options, err error) {
				assert.Nil(t, err)
				assert.NotNil(t, p)
				assert.NotNil(t, p.Page)
				assert.NotNil(t, p.Query)
				assert.NotNil(t, p.Size)
				assert.NotNil(t, p.Sort)
				assert.NotNil(t, p.Filters)
				assert.Equal(t, "search", *p.Query)
				assert.Equal(t, *p.Page, 2)
				assert.Equal(t, *p.Size, 15)
				assert.Len(t, p.Sort, 2)
				assert.Len(t, p.Filters, 1)
				assert.Equal(t, "name", p.Sort[0].GetName())
				assert.True(t, p.Sort[0].IsAscending())
				assert.Equal(t, "age", p.Sort[1].GetName())
				assert.False(t, p.Sort[1].IsAscending())
				assert.Equal(t, 1, p.Filters["active"])
				assert.Equal(t, "user", p.Includes[0])
				assert.Equal(t, "mobile", p.Includes[1])
				assert.Equal(t, "id", p.Fields["user"][0])
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := querybuilder.ParseUrl(tt.args.originUrl)
			tt.validate(t, got, err)
		})
	}
}
