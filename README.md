## SQLite Data Generator

This is just a super simple utility for creating SQLite databases
filled with random data.  It's not even a CLI, it just has three
variables at the top to change, then compile/run it:

  * output directory for the SQLite file
  * output filename for the SQLite file
  * number of data rows in each table

eg:

```
outputDir = "Databases"
fileName = "72mb.sqlite"
numRows = 100000
```

A row count of 100000 makes a 72MB file (taking ~4.6 seconds) on my
(old) Linux desktop.

Adjust the variables to suit your desired target file name and
size, then run it:

```
$ go run main.go
2017/08/21 20:01:39 A SQLite database appears to be there already... removing it
2017/08/21 20:01:39 Creating new SQLite database file '/home/jc/Databases/72mb.sqlite'
2017/08/21 20:01:39 Creating tables
2017/08/21 20:01:39 # of cpu's detected: 4.  Launching 4 data generation workers
2017/08/21 20:01:39 Adding data
2017/08/21 20:01:44 SQLite database generation completed
```


### Why?

I needed something to generate SQLite files of a specific size,
while tracking down a bug which looks related to the size of
files being transferred.
