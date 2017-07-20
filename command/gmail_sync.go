package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"

	"github.com/diogomonica/linked-air/contacts_base"

	"log"

	"github.com/fabioberger/airtable-go"
	"github.com/urfave/cli"
)

type contact struct {
	gmailID string
	date    string // retrieved from message header
	from    string
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(tokenCacheFile())
	if err != nil {
		tok = getTokenFromWeb(config)
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(tok)
		log.Printf("Save this token to /run/secrets/gmail_credentials:\n%s\n", b)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() string {
	gmailCredLocation := os.Getenv("GMAIL_CREDENTIALS_DIR")
	if gmailCredLocation == "" {
		gmailCredLocation = "/run/secrets/"
	}
	os.MkdirAll(gmailCredLocation, 0700)
	tokenCacheFile := filepath.Join(gmailCredLocation, "gmail_credentials")
	return tokenCacheFile
}

// clientSecretFile generates the secret file path/filename.
func clientSecretFile() string {
	gmailCredLocation := os.Getenv("GMAIL_CREDENTIALS_DIR")
	if gmailCredLocation == "" {
		gmailCredLocation = "/run/secrets/"
	}
	os.MkdirAll(gmailCredLocation, 0700)
	tokenCacheFile := filepath.Join(gmailCredLocation, "client_secret.json")
	return tokenCacheFile
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func setupGmailConfig() *gmail.Service {
	ctx := context.Background()

	b, err := ioutil.ReadFile(clientSecretFile())
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	return srv
}

func retrieveLastContacts(srv *gmail.Service) []contacts_base.Contact {
	contacts := []contacts_base.Contact{}
	user := "me"
	date := ""
	from := ""
	req := srv.Users.Messages.List(user)
	r, err := req.Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	log.Printf("Processing %v messages...\n", len(r.Messages))
	for _, m := range r.Messages {
		msg, err := srv.Users.Messages.Get(user, m.Id).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve message %v: %v", m.Id, err)
		}
		for _, h := range msg.Payload.Headers {
			if h.Name == "Date" {
				date = h.Value
			}
			if h.Name == "From" {
				from = h.Value
			}
		}

		emails, err := mail.ParseAddressList(from)
		if err != nil {
			log.Fatal(err)
		}
		if len(emails) < 1 {
			log.Fatal("no emails were found in FROM header")
		}

		splitName := strings.Split(emails[0].Name, " ")
		firstName := splitName[0]
		lastName := splitName[len(splitName)-1]
		email := emails[0].Address

		if email != "" && firstName != "" {
			log.Printf("Received email from: %s %s (%s) at %s\n", firstName, lastName, email, date)

			newContact := contacts_base.Contact{}
			newContact.Fields.First = firstName
			newContact.Fields.Last = lastName
			newContact.Fields.Email = email
			newContact.Fields.LastContact = time.Now().Format("2006-01-02")
			contacts = append(contacts, newContact)
		} else {
			log.Printf("Skipped processing received email: %s\n", emails[0])
		}
	}

	return contacts
}

func CmdGmailSync(c *cli.Context) error {
	for {
		srv := setupGmailConfig()
		contacts := retrieveLastContacts(srv)
		if len(contacts) < 1 {
			break
		}

		baseID := "appVa3p4rNePGxslA" // replace this with your airtable base's id
		airtableAPIKey := os.Getenv("AIRTABLE_API_KEY")
		if airtableAPIKey == "" {
			b, err := ioutil.ReadFile("/run/secrets/airtable-api-key")
			if err != nil {
				log.Fatalf("Error reading Airtable API from file: %v\n", err)
			}
			airtableAPIKey = string(b)
		}

		airClient, err := airtable.New(airtableAPIKey, baseID)
		if err != nil {
			log.Fatalf("Error accessing Airtable API: %v\n", err)
		}

		for _, contact := range contacts {
			listParams := airtable.ListParameters{
				FilterByFormula: "{Email} = \"" + contact.Fields.Email + "\"",
			}

			humanEmails := []contacts_base.HumanEmail{}
			err = airClient.ListRecords("All Emails", &humanEmails, listParams)
			if err != nil {
				log.Fatalf("Unable to access record in All Emails table: %v\n", err)
			}

			if len(humanEmails) == 0 {
				humanEmail := &contacts_base.HumanEmail{}
				humanEmail.Fields.Email = contact.Fields.Email

				err = airClient.CreateRecord("All Emails", &humanEmail)
				if err != nil {
					log.Fatalf("Unable to access record in All Emails table: %v\n", err)
				}
			} else {
				if humanEmails[0].Fields.Human == true && humanEmails[0].Fields.Ignore == false {
					foundContacts := []contacts_base.Contact{}
					err = airClient.ListRecords("Contacts", &foundContacts, listParams)
					if err != nil {
						log.Fatalf("Unable to access record in All Emails table: %v\n", err)
					}
					if len(foundContacts) == 0 {
						err = airClient.CreateRecord("Contacts", &contact)
						if err != nil {
							log.Fatalf("Unable to access record in Contacts table: %v\n", err)
						}

						log.Printf("Creating new Contact: %s\n", contact.Fields.Email)
					} else {
						log.Printf("Contact already exists: %s\n", contact.Fields.Email)
					}
				}
			}

		}
		time.Sleep(c.Duration("howlong"))
	}
	return nil
}
