package main

import (
	"os"
	"fmt"
	"log"
	"syscall"
	"github.com/cornfeedhobo/pflag"
	"github.com/markkurossi/tabulate"
)

func readable(b int64, d int) string {
	f := float64(b)
	s := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	for _, p := range s {
		if f <= 1024.0 {
			return fmt.Sprintf("%.[3]*[1]f%s", f, p, d)
		}

		f /= 1024.0
	}

	return fmt.Sprintf("%.2f", f)
}

func main() {
	var decimal = pflag.IntP("decimals", "d", 2, "Gets the number of decimal places")
	var ddir = pflag.StringP("dir", "p", "", "Define a path")
	var format = pflag.StringP("format", "f", "02-01-2006 15:04:05", "Formats the date")
	var totalBytes int64
	totalBytes = 0

	pflag.Parse()
	wd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	if *ddir != "" {
		wd = *ddir
	}

	dir, err := os.ReadDir(wd)

	if err != nil {
		log.Fatal(err)
	}

	tab := tabulate.New(tabulate.UnicodeLight)
	tab.Header("File")
	tab.Header("Size")
	tab.Header("Date Modified")
	tab.Header("User ID")
	tab.Header("Group ID")

	for _, e := range dir {
		row := tab.Row()
		row.Column(e.Name())
		if e.IsDir() {
			row := tab.Row()
			row.Column("Dir")
		} else {
			info, err := os.Stat(e.Name())
			if err != nil {
				log.Fatal(err)
			}

			stat := info.Sys().(*syscall.Stat_t)
			totalBytes += info.Size()
			row.Column(readable(info.Size(), *decimal))
			row.Column(info.ModTime().Format(*format))
			row.Column(fmt.Sprint(stat.Uid))
			row.Column(fmt.Sprint(stat.Uid))
		}
	}

	tab.Print(os.Stdout)
	fmt.Println("\033[1;37m" + readable(totalBytes, *decimal) + " \033[0;34mBytes Fetched")
	fmt.Printf("\033[0m")
}
