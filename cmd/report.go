package cmd

import (
	"fmt"
	"log"
	"os"
)

func Report() {
	emails, obEmails := (UniqueEmails())
	if len(OutputFile) > 0 {
		printToFile(emails, obEmails)
	} else {
		printToConsole(emails, obEmails)
	}

}
func printToConsole(emails, obEmails PairList) {
	if len(emails) > 0 {
		log.Println("Emails Found:")
		for _, email := range emails {
			fmt.Printf("%s (%v)\n", email.Key, email.Value)
		}
	} else {
		log.Println("No Email Found")
	}

	if len(obEmails) > 0 {
		log.Println("Potential Emails Found:")
		for _, email := range obEmails {
			fmt.Printf("%s (%v)\n", email.Key, email.Value)
		}
	}
}

func printToFile(emails, obEmails PairList) {
	fileOutput, err := os.Create(OutputFile)
	if err != nil {
		log.Fatalf("Error creating file: %v\n", err)
	}
	defer fileOutput.Close()
	if len(emails) > 0 {
		fileOutput.WriteString("Emails Found:\n")
		for _, email := range emails {
			fmt.Fprintf(fileOutput, "%s (%v)\n", email.Key, email.Value)
		}
	} else {
		fileOutput.WriteString("No Email Found\n")
	}
	if len(obEmails) > 0 {
		fileOutput.WriteString("Potential Emails Found:\n")
		for _, email := range obEmails {
			fmt.Fprintf(fileOutput, "%s (%v)\n", email.Key, email.Value)
		}
	}
	fileOutput.WriteString("Links visited:\n")
	for link := range Links {
		fmt.Fprintf(fileOutput, "%s (%v)\n", link, Links[link])
	}
	fileOutput.Sync()
	log.Printf("Results wrote to file: %s", OutputFile)
}
