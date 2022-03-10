package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/byteford/goproto/github.com/byteford/goproto/modules/click"
	"google.golang.org/protobuf/proto"
)

func lambdaHander(ctx context.Context, input []byte) {
	fmt.Println(input)
	clicked := &click.Click{}
	err := proto.Unmarshal(input, clicked)
	if err != nil {
		log.Fatalln("Failed to unencode: ", err)
	}
	fmt.Println(clicked.AmountClicked)
}
func root(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(body)
	w.Write([]byte("yay"))
}
func router() {
	http.HandleFunc("/", root)

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
