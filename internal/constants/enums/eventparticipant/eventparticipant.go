package eventparticipant

type TypeParticipantEnum string

const (
	TypeParticipantEnumRegistered TypeParticipantEnum = "JOINED"
	TypeParticipantEnumJoined TypeParticipantEnum = "REQUESTED"
	TypeParticipantEnumInvited TypeParticipantEnum = "INVITED"
)

func (t TypeParticipantEnum) String() string {
	return string(t)
}