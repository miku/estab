package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/belogik/goes"
)

func main() {

	host := flag.String("host", "localhost", "elasticsearch host")
	port := flag.String("port", "9200", "elasticsearch port")
	indicesString := flag.String("indices", "", "indices to search (or all)")
	fieldsString := flag.String("f", "content.245.a", "field or fields space separated")

	flag.Parse()

	var indices []string
	trimmed := strings.TrimSpace(*indicesString)
	if len(trimmed) > 0 {
		indices = strings.Fields(trimmed)
	}

	fields := strings.Fields(*fieldsString)
	conn := goes.NewConnection(*host, *port)
	var query = map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"fields": fields,
	}

	extraArgs := make(url.Values, 1)
	searchResults, err := conn.Search(query, indices, []string{""}, extraArgs)
	if err != nil {
		log.Fatalln(err)
	}
	for _, hit := range searchResults.Hits.Hits {
		var columns []string
		for _, f := range fields {
			switch value := hit.Fields[f].(type) {
			case []interface{}:
				var c []string
				for _, e := range value {
					c = append(c, e.(string))
				}
				columns = append(columns, strings.Join(c, "|"))
			default:
				log.Fatal("unknown field type in response")
			}
		}
		fmt.Println(strings.Join(columns, "\t"))
	}

}
