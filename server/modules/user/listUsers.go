package user

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func NewUsers() *Users {
	u := new(Users)
	u.List = make(map[string]*User)
	return u
}

func (list *Users) GetFromId(id string) (*User, error) {
	if val, ok := list.List[id]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("no user found")
}

func (list *Users) GetFromName(name string) (*User, error) {
	for _, user := range list.List {
		if name == user.Name {
			return user, nil
		}
	}
	return nil, fmt.Errorf("no user found")
}

func (list *Users) ToByte() ([]byte, error) {
	out, err := proto.Marshal(list)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (list *Users) AddUser(user *User) (*User, error) {
	_, err := list.GetFromName(user.Name)
	if err == nil {
		return nil, fmt.Errorf("user already exsists")
	}
	user.AmountClicked = 0
	user.Id = uuid.New().String()
	user.Pos = &Vector2{X: 0, Y: 0}

	list.List[user.Id] = user
	return user, nil

}
