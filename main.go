package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gwenn/gosqlite"
	"github.com/minio/go-homedir"
)

var (
	outputDir = "Databases"
	fileName = "72mb.sqlite"
	numRows = 100000 // 100000 makes a 72MB file (taking ~4.6 seconds) on my Linux desktop.  Adjust to suit your desired target file size
)

type oneRow struct {
	key_data int
	int_data int
	signed_data int
	float_data float32
	double_data float64
	decim_data string
	date_data string
	code_data string
	name_data string
	address_data string
}

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
		// No error occurred when looking for an existing file, which means something is there.  For now, we're just
		// going to blindly kill the existing thing without any kind of better safeguard
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

	// Enable WAL mode
	err = sdb.Select("PRAGMA journal_mode=WAL", func(s *sqlite.Stmt) error {
		return nil
	})
	if err != nil {
		log.Printf("Error when setting WAL mode: %s\n", err)
		return
	}

	// Turn off synchronous mode
	err = sdb.Exec("PRAGMA synchronous=OFF")
	if err != nil {
		log.Printf("Error when setting synchronous mode: %s\n", err)
		return
	}

	// Create the schema
	log.Println("Creating tables")
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

	// Launch a worker pool generating row data
	cpus := runtime.NumCPU()
	log.Printf("# of cpu's detected: %d.  Launching %d data generation workers\n", cpus, cpus)
	results := make(chan *oneRow, cpus * 5) // 5 seems ok, less than 5 seems slightly slower (not properly measured though!)
	for w := 0; w < cpus; w++ {
		go worker(results)
	}

	// Bulk insert row data (inside a single transaction per table)
	log.Println("Adding data")
	var r *oneRow
	for _, tbl := range tableNames {
		err = sdb.Begin()
		if err != nil {
			log.Printf("Error for Begin(): %s\n", err)
			return
		}

		// Prepare the insert statement
		dbQuery := sqlite.Mprintf(`
			INSERT into %w (col_key, col_int, col_signed, col_float, col_double, col_decim, col_date, col_code, col_name, col_address)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, tbl)
		stmt, err := sdb.Prepare(dbQuery)
		if err != nil {
			log.Printf("Error when preparing statement for inserts: %s\n", err)
			return
		}

		// Insert the data
		for i := 0; i < numRows; i++ {
			r = <- results
			err = stmt.Exec(r.key_data, r.int_data, r.signed_data, r.float_data, r.double_data, r.decim_data,
				r.date_data, r.code_data, r.name_data, r.address_data)
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
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = alphaNum[rand.Intn(len(alphaNum))]
	}
	return string(randomString)
}

// Goroutine which generates rows of test data
func worker(results chan <- *oneRow) {
	for {
		row := new(oneRow)
		row.key_data = rand.Int()
		row.int_data = rand.Int()
		row.signed_data = rand.Int()
		row.float_data = rand.Float32()
		row.double_data = rand.Float64()
		row.decim_data = fmt.Sprintf("%d.%d", rand.Intn(100000000000000000), rand.Intn(100))
		row.date_data = fmt.Sprintf("%d%d%d%d-%d%d-%d%d", rand.Intn(10), rand.Intn(10), rand.Intn(10),
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10))
		row.code_data = randomString(10)
		row.name_data = randomString(20)
		addLen := rand.Intn(80)
		if addLen < 8 {
			addLen = 8
		}
		row.address_data = randomString(addLen)
		results <- row
	}
}
