tagq
====

![Coverage Badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/gabstv/f6b88267c7dcdd8b8f0adb53441566c9/raw/tagq__heads_master.json)

```go
package main

import (
    "fmt"
    "time"

    "github.com/gabstv/tagq"
)

type Item struct {
    Name string `json:"jname" xml:"xname"`
    Score int
    ScoreHist []int
    Nested struct{
        Moment time.Time `json:"m_oment"`
        Title string
    }
}

func main() {
    item := &Item{}
    item.Name = "Example"
    item.Score = 100
    item.ScoreHist = []int{80, 5, 15}
    item.Nested.Moment = time.Now()
    item.Nested.Title = "Hello"

    fmt.Println(tagq.Q(item, "jname").Str()) // prints: Example
    fmt.Println(tagq.Q(item, "xname").Str()) // prints: Example
    fmt.Println(tagq.Q(item, "Name").Str()) // prints: Example
    fmt.Println(tagq.Q(item, "Name", "something").Str()) // prints: ""
    fmt.Println(tagq.Q(item, "Score").Int()) // prints: 100
    fmt.Println(tagq.Q(item, "ScoreHist", "last").Int()) // prints: 15
    fmt.Println(tagq.Q(item, "ScoreHist", "0").Int()) // prints: 80
    fmt.Println(tagq.Q(item, "Nested", "m_oment").Time())
    fmt.Println(tagq.Q(item, "Nested", "Title").Str()) // prints: Hello
}
```