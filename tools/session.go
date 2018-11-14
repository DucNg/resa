// Database session based system
package tools

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	_ "github.com/mattn/go-sqlite3"

	"github.com/DucNg/resa/modele"
)

// generateRandomString generate a random base64 encoded string using crypto/rand.
// It gets a 20 random bytes and encode them to base64.
func generateRandomString() (string, error) {
	c := 20
	b := make([]byte, c)

	_, err := rand.Read(b)
	if err != nil {
		err = errors.New("Error generating random")
		return "error", err
	}
	randomString := base64.StdEncoding.EncodeToString(b)
	return randomString, err
}

// CreateSession insert a random token linked to a user id in database.
// This token is used to authentificate the user.
// It's the server side session.
// Return the generated token in case of sucess.
func CreateSession(db *sql.DB, idUser int64) (string, error) {
	randomString, err := generateRandomString()

	if err != nil {
		return "error", err
	}
	//strconv.FormatInt(idUser,10)
	_, err = db.Exec("INSERT INTO Session VALUES (?,?)", randomString, idUser)

	return randomString, err // Return the inserted token
}

// VerifySession check if the token exist in database a return the associated invite in case of success.
// --> From a token, get an Invite
// Return an invite modele.
func VerifySession(db *sql.DB, token string) (modele.Invite, error) {
	var i modele.Invite

	result, err := db.Query("SELECT id_invite,nom,prenom,mail,mdp,numtel,parrain"+
		" FROM Invite,Session"+
		" WHERE id_invite = id_user AND token = ?",
		token)
	if err != nil {
		return i, err
	}
	defer result.Close()

	if result.Next() { // No iteration because token should be unique
		err = result.Scan( // Fill invite
			&i.Id,
			&i.Nom,
			&i.Prenom,
			&i.Mail,
			&i.Mdp, // Getting hashed password from database
			&i.Numtel,
			&i.Parrain,
		)
		return i, err
	}
	err = errors.New("Verify session: No user found") // Session has been deleted or incorrect token
	return i, err                                     // No user found for this token, i should be nil
}

// DeleteSession delete the session from the specified token.
func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec("DELETE FROM Session WHERE token = ?", token)
	return err
}

// CreateAdminSession equivalent to CreateSession but for admin.
func CreateAdminSession(db *sql.DB, idAdmin int64) (string, error) {
	randomString, err := generateRandomString()

	if err != nil {
		return "error", err
	}
	_, err = db.Exec("INSERT INTO AdminSession VALUES (?,?)", randomString, idAdmin)

	return randomString, err // Return the inserted token
}

// CreateAdminSession equivalent to VerifySession but for admin.
func VerifyAdminSession(db *sql.DB, token string) (bool, error) {
	result, err := db.Query("SELECT * FROM AdminSession WHERE token = ?", token)
	defer result.Close()

	return result.Next(), err
}
