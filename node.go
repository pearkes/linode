package linode

import (
	"fmt"
	"strconv"
)

type BasicNodeResponse struct {
	Data map[string]int64 `json:"DATA"`
}

// Returns the slug for the region
func (r *BasicNodeResponse) StringID() string {
	if _, ok := r.Data["LinodeID"]; ok {
		return strconv.FormatInt(r.Data["LinodeID"], 10)
	} else {
		return ""
	}
}

type NodesResponse struct {
	Responses []NodeResponse
}

// Returns a node from the Nodes Response
func (r *NodesResponse) Node(id string) (Node, error) {
	if len(r.Responses) > 3 {
		return Node{}, fmt.Errorf("Incorrect data returned from API: %#v", r.Responses)
	}

	return Node{}, nil
}

type NodeResponse struct {
	Data []map[string]interface{} `json:"data"`
}

// Node is used to represent a retrieved Node. All properties
// are set as strings.
type Node struct {
	DataCenterID   int64  `json:"DATACENTERID"`
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

	return statusMap[int(n.Status)]
}

// Returns the slug for the region
func (n *Node) StringID() string {
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
	create := make(map[string]string)

	create["DataCenterID"] = opts.DatacenterID
	create["PlanID"] = opts.PlanID

	if opts.PaymentTerm != "" {
		create["PaymentTerm"] = opts.PaymentTerm
	}

	create["api_action"] = "linode.create"

	req, err := c.NewRequest("POST", []map[string]string{create})

	if err != nil {
		return "", err
	}

	resp, err := checkResp(c.Http.Do(req))

	if err != nil {
		return "", fmt.Errorf("Error creating node: %s", err)
	}

	node := new(BasicNodeResponse)

	err = decodeBody(resp, &node)

	if err != nil {
		return "", fmt.Errorf("Error parsing node response: %s", err)
	}

	// The request was successful
	return node.StringID(), nil
}

// DestroyNode contains the request parameters to destroy a
// node.
type DestroyNode struct {
	LinodeID   string
	SkipChecks string // bool
}

// DestroyNode destroys a node by the ID specified and
// returns an error if it fails. If no error is returned,
// the Node was succesfully destroyed.
func (c *Client) DestroyNode(opts *DestroyNode) error {
	destroy := make(map[string]string)

	destroy["LinodeID"] = opts.LinodeID

	if opts.SkipChecks != "" {
		destroy["skipChecks"] = opts.SkipChecks
	} else {
		destroy["skipChecks"] = "false"
	}

	destroy["api_action"] = "linode.delete"

	req, err := c.NewRequest("POST", []map[string]string{destroy})

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
	as := []string{"linode.list", "linode.ip.list", "linode.disk.list"}
	actions := []map[string]string{}

	for _, v := range as {
		a := make(map[string]string)
		a["LinodeID"] = id
		a["api_action"] = v
		actions = append(actions, a)
	}

	req, err := c.NewRequest("GET", actions)

	if err != nil {
		return Node{}, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return Node{}, fmt.Errorf("Error retrieving node: %s", err)
	}

	node := new(NodesResponse)

	err = decodeBody(resp, node)

	if err != nil {
		return Node{}, fmt.Errorf("Error decoding node response: %s", err)
	}

	// The request was successful
	return node.Node(id)

}
