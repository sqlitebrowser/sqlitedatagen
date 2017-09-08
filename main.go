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
	numRows = 100000 // 100000 makes a 72MB file (taking ~4.5 seconds) on my Linux desktop.  Adjust to suit your desired target file size
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

	// Disable the journal
	err = sdb.Select("PRAGMA journal_mode=OFF", func(s *sqlite.Stmt) error {
		return nil
	})
	if err != nil {
		log.Printf("Error when disabling the journal: %s\n", err)
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
	results := make(chan *oneRow, cpus * 30) // 30 seems ok, less than 20 seems slightly slower (not properly measured though!)
	for w := 0; w < cpus; w++ {
		go worker(results)
	}

	// Bulk insert row data (inside a single transaction per table)
	log.Println("Adding data")
	//var r1 *oneRow
	//var r1, r2 *oneRow
	//var r1, r2, r3 *oneRow
	//var r1, r2, r3, r4 *oneRow
	var r1, r2, r3, r4, r5, r6, r7, r8, r9, r10 *oneRow
	for _, tbl := range tableNames {
		err = sdb.Begin()
		if err != nil {
			log.Printf("Error for Begin(): %s\n", err)
			return
		}

		// Prepare the insert statement
		dbQuery := sqlite.Mprintf(`
			INSERT into %w (col_key, col_int, col_signed, col_float, col_double, col_decim, col_date, col_code, col_name, col_address) VALUES
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			tbl)
		stmt, err := sdb.Prepare(dbQuery)
		if err != nil {
			log.Printf("Error when preparing statement for inserts: %s\n", err)
			return
		}

		// Insert the data
		numLoops := numRows / 10;
		for i := 0; i < numLoops; i++ {
			r1 = <- results
			r2 = <- results
			r3 = <- results
			r4 = <- results
			r5 = <- results
			r6 = <- results
			r7 = <- results
			r8 = <- results
			r9 = <- results
			r10 = <- results
			err = stmt.Exec(r1.key_data, r1.int_data, r1.signed_data, r1.float_data, r1.double_data, r1.decim_data, r1.date_data, r1.code_data, r1.name_data, r1.address_data,
				r2.key_data, r2.int_data, r2.signed_data, r2.float_data, r2.double_data, r2.decim_data, r2.date_data, r2.code_data, r2.name_data, r2.address_data,
				r3.key_data, r3.int_data, r3.signed_data, r3.float_data, r3.double_data, r3.decim_data, r3.date_data, r3.code_data, r3.name_data, r3.address_data,
				r4.key_data, r4.int_data, r4.signed_data, r4.float_data, r4.double_data, r4.decim_data, r4.date_data, r4.code_data, r4.name_data, r4.address_data,
				r5.key_data, r5.int_data, r5.signed_data, r5.float_data, r5.double_data, r5.decim_data, r5.date_data, r5.code_data, r5.name_data, r5.address_data,
				r6.key_data, r6.int_data, r6.signed_data, r6.float_data, r6.double_data, r6.decim_data, r6.date_data, r6.code_data, r6.name_data, r6.address_data,
				r7.key_data, r7.int_data, r7.signed_data, r7.float_data, r7.double_data, r7.decim_data, r7.date_data, r7.code_data, r7.name_data, r7.address_data,
				r8.key_data, r8.int_data, r8.signed_data, r8.float_data, r8.double_data, r8.decim_data, r8.date_data, r8.code_data, r8.name_data, r8.address_data,
				r9.key_data, r9.int_data, r9.signed_data, r9.float_data, r9.double_data, r9.decim_data, r9.date_data, r9.code_data, r9.name_data, r9.address_data,
				r10.key_data, r10.int_data, r10.signed_data, r10.float_data, r10.double_data, r10.decim_data, r10.date_data, r10.code_data, r10.name_data, r10.address_data)
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
