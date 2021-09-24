package main

import (
	"fmt"
	"net/url"
	"strings"
)

func encode(value string) string {
	encodedValue := fmt.Sprintf("%s", url.QueryEscape(value))

	return encodedValue
}

func main() {
	str := "#a-db-user:"
	fmt.Printf("Orig=%s, encoded=%s\n", str, encode(str))

	special := ":/?#[]@"
	fmt.Println(strings.ContainsAny(str, special))
	fmt.Println(strings.ContainsAny("password:", special))
	fmt.Println(strings.ContainsAny("password/", special))
	fmt.Println(strings.ContainsAny("password?", special))
	fmt.Println(strings.ContainsAny("password#", special))
	fmt.Println(strings.ContainsAny("password[", special))
	fmt.Println(strings.ContainsAny("password]", special))
	fmt.Println(strings.ContainsAny("password:/?", special))
	fmt.Println(strings.ContainsAny("password", special))

}
