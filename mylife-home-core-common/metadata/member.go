package metadata

type MemberType string

const (
	Action MemberType = "action"
	State             = "state"
)

type Member struct {
	name        string
	description string
	memberType  MemberType
	valueType   Type
}

func (this *Member) Name() string {
	return this.name
}

func (this *Member) Description() string {
	return this.description
}

func (this *Member) MemberType() MemberType {
	return this.memberType
}

func (this *Member) ValueType() Type {
	return this.valueType
}
