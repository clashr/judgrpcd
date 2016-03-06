package services

import (
	"log"
	"math"
)

func Grade(runs []RunDetails, tests []string) int {
	scores := make([]int, len(tests))
	for i, test := range tests {
		// For each case, we use a logistical model with respect to the
		// time each process uses and the size of the input data.
		k := 1.0 / (25.0 * math.Log(float64(len(test))))
		x := float64(runs[i].TTotal.Seconds())
		scores[i] = int(50000.0 / (1 + math.Exp(k*x)))
		log.Println(k, x, scores[i])
	}

	// We use a harmonic mean to highlight the poor scaling in algorithm
	// design or implementation, and to hide testing variability.
	final := hm(scores)

	log.Printf("Program recieved score of %f", final)
	return int(final)
}

func rms(collection []int) float64 {
	sumSquares := 0.0
	for _, value := range collection {
		sumSquares += math.Pow(float64(value), 2)
	}
	return math.Sqrt(sumSquares) / float64(len(collection))
}

func hm(collection []int) float64 {
	sumInvs := 0.0
	for _, value := range collection {
		sumInvs += 1.0 / float64(value)
	}
	return float64(len(collection)) / sumInvs
}
