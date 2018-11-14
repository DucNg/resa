// Package modele provides data structures for every objects.
// It also gives few functions to help building this structures.
// It links data structure from database to the code.
package modele

import (
	"regexp"
)

// Invite is the datastructure for invite. Mirror of invite on database.
// Links to parrain using parrain id isn't made her.
type Invite struct {
	Id      int64
	Nom     string
	Prenom  string
	Mail    string
	Mdp     string
	Numtel  string
	Parrain int64
	Voucher string
}

// CheckMail check email formatting using a regex.
func CheckMail(mail string) (bool, error) {
	regex := `^[^\W][a-zA-Z0-9_]+(\.[a-zA-Z0-9_]+)*\@[a-zA-Z0-9_]+(\.[a-zA-Z0-9_]+)*\.[a-zA-Z]{2,}$`
	matchFormat, err := regexp.MatchString(regex, mail)
	if err != nil {
		return false, err // Error checking regex
	}

	return matchFormat, nil // Tell if mail match regex or not
}

// GetParrainMail using a list of invite and a userId will return the parrain email.
// Only return the first match.
// This function is used to build the list of invite on the admin page.
func GetParrainMail(parrainId int64, inviteList []Invite) string {
	for _, element := range inviteList {
		if element.Id == parrainId {
			return element.Mail
		}
	}
	return ""
}
