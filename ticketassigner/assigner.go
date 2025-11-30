package ticketassigner

import (
	"fmt"
	"log"
	"ots/model"
	"ots/mongo/dbops"
	"ots/pipeline"
	"ots/settings"
	"ots/util"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func New() *TickerAssigner {
	return &TickerAssigner{
		timeoutMin: settings.MySettings.Get_TicketAssignTimeoutMin(),
	}
}

type TickerAssigner struct {
	timeoutMin time.Duration
}

func (t *TickerAssigner) Init() {
	log.Printf(`
		Ticket Assigner initialized
		Assigner will run every %v
	`, t.timeoutMin)

	for {
		t.Run()
		time.Sleep(t.timeoutMin)
	}
}

func (t *TickerAssigner) Run() {
	log.Printf("Assigner called at: %v", time.Now())

	// Get the current ticket pipeline of reservoir size
	// And create a reservoir
	tickets := pipeline.TicketPipeline.GetFirstOf(settings.MySettings.Get_ReservoirSize())
	r := pipeline.NewReservoir[*model.Ticket](-1) // use default size
	r.Fill(tickets)

	// Size of reservoir
	rsize := r.Size()
	log.Printf("Size of reservoir - %d", rsize)

	if rsize <= 0 {
		log.Println("Skipping..")
		return
	}

	for _, ticket := range r.Get_Queue() {
		// Assign resolver
		resolverId, err := t.Assign(ticket)
		if err != nil {
			log.Println(err)
			continue
		}

		// Create new tracker
		err = t.AddTicketTracker(ticket.ID, resolverId)
		if err != nil {
			log.Println(err)
			continue
		}

		// Remove the item from reservoir to discard bin
		r.QueueToBin()
	}

	// Remove the processed tickets from the pipeline
	log.Printf("Pipeline size before emptying: %d", pipeline.TicketPipeline.Size())
	log.Printf("Reservoir bin size now: %d", r.BinSize())

	pipeline.TicketPipeline.EmptyUpto(r.BinSize())

	log.Printf("Pipeline size after emptying: %d", pipeline.TicketPipeline.Size())
	// Empty the reservoir bin
	r.EmptyBin()

	log.Printf("Reservoir bin size now: %d", r.BinSize())
}

func (t *TickerAssigner) Assign(ticket *model.Ticket) (primitive.ObjectID, error) {
	// Query the DB to get ticket trackers
	trackers := dbops.GetTicketTrackers()
	log.Printf("Found Ticket Trackers: %d", len(trackers))
	if len(trackers) <= 0 {
		return primitive.NilObjectID, fmt.Errorf("no trackers found. aborting..")
	}

	// Sort the trackers by resolver ID (A-Z) to group multiple trackers
	// of same resolver
	sort.Slice(trackers, func(i, j int) bool { return trackers[i].ResolverID.Hex() < trackers[j].ResolverID.Hex() })

	// From the sorted slice, count the tracker with
	// least number of resolvers (resolver ID)
	occurences := util.GetOccurences(trackers, func(a, b *model.TicketTracker) bool { return a.ResolverID.Hex() == b.ResolverID.Hex() })
	// Sort occurences (A-Z)
	sort.Slice(occurences, func(i, j int) bool { return occurences[i].Count < occurences[j].Count })
	// log.Println(" **** Occurences ****")
	// for _, o := range occurences {
	// 	log.Println(*o)
	// }
	// log.Println(" **** Occurences ****")

	if len(occurences) <= 0 {
		return primitive.NilObjectID, fmt.Errorf("no resolver found to assign. aborting..")
	}

	resolverId := occurences[0].Data.ResolverID // resolver with least occurences (which means least tickets assigned)

	// Query the DB to assign this resolver ID to the ticket
	_, err := dbops.AssignResolverToTicket(ticket, resolverId)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("error assigning resolver with ID - %s to ticket with ID - %s: %v", resolverId, ticket.ID, err)
	}

	log.Printf("Ticket assigned to resolver - %s", resolverId.Hex())
	return resolverId, nil
}

func (t *TickerAssigner) AddTicketTracker(ticketId, resolverId primitive.ObjectID) error {
	// Insert new entry in ticket tracker
	tickettracker, err := dbops.AddTicketTracker(ticketId, resolverId)
	if err != nil {
		return err
	}

	log.Printf("Ticker tracker created: %v", *tickettracker)
	return nil
}
