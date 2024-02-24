package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func splitCoords(r *csv.Reader, coords [][]string) ([][]string, [][]string) {
	ggs := make([][]string, 0)
	vr := make([][]string, 0)
	for _, row := range coords {
		if row[0] == "ggs" {
			ggs = append(ggs, row)
		} else {
			vr = append(vr, row)
		}
	}
	return ggs, vr
}

func computeDiffs(ggsCoord string, vrCoord string) float64 {
	observed, _ := strconv.ParseFloat(vrCoord, 64)
	refPoint, _ := strconv.ParseFloat(ggsCoord, 64)
	return observed - refPoint
}

// strconv.FormatFloat(first-second, 'f', 3, 64)
func makeEquationTable(ggs [][]string, vr [][]string) [][]string {
	equationTable := make([][]string, 0)
	for _, point := range vr {
		for _, controlPoint := range ggs {
			cp := controlPoint[1:]
			p := point[1:]

			current := len(equationTable)
			equationTable = append(equationTable, []string{})
			equationTable[current] = append(equationTable[current], cp[0]+"-"+p[0])

			min := -0.006
			max := 0.006
			randSko1 := min + rand.Float64()*(max-min)
			randSko2 := min + rand.Float64()*(max-min)
			randSko3 := min + rand.Float64()*(max-min)

			dx := computeDiffs(cp[1], p[1]) + randSko1
			dy := computeDiffs(cp[2], p[2]) + randSko2
			dh := computeDiffs(cp[3], p[3]) + randSko3

			equationTable[current] = append(equationTable[current],
				strconv.FormatFloat(dx, 'f', 3, 64))
			equationTable[current] = append(equationTable[current],
				strconv.FormatFloat(dy, 'f', 3, 64))
			equationTable[current] = append(equationTable[current],
				strconv.FormatFloat(dh, 'f', 3, 64))

			equationTable[current] = append(equationTable[current],
				strconv.FormatFloat(math.Sqrt(randSko1*randSko1+randSko2*randSko2), 'f', 3, 64))
			equationTable[current] = append(equationTable[current],
				strconv.FormatFloat(math.Abs(randSko3), 'f', 3, 64))

		}
	}
	return equationTable
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("error! execution is \"./coord <path_to_file>\"")
		os.Exit(-1)
	}
	input := os.Args[1]
	f, err := os.Open(input)
	if err != nil {
		fmt.Println("error! no input file:", err)
		os.Exit(-1)
	}
	defer f.Close()
	r := csv.NewReader(f)
	coords, err := r.ReadAll()
	if err != nil {
		fmt.Println("error! something went wrong while reading file:", err)
		os.Exit(-1)
	}

	ggs, vr := splitCoords(r, coords)
	newName, _ := strings.CutSuffix(f.Name(), ".csv")
	newFile, err := os.Create(newName + "_equation_table.csv")
	if err != nil {
		fmt.Println("error! something went wrong while writing file")
		os.Exit(-1)
	}
	defer newFile.Close()
	et := makeEquationTable(ggs, vr)
	if len(et) == 0 {
		fmt.Println("No data given or wrong format")
		return
	}
	w := csv.NewWriter(newFile)
	w.Write([]string{"Имя", "dN (m)", "dE (m)", "dHt (m)", "СКО в плане (m)", "СКО по высоте (m)"})
	for _, row := range et {
		w.Write(row)
		w.Flush()
	}
}
