package utils

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

// JunitXML specifies the data structure of junit.xml
type JunitXML struct {
	XMLName    xml.Name    `xml:"testsuites"`
	Name       string      `xml:"name,attr"`
	Testsuites []Testsuite `xml:"testsuite"`
}

// Testsuite specifies the <testsuite> under <testsuites>
type Testsuite struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Name      string     `xml:"name,attr"`
	Time      string     `xml:"time,attr"`
	Testcases []Testcase `xml:"testcase"`
}

// Testcase specifies the <testcase> under <testsuite>
type Testcase struct {
	XMLName xml.Name `xml:"testcase"`
	Name    string   `xml:"name,attr"`
	Failure string   `xml:"failure"`
}

// ParseJunitText parse the xml file
func ParseJunitText(path string) (string, error) {
	xmlFile, _ := os.Open(path)
	defer xmlFile.Close()

	file, err := ioutil.ReadAll(xmlFile)
	var data JunitXML
	xml.Unmarshal(file, &data)

	var text string
	for _, testsuite := range data.Testsuites {
		for _, testcase := range testsuite.Testcases {
			if testcase.Failure != "" {
				text = text + "#### &#x1F198; " + testcase.Name + "\n```javascript\n" + testcase.Failure + "\n```\n"
			}
		}
	}
	return text, err
}
