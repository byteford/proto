package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/byteford/goproto/github.com/byteford/goproto/modules/click"
	"google.golang.org/protobuf/proto"
)

func toClicked(b []byte) (*click.Click, error) {
	clicked := &click.Click{}
	if err := proto.Unmarshal(b, clicked); err != nil {
		return nil, err
	}
	return clicked, nil
}

func errorToByte(error *click.Error) ([]byte, error) {
	out, err := proto.Marshal(error)
	if err != nil {
		return nil, err
	}
	return out, nil
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
}
func root(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	clicked, err := toClicked(body)
	if err != nil {
		fmt.Printf("error in unmarshel: %s", err)
	}
	fmt.Printf("name: %s \n", clicked.Name)
	w.Write([]byte(clicked.Name))
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
	router()
	// in, err := ioutil.ReadFile("./temp.txt")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// inclicked := &click.Click{}
	// if err := proto.Unmarshal(in, inclicked); err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println("in:", inclicked)
}
