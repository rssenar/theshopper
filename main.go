package main

import (
	"encoding/csv"
	"flag"
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
	RouteCount = RouteCount - mod
	BundlePerRoute := RouteCount / BundleCount
	LastBundle := RouteCount - (BundleCount * BundlePerRoute)

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
		mod = flag.Int("mod", 0, "Modify RouteCount qty")
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

	for counter := 0; ; counter++ {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		w := csv.NewWriter(os.Stdout)

		if counter == 0 {
			w.Write(header)
			w.Flush()
		} else {
			r, err := newRecord(rec, *mod)
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
}
