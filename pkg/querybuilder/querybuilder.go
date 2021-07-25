package querybuilder

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Sort struct {
	Name      string
	Ascending bool
}

func (s Sort) Direction() string {
	if s.Ascending {
		return "ASC"
	}

	return "DESC"
}

type Options struct {
	Query    *string                  `json:"q" form:"q" query:"q"`
	Page     *int                     `json:"page" form:"page" query:"page"`
	Size     *int                     `json:"size" form:"size" query:"size"`
	Filters  map[string]interface{} `json:"filter" form:"filter" query:"filter"`
	Sort     []Sort                   `json:"sort" form:"sort" query:"sort"`
	Includes []string                 `json:"includes" form:"include" query:"include"`
	Fields   map[string][]string      `json:"fields" form:"fields" query:"fields"`
	Errors   []error                  `json:"errors"`
	filterRegex *regexp.Regexp
	fieldsRegex *regexp.Regexp
}

func NewOptions() (*Options, error) {
	filterRegex, err := regexp.Compile(`filter\[(.+)\]`)
	if err != nil {
		return nil, err
	}
	fieldsRegex, err := regexp.Compile(`fields\[(.+)\]`)
	if err != nil {
		return nil, err
	}

	return &Options{
		filterRegex: filterRegex,
		fieldsRegex: fieldsRegex,
	}, nil
}

func (p *Options) setIncludes(queryParams url.Values) *Options {
	p.Includes = strings.Split(queryParams.Get("include"),",")
	return p
}

func (p *Options) setFilters(queryParams url.Values) *Options {
	if len(p.Filters) == 0 {
		p.Filters = make(map[string]interface{})
	}
	for k, val := range queryParams {
		result := p.filterRegex.FindStringSubmatch(k)
		if len(result) > 1 && len(val) > 0 {
			filterKey := result[1]
			p.Filters[filterKey] = p.simpleParseString(val[0])
		}
	}
	return p
}

func (p *Options) setFields(queryParams url.Values) *Options {
	if len(p.Fields) == 0 {
		p.Fields = make(map[string][]string)
	}
	for k, val := range queryParams {
		result := p.fieldsRegex.FindStringSubmatch(k)
		if len(result) > 1 {
			filterKey := result[1]
			if len(p.Fields[filterKey]) == 0 {
				p.Fields[filterKey] = []string{}
			}
			for _, item := range val {
				p.Fields[filterKey] = append(p.Fields[filterKey],strings.Split(item,",")...)
			}
		}
	}
	return p
}

func (p *Options) simpleParseString(item string) interface{} {
	num, err := strconv.Atoi(item)
	if err == nil {
		return num
	}

	temp := strings.TrimSpace(strings.ToLower(item))
	if  temp == "true" {
		return true
	}

	if temp == "false" {
		return false
	}

	return item
}

func (p *Options) setSort(queryParams url.Values) *Options {
	val := queryParams.Get("sort")
	if val != "" {
		p.addSort(val)
	}

	return p
}

func (p *Options) addSort(val string) {
	sortList := strings.Split(val, ",")
	for _, sortItem := range sortList {
		s := Sort{Ascending: true, Name: sortItem}
		if sortItem[:1] == "-" {
			s.Ascending = false
			s.Name = sortItem[1:]
		}
		p.Sort = append(p.Sort, s)
	}
}

func (p *Options) setSize(queryParams url.Values) *Options {
	val := queryParams.Get("size")
	if val != "" {
		sizeInt, err := strconv.Atoi(val)
		if err != nil {
			p.Errors = append(p.Errors, fmt.Errorf("size parse error: %w", err))
			return p
		}
		p.Size = &sizeInt
	}

	return p
}

func (p *Options) setPage(queryParams url.Values) *Options {
	val := queryParams.Get("page")
	if val != "" {
		pageInt, err := strconv.Atoi(val)
		if err != nil {
			p.Errors = append(p.Errors, fmt.Errorf("page parse error: %w", err))
			return p
		}
		p.Page = &pageInt
	}

	return p
}

func (p *Options) setQuery(queryParams url.Values) *Options {
	qVal := queryParams.Get("q")
	if qVal != "" {
		p.Query = &qVal
	}
	return p
}

func ParseUrl(originUrl string) (*Options, error) {
	uriParams, err := url.Parse(originUrl)
	if err != nil {
		return nil, err
	}
	p,err := NewOptions()
	if err != nil {
		return nil, err
	}
	queryParams := uriParams.Query()
	p.setQuery(queryParams)
	p.setPage(queryParams)
	p.setSize(queryParams)
	p.setSort(queryParams)
	p.setFilters(queryParams)
	p.setIncludes(queryParams)
	p.setFields(queryParams)

	return p, nil
}
