package auth

type UserAuth struct {
	Email  string
	Name   string
	AuthID string // [deprecated], this field currently is firebaseUID, move to userAuth.ID instead
	ID     string // hex string of models.User
}

// TODO: add interface for role-based auth
type Manager interface {
	Verify(jwt string) (*UserAuth, error)
}
