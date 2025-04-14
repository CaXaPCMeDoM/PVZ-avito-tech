package entity

import "errors"

type ReceptionsStatus string

const (
	InProgressStatus ReceptionsStatus = "in_progress"
	CloseStatus      ReceptionsStatus = "close"
)

var validReceptionsStatusMap = map[ReceptionsStatus]struct{}{
	InProgressStatus: {},
	CloseStatus:      {},
}

func (r ReceptionsStatus) IsValidReceptionsStatus() bool {
	_, exists := validReceptionsStatusMap[r]
	return exists
}

func (r ReceptionsStatus) ValidateReceptionsStatus() error {
	if !r.IsValidReceptionsStatus() {
		return errors.New("invalid city")
	}
	return nil
}
