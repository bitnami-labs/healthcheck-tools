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

// ActiveCertKeyPath contains paths of an active certificate-key path
type ActiveCertKeyPath struct {
	apacheConfPath string
	certPath string
	keyPath string
}

// getActivesslcertspath obtains the certificate-key pairs that are being used in each Apache configuration file
func getActiveSSLCertsPath(apacheConf map[string]string, apacheRoot string) (res []ActiveCertKeyPath){
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
			res = append(res, ActiveCertKeyPath{file, newCert, newKey})
		}
	}
	return res
}

// showActiveSSLCertsPath prints on screen the certificate-key pairs that are used in each Apache configuration file
func showActiveSSLCertsPath(certKeyPairs []ActiveCertKeyPath) {
	fmt.Println("\n--- Check: Included SSL Certificates in Apache ---")
	for index, keyPair := range certKeyPairs {
		fmt.Println("  Ocurrence #", index + 1)
		fmt.Println("  - Apache file: " + keyPair.apacheConfPath)
		fmt.Println("  - Certificate file: " + keyPair.certPath)
		fmt.Println("  - Key file: " + keyPair.keyPath)
	}
	return
}

// SSLCertInfo contains the following information about an active certificate
//     - Path to the certificate
//     - If the certificate can be opened
//     - If the certificate can be decoded
//     - The certificate domain name that uses
type SSLCertInfo struct {
	certPath string
	certOpened bool
    certDecoded bool
	certDomain string
}

// getSSLCertsInfo returns some the following information from the active certificate files
//     - Path to the certificate
//     - If the certificate can be opened
//     - If the certificate can be decoded
//     - The certificate domain name that uses
func getSSLCertsInfo(certKeyPairs []ActiveCertKeyPath) (res []SSLCertInfo) {
	for _, keyPair := range certKeyPairs {
		certOpened := false
		certDecoded := false
		certDomain := "Could not decode certificate"
		cert, err := ioutil.ReadFile(keyPair.certPath)
		if err == nil {
			certOpened = true
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
		res = append(res, SSLCertInfo{keyPair.certPath, certOpened, certDecoded, certDomain})
	}
	return
}

// showSSLCertsInfo prints on screen the following information from the active certificate files
//     - Path to the certificate
//     - If the certificate can be opened
//     - If the certificate can be decoded
//     - The certificate domain name that uses
func showSSLCertsInfo(certsInfo []SSLCertInfo) {
 	fmt.Println("\n--- Check: Decode SSL Certificate ---")
	for index, certInfo := range certsInfo {
		fmt.Println("  Ocurrence #", index + 1)
		fmt.Println("  - Certificate file: " + certInfo.certPath)
		fmt.Println("  - Certificate can be opened: ", certInfo.certOpened)
		fmt.Println("  - Certificate file can be decoded: ", certInfo.certDecoded)
		fmt.Println("  - Domain name: ", certInfo.certDomain)
	}
}

// CertKeyMatchInfo contains information about a certificate-key pair and whether they match or not
type CertKeyMatchInfo struct {
	certPath string
	keyPath string
	match bool
}

// getCertKeyMatchInfo returns, for each active certificate-key pair, whether they match or not
func getCertKeyMatchInfo(certKeyPairs []ActiveCertKeyPath) (res []CertKeyMatchInfo) {
	for _, keyPair := range certKeyPairs {
		match := true
		_ , err := tls.LoadX509KeyPair(keyPair.certPath, keyPair.keyPath)
		if err != nil {
			match = false
			log.Println(err)
		}
		res = append(res, CertKeyMatchInfo{keyPair.certPath, keyPair.keyPath, match})
	}
	return
}

// showCertKeyMatchInfo shows, for each active certificate-key pair, whether they match or not
func showCertKeyMatchInfo(certKeyMatches []CertKeyMatchInfo) {
 	fmt.Println("\n--- Check: Certificate and Key match ---")
	for index, certKeyMatchInfo := range certKeyMatches {
		fmt.Println("  Ocurrence #", index + 1)
		fmt.Println("  - Certificate file: " + certKeyMatchInfo.certPath)
		fmt.Println("  - Key file: ", certKeyMatchInfo.keyPath)
		fmt.Println("  - Certificate and key match: ", certKeyMatchInfo.match)
	}
}

// HTTPSConnectionInfo contains information about the HTTPS connection attempt to the server
//     - Hostname and port
//     - Whether the connection was successful or not
//     - The domain name of the server certificate
type HTTPSConnectionInfo struct {
	hostname string
	port string
	canConnect bool
	serverCertDomain string
}

// getHTTPSConnectionInfo attempts a HTTPS connection to the server and returns the following information
//     - Hostname and port
//     - Whether the connection was successful or not
//     - The domain name of the server certificate
func getHTTPSConnectionInfo(hostname string, port string) (res HTTPSConnectionInfo){
	conf := &tls.Config{
		InsecureSkipVerify: true,
    }

	res.hostname = hostname
	res.port = port
	res.canConnect = true
	res.serverCertDomain = "Could not obtain certificate"

    conn, err := tls.Dial("tcp", hostname + ":" + port, conf)
    if err != nil {
		res.canConnect = false
        log.Println(err)
        return
    }

	conn.Handshake()
	res.serverCertDomain = conn.ConnectionState().PeerCertificates[0].Subject.CommonName

	return
}

// getHTTPSConnectionInfo prints the results of the HTTPS connection attempt to the server
//     - Hostname and port
//     - Whether the connection was successful or not
//     - The domain name of the server certificate
func showHTTPSConnectionInfo(info HTTPSConnectionInfo) {
 	fmt.Println("\n--- Check: HTTPS Connection with server ---")
	fmt.Println("  - Hostname: ", info.hostname)
	fmt.Println("  - Port: ", info.port)
	fmt.Println("  - Can connect: ", info.canConnect)
	fmt.Println("  - Server certificate domain: ", info.serverCertDomain)
}

// RunChecks executes the health check suite and prints the results on screen
func RunChecks(confFile, apacheRoot, hostname, port string) {
	apacheConf := apache.OpenAllApacheConfigurationFiles(confFile, apacheRoot)
	sslCertPaths := getActiveSSLCertsPath(apacheConf, apacheRoot)
	if len(sslCertPaths) == 0 {
		fmt.Println("No SSL certificates found in the Apache configuration")
	} else {
	    showActiveSSLCertsPath(sslCertPaths)
		showSSLCertsInfo(getSSLCertsInfo(sslCertPaths))
		showCertKeyMatchInfo(getCertKeyMatchInfo(sslCertPaths))
		showHTTPSConnectionInfo(getHTTPSConnectionInfo(hostname, port))
	}
}
