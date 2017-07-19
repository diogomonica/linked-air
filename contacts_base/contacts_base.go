package contacts_base

var (
	// ContactsTableName is the name of the Airtable table containing Contact Records
	TasksTabContactsTableNameleName = "Contacts"
)

// Contact represents a single record in the `Contacts` Airtable table
type Contact struct {
	AirtableID string `json:"id,omitempty"`
	Fields     struct {
		Email    string   `json:"Email"`
		First    string   `json:"First"`
		Last     string   `json:"Last"`
		Title    string   `json:"Title"`
		Company  []string `json:"Company"`
		Linkedin string   `json:"Linkedin"`
		Comments string   `json:"Comments"`
	} `json:"fields"`
}

type LinkedinContact struct {
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	EmailAddress string `json:"emailaddress"`
	Company      string `json:"company"`
	Position     string `json:"position"`
}

// Company represents a single record in the `Contacts` Airtable table
type Company struct {
	AirtableID string `json:"id,omitempty"`
	Fields     struct {
		Name string `json:"Name"`
	} `json:"fields"`
}
