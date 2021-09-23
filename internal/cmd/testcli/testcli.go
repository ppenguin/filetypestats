package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ppenguin/filetypestats"
	"github.com/ppenguin/filetypestats/ftsdb"
	"github.com/ppenguin/filetypestats/types"
	utils "github.com/ppenguin/gogenutils"
)

func main() {
	pscandirs := flag.String("dirs", "./", "directories to scan, comma-separated")
	dbfile := flag.String("db", "scandb.sqlite", "database in which the scan result is stored")

	rm := flag.Bool("rm", false, "remove database if exists")
	flag.Parse()

	if len(flag.Args()) == 0 {
		usage()
	}

	scandirs := strings.Split(*pscandirs, ",")

	switch flag.Arg(0) {
	case "scan":
		if *rm {
			os.Remove(*dbfile)
		}
		scan(scandirs, *dbfile)
	case "show":
		show(scandirs, *dbfile)
	default:
		usage()
	}
}

func usage() {
	fmt.Printf(
		"Usage: %s [ --dirs=dir1,dir2 ] [ --db=scandb.sqlite ] [ scan | show ]\n"+
			"\tscan: scans all dirs given recursively and stores statistics per dir in scandb\n"+
			"\tshow: gets the totals from scandb for the given dirs.\n"+
			"\t\tTo show totals under a dir, use the special form --dir='/dir/to/%%' (remember quoting if necessary)\n", os.Args[0])
	os.Exit(0)
}

func exiterr(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s", err.Error())
	os.Exit(1)
}

func scan(dirs []string, file string) {
	fmt.Printf("Scanning %v to database %s...\n", dirs, file)
	ts := time.Now()
	if fstats, err := filetypestats.WalkFileTypeStatsDB(dirs, file); err != nil {
		exiterr(err)
	} else {
		fmt.Printf("Scanning took %s\n\n", time.Since(ts))
		fmt.Println("Scan totals:")
		printstats(fstats)
	}
}

func show(dirs []string, file string) {
	db, err := ftsdb.New(file, false)
	if err != nil {
		exiterr(err)
	}
	defer db.Close()
	ts := time.Now()
	fstats, err := db.FTStatsDirsSum(dirs)
	if err != nil {
		exiterr(err)
	}
	fmt.Printf("Query took %s\n\n", time.Since(ts))
	fmt.Println("Query totals:")
	printstats(fstats)
}

func printstats(fstats types.FileTypeStats) {
	totCount := uint(0)
	totSize := uint64(0)
	for k, v := range fstats {
		fmt.Printf("%d %s files taking %s of space\n", v.FileCount, k, utils.ByteCountSI(v.NumBytes))
		totCount += v.FileCount
		totSize += v.NumBytes
	}
	fmt.Printf("\nTotal %d files taking %s of space\n", totCount, utils.ByteCountSI(totSize))
}
