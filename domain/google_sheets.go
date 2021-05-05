package domain

// Entity
type GoogleSheets struct {
	Row        int
	Ticket     string
	Status     string
	CurrentBug int
}

// Usecase
type GoogleSheetsUsecase interface {
	Tracking() ([]*GoogleSheets, error)
	Update(item *GoogleSheets, bug int) error
	Notify(items []GoogleSheets) error
}

// Repository
type GoogleSheetsRepository interface {
	Fetch() ([]*GoogleSheets, error)
	Update(item *GoogleSheets, bug int) error
}
