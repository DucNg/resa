package modele

// Admin is the adminitrateur modele.
// Token is used to contrain the session token.
type Admin struct {
	IdAdmin int64
	Login   string
	Psw     string
	Token   string
}
