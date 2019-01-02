package pac

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	fileName = "common/pac/pac.js"
	pacPort  = 2600
)

func Start() {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, string(data))
	})
	log.Printf("PAC working on %d ...\n", pacPort)
	err = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(pacPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
