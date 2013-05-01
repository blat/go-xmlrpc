Go XML-RPC
==========

XML-PRC package for Go programming language.


Usage
----------

    package main

    import (
        "net/rpc/xmlrpc",
        "fmt",
        "log"
    }

    const (
        URL = "http://my.blog.com/wordpress/xmlrpc.php"
        USERNAME = "admin"
        PASSWORD = ""
        BLOG_ID = "0"
    )

    func main() {
        options := make(map[string]interface{})
        options["offset"] = 0
        options["number"] = 10
        response := xmlrpc.Request(URL, "getPosts", BLOG_ID, USERNAME, PASSWORD, options)

        for _, params := range response {
            for _, param := range params.([]interface{}) {
                log.Printf("%s", param.(map[string]interface{})["post_title"])
            }
        }
    }

This example displays the title of the 10 last posts in a WordPress.
