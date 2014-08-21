README
======

Export elasticsearch document fields into tab separated values.

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
    Tim red
    Alice   yellow
