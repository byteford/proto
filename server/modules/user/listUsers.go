package user

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func (list *Users) GetFromId(id string) (*User, error) {
	for _, user := range list.User {
		if id == user.Id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("no user found")
}

func (list *Users) GetFromName(name string) (*User, error) {
	for _, user := range list.User {
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
	for _, u := range list.User {
		if user.Name == u.Name {
			return nil, fmt.Errorf("user already exsists")
		}
	}
	user.AmountClicked = 0
	user.Id = uuid.New().String()
	user.Pos = &Vector2{X: 0, Y: 0}

	list.User = append(list.User, user)
	return user, nil

}
