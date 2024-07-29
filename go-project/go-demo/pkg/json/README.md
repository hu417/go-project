

# json序列化

json是一种轻量级的数据交换格式，易于人阅读和编写，同时也易于机器解析和生成，之所以json这么流行，是因为json的结构和多级结构体（对象）刚好能对应上，并且本身也十分易读。而前后端交互的时候后端通常会返回给前端一个多级的结构体，于是json慢慢开始流行了，且json是跨语言和跨平台的，自身也足够轻量级。

json格式
```go

一个标准的json数据
//每个key对应的是一个value
{
“k1": 1,
"k2": 2 //注意结尾的这个不能有逗号
}



json字符串
{
"k1": "1",
"k2": "2"
}


json数组
{
“k1”: [1,2],
“k2”: [3,4]
}


json对象
{
“k1”: {“1”: “haihai”},
“k2”: {“2”:”haihahai”}
}


json对象数组
{
“k1”: [
{“k11”: “hellohello”},
{“k12”: “badbad”}
]
}



json数组对象
{
“k2”: {
	“hello”: [1,2,3]
	}
}

所有的JSON数据都是由上述几种JSON数据组合而成

```


## 解析不定JSON数据

### 方式一: 使用接口类型

```go

// 请求其他服务   
jsonStr := xxx

var data interface{}

err := json.Unmarshal([]byte(jsonStr),&data)

fmt.Println(data)

```


### 方式二: 使用map

```go

jsonStr := xxx
var data map[string]interface{} 
err := json.Unmarshal([]byte(jsonStr),&data)

// 你想取的字段
fieldValue := data["field"]

// 类型断言
if value,ok := data["field"].(float64);ok {

} else if vluae,ok := data["field"].(int64); ok {

}


```

理论上所有的合法的JSON数据都可以被反序列化到map[string]interface{}中,但是实际应用中 可能会出现一些无法被map[string]interface{}解析的JSON数据
- JSON 数据中包含了多层嵌套的数据结构。在这种情况下，如果没有使用递归或者其他方式对嵌套数据进行处理，可能会导致反序列化失败。
- JSON 数据中包含了数组类型，但是数组元素类型不一致或者无法转换成相应的类型。在这种情况下，可能需要手动处理数组元素或者使用其他数据类型来保存数组数据。
- JSON 数据中包含了自定义数据类型或者复杂的数据结构，无法使用 map[string]interface{} 类型来反序列化。在这种情况下，需要定义相应的结构体或者使用其他适合的数据类型来反序列化。


### 方式三: 第三方库

如gjson和jsonparser，这些库通常提供了更加灵活和高效的JSON解析方式，可以根据具体的需求选择合适的库来使用


### 方式四: 使用其他包

json.RawMessage 是一个非常高效的数据类型，因为她不需要进行任何解析和类型转换，直接保存了未经处理的原始JSON数据，在反序列化的时候只需要将json.RawMessage转化为对应的数据类型即可，无需重新解析JSON数据
json.Number 表示JSON中的数字类型，可以用来保存任意精度的数字。这个数字可以特别大，可能会无法用Go中的整数或者浮点数来表示

```go

package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    jsonData := []byte(`{
        "id": 12345,
        "name": "John Doe",
        "score": 95.5,
        "is_student": true,
        "tags": ["tag1", "tag2", "tag3"],
        "extra": {
            "field1": "value1",
            "field2": 123
        }
    }`)

    var m map[string]json.RawMessage
    err := json.Unmarshal(jsonData, &m)
    if err != nil {
        panic(err)
    }

    var id int
    err = json.Unmarshal(m["id"], &id)
    if err != nil {
        panic(err)
    }
    fmt.Printf("id: %d\n", id)

    var name string
    err = json.Unmarshal(m["name"], &name)
    if err != nil {
        panic(err)
    }
    fmt.Printf("name: %s\n", name)

    var score float64
    err = json.Unmarshal(m["score"], &score)
    if err != nil {
        panic(err)
    }
    fmt.Printf("score: %f\n", score)

    var isStudent bool
    err = json.Unmarshal(m["is_student"], &isStudent)
    if err != nil {
        panic(err)
    }
    fmt.Printf("is_student: %v\n", isStudent)

    var tags []string
    err = json.Unmarshal(m["tags"], &tags)
    if err != nil {
        panic(err)
    }
    fmt.Printf("tags: %v\n", tags)

    var extra map[string]json.RawMessage
    err = json.Unmarshal(m["extra"], &extra)
    if err != nil {
        panic(err)
    }
    var field1 string
    err = json.Unmarshal(extra["field1"], &field1)
    if err != nil {
        panic(err)
    }
    fmt.Printf("extra.field1: %s\n", field1)

    var field2 int
    err = json.Unmarshal(extra["field2"], &field2)
    if err != nil {
        panic(err)
    }
    fmt.Printf("extra.field2: %d\n", field2)
}

// 不确定的类型
data := make(map[string]interface{})
if err := json.Unmarshal(rawData, &data); err != nil {
    log.Fatal(err)
}

if value, ok := data["age"].(float64); ok {
    // 处理年龄为浮点数的情况
} else if value, ok := data["age"].(int); ok {
    // 处理年龄为整数的情况
} else {
    // 处理年龄为其他类型或不存在的情况
}

```