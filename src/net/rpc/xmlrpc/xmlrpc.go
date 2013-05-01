package xmlrpc

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "reflect"
    "strconv"
    "time"
)

func Request(url string, method string, params ...interface{}) ([]interface{}) {
    request := Serialize(method, params)
    log.Printf("%s", request)
    buffer := bytes.NewBuffer([]byte(request))

    response, err := http.Post(url, "text/xml", buffer)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    return Unserialize(response.Body)
}

type MethodResponse struct {
    Params []Param `xml:"params>param"`
}

type Param struct {
    Value Value `xml:"value"`
}

type Value struct {
    List []Value `xml:"array>data>value"`
    Object []Member `xml:"struct>member"`
    String string `xml:"string"`
    Int string `xml:"int"`
    Boolean string `xml:"boolean"`
    DateTime string `xml:"dateTime.iso8601"`
}

type Member struct {
    Name string `xml:"name"`
    Value Value `xml:"value"`
}

func unserialize(value Value) (interface{}) {
    if value.List != nil {
        result := make([]interface{}, len(value.List))
        for i, v := range value.List {
            result[i] = unserialize(v)
        }
        return result

    } else if value.Object != nil {
        result := make(map[string]interface{}, len(value.Object))
        for _, member := range value.Object {
            result[member.Name] = unserialize(member.Value)
        }
        return result

    } else if value.String != "" {
        return fmt.Sprintf("%s", value.String)

    } else if value.Int != "" {
        result, _ := strconv.Atoi(value.Int)
        return result

    } else if value.Boolean != "" {
        return value.Boolean == "1"

    } else if value.DateTime != "" {
        var format = "20060102T15:04:05"
        result, _ := time.Parse(format, value.DateTime)
        return result
    }

    return nil
}

func Unserialize(buffer io.ReadCloser) ([]interface{}) {
    body, err := ioutil.ReadAll(buffer)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("%s", body)

    var response MethodResponse
    xml.Unmarshal(body, &response)

    result := make([]interface{}, len(response.Params))
    for i, param := range response.Params {
        result[i] = unserialize(param.Value)
    }

    return result
}

func Serialize(method string, params []interface{}) (string) {
    request := "<methodCall>"
    request += fmt.Sprintf("<methodName>wp.%s</methodName>", method)
    request += "<params>"

    for _, value := range params {
        request += "<param>"
        request += serialize(value)
        request += "</param>"
    }

    request += "</params></methodCall>"

    return request
}

func serialize(value interface{}) (string) {
    result := "<value>"
    switch value.(type) {
        case string:
            result += fmt.Sprintf("<string>%s</string>", value.(string))
            break
        case int:
            result += fmt.Sprintf("<int>%d</int>", value)
            break
        default:
            if reflect.ValueOf(value).Kind() == reflect.Map {
                result += "<struct>"
                for k, v := range value.(map[string]interface{}) {
                    result += "<member>"
                    result += fmt.Sprintf("<name>%s</name>", k)
                    result += serialize(v)
                    result += "</member>"
                }
                result += "</struct>"
            } else {
                log.Fatal(value)
            }

    }
    result += "</value>"
    return result
}
