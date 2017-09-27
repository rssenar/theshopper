package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type record struct {
	City           string
	State          string
	Zip            string
	Crrt           string
	RouteCount     int
	BundleCount    int
	BundlePerRoute int
	LastBundle     int
	SeqNum         int
	FinalBundles   int
	Count          string
}

func newRecord(r []string, mod int) (*record, error) {
	RouteCount, err := strconv.Atoi(r[4])
	BundleCount, err := strconv.Atoi(r[5])
	if err != nil {
		return nil, err
	}
	RouteCount = RouteCount + mod

	var BundlePerRoute int
	var LastBundle int

	if RouteCount%BundleCount == 0 {
		BundlePerRoute = RouteCount / BundleCount
		LastBundle = BundleCount
	} else {
		BundlePerRoute = (RouteCount / BundleCount) + 1
		LastBundle = RouteCount - (BundleCount * (BundlePerRoute - 1))
	}

	return &record{
		City:           r[0],
		State:          r[1],
		Zip:            r[2],
		Crrt:           r[3],
		RouteCount:     RouteCount,
		BundleCount:    BundleCount,
		BundlePerRoute: BundlePerRoute,
		LastBundle:     LastBundle,
	}, nil
}

func main() {
	var (
		mrc     = flag.Int("mrc", 0, "Modify RouteCount number (defaults to 0)")
		outfile = flag.String("output", "output.csv", "Filename to export the CSV results")
	)
	flag.Parse()

	header := []string{"City",
		"State",
		"Zip",
		"Crrt",
		"RouteCount",
		"BundleCount",
		"BundlePerRoute",
		"LastBundle",
		"SeqNum",
		"FinalBundles",
		"Count",
	}

	r := csv.NewReader(os.Stdin)

	of, err := os.Create(*outfile)
	if err != nil {
		log.Fatalln(err)
	}
	defer of.Close()
	w := csv.NewWriter(of)

	for counter := 0; ; counter++ {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		if counter == 0 {
			w.Write(header)
			w.Flush()
		} else {
			r, err := newRecord(rec, *mrc)
			if err != nil {
				log.Fatalln(err)
			}

			for bndl := 1; bndl < r.BundlePerRoute; bndl++ {
				r.SeqNum = bndl
				w.Write([]string{
					r.City,
					r.State,
					r.Zip,
					r.Crrt,
					strconv.Itoa(r.RouteCount),
					strconv.Itoa(r.BundleCount),
					strconv.Itoa(r.BundlePerRoute),
					strconv.Itoa(r.LastBundle),
					strconv.Itoa(r.SeqNum),
					strconv.Itoa(r.BundleCount),
					strconv.Itoa(r.SeqNum) + "_of_" + strconv.Itoa(r.BundlePerRoute),
				})
				w.Flush()
				r.SeqNum++
			}
			w.Write([]string{
				r.City,
				r.State,
				r.Zip,
				r.Crrt,
				strconv.Itoa(r.RouteCount),
				strconv.Itoa(r.BundleCount),
				strconv.Itoa(r.BundlePerRoute),
				strconv.Itoa(r.LastBundle),
				strconv.Itoa(r.SeqNum),
				strconv.Itoa(r.LastBundle),
				strconv.Itoa(r.SeqNum) + "_of_" + strconv.Itoa(r.BundlePerRoute),
			})
			w.Flush()
		}
	}
	fmt.Println("Job Completed!")
}
