package main

import (
	"flag"
	"log"
	"fmt"
	"os"
	"os/user"
	"code.google.com/p/gopass"
	"github.com/sudharsh/henchman/henchman"
)



func currentUsername() string {
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Couldn't get current username")
		return ""
	}
	return u.Username
}


func main() {
	username  := flag.String("user", currentUsername(), "User to run as")
	password  := flag.String("password", "", "Path to the private key file")
	planFile := flag.String("plan", "", "Path to the plan")

	flag.Parse()
	if *username == "" {
		os.Exit(1)
	}

	if *password == "" {
		var err error
		if *password, err = gopass.GetPass("Password:"); err != nil {
			fmt.Println(err)
		}
	}

	plan, err := henchman.ParsePlan(planFile)
	sem := make(chan int, 100)
	for _, hostname := range plan.Hosts {
		go func() {
			machine := henchman.Machine{*username, *password, hostname}
			for _, task := range plan.Tasks {
				machine.RunTask(task)
			}
			sem <-1
		}()
		<- sem
	}

	if err != nil {
		log.Fatalf("Couldn't read the plan: %s", err)
		os.Exit(1)
	}

}
