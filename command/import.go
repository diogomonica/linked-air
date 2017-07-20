package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/diogomonica/linked-air/contacts_base"

	"bufio"
	"encoding/csv"
	"io"
	"log"

	"github.com/fabioberger/airtable-go"
	"github.com/urfave/cli"
)

func CmdImportCompanies(c *cli.Context) error {
	// Write your code here

	airtableAPIKey := os.Getenv("AIRTABLE_API_KEY")
	baseID := "appVa3p4rNePGxslA" // replace this with your airtable base's id

	client, err := airtable.New(airtableAPIKey, baseID)
	if err != nil {
		panic(err)
	}

	csvFile, _ := os.Open(c.Args()[0])
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var linkedinContacts []contacts_base.LinkedinContact
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		actualLastNameList := strings.Split(line[1], " ")
		actualLastName := actualLastNameList[len(actualLastNameList)-1]

		linkedinContacts = append(linkedinContacts, contacts_base.LinkedinContact{
			Firstname:    line[0],
			Lastname:     actualLastName,
			EmailAddress: line[2],
			Company:      line[3],
			Position:     line[4],
		})
	}
	linkedinContactsJson, _ := json.Marshal(linkedinContacts)
	fmt.Println(string(linkedinContactsJson))

	for _, linkedInContact := range linkedinContacts {
		companies := []contacts_base.Company{}

		listParams := airtable.ListParameters{
			FilterByFormula: "{Name} = \"" + linkedInContact.Company + "\"",
		}

		err = client.ListRecords("Companies", &companies, listParams)
		if err != nil {
			panic(err)
		}

		fmt.Println("Filtering by: {Name} = \"" + linkedInContact.Company + "\"")
		fmt.Printf("%+v", companies)
		if len(companies) == 0 {
			newCompany := &contacts_base.Company{}
			newCompany.Fields.Name = linkedInContact.Company
			err = client.CreateRecord("Companies", &newCompany)
			if err != nil {
				panic(err)
			}

			fmt.Println("Creating new Company: ", linkedInContact.Company)
		}
	}

	return nil
}

func CmdImportContacts(c *cli.Context) error {
	// Write your code here

	airtableAPIKey := os.Getenv("AIRTABLE_API_KEY")
	baseID := "appVa3p4rNePGxslA" // replace this with your airtable base's id

	client, err := airtable.New(airtableAPIKey, baseID)
	if err != nil {
		panic(err)
	}

	csvFile, _ := os.Open(c.Args()[0])
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var linkedinContacts []contacts_base.LinkedinContact
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		actualLastNameList := strings.Split(line[1], " ")
		actualLastName := actualLastNameList[len(actualLastNameList)-1]

		linkedinContacts = append(linkedinContacts, contacts_base.LinkedinContact{
			Firstname:    line[0],
			Lastname:     actualLastName,
			EmailAddress: line[2],
			Company:      line[3],
			Position:     line[4],
		})
	}
	linkedinContactsJson, _ := json.Marshal(linkedinContacts)
	fmt.Println(string(linkedinContactsJson))

	for _, linkedInContact := range linkedinContacts {
		companies := []contacts_base.Company{}

		listParams := airtable.ListParameters{
			FilterByFormula: "{Name} = \"" + linkedInContact.Company + "\"",
		}

		err = client.ListRecords("Companies", &companies, listParams)
		if err != nil {
			panic(err)
		}

		contactCompanies := []string{companies[0].AirtableID}

		contacts := []contacts_base.Contact{}

		listParams = airtable.ListParameters{
			FilterByFormula: "{Email} = \"" + linkedInContact.EmailAddress + "\"",
		}

		err = client.ListRecords("Contacts", &contacts, listParams)
		if err != nil {
			panic(err)
		}

		if len(contacts) == 0 {
			newContact := &contacts_base.Contact{}
			newContact.Fields.First = linkedInContact.Firstname
			newContact.Fields.Last = linkedInContact.Lastname
			newContact.Fields.Email = linkedInContact.EmailAddress
			newContact.Fields.Title = linkedInContact.Position
			newContact.Fields.Company = contactCompanies

			err = client.CreateRecord("Contacts", &newContact)
			if err != nil {
				panic(err)
			}

			fmt.Println("Creating new Contact: ", linkedInContact.EmailAddress)
		}
	}

	return nil
}
func PrettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}

func main() {
}
