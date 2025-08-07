package regatta

import "encoding/xml"

type RunLog struct {
	XMLName     xml.Name  `xml:"runlog"`
	ContestID   string    `xml:"contest_id,attr"`
	StartTime   string    `xml:"start_time,attr"`
	CurrentTime string    `xml:"current_time,attr"`
	Name        string    `xml:"name"`
	Users       Users     `xml:"users"`
	Problems    Problems  `xml:"problems"`
	Languages   Languages `xml:"languages"`
	Headers     Headers   `xml:"userrunheaders"`
	Runs        Runs      `xml:"runs"`
}

type Users struct {
	XMLName xml.Name `xml:"users"`
	Users   []User   `xml:"user"`
}

type User struct {
	XMLName xml.Name `xml:"user"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

type Problems struct {
	XMLName  xml.Name     `xml:"problems"`
	Problems []XMLProblem `xml:"problem"`
}

type XMLProblem struct {
	XMLName   xml.Name `xml:"problem"`
	ID        string   `xml:"id,attr"`
	ShortName string   `xml:"short_name,attr"`
	LongName  string   `xml:"long_name,attr"`
}

type Languages struct {
	XMLName   xml.Name   `xml:"languages"`
	Languages []Language `xml:"language"`
}

type Language struct {
	XMLName   xml.Name `xml:"language"`
	ID        string   `xml:"id,attr"`
	ShortName string   `xml:"short_name,attr"`
	LongName  string   `xml:"long_name,attr"`
}

type Headers struct {
	XMLName xml.Name `xml:"userrunheaders"`
	Headers []Header `xml:"userrunheader"`
}

type Header struct {
	XMLName xml.Name `xml:"userrunheader"`
	UserID  string   `xml:"user_id,attr"`
}

type Runs struct {
	XMLName xml.Name `xml:"runs"`
	Runs    []Run    `xml:"run"`
}

type Run struct {
	XMLName    xml.Name `xml:"run"`
	RunID      string   `xml:"run_id,attr"`
	Time       int      `xml:"time,attr"`
	RunUUID    string   `xml:"run_uuid,attr"`
	Status     string   `xml:"status,attr"`
	UserID     int      `xml:"user_id,attr"`
	ProbID     int      `xml:"prob_id,attr"`
	LangID     string   `xml:"lang_id,attr"`
	Test       string   `xml:"test,attr"`
	Nsec       string   `xml:"nsec,attr"`
	PassedMode string   `xml:"passed_mode,attr"`
}
