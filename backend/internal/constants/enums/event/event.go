package event

type RegisterStatusEnum string

const (
	RegisterStatusEnumPending RegisterStatusEnum = "PENDING"
	RegisterStatusEnumAccepted RegisterStatusEnum = "ACCEPTED"
	RegisterStatusEnumFull RegisterStatusEnum = "FULL"
)