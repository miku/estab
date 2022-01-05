estab
=====

Please see [esdump](https://github.com/miku/esdump) for a alternative elasticsearch command line export tool.

----

Export elasticsearch document fields into tab separated values. `estab`
uses the [scan search type](http://www.elasticsearch.org/guide/en/elasticsearch/guide/current/scan-scroll.html)
and the scroll API, which help

> to retrieve large numbers of documents from Elasticsearch efficiently ...

[![Project Status: Inactive – The project has reached a stable, usable state but is no longer being actively developed; support/maintenance will be provided as time allows.](https://www.repostatus.org/badges/latest/inactive.svg)](https://www.repostatus.org/#inactive)

Note: For another command line export option, take a look at
[esdump](https://github.com/miku/esdump), which is can be combined with tools
like [jq](https://stedolan.github.io/jq/) to yield similar results as estab.

Installation
------------

    $ go get github.com/miku/estab/cmd/estab

Or if your system speaks `dpkg` or `rpm`, there is a [release](https://github.com/miku/estab/releases).

Usage
-----

    $ estab -h
    Usage of estab:
      -1=false: one value per line (works only with a single column in -f)
      -cpuprofile="": write cpu profile to file
      -delimiter="\t": column delimiter
      -f="_id _index": field or fields space separated
      -header=false: output header row with field names
      -host="localhost": elasticsearch host
      -indices="": indices to search (or all)
      -limit=-1: maximum number of docs to return (return all by default)
      -null="NOT_AVAILABLE": value for empty fields
      -port="9200": elasticsearch port
      -precision=0: precision for numeric output
      -query="": custom query to run
      -raw=false: stream out the raw json records
      -separator="|": separator to use for multiple field values
      -size=10000: scroll batch size
      -timeout="10m": scroll timeout
      -v=false: prints current program version
      -zero-as-null=false: treat zero length strings as null values

Example
-------

Assuming an elasticsearch is running on localhost:9200.

    $ curl -XPOST localhost:9200/test/default/ -d '{"name": "Tim", "color": "red"}'
    $ curl -XPOST localhost:9200/test/default/ -d '{"name": "Alice", "color": "yellow"}'
    $ curl -XPOST localhost:9200/test/default/ -d '{"name": "Brian", "color": "green"}'

    $ estab -indices "test" -f "name color"
    Brian   green
    Tim     red
    Alice   yellow

Specify multiple indices:

    $ curl -XPOST localhost:9200/test2/default/ -d '{"name": "Yang", "color": "white"}'
    $ curl -XPOST localhost:9200/test2/default/ -d '{"name": "Ying", "color": "black"}'

    $ estab -indices "test test2" -f "name color"
    Ying    black
    Yang    white
    Tim     red
    Alice   yellow
    Brian   green

Multiple values are packed into a single value:

    $ curl -XPOST localhost:9200/test/default/ \
           -d '{"name": "Meltem", "color": ["green", "white"]}'

    $ estab -indices "test" -f "name color"
    Brian   green
    Meltem  green|white
    Tim     red
    Alice   yellow

Missing values get a special value via `-null`, which defaults to `NOT_AVAILABLE`:

    $ curl -XPOST localhost:9200/test/default/ -d '{"name": "Jin"}'

    $ estab -indices "test" -f "name color"
    Brian   green
    Meltem  green|white
    Tim     red
    Alice   yellow
    Jin     NOT_AVAILABLE

In 0.2.0 a `-raw` flag was added that will stream out the full JSON documents:

    $ estab -indices test -raw
    {"_index":"test","_type":"default","_id":"...BvA4z","_score":0,"_source":{...}}
    {"_index":"test","_type":"default","_id":"...BvA40","_score":0,"_source":{...}}
    {"_index":"test","_type":"default","_id":"...BvA41","_score":0,"_source":{...}}

This can be fed into json processors like [jq](http://stedolan.github.io/jq/):

    $ estab -indices test -raw | jq --raw-output '._source.color'
    red
    yellow
    green

In 0.2.1 a `-1` flag was added:

    $ estab -indices "test" -f "color" -1
    green
    green
    white
    red
    yellow
    NOT_AVAILABLE
