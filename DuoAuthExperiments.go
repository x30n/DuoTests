package main

import (
	"fmt"
	"encoding/json"
	"os"
	"time"
	"github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/authapi"
)

func main() {

	ikey := os.Getenv("DUO_IKEY")
	skey := os.Getenv("DUO_SKEY")
	host := os.Getenv("DUO_HOST")
	if (ikey == "" or skey == "" or host == "") {
	   fmt.Println("Must set DUO_IKEY, DUO_SKEY and DUO_HOST environment variables")
	}
	userAgent := "Test User Agent"
	timeOut := duoapi.SetTimeout(10*time.Second)
	
	api := duoapi.NewDuoApi(ikey,skey,host,userAgent,duoapi.SetTimeout(10*time.Second))
	auth := authapi.NewAuthApi(*api)
	status, err := auth.Ping()
	if err != nil {
		fmt.Println("error:", err)
	}
	b, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)

	// fmt.Printf(strings.String(auth.Ping()))
}
