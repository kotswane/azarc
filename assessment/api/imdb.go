package api

import "azarc/models"

type ImdbAPI interface {
	QueryIMDB(title string) (models.Plot, error)
}
