package linode

import (
	"testing"

	. "github.com/motain/gocheck"
)

func TestNode(t *testing.T) {
	TestingT(t)
}

func (s *S) Test_CreateNode(c *C) {
	testServer.Response(202, nil, nodeCreateExample)

	opts := CreateNode{
		DatacenterID: "1",
		PlanID:       "1",
	}

	id, err := s.client.CreateNode(&opts)

	req := testServer.WaitRequest()

	c.Assert(req.Form["api_action"], DeepEquals, []string{"batch"})
	c.Assert(req.Form["api_requestArray"], DeepEquals, []string{"[{\"DataCenterID\":\"1\",\"PlanID\":\"1\",\"api_action\":\"linode.create\"}]"})
	c.Assert(err, IsNil)
	c.Assert(id, Equals, "8098")
}

func (s *S) Test_CreateNode_Bad(c *C) {
	testServer.Response(200, nil, nodeExampleError)

	opts := CreateNode{
		PlanID: "1",
	}

	id, err := s.client.CreateNode(&opts)

	req := testServer.WaitRequest()

	c.Assert(req.Form["api_action"], DeepEquals, []string{"batch"})
	c.Assert(req.Form["api_requestArray"], DeepEquals, []string{"[{\"DataCenterID\":\"\",\"PlanID\":\"1\",\"api_action\":\"linode.create\"}]"})
	c.Assert(err.Error(), Equals, "Errors from API: 8: PlanID is invalid. Check linode.plans.list")
	c.Assert(id, Equals, "")
}

func (s *S) Test_RetrieveNode(c *C) {
	testServer.Response(200, nil, nodeExample)

	node, err := s.client.RetrieveNode("586892")

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(node.ID, Equals, "586892")
}

func (s *S) Test_DestroyNode(c *C) {
	testServer.Response(200, nil, nodeExampleDelete)

	opts := DestroyNode{
		LinodeID:   "1",
		SkipChecks: "true",
	}

	err := s.client.DestroyNode(&opts)

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
}

var nodeExampleError = `{
  "ERRORARRAY": [
    {
      "ERRORCODE": 8,
      "ERRORMESSAGE": "PlanID is invalid. Check linode.plans.list"
    }
  ],
  "DATA": {},
  "ACTION": "linode.create"
}`

var nodeExampleDelete = `{
   "ERRORARRAY":[],
   "ACTION":"linode.delete",
   "DATA":{
      "LinodeID":8204
   }
}`

var nodeExampleUpdate = `{
   "ERRORARRAY":[],
   "ACTION":"linode.update",
   "DATA":{
      "LinodeID":8098
   }
}`

var nodeCreateExample = `{
   "ERRORARRAY":[],
   "ACTION":"linode.create",
   "DATA":{
      "LinodeID":8098
   }
}`

var nodeExample = `
[
  {
    "ERRORARRAY": [],
    "DATA": [
      {
        "UPDATE_DT": "2009-07-18 12:53:043.0",
        "DISKID": 55320,
        "LABEL": "256M Swap Image",
        "TYPE": "swap",
        "LINODEID": 98,
        "ISREADONLY": 0,
        "STATUS": 1,
        "CREATE_DT": "2008-04-04 10:08:06.0",
        "SIZE": 256
      }
    ],
    "ACTION": "linode.disk.list"
  },
  {
    "ERRORARRAY": [],
    "DATA": [
      {
        "ALERT_CPU_ENABLED": 1,
        "ALERT_BWIN_ENABLED": 1,
        "ALERT_BWQUOTA_ENABLED": 1,
        "BACKUPWINDOW": 0,
        "ALERT_DISKIO_THRESHOLD": 1000,
        "DISTRIBUTIONVENDOR": "",
        "WATCHDOG": 1,
        "DATACENTERID": 6,
        "STATUS": 0,
        "ALERT_DISKIO_ENABLED": 1,
        "CREATE_DT": "2014-07-26 08:23:51.0",
        "TOTALHD": 24576,
        "ALERT_BWQUOTA_THRESHOLD": 80,
        "TOTALRAM": 1024,
        "ALERT_BWIN_THRESHOLD": 5,
        "LINODEID": 586892,
        "ALERT_BWOUT_THRESHOLD": 5,
        "ALERT_BWOUT_ENABLED": 1,
        "BACKUPSENABLED": 0,
        "ALERT_CPU_THRESHOLD": 90,
        "PLANID": 1,
        "BACKUPWEEKLYDAY": 0,
        "LABEL": "linode586892",
        "LPM_DISPLAYGROUP": "",
        "TOTALXFER": 2000
      }
    ],
    "ACTION": "linode.list"
  },
  {
    "ERRORARRAY": [],
    "DATA": [
      {
        "IPADDRESSID": 89244,
        "RDNS_NAME": "li326-141.members.linode.com",
        "LINODEID": 586892,
        "ISPUBLIC": 1,
        "IPADDRESS": "66.228.45.141"
      }
    ],
    "ACTION": "linode.ip.list"
  }
]
`
