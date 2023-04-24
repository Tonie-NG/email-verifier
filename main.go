package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/mail"
	"os"
	"strings"
)

func main () {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Domain, hasMX, hasSPF, SPFRecord, hasDMARC, DMARCRecord\n")

	for scanner.Scan(){
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: Could not read from input: %v\n",err)
	}

}

func checkDomain(emailAddress string) {

	splitted := strings.Split(emailAddress, "@")

	domain := splitted[1]

	var hasMX, hasSPF, hasDMARC, validEmail bool

	var SPFRecord, DMARCRecord string

	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecord, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	
	for _, record := range txtRecord {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			SPFRecord = record
			break
		}
	}

	// Validate email addresses for the domain
	_, err = mail.ParseAddress(emailAddress)
    if err != nil {
        log.Printf("Error: Invalid email address: %v\n", err)
    }

	validEmail = true


	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC= true
			DMARCRecord = record
			break
		}
	}

	fmt.Printf("%v,%v,%v,%v,%v,%v, %v", domain, hasMX, hasSPF, SPFRecord, hasDMARC, DMARCRecord, validEmail)
}