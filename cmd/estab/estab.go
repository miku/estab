// estab exports elasticsearch fields as tab separated values
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/belogik/goes"
	"github.com/miku/estab"
)

func main() {

	host := flag.String("host", "localhost", "elasticsearch host")
	port := flag.String("port", "9200", "elasticsearch port")
	indicesString := flag.String("indices", "", "indices to search (or all)")
	fieldsString := flag.String("f", "_id _index", "field or fields space separated")
	timeout := flag.String("timeout", "10m", "scroll timeout")
	size := flag.Int("size", 10000, "scroll batch size")
	nullValue := flag.String("null", "NOT_AVAILABLE", "value for empty fields")
	separator := flag.String("separator", "|", "separator to use for multiple field values")
	delimiter := flag.String("delimiter", "\t", "column delimiter")
	limit := flag.Int("limit", -1, "maximum number of docs to return (return all by default)")
	version := flag.Bool("v", false, "prints current program version")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	queryString := flag.String("query", "", "custom query to run")
	raw := flag.Bool("raw", false, "stream out the raw json records")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *version {
		fmt.Println(estab.Version)
		os.Exit(0)
	}

	var query map[string]interface{}
	if *queryString == "" {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	} else {
		err := json.Unmarshal([]byte(*queryString), &query)
		if err != nil {
			log.Fatal(err)
		}
	}

	indices := strings.Fields(*indicesString)
	fields := strings.Fields(*fieldsString)

	if !*raw {
		query["fields"] = fields
	}

	conn := goes.NewConnection(*host, *port)
	scanResponse, err := conn.Scan(query, indices, []string{""}, *timeout, *size)
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	i := 0

	for {
		scrollResponse, err := conn.Scroll(scanResponse.ScrollId, *timeout)
		if err != nil {
			log.Fatal(err)
		}
		if len(scrollResponse.Hits.Hits) == 0 {
			break
		}
		for _, hit := range scrollResponse.Hits.Hits {
			if i == *limit {
				return
			}
			if *raw {
				b, err := json.Marshal(hit)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintln(w, string(b))
				continue
			}
			var columns []string
			for _, f := range fields {
				var c []string
				switch f {
				case "_id":
					c = append(c, hit.Id)
				case "_index":
					c = append(c, hit.Index)
				case "_type":
					c = append(c, hit.Type)
				case "_score":
					c = append(c, strconv.FormatFloat(hit.Score, 'f', 6, 64))
				default:
					switch value := hit.Fields[f].(type) {
					case nil:
						c = []string{*nullValue}
					case []interface{}:
						for _, e := range value {
							c = append(c, e.(string))
						}
					default:
						log.Fatalf("unknown field type in response: %+v\n", hit.Fields[f])
					}
				}
				columns = append(columns, strings.Join(c, *separator))
			}
			fmt.Fprintln(w, strings.Join(columns, *delimiter))
			i++
		}
	}
}
