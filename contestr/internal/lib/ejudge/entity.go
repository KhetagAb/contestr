package ejudge

import "encoding/xml"

type Submission struct {
	XMLName    xml.Name `xml:"run"`
	RunID      int      `xml:"run_id,attr"`
	Time       int      `xml:"time,attr"`
	RunUUID    string   `xml:"run_uuid,attr"`
	Status     string   `xml:"status,attr"`
	UserID     int      `xml:"user_id,attr"`
	ProbID     int      `xml:"prob_id,attr"`
	LangID     int      `xml:"lang_id,attr"`
	Test       int      `xml:"test,attr"`
	NSec       int      `xml:"nsec,attr"`
	PassedMode string   `xml:"passed_mode,attr"`
}
