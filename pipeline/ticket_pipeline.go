package pipeline

import (
	"ots/model"
	"ots/settings"
)

var TicketPipeline *Pipeline[*model.Ticket]

func GenerateTicketPipeline() {
	TicketPipeline = &Pipeline[*model.Ticket]{
		defsize: settings.MySettings.Get_PipelineSize(),
	}
}
