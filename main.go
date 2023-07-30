package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: wifionice <login|logout>")
	}

	action := strings.ToLower(os.Args[1])
	if _, ok := map[string]struct{}{"login": struct{}{}, "logout": struct{}{}}[action]; !ok {
		log.Fatal("Usage: wifionice <login|logout>")
	}

	ir := &IndexRequest{}
	resp, err := ir.Send()
	if err != nil {
		log.Fatalf("Failed to aquire CSRF: %+v", err)
	}

	lr := &LoginRequest{
		CSRFToken: resp.CSRFToken,
		Login:     action == "login",
		Logout:    action == "logout",
	}

	res, err := lr.Send()
	if err != nil {
		log.Fatalf("Failed to login/out: %#v", err)
	}

	if (res.StatusCode != 200) && (res.StatusCode != 302) {
		log.Fatalf("Failed to login/out! code(%d)", res.StatusCode)
	}

	fmt.Printf("Ok, %d\n", res.StatusCode)
}
