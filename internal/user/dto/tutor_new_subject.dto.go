package dto

import (
	"eduVix_backend/internal/common/types"
)

type NewSubjectDto struct {
	Price       int               `json:"price"`
	SubjectName types.SubjectName `json:"subject_name"`
}
