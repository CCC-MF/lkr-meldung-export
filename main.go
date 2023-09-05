/*
 * MIT License
 *
 * Copyright (c) 2023 Comprehensive Cancer Center Mainfranken
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"database/sql"
	"fmt"
	"github.com/alecthomas/kong"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/term"
	"log"
	"os"
	"syscall"
)

var (
	cli     *CLI
	context *kong.Context
	db      *sql.DB
)

type CLI struct {
	User     string `short:"U" help:"Database username" required:""`
	Password string `short:"P" help:"Database password"`
	Host     string `short:"H" help:"Database host" default:"localhost"`
	Port     int    `help:"Database port" default:"3306"`
	Database string `short:"D" help:"Database name" default:"onkostar"`
}

func initCLI() {
	cli = &CLI{}
	context = kong.Parse(cli,
		kong.Name("lkr-export"),
		kong.Description("Tool to export data from Onkostar table 'lkr_meldung_export'"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
}

func main() {

	initCLI()

	if len(cli.Password) == 0 {
		fmt.Print("Passwort: ")
		if bytePw, err := term.ReadPassword(int(syscall.Stdin)); err == nil {
			cli.Password = string(bytePw)
		}
		println()
	}

	if dbx, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=skip-verify", cli.User, cli.Password, cli.Host, cli.Port, cli.Database)); err == nil {
		if err := dbx.Ping(); err == nil {
			db = dbx
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {
					log.Println("Cannot close database connection")
				}
			}(db)
		} else {
			log.Fatalf("Cannot connect to Database: %s\n", err.Error())
		}
	} else {
		log.Fatalf("Cannot connect to Database: %s\n", err.Error())
	}

	query := `SELECT lkr_export, xml_daten FROM lkr_meldung_export WHERE xml_daten IS NOT NULL ORDER BY lkr_export`

	var lkrExport uint32
	var xmlDaten sql.NullString

	if rows, err := db.Query(query); err == nil {
		var xmlDataCollection = map[string][]AdtGekid{}
		for rows.Next() {
			if err := rows.Scan(&lkrExport, &xmlDaten); err == nil {
				if value, err := xmlDaten.Value(); err == nil && value != nil {
					if len(xmlDataCollection[fmt.Sprintf("%d", lkrExport)]) == 0 {
						xmlDataCollection[fmt.Sprintf("%d", lkrExport)] = []AdtGekid{}
					}

					if data, err := UnmarschallAdtGekid([]byte(fmt.Sprint(xmlDaten))); err == nil {
						content := append(xmlDataCollection[fmt.Sprintf("%d", lkrExport)], *data)
						xmlDataCollection[fmt.Sprintf("%d", lkrExport)] = content
					} else {
						log.Printf("Problematischer Datensatz bei Export: %v", lkrExport)
					}
				} else {
					log.Printf("Kann XML-Daten nicht lesen")
				}
			} else {
				log.Printf("Datenbankeintrag nicht lesen")
			}
		}

		for exportId, content := range xmlDataCollection {
			if fileContent, err := ToMergedString(content); err == nil {
				filename := fmt.Sprintf("ADTGEKID_Export_%s.xml", exportId)
				if err := os.WriteFile(filename, []byte(fileContent), 0664); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Printf("Überspringe Meldungen beim Zusammenführen von Export-ID %s: %s", exportId, err)
			}
		}

	}

}
