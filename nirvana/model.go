package nirvana

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

// --- TAG TYPES

// Area value of tag.type for area
const Area = "1"

// Context value of tag.type for context
const Context = "2"

// Contact value of tag.type for contact
const Contact = "3"

// --- TASK types

// Action value of task.type for task item
const Action = "0"

// ActionProject value of task.type for project tasks
const ActionProject = "1"

// Reference value of task.type for reference item
const Reference = "2"

// ReferenceProject value of task.type for project references
const ReferenceProject = "3"

// --- TASK / PROJECT sates

// StateNext Next action list (1)
const StateNext = "1"

// StateWaiting Waiting for action (2)
const StateWaiting = "2"

// StateScheduled Scheduled action / project (3)
const StateScheduled = "3"

// StateSomeday Someday action / project (4)
const StateSomeday = "4"

// StateLater Later action / project (5)
const StateLater = "5"

// StateTrashed Trashed action / project (6)
const StateTrashed = "6"

// StateLogged Logged action / project (7)
const StateLogged = "7"

// StateDeleted Deleted action / project (8)
const StateDeleted = "8"

// StateRecurring Recurring action (9)
const StateRecurring = "9"

// StateActiveProject active project (11)
const StateActiveProject = "11"

// --- JSON objects

// NirvanaResponse ...
type NirvanaResponse struct {
	Request APIRequest  `json:"request"`
	Results []APIResult `json:"results"`
}

// APIRequest Nirvana API request
type APIRequest struct {
	RequestID string `json:"requestid"`
}

// APIResult Nirvana API result
type APIResult struct {
	Error APIError `json:"error"`
	Auth  APIAuth  `json:"auth"`
	User  User     `json:"user"`
	Tag   Tag      `json:"tag"`
	Task  Task     `json:"task"`
}

// APIError Nirvana API error
type APIError struct {
	Code    int64      `json:"code"`
	Message string     `json:"message"`
	Request APIRequest `json:"request"`
}

// APIAuth Nirvana API authentification information
type APIAuth struct {
	Token string `json:"token"`
}

// UserJSON JSON User structure
type UserJSON struct {
	ID string `json:"id"`
}

// User structure
type User struct {
	// ID ...
	ID string
}

// UnmarshalJSON exported
func (user *User) UnmarshalJSON(data []byte) error {
	var res UserJSON
	if err := json.Unmarshal(data, &res); err != nil {
		log.Printf("unable to unmarshall tag %v", data)
		return err
	}
	user.ID = res.ID
	return nil
}

// TagJSON JSON Tag structure
type TagJSON struct {
	Key     string `json:"key"`
	Type    string `json:"type"`
	Email   string `json:"email"`
	Color   string `json:"color"`
	Meta    string `json:"meta"`
	Deleted string `json:"deleted"`
}

// Tag structure
type Tag struct {
	// Key Tag name
	Key string
	// Type Tag type (0 = Tag, 1 = Area, 2 = Context, 3 = Contact)
	Type string
	// Email of contact (blank for others)
	Email string
	// Color Colour to use in displaying tag
	Color string
	// Meta Unknown
	Meta string
	// Deleted Time tag was deleted (0 for not deleted)
	Deleted time.Time
	// IsDeleted true if tag is deleted
	IsDeleted bool
}

// UnmarshalJSON exported
func (tag *Tag) UnmarshalJSON(data []byte) error {
	var res TagJSON
	if err := json.Unmarshal(data, &res); err != nil {
		log.Printf("unable to unmarshall tag %v", data)
		return err
	}
	tag.Key = res.Key
	tag.Type = res.Type
	tag.Email = res.Email
	tag.Color = res.Color
	tag.Meta = res.Meta
	tag.Deleted = strUnixTime(res.Deleted, 0, "tag.deleted")
	tag.IsDeleted = res.Deleted != "0"
	if res.Type == "0" {
		tag.Type = Context
		/*
			if !tag.IsDeleted {
				log.Printf("Tag type changed to context for %v", tag)
			}
		*/
	}
	return nil
}

// TaskJSON JSON Task structure
type TaskJSON struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Sequential string `json:"ps"`
	State      string `json:"state"`
	ParentID   string `json:"parentid"`
	Seq        string `json:"seq"`
	Seqt       string `json:"seqt"`
	Seqp       string `json:"seqp"`
	Name       string `json:"name"`
	Tags       string `json:"tags"`
	Etime      string `json:"etime"`
	Energy     string `json:"energy"`
	Waitingfor string `json:"waitingfor"`
	Startdate  string `json:"startdate"`
	Duedate    string `json:"duedate"`
	Recurring  string `json:"recurring"`
	Note       string `json:"note"`
	Completed  string `json:"completed"`
	Cancelled  string `json:"cancelled"`
	Deleted    string `json:"deleted"`
	Updated    string `json:"updated"`
}

// Task Task structure
type Task struct {
	// ID uuid, Unique identifier of task
	ID string
	// Type of task (0 = Task, 1 = Project)
	Type string
	// Sequential Is project sequential (0 if task is parallel, 1 if sequential)
	Sequential bool
	// State Action state (0 = Inbox, 1 = Next, 2 = Waiting, 3 = Scheduled, 4 = Someday, 5 = Later, 6 = Trashed, 7 = Logged, 8 = Deleted, 9 = Recurring, 11 = Active project )
	State string
	// ParentID Parent project UUID
	ParentID string
	// Seq Order of task / project in project list
	Seq int64
	// Seqt Time task is focused (0 for not focused)
	Seqt time.Time
	// Seqp Order of task in task list
	Seqp int64
	// Name of task / project
	Name string
	// Tags Comma delimited list of tags
	Tags string
	// Etime Minutes task requires (whoa nelly = 600)
	Etime int64
	// Energy Effort for task / project
	Energy int64
	// Waitingfor Name of contact task / project is waiting for
	Waitingfor string
	// Startdate Start date of task / project in YYYYMMDD
	Startdate time.Time
	// Duedate Due date of task / project in YYYYMMDD
	Duedate time.Time
	// Recurring JsonString "{"paused":false,"freq":"daily","interval":1,"nextdate":"20120807","hasduedate":0}"
	Recurring string
	// Note Contents of task / project note
	Note string
	// Completed Date of completion
	Completed time.Time
	// Cancelled ?
	Cancelled bool
	// Deleted Date of deleted
	Deleted time.Time
	// Updated Date of last update
	Updated time.Time
	// IsDeleted true if task is delete
	IsDeleted bool
	// IsCompleted true if task is done
	IsCompleted bool
}

// UnmarshalJSON exported
func (task *Task) UnmarshalJSON(data []byte) error {
	var res TaskJSON
	if err := json.Unmarshal(data, &res); err != nil {
		log.Printf("unable to unmarshall task %v", data)
		return err
	}

	task.ID = res.ID
	task.Type = res.Type
	task.Sequential = parseBool(res.Sequential)
	task.State = res.State
	task.ParentID = res.ParentID
	task.Seq = strInt64(res.Seq, 0, "task.seq")
	task.Seqt = strUnixTime(res.Seqt, 0, "task.seqt")
	task.Seqp = strInt64(res.Seqp, 0, "task.seqp")
	task.Name = res.Name
	task.Tags = res.Tags
	task.Etime = strInt64(res.Etime, 0, "task.etime")
	task.Energy = strInt64(res.Energy, 0, "task.energy")
	task.Waitingfor = res.Waitingfor
	if date, found := parseYYYMMDD(res.Startdate, "task.startdate"); found {
		task.Startdate = date
	}
	if date, found := parseYYYMMDD(res.Duedate, "task.duedate"); found {
		task.Duedate = date
	}
	task.Recurring = res.Recurring
	task.Note = res.Note
	task.Completed = strUnixTime(res.Completed, 0, "task.completed")
	task.IsCompleted = res.Completed != "0"
	task.Cancelled = parseBool(res.Cancelled)
	task.Deleted = strUnixTime(res.Deleted, 0, "task.deleted")
	task.Updated = strUnixTime(res.Updated, 0, "task.updated")
	task.IsDeleted = res.Deleted != "0"
	return nil
}

// ResultError returns response error or nil
func (p *NirvanaResponse) ResultError() error {
	if len(p.Results) > 0 && p.Results[0].Error.Code > 0 {
		return fmt.Errorf("Nirvana response error: %d %s", p.Results[0].Error.Code, p.Results[0].Error.Message)
	}
	return nil
}

// AuthToken returns authentication token or nil
func (p *NirvanaResponse) AuthToken() (string, bool) {
	if len(p.Results) > 0 && len(p.Results[0].Auth.Token) > 0 {
		return p.Results[0].Auth.Token, true
	}
	return "", false
}

func strInt64(val string, defaultVal int64, name string) int64 {
	if val == "" {
		return defaultVal
	}
	ret, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Printf("convertion error : field '%s' value: '%s' error: '%v'", name, val, err)
		ret = defaultVal
	}
	return ret
}

func strUnixTime(val string, defaultVal int64, name string) time.Time {
	return time.Unix(strInt64(val, defaultVal, name), 0)
}

const nirvanaDateLayout = "20060102"

var noDate = time.Unix(0, 0)

func parseYYYMMDD(val string, name string) (time.Time, bool) {
	if val == "" {
		return noDate, false
	}
	ret, err := time.Parse(nirvanaDateLayout, val)
	if err != nil {
		ret = noDate
		log.Printf("convertion error : field '%s' value: '%s' error: '%v'", name, val, err)
	}
	return ret, true
}

func parseBool(val string) bool {
	if val == "" || val == "0" {
		return false
	}
	return true
}
