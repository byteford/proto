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

func clickedToObj(b []byte) (*click.Click, error) {
	clicked := &click.Click{}
	if err := proto.Unmarshal(b, clicked); err != nil {
		return nil, err
	}
	return clicked, nil
}
func clickedToByte(clicked *click.Click) ([]byte, error) {
	out, err := proto.Marshal(clicked)
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

func getUserFromId(id string) (*click.Click, error) {
	for _, user := range Users.User {
		if id == user.Id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("No user found")
}

func writeUser(w http.ResponseWriter, user *click.Click) error {
	out, err := clickedToByte(user)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func newUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errobj := click.Error{
			Name: "wrong method",
		}
		out, err := errorToByte(&errobj)
		if err != nil {
			fmt.Printf("error Marsheling error: %s \n", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write(out)
		if err != nil {
			fmt.Printf("error sending error: %s \n", err)
		}
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	clicked, err := clickedToObj(body)
	if err != nil {
		fmt.Printf("error in unmarshel: %s", err)
	}
	clicked.AmountClicked = 0
	clicked.Id = uuid.New().String()
	Users.User = append(Users.User, clicked)
	err = writeUser(w, clicked)
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
		user, err := getUserFromId(params.Get("id"))
		if err != nil {
			fmt.Printf("error finding user: %s \n", user)
			errobj := click.Error{
				Name: "no User found",
			}
			out, err := errorToByte(&errobj)
			if err != nil {
				fmt.Printf("error Marsheling error: %s \n", err)
			}
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write(out)
			if err != nil {
				fmt.Printf("error sending error: %s \n", err)
			}
		}
		writeUser(w, user)
	}
}
func router() {
	http.HandleFunc("/", root)
	http.HandleFunc("/new", newUser)
	fmt.Println("Starting server")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
func main() {
	Users = &click.Users{}
	router()
}
