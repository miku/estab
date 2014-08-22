README
======

Export elasticsearch document fields into tab separated values.

Installation
------------

    $ go get github.com/miku/estab/cmd/estab

Or if your system speaks `dpkg` or `rpm`, there is a [release](https://github.com/miku/estab/releases).

Usage
-----

    $ estab -h
    Usage of ./estab:
      -cpuprofile="": write cpu profile to file
      -delimiter="\t": column delimiter
      -f="content.245.a": field or fields space separated
      -host="localhost": elasticsearch host
      -indices="": indices to search (or all)
      -null="NOT_AVAILABLE": value for empty fields
      -port="9200": elasticsearch port
      -separator="|": separator to use for multiple field values
      -size=10000: scroll batch size
      -timeout="10m": scroll timeout
      -v=false: prints current program version

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
    Jin NOT_AVAILABLE
