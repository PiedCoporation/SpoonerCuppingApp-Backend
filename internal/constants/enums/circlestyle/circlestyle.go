package circlestyle

type CircleStyleEnum string

const (
	Minimal CircleStyleEnum = "MINIMAL"
	Expert  CircleStyleEnum = "EXPERT"
)

func (e CircleStyleEnum) IsValid() bool {
	switch e {
	case Minimal, Expert:
		return true
	}
	return false
}
