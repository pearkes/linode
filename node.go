package linode

import (
	"fmt"
	"strconv"
	"strings"
)

type BasicNodeResponse struct {
	Data   map[string]int64 `json:"DATA"`
	Errors []LinodeError    `json:"ERRORARRAY"`
}

// Returns the slug for the region
func (r *BasicNodeResponse) StringID() string {
	if _, ok := r.Data["LinodeID"]; ok {
		return strconv.FormatInt(r.Data["LinodeID"], 10)
	} else {
		return ""
	}
}

func (r *BasicNodeResponse) Error() error {
	if len(r.Errors) > 0 {
		errs := []string{}
		for _, v := range r.Errors {
			errs = append(errs, v.ErrorMessage())
		}
		return fmt.Errorf("Errors from API: %s", strings.Join(errs, ", "))
	} else {
		return nil
	}
}

type NodesResponse struct {
	Responses []NodeResponse
}

type NodeResponse struct {
	Data   []map[string]interface{} `json:"DATA"`
	Errors []LinodeError            `json:"ERRORARRAY"`
	Action string                   `json:"ACTION"`
}

// Constructs a node from a batch retrieve response.
func (r *NodesResponse) Node() (Node, error) {
	// Check for an error in the response first and
	// return if we have one
	if r.Error() != nil {
		return Node{}, r.Error()
	}

	// We should only have 3 responses to the batch actions,
	// if not something is wrong.
	if len(r.Responses) != 3 {
		return Node{}, fmt.Errorf("Incorrect data returned from API: %#v", r.Responses)
	}

	// The node we're building
	node := Node{}

	// Iterate over the data we received and, depending on the data,
	// build our Node object.
	for _, resp := range r.Responses {
		// We should only have one node, we will access later with a 0 index
		if len(resp.Data) != 1 {
			return Node{}, fmt.Errorf("Incorrect data returned from API: %#v", resp.Data)
		}

		nodeData := resp.Data[0]

		// What type of response it is.
		switch resp.Action {
		case "linode.list":
			node.DataCenterID = strconv.FormatFloat(nodeData["DATACENTERID"].(float64), 'f', 0, 64)
			node.Dist = nodeData["DISTRIBUTIONVENDOR"].(string)
			node.ID = strconv.FormatFloat(nodeData["LINODEID"].(float64), 'f', 0, 64)
			node.Label = nodeData["LABEL"].(string)
			node.Status = strconv.FormatFloat(nodeData["STATUS"].(float64), 'f', 0, 64)
			node.TotalHD = strconv.FormatFloat(nodeData["TOTALHD"].(float64), 'f', 0, 64)
		case "linode.disk.list":
			node.DiskLabel = nodeData["LABEL"].(string)
			node.DiskType = nodeData["TYPE"].(string)
			node.DiskStatus = strconv.FormatFloat(nodeData["STATUS"].(float64), 'f', 0, 64)
			node.DiskSize = strconv.FormatFloat(nodeData["SIZE"].(float64), 'f', 0, 64)
		case "linode.ip.list":
			node.IPAddress = nodeData["IPADDRESS"].(string)
			node.DNSName = nodeData["RDNS_NAME"].(string)
		default:
			return Node{}, fmt.Errorf("Unexpected action returned from API: %s", resp.Action)
		}
	}

	// We build the node without errors
	return node, nil
}

func (r *NodesResponse) Error() error {
	errs := []string{}

	// Iterate over the responses we received and
	// build a list of errors, if they exist.
	for _, v := range r.Responses {
		if len(v.Errors) > 0 {
			for _, e := range v.Errors {
				errs = append(errs, e.ErrorMessage())
			}
		}
	}

	// If we have any errors, return them as an error type, if not,
	// return nil
	if len(errs) > 0 {
		return fmt.Errorf("Errors from API: %s", strings.Join(errs, ", "))
	} else {
		return nil
	}
}

// Node is used to represent a retrieved Node. All properties
// are set as strings.
type Node struct {
	DataCenterID string
	Dist         string
	ID           string
	Label        string
	Status       string
	TotalHD      string
	IPAddress    string
	DNSName      string
	DiskLabel    string
	DiskType     string
	DiskStatus   string
	DiskSize     string
}

// Returns the slug for the region
func (n *Node) StringStatus() string {
	statusMap := map[string]string{
		"-2": "boot failed",
		"-1": "being created",
		"0":  "brand new",
		"1":  "running",
		"2":  "powered off",
		"3":  "shutting down",
		"4":  "saved to disk",
	}

	return statusMap[n.Status]
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

	if node.Error() != nil {
		return "", node.Error()
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

	resp, err := checkResp(c.Http.Do(req))

	if err != nil {
		return fmt.Errorf("Error destroying node: %s", err)
	}

	node := new(BasicNodeResponse)

	err = decodeBody(resp, &node)

	if err != nil {
		return fmt.Errorf("Error parsing node response: %s", err)
	}

	if node.Error() != nil {
		return node.Error()
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
	resps := new([]NodeResponse)

	err = decodeBody(resp, resps)

	node.Responses = *resps

	if err != nil {
		return Node{}, fmt.Errorf("Error decoding node response: %s", err)
	}

	// The request was successful
	return node.Node()

}
