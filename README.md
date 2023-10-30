# gojson

![](https://github.com/rhaeguard/gojson/actions/workflows/go.yml/badge.svg)

A JSON parser written in Go using Shift-Reduce Parsing technique without any parsing table.

### to add the dependency:

```shell
go get github.com/rhaeguard/gojson
```

### how to use:
```go
import (
    "fmt"
    "github.com/rhaeguard/gojson"
)

type WebConfig struct {
    Hostname string
    Port     int
    IsActive bool
}

func main() {
    json, err := gojson.ParseJson(
        `{
              "Hostname": "localhost",
              "Port": 8282,
              "IsActive": true
        }`,
    )
    
    if err != nil {
        // error occurred
    }
    
	// we can directly try to get the value from the JsonValue object
    objectFields := json.Value.(map[string]gojson.JsonValue)
    port := int(objectFields["Port"].Value.(float64)) // all numeric values are converted to float64
    fmt.Printf("%d\n", port)
    
	// we can also try to map the values to a struct
    var wc WebConfig
    if uerr := json.Unmarshal(&wc); uerr != nil {
        // error occurred
    } else {
        fmt.Printf("%-v\n", wc.Port)
    }
}

```