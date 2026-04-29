package shogi

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

func ReadRecords(file_name string) []Record {
	var data = make([]Record, 0)
	fp, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	reader := csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			var r Record
			var l = len(line)
			r.result = line[0]
			a, _ := strconv.Atoi(line[1])
			r.ply = a
			for i := 2; i < l; i++ {
				r.str_moves = append(r.str_moves, line[i])
			}
			data = append(data, r)
		}
	}
	return data
}
