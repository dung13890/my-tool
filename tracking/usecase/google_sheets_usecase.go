package usecase

import (
	"github.com/dung13890/my-tool/domain"
)

type googleSheetsUsecase struct {
}

func NewGoogleSheetsUsecase() domain.GoogleSheetsUsecase {
	return &pherusaUsecase{}
}
