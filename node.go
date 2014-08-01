package linode

import (
	"fmt"
	"strconv"
)

type BasicNodeResponse struct {
	Id map[string]map[string]int64 `json:"LinodeID"`
}

// Returns the slug for the region
func (r *BasicNodeResponse) StringId() string {
	if _, ok := r.Id["DATA"]["LINODEID"]; ok {
		return strconv.FormatInt(r.Id["DATA"]["LINODEID"], 10)
	} else {
		return ""
	}
}

type NodesResponse struct {
	Responses []NodeResponse `json:""`
}

type NodeResponse struct {
	Data []map[string]interface{} `json:"data"`
}

// Node is used to represent a retrieved Node. All properties
// are set as strings.
type Node struct {
	DataCenterId   int64  `json:"DATACENTERID"`
	Dist           string `json:"DISTRIBUTIONVENDOR"`
	Id             int64  `json:"LINODEID"`
	Label          string `json:"LABEL"`
	Status         int64  `json:"STATUS"`
	TotalDiskspace int64  `json:"TOTALHD"`
	IPAddress      string `json:"IPADDRESS"`
	DNSName        string `json:"RDNS_NAME"`
	Disk           Disk
}

// Represents a disk on a node
type Disk struct {
	Label  string `json:"LABEL"`
	Type   string `json:"TYPE"`
	Status string `json:"STATUS"`
	Size   string `json:"SIZE"`
}

// Returns the slug for the region
func (n *Node) StringStatus() string {
	statusMap := map[int]string{
		-2: "boot failed",
		-1: "being created",
		0:  "brand new",
		1:  "running",
		2:  "powered off",
		3:  "shutting down",
		4:  "saved to disk",
	}

	return statusMap[d.Status]
}

// Returns the slug for the region
func (n *Node) StringId() string {
	return strconv.FormatInt(n.Id, 10)
}

// CreateNode contains the request parameters to create a new
// node.
type CreateNode struct {
	DatacenterID string
	PlanID       string
	PaymentTerm  string
}

// CreateNode creates a node from the parameters specified and
// returns an error if it fails. If no error and an ID is returned,
// the Node was succesfully created. Sometimes, it can return an ID
// and an error – this is if the provisioning of the disk or IP
// failed later on.
func (c *Client) CreateNode(opts *CreateNode) (string, error) {
	// Make the request parameters
	params := make(map[string]string)

	params["DataCenterID"] = opts.DatacenterID
	params["PlanID"] = opts.PlanID

	if opts.PaymentTerm != "" {
		params["PaymentTerm"] = opts.PaymentTerm
	}

	req, err := c.NewRequest(params, "POST", "linode.create")

	if err != nil {
		return "", err
	}

	resp, err := checkResp(c.Http.Do(req))

	if err != nil {
		return "", fmt.Errorf("Error creating node: %s", err)
	}

	node := new(NodeResponse)

	err = decodeBody(resp, &node)

	if err != nil {
		return "", fmt.Errorf("Error parsing node response: %s", err)
	}

	// The request was successful
	return node.Node.StringId(), nil
}

// DestroyNode destroys a node by the ID specified and
// returns an error if it fails. If no error is returned,
// the Node was succesfully destroyed.
func (c *Client) DestroyNode(id string) error {
	req, err := c.NewRequest(map[string]string{}, "POST", "linode.destroy")

	if err != nil {
		return err
	}

	_, err = checkResp(c.Http.Do(req))

	if err != nil {
		return fmt.Errorf("Error destroying node: %s", err)
	}

	// The request was successful
	return nil
}

// RetrieveNode gets  a node by the ID specified and
// returns a Node and an error. An error will be returned for failed
// requests with a nil Node.
func (c *Client) RetrieveNode(id string) (Node, error) {
	actions := []string{"linode.list", "linode.ip.list", "linode.disk.list"}
	req, err := c.NewRequest(map[string]string{}, "GET", actions)

	if err != nil {
		return Node{}, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return Node{}, fmt.Errorf("Error destroying node: %s", err)
	}

	node := new(NodeResponse)

	err = decodeBody(resp, node)

	if err != nil {
		return Node{}, fmt.Errorf("Error decoding node response: %s", err)
	}

	// The request was successful
	return node.Node, nil
}
