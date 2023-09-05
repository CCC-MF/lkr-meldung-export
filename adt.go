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
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"
)

// XML Data Structures

type AdtGekid struct {
	SchemaVersion string       `xml:"Schema_Version,attr"`
	Absender      string       `xml:",innerxml"`
	MengePatient  MengePatient `xml:"Menge_Patient"`
	MengeMelder   MengeMelder  `xml:"Menge_Melder"`
}

func UnmarschallAdtGekid(content []byte) (*AdtGekid, error) {
	data := AdtGekid{}
	if err := xml.Unmarshal(content, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func ToMergedString(content []AdtGekid) (string, error) {
	var schema_version = ""
	var absender = ""

	var patienten = []string{}
	var melder = []string{}

	for _, entry := range content {
		patchedAbsenderString := entry.ToAbsenderString()
		reSoftwareID := regexp.MustCompile("SOFTWARE_ID=\"[^\"]+\"")
		patchedAbsenderString = reSoftwareID.ReplaceAllString(patchedAbsenderString, "SOFTWARE_ID=\"lkr-meldung-export\"")
		reInstallationID := regexp.MustCompile("Installations_ID=\"[^\"]+\"")
		patchedAbsenderString = reInstallationID.ReplaceAllString(patchedAbsenderString, "Installations_ID=\"1.0.0\"")

		if schema_version == "" {
			schema_version = entry.ToSchemaVersionString()
		}
		if absender == "" {
			absender = patchedAbsenderString
		}
		if schema_version != entry.ToSchemaVersionString() {
			return "", fmt.Errorf("Verschiedene Schema-Versionen")
		}

		if absender != patchedAbsenderString {
			return "", fmt.Errorf("Verschiedene Absender")
		}

		if !contains(patienten, entry.ToPatientString()) {
			patienten = append(patienten, entry.ToPatientString())
		}
		if !contains(melder, entry.ToMelderString()) {
			melder = append(melder, entry.ToMelderString())
		}
	}

	result := "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n"
	result += fmt.Sprintf("<ADT_GEKID xmlns=\"http://www.gekid.de/namespace\" Schema_Version=\"%s\">", schema_version)
	result += fmt.Sprintf("%s\n", absender)
	result += "    <Menge_Patient>"
	result += strings.Join(patienten, "")
	result += "</Menge_Patient>\n"
	result += "    <Menge_Melder>"
	result += strings.Join(melder, "")
	result += "</Menge_Melder>\n"
	result += "</ADT_GEKID>"

	re, _ := regexp.Compile("\n\\s*\n")
	result = re.ReplaceAllString(result, "\n")

	return result, nil
}

func contains(hey []string, needle string) bool {
	for _, current := range hey {
		if current == needle {
			return true
		}
	}
	return false
}

func (adtGekid *AdtGekid) ToSchemaVersionString() string {
	return adtGekid.SchemaVersion
}

func (adtGekid *AdtGekid) ToAbsenderString() string {
	x := strings.Split(adtGekid.Absender, "</Absender>")
	if len(x) > 0 {
		return x[0] + "</Absender>"
	}
	return ""
}

func (adtGekid *AdtGekid) ToPatientString() string {
	return adtGekid.MengePatient.Value
}

func (adtGekid *AdtGekid) ToMelderString() string {
	return adtGekid.MengeMelder.Value
}

type MengePatient struct {
	Value string `xml:",innerxml"`
}

type MengeMelder struct {
	Value string `xml:",innerxml"`
}
