package controllers

import "math"

type Metadata struct {
	CurrentPage  int64 `json:"page"`
	Limit        int64 `json:"limit"`
	FirstPage    int64 `json:"first_page"`
	LastPage     int64 `json:"last_page"`
	TotalRecords int64 `json:"total_records"`
}

func computeMetadata(totalRecords, page, limit int64) *Metadata {
	if totalRecords == 0 {
		return &Metadata{}
	}
	return &Metadata{
		CurrentPage:  page,
		Limit:        limit,
		FirstPage:    1,
		LastPage:     int64(math.Ceil(float64(totalRecords) / float64(limit))),
		TotalRecords: totalRecords,
	}
}
