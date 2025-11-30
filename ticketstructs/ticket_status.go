package ticketstructs

type TicketStatus struct {
	Created    string
	Open       string
	InProgress string
	Closed     string
}

func GenerateTicketStatus() *TicketStatus {
	return &TicketStatus{
		Created:    "created",
		Open:       "open",
		InProgress: "in-progress",
		Closed:     "closed",
	}
}
