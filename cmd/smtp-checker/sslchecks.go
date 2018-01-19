package main

import (
	"crypto/x509"
	"crypto/tls"
	"encoding/pem"
	"path"
	"io/ioutil"
  	"regexp"
	"log"
	"fmt"
	"github.com/bitnami-labs/healthcheck-tools/pkg/apache"
)

func getActiveSSLCertsPath(apacheConf map[string]string, apacheRoot string) (res [][]string){
	for file, content := range apacheConf {
		sslcertRe, err := regexp.CompilePOSIX("^[[:space:]]*SSLCertificateFile [\"]?([^\n\"]+)[\"]?")
		if err != nil {
			log.Fatal(err)
		}
		sslcert := sslcertRe.FindAllStringSubmatch(content, -1)

		sslcertKeyRe, err := regexp.CompilePOSIX("^[[:space:]]*SSLCertificateKeyFile [\"]?([^\n\"]+)[\"]?")
		if err != nil {
			log.Fatal(err)
		}
		sslcertKey := sslcertKeyRe.FindAllStringSubmatch(content, -1)

		for index, element := range sslcert {
			newCert := element[1]
			if !path.IsAbs(newCert) {
				newCert = path.Join(apacheRoot, newCert)
			}
			newKey := ""
			if len(sslcertKey) > index - 1 {
				newKey = sslcertKey[index][1]
				if !path.IsAbs(newKey) {
					newKey = path.Join(apacheRoot, newKey)
				}
			}
			res = append(res, []string{file, newCert, newKey})
		}
	}
	return res
}

type checkIncludedSSLCertsResult struct {
	apacheFile string
	certFile string
	keyFile string
}

func checkIncludedSSLCertsInApache(sslCertPaths [][]string) (res []checkIncludedSSLCertsResult){
	fmt.Println("\n--- Check: Included SSL Certificates in Apache ---")
	for index, element := range sslCertPaths {
		fmt.Print("Ocurrence #")
		fmt.Println(index + 1)
		fmt.Println("  - Apache file: " + element[0])
		fmt.Println("  - Certificate file: " + element[1])
		fmt.Println("  - Key file: " + element[2])
		res = append(res, checkIncludedSSLCertsResult{element[0], element[1], element[2]})
	}
	return
}

type checkSSLCertificateResult struct {
	certFile string
	certExists bool
    certDecoded bool
	certDomain string
}

func checkSSLCertificate(sslCertPaths [][]string) (res []checkSSLCertificateResult) {
 	fmt.Println("\n--- Check: Decode SSL Certificate ---")
	for index, element := range sslCertPaths {
		fmt.Print("Ocurrence #")
		fmt.Println(index + 1)
		cert, err := ioutil.ReadFile(element[1])
		certExists := false
		certDecoded := false
		certDomain := "Could not decode certificate"
		if err == nil {
			certExists = true
			block, _ := pem.Decode(cert)
			parsedCert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				certDecoded = true
				certDomain = parsedCert.Subject.CommonName
			} else {
				log.Println(err)
			}
		 } else {
			 log.Println(err)
		 }

		fmt.Println("  - Certificate file: " + element[1])
		fmt.Println("  - Certificate can be opened: ", certExists)
		fmt.Println("  - Certificate file can be decoded: ", certDecoded)
		fmt.Println("  - Domain name: " + certDomain)
		res = append(res, checkSSLCertificateResult{element[1], certExists, certDecoded, certDomain})
	}
	return
}

type checkSSLMatchResult struct {
	certFile string
	keyFile string
	match bool
}

func checkSSLMatch(sslCertPaths [][]string) (res []checkSSLMatchResult) {
 	fmt.Println("\n--- Check: Certificate and Key match ---")
	for index, element := range sslCertPaths {
		fmt.Print("Ocurrence #")
		fmt.Println(index + 1)

		match := true
		_ , err := tls.LoadX509KeyPair(element[1], element[2])
		if err != nil {
			match = false
			log.Println(err)
		}
		fmt.Println("  - Certificate file: " + element[1])
		fmt.Println("  - Certificate key: " + element[2])
		fmt.Println("  - Certificate and Key match: ", match)
		res = append(res, checkSSLMatchResult{element[1], element[2], match})
	}
	return
}

func RunChecks(confFile, apacheRoot string) {
	apacheConf := apache.OpenAllApacheConfigurationFiles(confFile, apacheRoot)
	sslCertPaths := getActiveSSLCertsPath(apacheConf, apacheRoot)
	if len(sslCertPaths) == 0 {
		fmt.Println("No SSL certificates found in the Apache configuration")
	} else {
		checkIncludedSSLCertsInApache(sslCertPaths)
		checkSSLCertificate(sslCertPaths)
		checkSSLMatch(sslCertPaths)
	}
}
