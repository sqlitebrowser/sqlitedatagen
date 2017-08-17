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
fileName = "gen.sqlite"
numRows = 10000
```

A row count of 10000 makes for a 5.3MB file on my Linux desktop.

Adjust to suit your desired target file size.


### Why?

I needed something to generate SQLite files of a specific size,
while tracking down a bug which looks related to the size of
files being transferred.
