package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/m3u/helper"
	"github.com/hoshsadiq/m3ufilter/net"
	"io"
)

const (
	HeaderCuid       = "cuid"
	HeaderSearchName = "search-name"
	HeaderNumber     = "chno"
	HeaderId         = "tvg-id"
	HeaderGroup      = "group-title"
	HeaderName       = "tvg-name"
	HeaderShift      = "tvg-shift"
	HeaderLogo       = "tvg-logo"
	HeaderUri        = "uri"
)

var log = logger.Get()

type StreamData struct {
	ChNo  string `yaml:"chno"`
	Name  string
	Id    string
	Logo  string
	Group string
	Shift string
}

var allowedColumns = map[string]bool{
	HeaderCuid:       true,
	HeaderSearchName: true,
	HeaderNumber:     true,
	HeaderId:         true,
	HeaderGroup:      true,
	HeaderName:       true,
	HeaderShift:      true,
	HeaderLogo:       true,
	HeaderUri:        true,
}

var defaultColumns = []string{
	HeaderSearchName,
	HeaderNumber,
	HeaderId,
	HeaderGroup,
	HeaderName,
	HeaderShift,
	HeaderLogo,
}

func getCsvColumnValue(column string, columnMapping map[string]int, record []string, defaultValue string) string {
	colIndex, ok := columnMapping[column]
	if !ok {
		return defaultValue
	}

	return record[colIndex]
}

func GetCsvMapping(csvFile string) (map[string]*StreamData, error) {
	csvResp, err := net.GetUri(csvFile)
	if err != nil {
		return nil, err
	}
	defer helper.Close(csvResp.Body, fmt.Sprintf("for csv %s", csvFile))()

	r := csv.NewReader(csvResp.Body)
	r.FieldsPerRecord = 0
	r.Comment = '#'

	result := map[string]*StreamData{}
	columnMapping := make(map[string]int, len(allowedColumns))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Errorf("failed to read csv line in file %s due to err %s", csvFile, err)
			continue
		}

		if len(result) == 0 && len(columnMapping) == 0 {
			var hasHeader = true
			for i, value := range record {
				if _, ok := allowedColumns[value]; !ok {
					hasHeader = false
					break
				}
				columnMapping[value] = i
			}
			if hasHeader {
				if _, ok := columnMapping[HeaderName]; !ok {
					return nil, errors.New("couldn't figure out which column is the name")
				}
				continue
			}

			for i, column := range defaultColumns {
				columnMapping[column] = i
			}
		}

		result[record[columnMapping[HeaderSearchName]]] = &StreamData{
			ChNo:  getCsvColumnValue(HeaderNumber, columnMapping, record, ""),
			Name:  getCsvColumnValue(HeaderName, columnMapping, record, ""),
			Id:    getCsvColumnValue(HeaderId, columnMapping, record, ""),
			Logo:  getCsvColumnValue(HeaderLogo, columnMapping, record, ""),
			Group: getCsvColumnValue(HeaderGroup, columnMapping, record, ""),
			Shift: getCsvColumnValue(HeaderShift, columnMapping, record, ""),
		}
	}

	return result, nil
}
