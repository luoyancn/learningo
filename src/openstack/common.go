package openstack

type _domain struct {
	Name string `json:"name"`
	Id   string `json:"id,omitempty"`
}

type _project struct {
	Domain _domain `json:"domain"`
	Name   string  `json:"name"`
	Id     string  `json:"id,omitempty"`
}

type _user struct {
	Domain   _domain `json:"domain"`
	Id       string  `json:"id,omitempty"`
	Name     string  `json:"name"`
	Password string  `json:"password,omitempty"`
	Expires  string  `json:"password_expires_at,omitempty"`
}

type _scope struct {
	Project _project `json:"project"`
}

type _password struct {
	User _user `json:"user"`
}

type _identity struct {
	Methods  []string  `json:"methods"`
	Password _password `json:"password"`
}

type _auth struct {
	Identity _identity `json:"identity"`
	Scope    _scope    `json:"scope"`
}

type _credAuth struct {
	Auth _auth `json:"auth"`
}

type _role struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type _endpoint struct {
	Id        string `json:"id"`
	Interface string `json:"interface"`
	Region    string `json:"region"`
	RegionId  string `json:"region_id"`
	Url       string `json:"url"`
}

type _catalog struct {
	Endpoints []_endpoint `json:"endpoints"`
	Id        string      `json:"id"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
}

type _token struct {
	AuditIds []string   `json:"audit_ids"`
	Catalog  []_catalog `json:"catalog"`
	Expires  string     `json:"expires_at"`
	IsDomain bool       `json:"is_domain"`
	IssuedAt string     `json:"issued_at"`
	Methods  []string   `json:"methods"`
	Project  _project   `json:"project"`
	Role     []_role    `json:"roles"`
	User     _user      `json:"user"`
}

type _network struct {
	Uuid string `json:"uuid"`
}

type _block_mapping struct {
	SourceType          string `json:"source_type"`
	BootIndex           int    `json:"boot_index,string"`
	Uuid                string `json:"uuid"`
	VolumeSize          int    `json:"volume_size,string"`
	DestinationType     string `json:"destination_type"`
	DeleteOnTermination bool   `json:"delete_on_termination,omitempty"`
}

type _server struct {
	Name         string           `json:"name"`
	ImageRef     string           `json:"imageRef,omitempty"`
	FlavorRef    string           `json:"flavorRef"`
	Max          int              `json:"max_count"`
	Min          int              `json:"min_count"`
	Networks     []_network       `json:"networks"`
	BlockDevices []_block_mapping `json:"block_device_mapping_v2,omitempty"`
}

type _reqServer struct {
	Server _server `json:"server"`
}

type _secGroup struct {
	Name string `json:"name"`
}

type _link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type Image struct {
	Id    string  `json:"id,omitempty"`
	Links []_link `json:"links,omitempty"`
}

type _flavor struct {
	Id    string  `json:"id"`
	Links []_link `json:"links"`
}

type _volume struct {
	Id                  string `json:"id"`
	DeleteOnTermination bool   `json:"delete_on_termination"`
}

type _addinfo struct {
	Mac     string `json:"OS-EXT-IPS-MAC:mac_addr,omitempty"`
	Type    string `json:"OS-EXT-IPS:type,omitempty"`
	Addr    string `json:"addr,omitempty"`
	Version int    `json:"version,omitempty"`
}

type _vmserver struct {
	SecurityGroups []_secGroup `json:"security_groups"`
	DiskConfig     string      `json:"OS-DCF:diskConfig"`
	Id             string      `json:"id"`
	Links          []_link     `json:"links"`
	AdminPass      string      `json:"adminPass"`

	Volumes []_volume `json:"os-extended-volumes:volumes_attached,omitempty"`
	Status  string    `json:"status,omitempty"`
	Tags    []string  `json:"tags,omitempty"`

	Address     map[string][]_addinfo `json:"addresses,omitempty"`
	Metadata    map[string]string     `json:"metadata,omitempty"`
	Update      string                `json:"updated,omitempty"`
	HostId      string                `json:"hostId,omitempty"`
	KeyName     string                `json:"key_name,omitempty"`
	Image       interface{}           `json:"image,omitempty"`
	Flavor      _flavor               `json:"flavor,omitempty"`
	Locked      bool                  `json:"locked,omitempty"`
	Description string                `json:"description,omitempty"`
	HostStatus  string                `json:"host_status,omitempty"`
	UserId      string                `json:"user_id,omitempty"`
	Name        string                `json:"name,omitempty"`
	Created     string                `json:"created,omitempty"`
	ProjectId   string                `json:"tenant_id,omitempty"`
	AccessIPv4  string                `json:"accessIPv4,omitempty"`
	AccessIPv6  string                `json:"accessIPv6,omitempty"`
	Process     int                   `json:"process,omitempty"`
	ConfigDrive string                `json:"config_drive,omitempty"`

	AZ            string `json:"OS-EXT-AZ:availability_zone,omitempty"`
	UserData      string `json:"OS-EXT-SRV-ATTR:user_data,omitempty"`
	TaskState     string `json:"OS-EXT-STS:task_state,omitempty"`
	VmState       string `json:"OS-EXT-STS:vm_state,omitempty"`
	InstanceName  string `json:"OS-EXT-SRV-ATTR:instance_name,omitempty"`
	RootDevice    string `json:"OS-EXT-SRV-ATTR:root_device_name,omitempty"`
	LaunchedAt    string `json:"OS-SRV-USG:launched_at,omitempty"`
	KernelId      string `json:"OS-EXT-SRV-ATTR:kernel_id,omitempty"`
	LaunchedIndex int    `json:"OS-EXT-SRV-ATTR:launch_index,omitempty"`
	Host          string `json:"OS-EXT-SRV-ATTR:host,omitempty"`
	RamdiskId     string `json:"OS-EXT-SRV-ATTR:ramdisk_id,omitempty"`
	ReservationId string `json:"OS-EXT-SRV-ATTR:reservation_id,omitempty"`
	HostName      string `json:"OS-EXT-SRV-ATTR:hostname,omitempty"`
	Hypervisor    string `json:"OS-EXT-SRV-ATTR:hypervisor_hostname,omitempty"`
	PowerState    int    `json:"OS-EXT-STS:power_state,omitempty"`
	TerminatedAt  string `json:"OS-SRV-USG:terminated_at,omitempty"`
}

type _respServer struct {
	Server _vmserver `json:"server"`
}

type RespToken struct {
	Token _token `json:"token"`
}

type NULL struct{}
