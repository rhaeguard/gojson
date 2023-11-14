# gojson

![](https://github.com/rhaeguard/gojson/actions/workflows/go.yml/badge.svg)

A JSON parser written in Go using Shift-Reduce Parsing technique without any parsing table.

read the [article](https://rhaeguard.github.io/posts/json-parsing-shift-reduce/)

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
    inputJson := `{
              "Hostname": "localhost",
              "Port": 8282,
              "IsActive": true
        }`
    json, err := gojson.Parse(inputJson)
    
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

	// we can unmarshall in one step as well
    if uerr := gojson.Unmarshal(inputJson, &wc); uerr != nil {
        // error occurred
    } else {
        fmt.Printf("%-v\n", wc.Port)
    }
}
```
