package qr

import "github.com/exexeeeex/fiks-api/internal/domain"

type generateQrRequest struct {
	StudentId string `json:"student_id"`
}

type ReadQrResponse struct {
	Id    string       `json:"id"`
	Group domain.Group `json:"group"`
	Key   string       `json:"key"`
}
