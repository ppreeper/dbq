package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ppreeper/dbq/database"
	"github.com/ppreeper/pad"
)

func main() {
	// flags
	var dbase, stmt string
	var timer bool

	flag.StringVar(&dbase, "db", "", "database")
	flag.StringVar(&stmt, "q", "", "sql query")
	flag.BoolVar(&timer, "t", false, "sql timer")
	flag.Parse()

	if dbase == "" {
		fmt.Println("no database specified")
		os.Exit(1)
	}
	if stmt == "" {
		fmt.Println("no query specified")
		os.Exit(2)
	}

	// Config File
	userConfigDir, err := os.UserConfigDir()
	checkErr(err)
	var c Conf
	c.getConf(userConfigDir + "/dbq/config.yml")

	src := c.getDB(dbase)

	// connect to source database
	// open database connection
	sdb, err := database.OpenDatabase(src)
	checkErr(err)
	defer func() {
		if err := sdb.Close(); err != nil {
			checkErr(err)
		}
	}()
	checkErr(err)

	start := time.Now()
	colNames, dataSet := queryData(sdb, stmt)
	elapsed := time.Since(start)

	printData(&colNames, &dataSet)
	if timer {
		fmt.Printf("query time: %s\n", elapsed.String())
	}
}

func queryData(sdb *database.Database, stmt string) (colNames []string, dataSet []interface{}) {
	rows, err := sdb.DB.Queryx(stmt)
	fatalErr(err)
	defer rows.Close()

	colNames, err = rows.Columns()
	checkErr(err)
	cols := make([]interface{}, len(colNames))
	colPtrs := make([]interface{}, len(colNames))
	for i := 0; i < len(colNames); i++ {
		colPtrs[i] = &cols[i]
	}

	for rows.Next() {
		var rowMap = make(map[string]interface{})
		err = rows.Scan(colPtrs...)
		fatalErr(err)
		for i, col := range cols {
			rowMap[colNames[i]] = col
		}
		dataSet = append(dataSet, rowMap)
	}
	return
}

func printData(colNames *[]string, dataSet *[]interface{}) {
	colLens := make([]int, len(*colNames))
	for k, v := range *colNames {
		if len(v) > colLens[k] {
			colLens[k] = len(v)
		}
	}
	for _, v := range *dataSet {
		for k, c := range *colNames {
			vs := fmt.Sprintf("%v", v.(map[string]interface{})[c])
			if len(vs) > colLens[k] {
				colLens[k] = len(vs)
			}
		}
	}
	for k, v := range *colNames {
		fmt.Printf("%v", pad.LJustLen(v, colLens[k]))
		if k < len(*colNames)-1 {
			fmt.Printf(";")
		}
	}
	fmt.Println()
	for _, v := range *dataSet {
		for k, c := range *colNames {
			vs := fmt.Sprintf("%v", v.(map[string]interface{})[c])
			fmt.Printf("%v", pad.LJustLen(vs, colLens[k]))
			if k < len(*colNames)-1 {
				fmt.Printf(";")
			}
		}
		fmt.Println()
	}
}
