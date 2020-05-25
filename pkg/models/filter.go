package models

import (
	"context"
	"net/url"
	"strconv"
	"time"

	grpc "google.golang.org/grpc/metadata"

	"github.com/jinzhu/gorm"
)

type FilterContextKey string

const (
	ErrBadFilter = ErrorBadRequest("Bad Filter Provided")
)

type Filter struct {
	DateFrom time.Time
	DateTo   time.Time
	Type     string
	Offset   uint
	Limit    uint
}

func getFirst(list []string) (string, error) {
	if len(list) != 1 {
		return "", ErrBadFilter
	}
	return list[0], nil
}

func NewFilterFromQueryString(queries url.Values) (*Filter, error) {
	f := &Filter{}
	if err := f.WithDateFrom(queries.Get("dateFrom")); err != nil {
		return nil, err
	}
	if err := f.WithDateTo(queries.Get("dateTo")); err != nil {
		return nil, err
	}
	f.WithType(queries.Get("type"))
	if err := f.WithOffset(queries.Get("offset")); err != nil {
		return nil, err
	}
	if err := f.WithLimit(queries.Get("limit")); err != nil {
		return nil, err
	}
	return f, nil
}

func NewFilterFromGRPCMetatdata(md grpc.MD) (*Filter, error) {
	f := &Filter{}
	dateFrom, err := getFirst(md.Get("dateFrom"))
	if err != nil {
		return nil, err
	}
	if err := f.WithDateFrom(dateFrom); err != nil {
		return nil, err
	}
	dateTo, err := getFirst(md.Get("dateTo"))
	if err != nil {
		return nil, err
	}
	if err := f.WithDateTo(dateTo); err != nil {
		return nil, err
	}
	Type, err := getFirst(md.Get("type"))
	if err != nil {
		return nil, err
	}
	f.WithType(Type)
	offset, err := getFirst(md.Get("offset"))
	if err != nil {
		return nil, err
	}
	if err := f.WithOffset(offset); err != nil {
		return nil, err
	}
	limit, err := getFirst(md.Get("limit"))
	if err != nil {
		return nil, err
	}
	if err := f.WithLimit(limit); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *Filter) WithDateFrom(dateFrom string) error {
	if dateFrom == "" {
		return nil
	}
	var err error
	f.DateFrom, err = time.Parse(time.RFC3339, dateFrom)
	if err != nil {
		return err
	}
	return nil
}
func (f *Filter) WithDateTo(dateTo string) error {
	if dateTo == "" {
		return nil
	}
	var err error
	f.DateTo, err = time.Parse(time.RFC3339, dateTo)
	if err != nil {
		return err
	}
	return nil
}

func (f *Filter) WithType(filterType string) {
	f.Type = filterType
}
func (f *Filter) WithOffset(offset string) error {
	if offset == "" {
		return nil
	}
	o, err := strconv.ParseUint(offset, 10, 64)
	if err != nil {
		return err
	}
	f.Offset = uint(o)
	return nil

}
func (f *Filter) WithLimit(limit string) error {
	if limit == "" {
		f.Limit = 100
		return nil
	}
	l, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return err
	}

	// Don't allow pages larger than 100
	if l > 100 {
		l = 100
	}

	f.Offset = uint(l)
	return nil
}

func (f Filter) ApplyDateFrom(db *gorm.DB) *gorm.DB {
	if !f.DateFrom.IsZero() {
		return db.Where("create_at > ?", f.DateFrom)
	}
	return db
}
func (f Filter) ApplyDateTo(db *gorm.DB) *gorm.DB {
	if !f.DateTo.IsZero() {
		return db.Where("create_at < ?", f.Type)
	}
	return db
}
func (f Filter) ApplyType(db *gorm.DB) *gorm.DB {
	if f.Type != "" {
		return db.Where("type = ?", f.Type)
	}
	return db
}
func (f Filter) ApplyOffset(db *gorm.DB) *gorm.DB {
	if f.Offset != 0 {
		return db.Offset(f.Offset)
	}
	return db
}
func (f Filter) ApplyLimit(db *gorm.DB) *gorm.DB {
	if f.Limit != 0 {
		return db.Limit(f.Limit)
	}
	// Always want a limit of 100
	return db.Limit(100)
}

func (f Filter) ApplyAll(db *gorm.DB) *gorm.DB {
	return db.Scopes(f.ApplyDateFrom,
		f.ApplyDateTo,
		f.ApplyType,
		f.ApplyOffset,
		f.ApplyLimit)
}

func ExtractFilterClaims(ctx context.Context) (*Filter, error) {
	filter, ok := ctx.Value(FilterContextKey("Filter")).(*Filter)
	if ok {
		return filter, nil
	}
	return nil, ErrBadFilter
}
