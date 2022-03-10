package msgError

import (
	"fmt"
	"net/http"

	"google.golang.org/protobuf/proto"
)

//ToByte Marshals Error in to bytes
func (msgError *Error) ToByte() ([]byte, error) {
	out, err := proto.Marshal(msgError)
	if err != nil {
		return nil, err
	}
	return out, nil
}

//Write converts Error to bytes and writes in in to the provided writer
func (msgError *Error) Write(w http.ResponseWriter) error {
	out, err := msgError.ToByte()
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
