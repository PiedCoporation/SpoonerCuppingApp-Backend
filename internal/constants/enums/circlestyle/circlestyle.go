package circlestyle

type CircleStyleEnum string

const (
	Default CircleStyleEnum = "DEFAULT"
	Second  CircleStyleEnum = "SECOND"
)

func (e CircleStyleEnum) IsValid() bool {
	switch e {
	case Default, Second:
		return true
	}
	return false
}
