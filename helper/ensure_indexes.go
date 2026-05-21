package helper

import "ots/mongo"

func EnsureAllIndexes() []error {
	idxs := map[int]func() error{
		0: mongo.EnsureTicketIndexes,
		1: mongo.EnsureTicketTrackerIndexes,
		2: mongo.EnsureResolverIndexes,
		3: mongo.EnsureAdminIndexes,
	}

	errors := make([]error, 0, len(idxs))
	for _, fn := range idxs {
		if err := fn(); err != nil {
			errors = append(errors, err)
			continue
		}
	}

	return errors
}
