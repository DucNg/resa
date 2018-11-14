// Package tools provide function to work with the database.
// It gets structs from modele and fill them or insert in database using them.
package tools

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/DucNg/resa/config"
	"github.com/DucNg/resa/modele"
)

// Connect to the database using the file provided in config and return a DB object.
// It needs to be used before every request to the database.
// You can't use any function in this file without getting the DB object so you need to use this function.
// Changing the SGBD should only be done here.
func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", *config.DbFile) // Connect to the sqlite database
	if err != nil {
		return nil, err
	}

	return db, nil // Return pointer to the connection
}

// Disconnect close the connection to the database.
// Needs to be called after every request. Using defer is highly recommanded.
// It only contrain one line for the moment but it could grow bigger in the futur if needed so keep using this to disconnect.
func Disconnect(db *sql.DB) {
	db.Close()
}

// HashPassword get password as string and return hashed password as string using bcrypt.
// bcrypt provide a highly secure way to generate password, using hashing and salt.
// Playing with password should only be done using this.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash check the password validity using bcrypt.
// Complementatry to HashPassword().
// Playing with password should only be done using this.
func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser use a Invite struct from modele to insert the invite into the database.
// It hash the password provided using HashPassword()
// Provided informations can be **empty** but **not nil**!!!
func CreateUser(db *sql.DB, i *modele.Invite) (int64, error) { // Create user, return user id or error
	tx, err := db.Begin() // Start transaction
	if err != nil {
		return -1, err
	}

	defer tx.Rollback() // Close transaction no matter what
	stmt, err :=
		tx.Prepare("INSERT INTO Invite(id_invite,nom,prenom,mail,mdp,numtel,parrain)" +
			" VALUES (NULL,?,?,?,?,?,?)") // Insert into Invite
	if err != nil {
		return -1, err
	}
	defer stmt.Close() // Close the statement no matter what

	hashedPsw, err := HashPassword(i.Mdp) // Hashing the password before sending to database
	if err != nil {
		return -1, err
	}

	result, err := stmt.Exec( // Fill placeholders
		i.Nom,
		i.Prenom,
		i.Mail,
		hashedPsw, // The password is hashed using bcrypt
		i.Numtel,
		i.Parrain,
	)
	if err != nil {
		return -1, err
	}

	err = tx.Commit() // Commit changes to database
	if err != nil {
		return -1, err
	}

	return result.LastInsertId() // Return the id of the created user or error
}

// UniqueMail tell if the provided email is unique in database or not
// This function doesn't use a modele, it could be merged with CreateUser somehow.
func UniqueMail(db *sql.DB, mail string) (bool, error) {
	tx, err := db.Begin() // Start transaction
	if err != nil {
		return false, err
	}

	defer tx.Rollback() // Close transaction no matter what
	stmt, err :=
		tx.Prepare("SELECT COUNT(*) FROM Invite" +
			" WHERE mail = ?") // Count line with selected email
	if err != nil {
		return false, err
	}
	defer stmt.Close() // Close the statement no matter what

	result, err := stmt.Query(mail) // Fill placeholder and execute query
	defer result.Close()

	result.Next() // No iteration because count * should only return one line
	var numOccurences int
	err = result.Scan(&numOccurences) // Getting result
	if err != nil {
		return false, err
	}
	return numOccurences <= 0, nil // Expect 0 if mail is unique
}

// CheckVoucher check the validity of a voucher. It check existance and expiration time.
// Return true and nil in case of sucess, return false and specify why in err if failed
func CheckVoucher(db *sql.DB, code string) (bool, error) {
	result, err := db.Query("SELECT id_voucher,code,expiration,proprietaire"+
		" FROM Voucher WHERE code = ?", code)
	if err != nil {
		return false, err
	}
	defer result.Close()

	if result.Next() { // No iteration because voucher should be unique
		var voucher modele.Voucher
		voucher = modele.Voucher{}
		err = result.Scan(
			&voucher.Id,
			&voucher.Code,
			&voucher.Expiration,
			&voucher.Prop,
		)
		if voucher.Expiration.After(time.Now()) {
			return true, nil
		}
		return false, errors.New("Voucher expired") // A voucher was found but expired
	}
	return false, errors.New("Voucher doesn't exist") // Nothing was found

}

// GetParrain return the parrain id for a voucher and check voucher validity.
// This function is used to link Invite to his parrain on registration.
// Doesn't use a modele, should be merged with CreateUser() somehow.
// Return parrain user id or -1 if voucher is invalid
func GetParrain(db *sql.DB, code string) (int64, error) {
	tx, err := db.Begin() // Start transaction
	if err != nil {
		return -1, err
	}

	defer tx.Rollback() // Close transaction no matter what
	stmt, err :=
		tx.Prepare("SELECT id_invite FROM Invite, Voucher" +
			" WHERE proprietaire = id_invite" +
			" AND code = ?") // Count line with selected email
	if err != nil {
		return -1, err
	}
	defer stmt.Close() // Close the statement no matter what

	// Check voucher validity
	voucherValid, err := CheckVoucher(db, code)
	if !voucherValid {
		return -1, err
	}

	result, err := stmt.Query(code) // Fill placeholder and execute query
	defer result.Close()

	if result.Next() { // No iteration because voucher should be unique
		var idParrain int64
		err = result.Scan(&idParrain)
		return idParrain, err
	}
	return -1, nil // No error but id parrain is -1 mean nothing found
}

// ConnectUser fill the Invite modele provided using his email as a filter.
// It check if email and password match.
// Info from database can be **empty** but **can't be nil**!!
func ConnectUser(db *sql.DB, i *modele.Invite) error { // Receive result in i by reference
	var notHashedPsw string = i.Mdp // Getting not hashed passowrd from form
	var hashedPsw string

	tx, err := db.Begin() // Start transaction
	if err != nil {
		return err
	}

	defer tx.Rollback() // Close transaction no matter what
	stmt, err :=
		tx.Prepare("SELECT * FROM Invite" +
			" WHERE mail = ?") // Select invite from his mail address
	if err != nil {
		return err
	}
	defer stmt.Close() // Close the statement no matter what

	result, err := stmt.Query(i.Mail) // Fill placeholder and execute query
	defer result.Close()

	if result.Next() { // No iteration because mail should be unique (PRIMARY KEY)
		err = result.Scan( // Fill invite
			&i.Id,
			&i.Nom,
			&i.Prenom,
			&i.Mail,
			&hashedPsw, // Getting hashed password from database
			&i.Numtel,
			&i.Parrain,
		)

		// Check password
		if CheckPasswordHash(notHashedPsw, hashedPsw) {
			return err // Password correct return invite
		}
		err = errors.New("Connect: Incorrect password")
		return err

	}
	err = errors.New("Connect: Incorrect mail")
	return err
}

// ListInvite fill the slice with every invite in database.
// TODO This should be improve with paging to avoid crash/lag/slowing/instability with heavy database.
// It should still work very fast if the number of registration is < 200
// Info from database can be **empty** but **can't be nil**!!
func ListInvite(db *sql.DB, listI *[]modele.Invite) error {
	result, err := db.Query("SELECT id_invite,nom,prenom,mail,numtel,parrain" +
		" FROM Invite ORDER BY nom")
	if err != nil {
		return err
	}
	defer result.Close()

	var errL string     // Could have multiple errors
	for result.Next() { // Iterate on each invite to fill the slice
		var inviteTmp modele.Invite // Fill this one first and then append it to slice
		inviteTmp = modele.Invite{} // Initialize it
		err = result.Scan(          // Fill invite
			&inviteTmp.Id,
			&inviteTmp.Nom,
			&inviteTmp.Prenom,
			&inviteTmp.Mail,
			&inviteTmp.Numtel,
			&inviteTmp.Parrain,
		)
		if err != nil { // If something goes wrong during iteration don't screw up everything, keep going and keep errors for later
			errL += err.Error() // Handle multiple errors
		}

		*listI = append(*listI, inviteTmp)
	}
	if errL != "" {
		return errors.New(errL) // get any error encountered during iteration
	}
	return err
}

// GetVouchers extract all vouchers from database in an HashMap associating userId with voucher code.
// TODO This should be improve with paging to avoid crash/lag/slowing/instability with heavy database.
// This isn't much of an issue because hashmap is fast. Needs testing.
// Info from database can be **empty** but **can't be nil**!!
func GetVouchers(db *sql.DB, vouchers map[int64]modele.Voucher) error {
	result, err := db.Query("SELECT id_invite,id_voucher,code,expiration,proprietaire" +
		" FROM Voucher,Invite" +
		" WHERE id_invite = proprietaire")
	if err != nil {
		return err
	}
	defer result.Close()

	var errL string // Could have multiple errors
	for result.Next() {
		var id int64
		var voucherTmp modele.Voucher
		err = result.Scan(
			&id,
			&voucherTmp.Id,
			&voucherTmp.Code,
			&voucherTmp.Expiration,
			&voucherTmp.Prop,
		)
		vouchers[id] = voucherTmp // Build the map with every vouchers, associate with id_invite
		if err != nil {
			errL += err.Error() // Handle multiple errors
		}
	}
	if errL != "" {
		return errors.New(errL) // get any error encountered during iteration
	}
	return err
}

// AddVoucher add a voucher in database using a modele.
// Values can be empty but can't be nil or it will troublesome when getting them.
func AddVoucher(db *sql.DB, voucher modele.Voucher) error {
	_, err := db.Exec("INSERT INTO Voucher(code,expiration,proprietaire)"+
		" VALUES (?,?,?)",
		voucher.Code, voucher.Expiration, voucher.Prop)
	return err
}

// DisableVoucher disable a voucher in database using a userId.
// Disable a voucher means set his expiration date to Thu Jan 01 00:00:00 1970 UTC (UNIX time 0)
func DisableVoucher(db *sql.DB, idInvite int64) error {
	disableTime := time.Unix(0, 0) // Disable a voucher means set his expiration date to Thu Jan 01 00:00:00 1970 UTC

	_, err := db.Exec("UPDATE Voucher SET expiration = ? WHERE proprietaire = ?", disableTime, idInvite)
	return err
}

// GetInvite return a Invite modele using a id_invite.
// Issue the request on database and then fill the invite using info from database. Then return it.
// Improvement: could be merge with ListInvite() since they're quiet similar.
func GetInvite(db *sql.DB, id_invite int64) (modele.Invite, error) {
	var invite modele.Invite = modele.Invite{}
	result, err := db.Query("SELECT id_invite,nom,prenom,mail,numtel,parrain"+
		" FROM Invite WHERE id_invite = ?", id_invite)
	if err != nil {
		return invite, err
	}
	defer result.Close()
	result.Next()
	err = result.Scan( // Fill invite
		&invite.Id,
		&invite.Nom,
		&invite.Prenom,
		&invite.Mail,
		&invite.Numtel,
		&invite.Parrain,
	)
	return invite, err
}

// ConnectAdmin check login password provided in the admin modele and fill the same struct on sucess
func ConnectAdmin(db *sql.DB, admin *modele.Admin) error {
	var notHashedPsw string = admin.Psw
	var hashedPsw string

	result, err := db.Query("SELECT id_admin,login,mdp FROM Administrateur WHERE login = ?", admin.Login)
	defer result.Close()

	if !result.Next() {
		err := errors.New("Connect admin: Incorrect login")
		return err
	}
	err = result.Scan(
		&admin.IdAdmin,
		&admin.Login,
		&hashedPsw,
	)
	// Check password
	if CheckPasswordHash(notHashedPsw, hashedPsw) {
		admin.Psw = hashedPsw
		return err // Password correct return invite
	}
	err = errors.New("Connect admin: Incorrect password")
	return err
}
