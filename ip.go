package egoscale

import (
	"encoding/json"
)

// IpAddress represents an IP Address
type IpAddress struct {
	Id                        string         `json:"id"`
	Account                   string         `json:"account,omitempty"`
	AllocatedAt               string         `json:"allocated,omitempty"`
	AssociatedNetworkId       string         `json:"associatednetworkid,omitempty"`
	AssociatedNetworkName     string         `json:"associatednetworkname,omitempty"`
	DomainId                  string         `json:"domainid,omitempty"`
	DomainName                string         `json:"domainname,omitempty"`
	ForDisplay                bool           `json:"fordisplay,omitempty"`
	ForVirtualNetwork         bool           `json:"forvirtualnetwork,omitempty"`
	IpAddress                 string         `json:"ipaddress"`
	IsElastic                 bool           `json:"iselastic,omitempty"`
	IsPortable                bool           `json:"isportable,omitempty"`
	IsSourceNat               bool           `json:"issourcenat,omitempty"`
	IsSystem                  bool           `json:"issystem,omitempty"`
	NetworkId                 string         `json:"networkid,omitempty"`
	PhysicalNetworkId         string         `json:"physicalnetworkid,omitempty"`
	Project                   string         `json:"project,omitempty"`
	ProjectId                 string         `json:"projectid,omitempty"`
	Purpose                   string         `json:"purpose,omitempty"`
	State                     string         `json:"state,omitempty"`
	VirtualMachineDisplayName string         `json:"virtualmachinedisplayname,omitempty"`
	VirtualMachineId          string         `json:"virtualmachineid,omitempty"`
	VirtualMachineName        string         `json:"virtualmachineName,omitempty"`
	VlanId                    string         `json:"vlanid,omitempty"`
	VlanName                  string         `json:"vlanname,omitempty"`
	VmIpAddress               string         `json:"vmipaddress,omitempty"`
	VpcId                     string         `json:"vpcid,omitempty"`
	ZoneId                    string         `json:"zoneid,omitempty"`
	ZoneName                  string         `json:"zonename,omitempty"`
	Tags                      []*ResourceTag `json:"tags,omitempty"`
	JobId                     string         `json:"jobid,omitempty"`
	JobStatus                 JobStatusType  `json:"jobstatus,omitempty"`
}

// AssociateIpProfileRequest represents the IP creation
type AssociateIpAddressRequest struct {
	Account    string `json:"account,omitempty"`
	DomainId   string `json:"domainid,omitempty"`
	ForDisplay bool   `json:"fordisplay,omitempty"`
	IsPortable bool   `json:"isportable,omitempty"`
	NetworkdId string `json:"networkid,omitempty"`
	ProjectId  string `json:"projectid,omitempty"`
	RegionId   string `json:"regionid,omitempty"`
	VpcId      string `json:"vpcid,omitempty"`
	ZoneId     string `json:"zoneid,omitempty"`
}

// Command returns the CloudStack API command
func (*AssociateIpAddressRequest) Command() string {
	return "associateIpAddressRequest"
}

// AssociateIpAddressResponse represents the response to the creation of an IpAddress
type AssociateIpAddressResponse struct {
	IpAddress *IpAddress `json:"ipaddress"`
}

// DisassociateIpAddressRequest represents the IP deletion
type DisassociateIpAddressRequest struct {
	Id string `json:"id"`
}

// Command returns the CloudStack API command
func (*DisassociateIpAddressRequest) Command() string {
	return "disassociateIpAddressRequest"
}

// ListPublicIpAddressesRequest represents a search for public IP addresses
type ListPublicIpAddressesRequest struct {
	Account            string         `json:"account,omitempty"`
	AllocatedOnly      bool           `json:"allocatedonly,omitempty"`
	AllocatedNetworkId string         `json:"allocatednetworkid,omitempty"`
	DomainId           string         `json:"domainid,omitempty"`
	ForDisplay         bool           `json:"fordisplay,omitempty"`
	ForLoadBalancing   bool           `json:"forloadbalancing,omitempty"`
	ForVirtualNetwork  string         `json:"forvirtualnetwork,omitempty"`
	Id                 string         `json:"id,omitempty"`
	IpAddress          string         `json:"ipaddress,omitempty"`
	IsElastic          bool           `json:"iselastic,omitempty"`
	IsRecursive        bool           `json:"isrecursive,omitempty"`
	IsSourceNat        bool           `json:"issourcenat,omitempty"`
	IsStaticNat        bool           `json:"isstaticnat,omitempty"`
	Keyword            string         `json:"keyword,omitempty"`
	ListAll            bool           `json:"listall,omiempty"`
	Page               int            `json:"page,omitempty"`
	PageSize           int            `json:"pagesize,omitempty"`
	PhysicalNetworkId  string         `json:"physicalnetworkid,omitempty"`
	ProjectId          string         `json:"projectid,omitempty"`
	Tags               []*ResourceTag `json:"tags,omitempty"`
	VlanId             string         `json:"vlanid,omitempty"`
	VpcId              string         `json:"vpcid,omitempty"`
	ZoneId             string         `json:"zoneid,omitempty"`
}

// Command returns the CloudStack API command
func (*ListPublicIpAddressesRequest) Command() string {
	return "listPublicIpAddresses"
}

// ListPublicIpAddressesResponse represents a list of public IP addresses
type ListPublicIpAddressesResponse struct {
	Count           int          `json:"count"`
	PublicIpAddress []*IpAddress `json:"publicipaddress"`
}

// AssociateIpAddress acquires and associates a public IP to a given zone
//
// https://doc.internal.exoscale.ch/others/cs/user/associateIpAddress.html
func (exo *Client) AssociateIpAddress(req AssociateIpAddressRequest, async AsyncInfo) (*IpAddress, error) {
	params, err := prepareValues(req)
	if err != nil {
		return nil, err
	}
	resp, err := exo.AsyncRequest("associateIpAddress", *params, async)
	if err != nil {
		return nil, err
	}

	var r AssociateIpAddressResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return r.IpAddress, nil
}

// DisassociateIpAddress disassociates a public IP from the account
//
// https://doc.internal.exoscale.ch/others/cs/user/disassociateIpAddress.html
func (exo *Client) DisassociateIpAddress(req *DisassociateIpAddressRequest, async AsyncInfo) error {
	return exo.BooleanAsyncRequest(req, async)
}

// ListPublicIpAddresses lists the public ip addresses
func (exo *Client) ListPublicIpAddresses(req *ListPublicIpAddressesRequest) ([]*IpAddress, error) {
	var r ListPublicIpAddressesResponse
	err := exo.Request(req, &r)
	if err != nil {
		return nil, err
	}

	return r.PublicIpAddress, nil
}

// XXX UpdateIpAddress
//
// https://doc.internal.exoscale.ch/others/cs/user/updateIpAddress.html
