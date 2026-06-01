package user

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

func (r Role) IsValid() bool {
	return r == RoleUser || r == RoleAdmin
}
