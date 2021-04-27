package entities

type Ticket struct {
	Url      string
	Title    string
	BugCount int
	Bugs     []string
	Status   string
}
