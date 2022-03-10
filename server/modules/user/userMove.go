package user

import "google.golang.org/protobuf/proto"

func (move *MoveUser) ToObj(b []byte) error {
	err := proto.Unmarshal(b, move)
	if err != nil {
		return err
	}
	return nil
}
