package nirvana

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestUnmarshall(t *testing.T) {
	rawData, err := ioutil.ReadFile("dump_testdata.json")
	if err != nil {
		t.Error(err)
	}
	var data NirvanaResponse

	err = json.Unmarshal(rawData, &data)
	if err != nil {
		t.Error(err)
	}

	for _, item := range data.Results {
		if item.User.ID == "139026" {
			return
		}
	}
	t.Error()
}
