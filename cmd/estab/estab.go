package main

import (
	"flag"
	"fmt"

	"github.com/belogik/goes"
)

func main() {

	host := flag.String("host", "localhost", "elasticsearch host")
	port := flag.String("port", "9200", "elasticsearch port")

	flag.Parse()

	conn := goes.NewConnection(*host, *port)
	fmt.Println(conn)

}
