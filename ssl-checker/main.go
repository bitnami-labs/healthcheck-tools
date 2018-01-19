package main

import (
	"flag"
	"fmt"
)

func main() {
	var apacheRoot = flag.String("apacheRoot", "/opt/bitnami/apache2/", "Root of Apache installation")
	var apacheConf = flag.String("apacheConf", "/opt/bitnami/apache2/conf/httpd.conf",
		"Path to the root Apache configuration file")
	flag.Parse()
	fmt.Println("======================================")
	fmt.Println("SSL CHECKS")
	fmt.Println("======================================")
	fmt.Println("Starting checks with these parameters:")
	fmt.Println("  - Apache Root: " + *apacheRoot)
	fmt.Println("  - Apache Root configuration: " + *apacheConf)
	fmt.Println("======================================")
	RunChecks(*apacheConf, *apacheRoot)
	fmt.Println("\nSSL CHECKS FINISHED")

}
