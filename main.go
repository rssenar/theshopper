package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type field struct {
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

func (f *field) valdiate() {
	if !(f.BundleCount*(f.BundlePerRoute-1)+f.LastBundle == f.RouteCount) {
		log.Fatalln("Invalid Formula")
	}
}

func parseField(r []string) (*field, error) {
	RouteCount, err := strconv.Atoi(r[4])
	BundleCount, err := strconv.Atoi(r[5])
	BundlePerRoute, err := strconv.Atoi(r[6])
	LastBundle, err := strconv.Atoi(r[7])
	if err != nil {
		return nil, err
	}
	return &field{
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

// helper function to check and handle errors
func check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func main() {
	header := []string{"City", "State", "Zip", "Crrt", "RouteCount", "BundleCount", "BundlePerRoute", "LastBundle", "SeqNum", "FinalBundles", "Count"}

	r := csv.NewReader(os.Stdin)

	for counter := 0; ; counter++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		check(err)

		w := csv.NewWriter(os.Stdout)

		if counter == 0 {
			w.Write(header)
			w.Flush()
		} else {
			r, err := parseField(record)
			check(err)
			r.valdiate()

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
				strconv.Itoa(r.SeqNum + 1),
				strconv.Itoa(r.LastBundle),
				strconv.Itoa(r.SeqNum+1) + "_of_" + strconv.Itoa(r.BundlePerRoute),
			})
			w.Flush()
		}
	}
}
