package eventparticipant

type TypeParticipantEnum string

const (
	TypeParticipantEnumRequested TypeParticipantEnum = "REQUESTED"
	TypeParticipantEnumJoined TypeParticipantEnum = "JOINED"
	TypeParticipantEnumInvited TypeParticipantEnum = "INVITED"
)

func (t TypeParticipantEnum) String() string {
	return string(t)
}