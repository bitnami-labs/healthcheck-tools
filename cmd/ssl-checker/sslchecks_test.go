package main

import (
	"testing"
	"github.com/bitnami-labs/healthcheck-tools/pkg/apache"
)

func TestSSLChecker(t *testing.T) {
	apacheConf := "testdata/apache2/conf/httpd.conf"
	apacheRoot := "testdata/apache2/"

	openedApacheFiles := apache.OpenAllApacheConfigurationFiles(apacheConf, apacheRoot)
	detectedCerts := getActiveSSLCertsPath(openedApacheFiles, apacheRoot)
	t.Run("Check Detected SSL files", func(t *testing.T) {
		if len(detectedCerts) != 1 {
			t.Errorf("Incorrect certs detected, expected: 1, got: %d", len(detectedCerts))
 		}
		if detectedCerts[0].apacheConfPath != "testdata/apache2/conf/bitnami/bitnami.conf" {
			t.Errorf("Incorrect file detected, expected: testdata/apache2/conf/bitnami/bitnami.conf, got: %s",
				detectedCerts[0].apacheConfPath)
		}
	})
    t.Run("Check Detected domain", func(t *testing.T) {
		checkResult := getSSLCertsInfo(detectedCerts)
		if (checkResult[0].certDomain != "www.example.com") {
			t.Errorf("Incorrect domain detected, expected: www.example.com, got: %s", checkResult[0].certDomain)
		}
	})
    t.Run("Check Cert and Key Match", func(t *testing.T) {
		checkResult := getCertKeyMatchInfo(detectedCerts)
		if !checkResult[0].match {
			t.Errorf("Incorrect certificate match detected, expected: true, got: %t", checkResult[0].match)
		}
	})
    t.Run("Check HTTPS Connection", func(t *testing.T) {
		checkResult := getHTTPSConnectionInfo("bitnami.com", "443")
		if checkResult.serverCertDomain != "www.bitnami.com" {
			t.Errorf("Incorrect HTTPS Server certificate detected, expected: www.bitnami.com, got: %s",
				checkResult.serverCertDomain)
		}
	})
}
