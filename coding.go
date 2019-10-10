//	Copyright(c) 2019 Zdravko I. Genov
//	This software code is licenced under the terms of MIT licence.
//	For more information read the LICENCE.txt file in this repository.

package main

import (
	"fmt"
	"strconv"
	"time"
)

//int32 is the set of all signed 32-bit integers. Range: -2147483648 through => 2147483647,  32 most significant bits = 0
const divisor = 2147483647 //	2147483648 => 2^31 (on powers of 31)	=> So, I can use int32 as Results?

var numberofPairs = 40000000 //	5	=>	int for the loops

func main() { //   Commented sections are ONLY for comparison with the original definition of the problem

	//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	//--------------------------------------------------		1.1)	=>	Initial Parameters HERE:
	const startA = 65
	const startB = 8921

	const factorA = 16807
	const factorB = 48271
	//--------------------------------------------------		1.2)	=>	Generating Values HERE:

	startTotalSeq := time.Now()
	fmt.Println("=================================================")

	resultA := generator(startA, factorA)
	resultB := generator(startB, factorB)
	/*
		fmt.Println("resultA(first 5) = ", len(resultA), resultA[:5])
		fmt.Println("resultB(first 5) = ", len(resultB), resultB[:5])
		fmt.Println("=====================================================")
	*/
	//--------------------------------------------------		1.3)	=>	Finding Matching Pairs HERE:
	numberofMatches := comparator(resultA, resultB)
	//--------------------------------------------------		1.4)	=>	Printing Matching results HERE:

	fmt.Println("--------------------------------------------------")
	TotalTimeSeq := time.Since(startTotalSeq) //	.Nanoseconds()
	fmt.Printf("Total Time is = %v (%vns per value)\n", TotalTimeSeq, float64(TotalTimeSeq)/float64(numberofPairs))

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Total Number of Matching Pairs is: => %d\n", numberofMatches)
	fmt.Println("====================== DONE ======================" + "\n")
	//--------------------------------------------------		1.5)	=>	Total of 3 lines of code (2 generators + 1 comparator) + some printing
	time.Sleep(time.Duration(1) * time.Second)

	//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	//--------------------------------------------------		2.1)	Optimizing performance by using 2 goroutines with channels
	valuesA := make(chan []int64)
	valuesB := make(chan []int64)
	//--------------------------------------------------		2.2)	=>	Generating Values in parallel HERE:
	startTotalPar := time.Now()

	go gogen(startA, factorA, valuesA)
	go gogen(startB, factorB, valuesB)
	//--------------------------------------------------		2.3)	=>	Collecting comming results HERE:
	resultA = <-valuesA
	resultB = <-valuesB
	//--------------------------------------------------		2.4)	=>	Finding Matching Pairs HERE:
	numberofMatches = comparator(resultA, resultB)
	//--------------------------------------------------		2.5)	=>	Printing Matching results HERE:
	fmt.Println("--------------------------------------------------")
	TotalTimePar := time.Since(startTotalPar) //	.Nanoseconds()
	fmt.Printf("Total Time is = %v (%vns per value)\n", TotalTimePar, float64(TotalTimePar)/float64(numberofPairs))

	fmt.Println("--------------------------------------------------")
	fmt.Printf("Total Number of Matching Pairs is: => %d\n", numberofMatches)
	fmt.Println("====================== DONE ======================" + "\n")
	//--------------------------------------------------		2.6)	=>	Total of 7 lines of code (3 x 2 chan, gogen and results + 1 comparator) + some printing
	time.Sleep(time.Duration(1) * time.Second)
	//--------------------------------------------------		2.7)	=>	Printing performance results HERE:

	fmt.Printf("Parallel(goroutines) execution is %v %% faster than Sequential\n\n", int(100*(float64(TotalTimeSeq)/float64(TotalTimePar)-1)))
	fmt.Println("====================== END ======================" + "\n")
	//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
}

//	Generator of values, according to the business rules logic	=>	7 lines of code
func generator(startingValue int64, factor int64) (Result []int64) {

	defer measureTime(time.Now(), "generator"+strconv.FormatInt(startingValue, 10))
	previousValue := startingValue //	int64(startingValue)

	for ig := 1; ig <= numberofPairs; ig++ {

		productValue := previousValue * factor     //	Multiplication is resulting in more than uint32	=>	uint64
		nextValue := productValue % int64(divisor) //	Remainder of x / y	=> Guaranteed to be 32 bits

		Result = append(Result, nextValue) //	Populating slice of resulting values
		previousValue = nextValue          //	Preparing for the next run of the loop

	}
	return Result
}

//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

//	PS. => Performance optimization	=> concurent goroutines with return channels	=>	10 lines of code - 3 more for handling channels
func gogen(startingValue int64, factor int64, result chan []int64) {

	defer measureTime(time.Now(), "gogen"+strconv.FormatInt(startingValue, 10))

	genValues := []int64{}         //	Empty Slice
	previousValue := startingValue //	uint64(startingValue)
	defer close(result)            //	closing channel on return

	for ig := 1; ig <= numberofPairs; ig++ {

		productValue := previousValue * factor     //	Multiplication is resulting in more than uint32	=>	uint64(factor)
		nextValue := productValue % int64(divisor) //	Remainder of x / y	=> Guaranteed to be 32 bits

		genValues = append(genValues, nextValue) //	Populating slice of resulting values
		previousValue = nextValue                //	Preparing for the next run of the loop
	}

	result <- genValues //	Done sending...
	return              //	cleaning everything...
}

//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

//	Comparator of pairs of values, according to the business rules logic	=>	7 lines of code (reminds me of mining, nonces and difficulty...)
func comparator(resultA []int64, resultB []int64) (numberofMatches int) {

	defer measureTime(time.Now(), "comparator")
	numberofMatches = 0

	for ic := 0; ic < numberofPairs; ic++ { //	Comparator Code here =>

		resultA16LSB := resultA[ic] & 0x0000FFFF
		resultB16LSB := resultB[ic] & 0x0000FFFF
		/*
			if ic < 5 {
				fmt.Printf("resultA bits: => %032b\n", resultA[ic])
				fmt.Printf("resultB bits: => %032b\n", resultB[ic])
				fmt.Println("--------------------------------------------------")

				fmt.Printf("resultA16LSB: => %032b\n", resultA16LSB)
				fmt.Printf("resultB16LSB: => %032b\n", resultB16LSB)
				fmt.Println("--------------------------------------------------")
			}
		*/
		if resultB16LSB == resultA16LSB {
			numberofMatches++
			//	fmt.Println(numberofMatches)
		}
	}
	return numberofMatches
}

//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func measureTime(startTime time.Time, nameID string) {
	endTime := time.Since(startTime)
	fmt.Printf("%s took %v (%vns per value)\n", nameID, endTime, float64(endTime)/float64(numberofPairs))
}

//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
