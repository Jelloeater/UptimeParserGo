# SNMP based up time parser

* Takes command line args and output XML for PRTG

## Build
```
go get -d ./...
go build
```

## Usage
```
USAGE:
   UptimeParserGo.exe [global options] command [command options] [arguments...]

COMMANDS:
     help, h  Shows a list of commands or help for one command

   output:
     xml, x   export as XML
     json, j  export as JSON

GLOBAL OPTIONS:
   --debug, -d             enable debug mode
   --ip value, -i value    IP address to scan
   --snmp value, -s value  SNMP community string
   --help, -h              show help

```

Example

`UptimeParserGo -ip 192.168.1.1/24 -snmp "public" x`

## Notes
No JSON out yet, only XML