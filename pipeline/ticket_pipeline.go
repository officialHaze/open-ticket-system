package pipeline

import (
	"ots/model"
	"ots/settings"
	"sync"
)

var TicketPipeline *Pipeline[*model.Ticket]

func GenerateTicketPipeline() {
	TicketPipeline = &Pipeline[*model.Ticket]{
		defsize: settings.MySettings.Get_PipelineSize(),
		mu:      &sync.Mutex{},
	}
}
