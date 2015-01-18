package main

import (
	//"flag"
	"github.com/ProfessorBeekums/PbStockResearcher/config"
	"github.com/ProfessorBeekums/PbStockResearcher/filings"
	"github.com/ProfessorBeekums/PbStockResearcher/log"
	"github.com/ProfessorBeekums/PbStockResearcher/parser"
	"github.com/ProfessorBeekums/PbStockResearcher/persist"
	//"github.com/ProfessorBeekums/PbStockResearcher/scraper"
	"github.com/ProfessorBeekums/PbStockResearcher/tmpStore"
	"strings"
)

//func init() {
//	flag.IntVar(&year, "year", 0, "The year to scrape")
//	flag.IntVar(&quarter, "quarter", 0, "The quarter to scrape")
//}

func main() {
	log.Println("Starting program")

	//flag.Parse()

	c := config.NewConfig("/home/beekums/Projects/stockResearch/config")

	log.Println("Loaded config: ", c)

	tmpStore.NewTempStore(c.TmpDir)

	//companyPersister := persist.NewMongoDbCompany(c.MongoHost, c.MongoDb)
	reportPersister := persist.NewMongoDbReportFiles(c.MongoHost, c.MongoDb)
	rawReportPersister :=
		persist.NewMongoDbFinancialReportsRaw(c.MongoHost, c.MongoDb)

	var batchLimit int64 = 20

	done := false
	var totalNumValid int64 = 0
	var totalNumInvalid int64 = 0

	for done == false {
		unparsedFiles := reportPersister.GetNextUnparsedFiles(batchLimit)

		if len(*unparsedFiles) == 0 {
			break
		}

		/**
		TODO loop through the entire batch and count the number of valid and unvalid entries
		*/

		var numValid int64 = 0
		var numInvalid int64 = 0
		for _, reportFile := range *unparsedFiles {
			filePath := reportFile.Filepath
			if !strings.Contains(filePath, "10-Q") &&
				!strings.Contains(filePath, "10-K") &&
				!strings.Contains(filePath, "10-K_A") &&
				!strings.Contains(filePath, "10-Q_A") {
				reportFile.Parsed = true

				reportPersister.InsertUpdateReportFile(&reportFile)
				continue
			}

			rawReport := &filings.FinancialReportRaw{CIK: reportFile.CIK, Year: reportFile.Year, Quarter: reportFile.Quarter}
			// TODO this is not optimal
			rawReport.RawFields = make(map[string]int64)

			frp := parser.NewFinancialReportParser(reportFile.Filepath,
				rawReport, rawReportPersister, &filings.BasicRawFieldNameList{})

			frp.Parse()

			fr := frp.GetFinancialReport()

			frValid := fr.IsValid()

			if frValid == nil {
				//log.Println("Parsed report for CIK <", fr.CIK, "> year <", fr.Year, "> quarter <", fr.Quarter, ">")
				numValid++
			} else {
				//log.Error("Invalid financial report <", reportFile, "> with error: ", frValid)

				// TODO temporary while I figure out what my parsing code is missing
				//break
				numInvalid++
			}

			reportFile.Parsed = true

			reportPersister.InsertUpdateReportFile(&reportFile)
		}

		totalNumInvalid += numInvalid
		totalNumValid += numValid

		log.Println("@@@@ Batch had <", numValid, "> valid and invalid: ", numInvalid)
		log.Println("@@@@ Total is <", totalNumValid, "> valid and invalid: ", totalNumInvalid)
	}
	log.Println("Ending program")
}
