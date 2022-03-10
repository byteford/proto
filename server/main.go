package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/byteford/goproto/github.com/byteford/goproto/modules/click"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var Users *click.Users

func UserMoveToObj(b []byte) (*click.MoveUser, error) {
	move := &click.MoveUser{}
	if err := proto.Unmarshal(b, move); err != nil {
		return nil, err
	}
	return move, nil
}

func UserToObj(b []byte) (*click.User, error) {
	User := &click.User{}
	if err := proto.Unmarshal(b, User); err != nil {
		return nil, err
	}
	return User, nil
}

func UserToByte(User *click.User) ([]byte, error) {
	out, err := proto.Marshal(User)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func usersToByte(users *click.Users) ([]byte, error) {
	out, err := proto.Marshal(users)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func errorToByte(error *click.Error) ([]byte, error) {
	out, err := proto.Marshal(error)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func getUserFromId(id string) (*click.User, error) {
	for _, user := range Users.User {
		if id == user.Id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("No user found")
}

func getUserFromName(name string) (*click.User, error) {
	for _, user := range Users.User {
		if name == user.Name {
			return user, nil
		}
	}
	return nil, fmt.Errorf("No user found")
}

func writeUser(w http.ResponseWriter, user *click.User) error {
	out, err := UserToByte(user)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func writeError(w http.ResponseWriter, name string) error {
	errobj := click.Error{
		Name: name,
	}
	out, err := errorToByte(&errobj)
	if err != nil {
		return fmt.Errorf("error Marsheling error: %s \n", err)
	}
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(out)
	if err != nil {
		return fmt.Errorf("error sending error: %s \n", err)
	}
	return nil
}

func addUser(user *click.User) (*click.User, error) {
	for _, u := range Users.User {
		if user.Name == u.Name {
			return nil, fmt.Errorf("User already exsists")
		}
	}
	user.AmountClicked = 0
	user.Id = uuid.New().String()
	user.Pos = &click.Position{X: 0, Y: 0}

	Users.User = append(Users.User, user)
	return user, nil

}
func newUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		err := writeError(w, "wrong method")
		if err != nil {
			fmt.Println(err)
		}
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	NewUser, err := UserToObj(body)
	if err != nil {
		fmt.Printf("error in unmarshel: %s", err)
	}
	user, err := addUser(NewUser)
	if err != nil {
		err := writeError(w, err.Error())
		if err != nil {
			fmt.Println(err)
		}
	}
	err = writeUser(w, user)
	if err != nil {
		fmt.Printf("error sending user: %s", err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if len(params) == 0 {
		out, err := usersToByte(Users)
		if err != nil {
			fmt.Printf("error in UsersTobyte %s", err)
		}
		w.Write(out)
	} else {
		if id, ok := params["id"]; ok {
			user, err := getUserFromId(id[0])
			if err != nil {
				fmt.Printf("error finding user: %s \n", user)
				err := writeError(w, "no user found")
				if err != nil {
					fmt.Println(err)
				}
			}
			writeUser(w, user)
		}
		if name, ok := params["name"]; ok {
			user, err := getUserFromName(name[0])
			if err != nil {
				fmt.Printf("error finding user: %s \n", user)
				err := writeError(w, "no user found")
				if err != nil {
					fmt.Println(err)
				}
			}
			writeUser(w, user)
		}

	}
}
func moveUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(body)
		move, err := UserMoveToObj(body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(move)
		user, err := getUserFromId(move.UserId)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(user)
		user.Pos.X = user.Pos.X + move.Pos.X
		user.Pos.Y = user.Pos.Y + move.Pos.Y
		fmt.Println(user)
		writeUser(w, user)
	}
}
func router() {
	http.HandleFunc("/", root)
	http.HandleFunc("/new", newUser)
	http.HandleFunc("/move", moveUser)
	fmt.Println("Starting server")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
func main() {
	Users = &click.Users{}
	router()
}
