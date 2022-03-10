package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/byteford/proto/modules/msgError"
	"github.com/byteford/proto/modules/user"
)

var Users *user.Users

//newUser is called when you go to url /new with a POST request.
//adds a user to the map
func newUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errmsg := msgError.Error{Name: "wrong method"}
		err := errmsg.Write(w)
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
		errmsg := msgError.Error{Name: err.Error()}
		err := errmsg.Write(w)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = user.Write(w)
	if err != nil {
		fmt.Printf("error sending user: %s", err)
	}
}

//root is called when you go to the url /
//if no params returns list of users. or returns user bases on id or name
func root(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	//get all users
	if len(params) == 0 {
		out, err := Users.ToByte()
		if err != nil {
			fmt.Printf("error in UsersTobyte %s", err)
		}
		w.Write(out)
	} else {
		//get user by id
		if id, ok := params["id"]; ok {
			user, err := Users.GetFromId(id[0])
			if err == nil {
				user.Write(w)
				return
			}
			//get user by name
		} else if name, ok := params["name"]; ok {
			user, err := Users.GetFromName(name[0])
			if err == nil {
				user.Write(w)
				return
			}
		}
		//return error if cant find user
		fmt.Printf("error finding user: %s \n", params)
		errmsg := msgError.Error{Name: "No user found"}
		err := errmsg.Write(w)
		if err != nil {
			fmt.Println(err)
		}

	}
}

//moveUser is called when you go to the url /move
//triggures to move user by given amount
func moveUser(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		move := &user.MoveUser{}
		err = move.ToObj(body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(move)
		user, err := Users.GetFromId(move.UserId)
		if err != nil {
			fmt.Println(err)
			errmsg := msgError.Error{Name: "No user found"}
			err := errmsg.Write(w)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		user.Move(move.Amount)

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
	Users = user.NewUsers()
	router()
}
