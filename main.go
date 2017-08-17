package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"

	"github.com/gwenn/gosqlite"
	"github.com/minio/go-homedir"
	"os"
)

var (
	outputDir = "Databases"
	fileName = "gen.sqlite"
	numRows = 10000 // 10000 makes for a 5.3MB file on my Linux desktop.  Adjust to suit your desired target file size
)

func main() {
	// Determine full path to target file
	userHome, err := homedir.Dir()
	if err != nil {
		log.Printf("User home directory couldn't be determined: %s", "\n")
		return
	}
	fn := filepath.Join(userHome, outputDir, fileName)

	// If the database file already exists, nuke the file
	_, err = os.Stat(fn)
	if err == nil {
		// No error occurred when looking for an exists file, which means something is there.  For now, we're just
		// going to blindly kill the existing thing without an kind of better safeguards
		log.Printf("A SQLite database appears to be there already... removing it")
		os.Remove(fn)
	}

	// Create empty SQLite database
	log.Printf("Creating new SQLite database file '%s'\n", fn)
	sdb, err := sqlite.Open(fn, sqlite.OpenCreate | sqlite.OpenReadWrite)
	if err != nil {
		log.Printf("Couldn't open database: %s", err)
		return
	}
	defer sdb.Close()

	// Generate schema
	log.Println("Creating tables")

	// Create the tables
	tableNames := []string{"uniques", "updates", "hundred", "tenpct", "tiny"}
	var dbQuery string
	for _, tbl := range tableNames {
		dbQuery = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				'col_key'     INTEGER NOT NULL,
				'col_int'     INTEGER NOT NULL,
				'col_signed'  INTEGER NOT NULL,
				'col_float'   REAL NOT NULL,
				'col_double'  REAL NOT NULL,
				'col_decim'   NUMERIC NOT NULL,
				'col_date'    TEXT NOT NULL,
				'col_code'    TEXT NOT NULL,
				'col_name'    TEXT NOT NULL,
				'col_address' TEXT NOT NULL
			)`, tbl)
		err := sdb.Exec(dbQuery)
		if err != nil {
			log.Printf("Error when creating table '%s': %s\n", tbl, err)
			return
		}
	}

	// Bulk insert row data (inside a single transaction)
	log.Println("Adding data")
	for _, tbl := range tableNames {
		err = sdb.Begin()
		if err != nil {
			log.Printf("Error for Begin(): %s\n", err)
			return
		}

		// Prepare the insert statement
		dbQuery = fmt.Sprintf(`
		INSERT into %s (col_key, col_int, col_signed, col_float, col_double, col_decim, col_date, col_code, col_name, col_address)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, tbl)
		stmt, err := sdb.Prepare(dbQuery)
		if err != nil {
			log.Printf("Error when preparing statement for inserts: %s\n", err)
			return
		}

		for i := 0; i < numRows; i++ {
			// Generate test data
			key_data := rand.Int()
			int_data := rand.Int()
			signed_data := rand.Int()
			float_data := rand.Float32()
			double_data := rand.Float64()
			decim_data := fmt.Sprintf("%d.%d", rand.Intn(100000000000000000), rand.Intn(100))
			date_data := fmt.Sprintf("%d%d%d%d-%d%d-%d%d", rand.Intn(10), rand.Intn(10), rand.Intn(10),
				rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
			code_data := randomString(10)
			name_data := randomString(20)
			address_data := randomString(rand.Intn(80))

			// Insert the data
			err = stmt.Exec(key_data, int_data, signed_data, float_data, double_data, decim_data, date_data,
				code_data, name_data, address_data)
			if err != nil {
				log.Printf("Error when inserting data: %s\n", err)
				return
			}
		}

		// Clean up for this loop
		stmt.Finalize()

		// Commit the transaction
		sdb.Commit()
	}

	// TODO: Create indexes ?

	// Let the user know the program completed ok
	log.Printf("SQLite database generation completed")
}

// Generate a random string
func randomString(length int) string {
	const alphaNum = "abcdefghijklmnopqrstuvwxyz0123456789"
	randomString := make([]byte, 10)
	for i := range randomString {
		randomString[i] = alphaNum[rand.Intn(len(alphaNum))]
	}

	return string(randomString)
}
