package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/byteford/proto/modules/user"
	"google.golang.org/protobuf/proto"
)

var Users *user.Users

func UserMoveToObj(b []byte) (*user.MoveUser, error) {
	move := &user.MoveUser{}
	if err := proto.Unmarshal(b, move); err != nil {
		return nil, err
	}
	return move, nil
}

func errorToByte(error *user.Error) ([]byte, error) {
	out, err := proto.Marshal(error)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func writeError(w http.ResponseWriter, name string) error {
	errobj := user.Error{
		Name: name,
	}
	out, err := errorToByte(&errobj)
	if err != nil {
		return fmt.Errorf("error Marsheling error: %s", err)
	}
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(out)
	if err != nil {
		return fmt.Errorf("error sending error: %s", err)
	}
	return nil
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
	NewUser := &user.User{}
	err = NewUser.ToObj(body)
	if err != nil {
		fmt.Printf("error in unmarshel: %s", err)
	}
	user, err := Users.AddUser(NewUser)
	if err != nil {
		err := writeError(w, err.Error())
		if err != nil {
			fmt.Println(err)
		}
	}
	err = user.Write(w)
	if err != nil {
		fmt.Printf("error sending user: %s", err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if len(params) == 0 {
		out, err := Users.ToByte()
		if err != nil {
			fmt.Printf("error in UsersTobyte %s", err)
		}
		w.Write(out)
	} else {
		if id, ok := params["id"]; ok {
			user, err := Users.GetFromId(id[0])
			if err != nil {
				fmt.Printf("error finding user: %s \n", user)
				err := writeError(w, "no user found")
				if err != nil {
					fmt.Println(err)
				}
			}
			user.Write(w)
		}
		if name, ok := params["name"]; ok {
			user, err := Users.GetFromName(name[0])
			if err != nil {
				fmt.Printf("error finding user: %s \n", user)
				err := writeError(w, "no user found")
				if err != nil {
					fmt.Println(err)
				}
			}
			user.Write(w)
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
		user, err := Users.GetFromId(move.UserId)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(user)
		user.Pos.X = user.Pos.X + move.Pos.X
		user.Pos.Y = user.Pos.Y + move.Pos.Y
		fmt.Println(user)
		user.Write(w)
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
	Users = &user.Users{}
	router()
}
