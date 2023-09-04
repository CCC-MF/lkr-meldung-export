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
	"os"
	"strings"
	"testing"
)

func TestShouldReturnSchemaVersion(t *testing.T) {
	content, _ := os.ReadFile("test_data/test_data1.xml")
	data, err := UnmarschallAdtGekid(content)

	if err != nil {
		t.Logf("Could not unmarschal content'")
		t.Fail()
	}

	actual := data.SchemaVersion

	if actual != "2.1.1" {
		t.Logf("Result did not match expected value")
		t.Fail()
	}
}

func TestShouldReturnAbsenderPart(t *testing.T) {
	content, _ := os.ReadFile("test_data/test_data1.xml")
	data, err := UnmarschallAdtGekid(content)

	if err != nil {
		t.Logf("Could not unmarschal content'")
		t.Fail()
	}

	actual := data.ToAbsenderString()

	if !strings.HasPrefix(strings.TrimSpace(actual), "<Absender ") {
		t.Logf("Result did not start with '<Absender ' as expected")
		t.Fail()
	}

	if !strings.HasSuffix(strings.TrimSpace(actual), "</Absender>") {
		t.Logf("Result did not end with '</Absender>' as expected")
		t.Fail()
	}
}

func TestShouldReturnPatientPart(t *testing.T) {
	content, _ := os.ReadFile("test_data/test_data1.xml")
	data, err := UnmarschallAdtGekid(content)

	if err != nil {
		t.Logf("Could not unmarschal content'")
		t.Fail()
	}

	actual := data.ToPatientString()

	if !strings.HasPrefix(strings.TrimSpace(actual), "<Patient") {
		t.Logf("Result did not start with '<Patient>' as expected")
		t.Fail()
	}

	if !strings.HasSuffix(strings.TrimSpace(actual), "</Patient>") {
		t.Logf("Result did not end with '</Patient>' as expected")
		t.Fail()
	}
}

func TestShouldReturnMelderPart(t *testing.T) {
	content, _ := os.ReadFile("test_data/test_data1.xml")
	data, err := UnmarschallAdtGekid(content)

	if err != nil {
		t.Logf("Could not unmarschal content'")
		t.Fail()
	}

	actual := data.ToMelderString()

	if !strings.HasPrefix(strings.TrimSpace(actual), "<Melder") {
		t.Logf("Result did not start with '<Melder>' as expected")
		t.Fail()
	}

	if !strings.HasSuffix(strings.TrimSpace(actual), "</Melder>") {
		t.Logf("Result did not end with '</Melder>' as expected")
		t.Fail()
	}
}

func TestShouldMergeAdtFiles(t *testing.T) {
	content1, _ := os.ReadFile("test_data/test_data1.xml")
	content2, _ := os.ReadFile("test_data/test_data2.xml")
	content3, _ := os.ReadFile("test_data/test_data3.xml")

	data1, _ := UnmarschallAdtGekid(content1)
	data2, _ := UnmarschallAdtGekid(content2)

	expected := string(content3)
	actual, err := ToMergedString([]AdtGekid{*data1, *data2})

	if err != nil {
		t.Logf("Could not combine contents'")
		t.Fail()
	}

	if expected != actual {
		t.Logf("Result did not match expected value")
		t.Fail()
	}

}
