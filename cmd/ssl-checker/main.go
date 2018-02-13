package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var apacheRoot = flag.String("apacheRoot", "/opt/bitnami/apache2/", "Root of Apache installation")
	var apacheConf = flag.String("apacheConf", "/opt/bitnami/apache2/conf/httpd.conf",
		"Path to the root Apache configuration file")
	var hostname = flag.String("hostname", "", "Web application hostname")
	var port = flag.String("port", "443", "Web application port")
	flag.Parse()
	if *hostname == "" {
		log.Fatal("-hostname flag must be set")
	}
	fmt.Println("======================================")
	fmt.Println("SSL CHECKS")
	fmt.Println("======================================")
	fmt.Println("Starting checks with these parameters:")
	fmt.Println("  - Apache Root: " + *apacheRoot)
	fmt.Println("  - Apache Root configuration: " + *apacheConf)
	fmt.Println("  - Hostname: " + *hostname)
	fmt.Println("  - Port: " + *port)
	fmt.Println("======================================")
	RunChecks(*apacheConf, *apacheRoot, *hostname, *port)
	fmt.Println("\nSSL CHECKS FINISHED")

}
