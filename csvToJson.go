package csvToJson

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadCsvFile(path *string) ([]byte, string) {
	csvFile, errorLoadFile := os.Open(*path)

	if errorLoadFile != nil {
		log.Fatal("!!!ERROR!!! the file is not found")
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvContent, _ := csvReader.ReadAll()

	if len(csvContent) < 1 {
		log.Fatal("!!!ERROR!!! file empty or length of the lines are not the same")
	}

	csvHeadersArray := make([]string, 0)
	for _, headers := range csvContent[0] {
		csvHeadersArray = append(csvHeadersArray, headers)
	}

	//Remove all headers from content
	csvContent = csvContent[1:]

	var objectBufferResult bytes.Buffer
	objectBufferResult.WriteString("[")
	for i, d := range csvContent {
		objectBufferResult.WriteString("{")
		for j, y := range d {
			objectBufferResult.WriteString(`"` + csvHeadersArray[j] + `":`)
			_, fErr := strconv.ParseFloat(y, 32)
			_, bErr := strconv.ParseBool(y)
			if fErr == nil {
				objectBufferResult.WriteString(y)
			} else if bErr == nil {
				objectBufferResult.WriteString(strings.ToLower(y))
			} else {
				objectBufferResult.WriteString((`"` + y + `"`))
			}
			//end of property
			if j < len(d)-1 {
				objectBufferResult.WriteString(",")
			}

		}
		//end of object of the array
		objectBufferResult.WriteString("}")
		if i < len(csvContent)-1 {
			objectBufferResult.WriteString(",")
		}
	}

	objectBufferResult.WriteString(`]`)
	rawMessage := json.RawMessage(objectBufferResult.String())
	x, _ := json.MarshalIndent(rawMessage, "", "  ")
	newFileName := filepath.Base(*path)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + ".json"
	r := filepath.Dir(*path)
	return x, filepath.Join(r, newFileName)
}
