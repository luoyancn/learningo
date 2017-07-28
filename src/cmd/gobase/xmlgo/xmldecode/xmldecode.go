package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type memory struct {
	Unit  string `xml:"unit,attr"`
	Value int    `xml:",chardata"`
}

type curmemory struct {
	Unit  string `xml:"unit,attr"`
	Value int    `xml:",chardata"`
}

type vcpu struct {
	Placement string `xml:"placement,attr"`
	Value     int    `xml:",chardata"`
}

type sysinfo struct {
	Type string `xml:"type,attr"`
}

type clock struct {
	OffSet string `xml:"offset,attr"`
}

type lable struct {
	Value interface{} `xml:",chardata"`
}

type imagelable struct {
	Value interface{} `xml:",chardata"`
}

type seclabel struct {
	Type    string `xml:"type,attr"`
	Model   string `xml:"model,attr"`
	Relabel string `xml:"relabel,attr,omitempty"`
	//label      lable      `xml:label,tag`
	//imageLabel imagelable `xml:imagelabel,tag`
}

type novaname struct {
	Value string `xml:",chardata"`
}

type novainstance struct {
	Xmlns string   `xml:"xmlns:nova,attr"`
	Name  novaname `xml:"nova:name,tag"`
}
type novametadata struct {
	Instance novainstance `xml:"nova:instance,tag"`
}

type domain struct {
	// 根节点一定要是XMLName，否则可能解析出现问题
	XMLName   xml.Name     `xml:"domain"`
	Name      string       `xml:"name,tag"`
	Uuid      string       `xml:"uuid,tag"`
	Metadata  novametadata `xml:"metadata,tag"`
	Memory    memory       `xml:"memory,tag"`
	CurMemory curmemory    `xml:"currentMemory,tag"`
	Vcpu      vcpu         `xml:"vcpu,tag"`
	Clock     clock        `xml:"clock,tag"`
	PowerOff  string       `xml:"on_poweroff,tag"`
	Restart   string       `xml:"on_reboot,tag"`
	Crash     string       `xml:"on_crash,tag"`
	Seclabel  []seclabel   `xml:"seclabel,tag"`
}

func main() {

	xml_file, err := os.Open("libvirt.xml")
	if nil != err {
		fmt.Fprintf(os.Stderr,
			"An error occured when openging the xml file: %v\n", err)
		os.Exit(2)
	}

	defer xml_file.Close()
	data, err := ioutil.ReadAll(xml_file)
	if nil != err {
		fmt.Fprintf(os.Stderr,
			"An error occured when reading the xml file: %v\n", err)
		os.Exit(2)
	}

	var libvirt domain
	err = xml.Unmarshal(data, &libvirt)
	if nil != err {
		fmt.Fprintf(os.Stderr,
			"An error occured when parse the xml file: %v\n", err)
		os.Exit(2)
	}
	fmt.Printf("%#v\n", libvirt)
	output, err := xml.Marshal(libvirt)
	fmt.Printf("%s\n", string(output))
}
