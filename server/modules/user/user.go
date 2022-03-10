package user

import (
	"net/http"

	"google.golang.org/protobuf/proto"
)

func (user *User) ToObj(b []byte) error {
	if err := proto.Unmarshal(b, user); err != nil {
		return err
	}
	return nil
}

func (user *User) ToByte() ([]byte, error) {
	out, err := proto.Marshal(user)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (user *User) Write(w http.ResponseWriter) error {
	out, err := user.ToByte()
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}
