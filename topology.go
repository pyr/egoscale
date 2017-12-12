package egoscale

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

func (exo *Client) GetSecurityGroups(params url.Values) ([]*SecurityGroup, error) {
	resp, err := exo.Request("listSecurityGroups", params)
	if err != nil {
		return nil, err
	}

	var r ListSecurityGroupsResponse
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	}

	return r.SecurityGroups, nil
}

func (exo *Client) GetAllSecurityGroups() (map[string]SecurityGroup, error) {
	var sgs map[string]SecurityGroup
	params := url.Values{}
	securityGroups, err := exo.GetSecurityGroups(params)

	if err != nil {
		return nil, err
	}

	sgs = make(map[string]SecurityGroup)
	for _, sg := range securityGroups {
		sgs[sg.Name] = *sg
	}
	return sgs, nil
}

func (exo *Client) GetSecurityGroupId(name string) (string, error) {
	params := url.Values{}
	params.Set("name", name)
	securityGroups, err := exo.GetSecurityGroups(params)
	if err != nil {
		return "", err
	}

	for _, sg := range securityGroups {
		if sg.Name == name {
			return sg.Id, nil
		}
	}

	return "", nil
}

func (exo *Client) GetZones(params url.Values) ([]*Zone, error) {
	resp, err := exo.Request("listZones", params)

	if err != nil {
		return nil, err
	}

	var r ListZonesResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return r.Zones, nil
}

func (exo *Client) GetAllZones() (map[string]string, error) {
	var zones map[string]string
	params := url.Values{}
	response, err := exo.GetZones(params)
	if err != nil {
		return zones, err
	}

	zones = make(map[string]string)
	for _, zone := range response {
		zones[zone.Name] = zone.Id
	}
	return zones, nil
}

func (exo *Client) GetProfiles() (map[string]string, error) {

	var profiles map[string]string
	params := url.Values{}
	resp, err := exo.Request("listServiceOfferings", params)

	if err != nil {
		return nil, err
	}

	var r ListServiceOfferingsResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	profiles = make(map[string]string)
	for _, offering := range r.ServiceOfferings {
		profiles[strings.ToLower(offering.Name)] = offering.Id
	}

	return profiles, nil
}

func (exo *Client) GetKeypairs() ([]SSHKeyPair, error) {
	params := url.Values{}

	resp, err := exo.Request("listSSHKeyPairs", params)

	if err != nil {
		return nil, err
	}

	var r ListSSHKeyPairsResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	var keypairs = make([]SSHKeyPair, r.Count, r.Count)
	for i, keypair := range r.SSHKeyPairs {
		keypairs[i] = *keypair
	}
	return keypairs, nil
}

func (exo *Client) GetAffinityGroups() (map[string]string, error) {
	var affinitygroups map[string]string
	params := url.Values{}

	resp, err := exo.Request("listAffinityGroups", params)

	if err != nil {
		return nil, err
	}

	var r ListAffinityGroupsResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	affinitygroups = make(map[string]string)
	for _, affinitygroup := range r.AffinityGroups {
		affinitygroups[affinitygroup.Name] = affinitygroup.Id
	}
	return affinitygroups, nil

}

// GetImages list the available featured images and group them by name, then size.
func (exo *Client) GetImages() (map[string]map[int]string, error) {
	var images map[string]map[int]string
	images = make(map[string]map[int]string)

	params := url.Values{}
	params.Set("templatefilter", "featured")

	resp, err := exo.Request("listTemplates", params)

	if err != nil {
		return nil, err
	}

	var r ListTemplatesResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^Linux (?P<name>.+?) (?P<version>[0-9.]+)\b`)
	for _, template := range r.Templates {
		size := int(template.Size / (1024 * 1024 * 1024))

		fullname := strings.ToLower(template.Name)

		if _, present := images[fullname]; !present {
			images[fullname] = make(map[int]string)
		}
		images[fullname][size] = template.Id

		submatch := re.FindStringSubmatch(template.Name)
		if len(submatch) > 0 {
			name := strings.Replace(strings.ToLower(submatch[1]), " ", "-", -1)
			version := submatch[2]
			image := fmt.Sprintf("%s-%s", name, version)

			if _, present := images[image]; !present {
				images[image] = make(map[int]string)
			}
			images[image][size] = template.Id
		}
	}
	return images, nil
}

func (exo *Client) GetTopology() (*Topology, error) {

	zones, err := exo.GetAllZones()
	if err != nil {
		return nil, err
	}
	images, err := exo.GetImages()
	if err != nil {
		return nil, err
	}
	securityGroups, err := exo.GetAllSecurityGroups()
	if err != nil {
		return nil, err
	}
	groups := make(map[string]string)
	for k, v := range securityGroups {
		groups[k] = v.Id
	}

	keypairs, err := exo.GetKeypairs()
	if err != nil {
		return nil, err
	}

	/* Convert the ssh keypair to contain just the name */
	keynames := make([]string, len(keypairs))
	for i, k := range keypairs {
		keynames[i] = k.Name
	}

	affinitygroups, err := exo.GetAffinityGroups()
	if err != nil {
		return nil, err
	}
	profiles, err := exo.GetProfiles()
	if err != nil {
		return nil, err
	}

	topo := &Topology{
		Zones:          zones,
		Profiles:       profiles,
		Images:         images,
		Keypairs:       keynames,
		AffinityGroups: affinitygroups,
		SecurityGroups: groups,
	}

	return topo, nil
}
