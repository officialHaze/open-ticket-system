package ticketstructs

import "reflect"

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

func (t *TicketStatus) IsValidStatus(status string) bool {
	v := reflect.ValueOf(*t)

	for i := 0; i < v.NumField(); i++ {
		value := v.Field(i)
		v, ok := value.Interface().(string)
		if !ok {
			continue
		}

		if status == v {
			return true
		}
	}

	return false
}
