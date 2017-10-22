package main

import (
	"encoding/csv"
	"flag"
	"fmt"
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

func csvOutput(rec <-chan *record, outfile string) error {
	outputFile, err := os.Create(outfile)
	if err != nil {
		return fmt.Errorf("Cannot create output file: %v", err)
	}
	defer outputFile.Close()
	w := csv.NewWriter(outputFile)

	header := []string{"City", "State", "Zip", "Crrt", "RouteCount", "BundleCount", "BundlePerRoute", "LastBundle", "SeqNum", "FinalBundles", "Count"}
	if err := w.Write(header); err != nil {
		return fmt.Errorf("error writing header record to csv: %v", err)
	}

	for r := range rec {
		var finalBundle string
		for bundle := 1; bundle <= r.BundlePerRoute; bundle++ {
			switch {
			case bundle < r.BundlePerRoute:
				finalBundle = fmt.Sprint(r.BundleCount)
			case bundle == r.BundlePerRoute:
				finalBundle = fmt.Sprint(r.LastBundle)
			}
			w.Write([]string{r.City, r.State, r.Zip, r.Crrt,
				fmt.Sprint(r.RouteCount),
				fmt.Sprint(r.BundleCount),
				fmt.Sprint(r.BundlePerRoute),
				fmt.Sprint(r.LastBundle),
				fmt.Sprint(bundle),
				finalBundle,
				fmt.Sprintf("%v_of_%v", bundle, r.BundlePerRoute)})
		}
	}
	w.Flush()
	return nil
}

func remSep(p string) string {
	sep := []string{"-", ".", "*", "(", ")", ",", "\"", " "}
	for _, v := range sep {
		p = strings.Replace(p, v, "", -1)
	}
	return p
}

func main() {
	var (
		mrc     = flag.Int("mrc", 0, "Modify RouteCount number (defaults to 0)")
		outfile = flag.String("output", "output.csv", "Filename to export the CSV results")
	)
	flag.Parse()

	var counter int
	records := make(chan *record)

	go func() {
		input, err := os.Open(os.Args[len(os.Args)-1])
		if err != nil {
			log.Fatalln(err)
		}
		r := csv.NewReader(input)

		allRecs, err := r.ReadAll()
		if err != nil {
			log.Fatalln(err)
		}
		for idx, rec := range allRecs {
			if idx == 0 {
				continue
			}
			record, err := newRecord(rec, *mrc)
			if err != nil {
				log.Fatalf("%v on row [%v]\n", err, counter)
			}
			records <- record
		}
		close(records)
	}()

	if err := csvOutput(records, *outfile); err != nil {
		log.Printf("could not write to %s: %v", *outfile, err)
	}
	fmt.Println("Done!")
}
