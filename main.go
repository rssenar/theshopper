package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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

func remSep(p string) string {
	sep := []string{"-", ".", "*", "(", ")", ",", "\"", " "}
	for _, v := range sep {
		p = strings.Replace(p, v, "", -1)
	}
	return p
}

func newRecord(r []string, mrc int) (*record, error) {
	RouteCount, err := strconv.Atoi(remSep(r[4]))
	if err != nil {
		return nil, fmt.Errorf("Error parsing RouteCount :: %s", err)
	}

	BundleCount, err := strconv.Atoi(remSep(r[5]))
	if err != nil {
		return nil, fmt.Errorf("Error parsing BundleCount :: %s", err)
	}

	RouteCount = RouteCount + mrc
	if RouteCount <= 0 {
		return nil, fmt.Errorf("RouteCount Cannot be <= 0")
	}

	var BundlePerRoute int
	var LastBundle int

	switch {
	case RouteCount%BundleCount == 0:
		BundlePerRoute = RouteCount / BundleCount
		LastBundle = BundleCount
	default:
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
				log.Fatalf("%v on row [%v]", err, counter)
			}

			for bndl := 1; bndl < r.BundlePerRoute; bndl++ {
				r.SeqNum = bndl
				w.Write([]string{
					r.City,
					r.State,
					r.Zip,
					r.Crrt,
					fmt.Sprint(r.RouteCount),
					fmt.Sprint(r.BundleCount),
					fmt.Sprint(r.BundlePerRoute),
					fmt.Sprint(r.LastBundle),
					fmt.Sprint(r.SeqNum),
					fmt.Sprint(r.BundleCount),
					fmt.Sprintf("%v_of_%v", r.SeqNum, r.BundlePerRoute),
				})
				w.Flush()
				r.SeqNum++
			}
			w.Write([]string{
				r.City,
				r.State,
				r.Zip,
				r.Crrt,
				fmt.Sprint(r.RouteCount),
				fmt.Sprint(r.BundleCount),
				fmt.Sprint(r.BundlePerRoute),
				fmt.Sprint(r.LastBundle),
				fmt.Sprint(r.SeqNum),
				fmt.Sprint(r.LastBundle),
				fmt.Sprintf("%v_of_%v", r.SeqNum, r.BundlePerRoute),
			})
			w.Flush()
		}
	}
	fmt.Println("Job Completed!")
}
