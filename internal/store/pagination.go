package store

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"alpha,oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}

		fq.Offset = o
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		fq.Since = parseTime(since)
	}

	until := qs.Get("until")
	if until != "" {
		fq.Until = parseTime(until)
	}

	return fq, nil
}

// func parseTime(s string) string {
// 	t, err := time.Parse("2006-01-02 15:04:05-07", s)
// 	if err != nil {
// 		return ""
// 	}
// 	return t.Format(time.RFC3339)

// 	// return t.Format(time.DateOnly)
// }

func parseTime(s string) string {

	format := "2006-01-02 15:04:05-07"
	// Try to parse with date and time first
	t, err := time.Parse("2006-01-02 15:04:05-07", s)
	if err == nil {
		// return t.Format(time.RFC3339)
		return t.Format(format)
	}

	// Try to parse with date and time only, without timezone
	t, err = time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		// Default timezone to UTC if no timezone is provided in input
		t = t.UTC()
		return t.Format(format)
	}

	// If the above fails, try parsing as a date only (without time)
	t, err = time.Parse("2006-01-02", s)
	if err == nil {
		// Default time to midnight (00:00:00) in UTC
		t = t.UTC()
		// return t.Format(time.RFC3339)
		return t.Format(format)
	}

	// If both parsing attempts fail, return an empty string or handle error
	return ""
}

func (fq *PaginatedFeedQuery) ValidateDates() error {
	// Define the expected layout
	layout := "2006-01-02 15:04:05-07"

	// Parse the since date
	sinceDate, err := time.Parse(layout, fq.Since)
	if err != nil {
		return fmt.Errorf("invalid since date format: %v", err)
	}

	// Parse the until date
	untilDate, err := time.Parse(layout, fq.Until)
	if err != nil {
		return fmt.Errorf("invalid until date format: %v", err)
	}

	// Check if since is before until
	if sinceDate.After(untilDate) {
		return fmt.Errorf("since date (%s) must be before until date (%s)", fq.Since, fq.Until)
	}

	return nil
}
