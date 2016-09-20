package main

import (
	"bufio"
	"fmt"
	// "encoding/json"
	"log"
	"os"
	"time"
	"github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/authapi"
)

func getUserInput(msg string) string {
	fmt.Printf("%s ", msg)
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	return scan.Text()
}
	

func main() {

	ikey := os.Getenv("DUO_IKEY")
	skey := os.Getenv("DUO_SKEY")
	host := os.Getenv("DUO_HOST")
	if ikey == "" || skey == "" || host == "" {
		fmt.Println("Must set DUO_IKEY, DUO_SKEY and DUO_HOST environment variables")
		os.Exit(2)
	}
	userAgent := "Test User Agent"
	timeOut := duoapi.SetTimeout(10*time.Second)
	
	api := duoapi.NewDuoApi(ikey,skey,host,userAgent,timeOut)
	auth := authapi.NewAuthApi(*api)


	/* Ping() Test */
	pingStatus, err := auth.Ping()
	if err != nil {
		fmt.Println("error:", err)
	}
	if pingStatus.Stat != "OK" {
		log.Fatal("Duo Ping() failed")
	}

	/* Check() */
	statStatus, err := auth.Check()
	if err != nil {
		fmt.Println("error:", err)
	}
	if statStatus.Stat != "OK" {
		log.Fatal("Duo Status() failed")
	}

	userName := getUserInput("Enter Username:")

	// Testing Preauth()
	statPreauth, err := auth.Preauth(authapi.PreauthUsername(userName))
	if err != nil || statPreauth.Stat != "OK" {
		log.Fatal("Duo Preauth() failed")
	}

	if statPreauth.Response.Result == "enroll" {
		devType := ""
		firstRun := true
		for (devType != "1" && devType != "2") {
			fmt.Printf("%s\n", statPreauth.Response.Status_Msg)
			devType = getUserInput("Please select the type of device you have:\n1. Smart Phone\n2. Other\n->")
			if firstRun != true {
				fmt.Println("Must select 1 or 2")
			}
			firstRun = false
		}
		if devType == "1" {
			statEnroll, err := auth.Enroll(authapi.AuthUsername(userName))
			if err != nil || statEnroll.Stat != "OK" {
				log.Fatal("Duo Enroll() failed")
			}
			fmt.Printf("Please scan the image here %s with the Duo App to enroll.\n", statEnroll.Response.Activation_Barcode)
		} else {
			fmt.Printf("Please click here: %s to enroll, and try again.\n",
				statPreauth.Response.Enroll_Portal_Url)
		}
	}
	
	if statPreauth.Response.Result == "auth" {
		if len(statPreauth.Response.Devices) > 1 {
			fmt.Println("Which device would you like us to use to verify your identity?")
		// for i, device := range statPreauth.Response.Devices {
			fmt.Printf("%d \"%s\" type: %s (%s)", 1, statPreauth.Response.Devices[0].Device, statPreauth.Response.Devices[0].Type, statPreauth.Response.Devices[0].Number)
		} else {
			fmt.Println("Authenticating out of band. Please wait...")
			statAuth, err := auth.Auth("auto", authapi.AuthUsername(userName), authapi.AuthDevice("auto"))
			if err != nil || statAuth.Stat != "OK" {
				log.Fatal("Duo Auth() failed.")
			}

			if statAuth.Response.Result == "allow" {
				fmt.Println("Yay! :)")
			} else {
				fmt.Println("Boo :(")
			}
			fmt.Println(statAuth.Response.Status_Msg)
		}
	}
	// }
		

	// b, err := json.MarshalIndent(statPreauth, "", "  ")
	// if err != nil {
		// log.Fatal(err)
		// fmt.Println("error:", err)
	// }
	// os.Stdout.Write(b)

}
