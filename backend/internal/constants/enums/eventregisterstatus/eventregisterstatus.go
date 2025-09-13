package eventregisterstatus

// RegisterStatusEnum represents the registration status of an event
// @Description Registration status of an event
type RegisterStatusEnum string

const (
	RegisterStatusEnumPending RegisterStatusEnum = "PENDING"
	RegisterStatusEnumAccepted RegisterStatusEnum = "ACCEPTED"
	RegisterStatusEnumFull RegisterStatusEnum = "FULL"
)