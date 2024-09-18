package services

import "fmt"

func HealthCheck(host string) {
	fmt.Printf("listening server with host %s", host)
}
