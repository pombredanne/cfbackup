package bosh

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

type Task struct {
	Id          int    `json:"id"`
	State       string `json:"state"`
	Description string `json:"description"`
	Result      string `json:"result"`
}

const (
	ERROR      int = 0
	PROCESSING int = 1
	DONE       int = 2
)

var TASKRESULT map[string]int = map[string]int{"error": ERROR, "processing": PROCESSING, "done": DONE}

func retrieveTaskId(resp *http.Response) (taskId int, err error) {
	if resp.StatusCode != 302 {
		err = TaskRedirectStatusCodeError
		return
	}
	redirectUrls := resp.Header["Location"]
	if redirectUrls == nil || len(redirectUrls) < 1 {
		err = errors.New("Could not find redirect url for bosh tasks")
		return
	}
	regex := regexp.MustCompile(`^.*tasks/`)
	idString := regex.ReplaceAllString(redirectUrls[0], "")
	return strconv.Atoi(idString)
}

func retrieveTaskStatus(resp *http.Response) (task *Task, err error) {
	if resp.StatusCode != 200 {
		err = TaskStatusCodeError
		return
	}
	task = &Task{}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, task)
	if err != nil {
		return
	}
	return

}