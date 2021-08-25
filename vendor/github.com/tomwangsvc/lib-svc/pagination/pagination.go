package pagination

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_time "github.com/tomwangsvc/lib-svc/time"
)

// Constants for data pagination
const (
	DefaultLimit     = 20
	DefaultOffset    = 0
	DefaultOrder     = "DESC"
	DefaultStaleness = time.Minute * 30

	MaxLimit = 100
)

// Pagination data model
type Pagination struct {
	// Provided by the client
	Cursor  string `json:"cursor"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Order   string `json:"order"`
	OrderBy string `json:"order_by"`

	// Populated by the search
	ReadTimestamp *time.Time `json:"read_timestamp,omitempty"`

	// Populated by the search
	Total *int64 `json:"total,omitempty"`
}

// NewPagination passes struct of pagination
// -> provide acceptableStaleness duration if shouldn't be default, e.g. when Csr searching new claims
func newPagination(r *http.Request, acceptableStaleness *time.Duration, allowedLimit *int) (*Pagination, error) {

	p := Default()
	p.Cursor = strings.TrimSpace(r.URL.Query().Get("cursor"))
	p.OrderBy = strings.TrimSpace(r.URL.Query().Get("order_by"))

	if err := p.setLimit(strings.TrimSpace(r.URL.Query().Get("limit")), allowedLimit); err != nil {
		return nil, lib_errors.Wrap(err, "Failed setting pagination limit")
	}
	if err := p.setOffset(strings.TrimSpace(r.URL.Query().Get("offset"))); err != nil {
		return nil, lib_errors.Wrap(err, "Failed setting pagination offset")
	}
	if err := p.setOrder(strings.TrimSpace(r.URL.Query().Get("order"))); err != nil {
		return nil, lib_errors.Wrap(err, "Failed setting pagination order")
	}
	if err := p.checkStaleness(strings.TrimSpace(r.URL.Query().Get("read_timestamp")), acceptableStaleness); err != nil {
		return nil, lib_errors.Wrap(err, "Failed checking pagination staleness")
	}

	return p, nil
}

// NewPagination passes struct of pagination
// -> provide acceptableStaleness duration if shouldn't be default, e.g. when Csr searching new claims
func NewPagination(r *http.Request, acceptableStaleness *time.Duration) (*Pagination, error) {

	p, err := newPagination(r, acceptableStaleness, nil)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed creating new pagination")
	}

	return p, nil
}

// NewPagination passes struct of pagination with allowed limit
// -> provide acceptableStaleness duration if shouldn't be default, e.g. when Csr searching new claims
func NewPaginationWithAllowedLimit(r *http.Request, acceptableStaleness *time.Duration, allowedLimit int) (*Pagination, error) {

	p, err := newPagination(r, acceptableStaleness, &allowedLimit)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed creating new pagination with allowed limit")
	}

	return p, nil
}

// Default returns the default pagination object
func Default() *Pagination {

	return &Pagination{
		Limit:  DefaultLimit,
		Offset: DefaultOffset,
		Order:  DefaultOrder,
	}
}

func (p *Pagination) setLimit(limit string, allowedLimit *int) error {

	if limit != "" {
		l, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			return lib_errors.NewCustomf(http.StatusBadRequest, "Not recognized: limit = %s", limit)
		}
		if l >= 0 {
			if allowedLimit == nil {
				if l > MaxLimit {
					l = MaxLimit
				}
			} else {
				if l > int64(*allowedLimit) {
					l = int64(*allowedLimit)
				}
			}
			p.Limit = int(l)
		}
	}
	return nil
}

func (p *Pagination) setOffset(offset string) error {

	if offset != "" {
		o, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			return lib_errors.NewCustomf(http.StatusBadRequest, "Not recognized: o = %s", offset)
		}
		if o > 0 {
			p.Offset = int(o)
		}
	}
	return nil
}

func (p *Pagination) setOrder(order string) error {

	if order != "" {
		if order != "ASC" && order != "DESC" {
			return lib_errors.NewCustomf(http.StatusBadRequest, "Not recognized: order = %s, expected ASC or DESC", order)
		}
		p.Order = order
	}
	return nil
}

func (p *Pagination) checkStaleness(readTimestamp string, acceptableStaleness *time.Duration) error {

	if readTimestamp != "" {
		t, err := time.Parse(time.RFC3339, readTimestamp)
		if err != nil {
			return lib_errors.NewCustomf(http.StatusBadRequest, "Not recognized: read_timestamp = %s", readTimestamp)
		}
		staleIfBefore := time.Now().Add(-DefaultStaleness)
		if acceptableStaleness != nil {
			staleIfBefore = time.Now().Add(-*acceptableStaleness)
		}
		if lib_time.IsFirstTimeBeforeSecondTime(t, staleIfBefore) {
			p.resetStale()
		} else {
			p.ReadTimestamp = &t
		}
	}
	return nil
}

func (p *Pagination) resetStale() {

	p.Offset = 0
	p.ReadTimestamp = nil
	p.Total = nil
}
