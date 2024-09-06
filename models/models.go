package models

type UserId = int
type AliasId = int
type ResourceId = int

type Alias struct {
	Id         AliasId    `json:"id"`
	CreatorId  UserId     `json:"creatorId"`
	ResourceId ResourceId `json:"resourceId"`
	Name       string     `json:"name"`
}
type User struct {
	Id       UserId `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

const (
	HiddenPassword = "****"
)

func (u *User) HidePassword() {
	u.Password = HiddenPassword
}

func (u *User) ValidatePasssword(password string) bool {
	// TODO: implement hashing
	return u.Password == password
}

type Resource struct {
	Id      ResourceId `json:"id"`
	UserId  UserId     `json:"userId"`
	Name    string     `json:"name"`
	Content string     `json:"value"`
}
