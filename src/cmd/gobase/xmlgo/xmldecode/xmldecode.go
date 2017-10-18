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

type cputune struct {
	Shares int `xml:"shares,tag"`
}

type resource struct {
	Partition string `xml:"partition,tag"`
}

type entry struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type system struct {
	Entry []entry `xml:"entry,tag"`
}

type sysinfo struct {
	Type   string `xml:"type,attr"`
	System system `xml:"system,tag"`
}

type os_type struct {
	Arch    string `xml:"arch,attr"`
	Machine string `xml:"machine,attr"`
	Value   string `xml:",chardata"`
}

type boot struct {
	Dev string `xml:"dev,attr"`
}

type smbios struct {
	Mode string `xml:"mode,attr"`
}

type os_kvm struct {
	OsType os_type `xml:"type,tag"`
	Boot   []boot  `xml:"boot,tag"`
	Smbios smbios  `xml:"smbios,tag"`
}

type features struct {
	Acpi string `xml:"acpi,tag"`
	Apic string `xml:"apci,tag"`
}

type cpu_topo struct {
	Sockets string `xml:"sockets,attr"`
	Cores   string `xml:"cores,attr"`
	Threads string `xml:"threads,attr"`
}

type cpu struct {
	CpuTopo cpu_topo `xml:"topology,tag"`
}

type clock struct {
	OffSet string `xml:"offset,attr"`
}

type disk_host struct {
	Name string `xml:"name,attr"`
	Port string `xml:"port,attr"`
}

type disk_source struct {
	Protocol string      `xml:"protocol,attr"`
	Name     string      `xml:"name,attr"`
	Host     []disk_host `xml:"host,tag"`
}

type disk_secret struct {
	Type string `xml:"type,attr"`
	Uuid string `xml:"uuid,attr"`
}

type disk_auth struct {
	Username string      `xml:"username,attr"`
	Secret   disk_secret `xml:"secret,tag"`
}

type disk_driver struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Cache string `xml:"cache,attr"`
}

type disk_target struct {
	Dev string `xml:"dev,attr"`
	Bus string `xml:"bus,attr"`
}

type alias struct {
	Name string `xml:"name,attr"`
}

type address struct {
	Type       string `xml:"type,attr,omitempty"`
	Controller string `xml:"controller,attr,omitempty"`
	Bus        string `xml:"bus,attr,omitempty"`
	Unit       string `xml:"unit,attr,omitempty"`
	Domain     string `xml:"domain,attr,omitempty"`
	Slot       string `xml:"slot,attr,omitempty"`
	Function   string `xml:"function,attr,omitempty"`
	Port       string `xml:"port,attr,omitempty"`
}

type disk_entry struct {
	Type         string      `xml:"type,attr"`
	Devices      string      `xml:"device,attr"`
	Driver       disk_driver `xml:"driver,tag"`
	Auth         disk_auth   `xml:"auth,tag"`
	Source       disk_source `xml:"source,tag"`
	BackingStore string      `xml:"backingStore,tag"`
	Target       disk_target `xml:"target,tag"`
	Serial       string      `xml:"serial,tag,omitempty"`
	ReadOnly     string      `xml:"readonly,tag,omitempty"`
	Alias        alias       `xml:"alias,tag"`
	Address      address     `xml:"address,tag"`
}

type device_controller struct {
	Type    string  `xml:"type,attr"`
	Index   string  `xml:"index,attr"`
	Model   string  `xml:"model,attr,omitempty"`
	Alias   alias   `xml:"alias,tag"`
	Address address `xml:"address,tag,omitempty"`
}

type interface_mac struct {
	Address string `xml:"address,attr"`
}

type interface_brige struct {
	Bridge string `xml:"bridge,attr"`
}

type interface_target struct {
	Dev string `xml:"dev,attr"`
}

type interface_model struct {
	Type string `xml:"type,attr"`
}

type interface_driver struct {
	Name string `xml:"name,attr"`
}

type interface_device struct {
	Type    string           `xml:"type,attr"`
	Mac     interface_mac    `xml:"mac,tag"`
	Source  interface_brige  `xml:"source,tag"`
	Traget  interface_target `xml:"target,tag"`
	Model   interface_model  `xml:"model,tag"`
	Driver  interface_driver `xml:"driver,tag"`
	Alias   alias            `xml:"alias,tag"`
	Address address          `xml:"address,tag"`
}

type serial_source struct {
	Path string `xml:"path,attr"`
}

type serial_target struct {
	Type string `xml:"type,attr,omitempty"`
	Port string `xml:"port,attr"`
}

type serial_device struct {
	Type   string        `xml:"type,attr"`
	Source serial_source `xml:"source,tag"`
	Target serial_target `xml:"target,tag"`
	Alias  alias         `xml:"alias,tag"`
}

type console_source serial_source
type console_target serial_target

type console_device struct {
	Type   string         `xml:"type,attr"`
	Source console_source `xml:"source,tag"`
	Target console_target `xml:"target,tag"`
	Alias  alias          `xml:"alias,tag"`
}

type input_device struct {
	Type    string  `xml:"type,attr"`
	Bus     string  `xml:"bus,attr"`
	Address address `xml:"address,tag,omitempty"`
	Alias   alias   `xml:"alias,tag"`
}

type devices struct {
	Emulator   string              `xml:"emulator,tag"`
	Disk       []disk_entry        `xml:"disk,tag"`
	Controller []device_controller `xml:"controller,tag"`
	Interface  []interface_device  `xml:"interface,tag"`
	Serial     []serial_device     `xml:"serial,tag"`
	Console    []console_device    `xml:"console,tag"`
	Input      []input_device      `xml:"input,tag"`
}

type seclabel struct {
	Type    string `xml:"type,attr"`
	Model   string `xml:"model,attr"`
	Relabel string `xml:"relabel,attr,omitempty"`
	//Label      string `xml:label,tag,omitempty`
	//ImageLabel string `xml:imagelabel,tag,omitempty`
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
	XMLName xml.Name `xml:"domain"`
	Type    string   `xml:"type,attr"`
	Id      string   `xml:"id,attr"`

	Name      string       `xml:"name,tag"`
	Uuid      string       `xml:"uuid,tag"`
	Metadata  novametadata `xml:"metadata,tag"`
	Memory    memory       `xml:"memory,tag"`
	CurMemory curmemory    `xml:"currentMemory,tag"`
	Vcpu      vcpu         `xml:"vcpu,tag"`
	Cputune   cputune      `xml:"cputune,tag"`
	Resource  resource     `xml:"resource,tag"`
	Sysinfo   sysinfo      `xml:"sysinfo,tag"`
	Os        os_kvm       `xml:"os,tag"`
	Features  features     `xml:"features,tag"`
	Cpu       cpu          `xml:"cpu,tag"`
	Clock     clock        `xml:"clock,tag"`
	PowerOff  string       `xml:"on_poweroff,tag"`
	Restart   string       `xml:"on_reboot,tag"`
	Crash     string       `xml:"on_crash,tag"`

	Devices devices `xml:"devices,tag"`

	Seclabel []seclabel `xml:"seclabel,tag"`
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
