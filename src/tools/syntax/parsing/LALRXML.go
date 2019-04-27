package parsing

import "encoding/xml"

type xmlLALR struct {
	XMLName     xml.Name `xml:"root"`
	Productions []xmlProduction
	States      []xmlState
}
type xmlProduction struct {
	XMLName xml.Name `xml:"production"`
	Head    string   `xml:"head,attr"`
	Nodes   string   `xml:"nodes,attr"`
	Script  string   `xml:"script,attr"`
	Len     int      `xml:"len,attr"`

	Id int `xml:"id,attr"`
}
type xmlState struct {
	XMLName xml.Name `xml:"state"`
	Gotos   []*xmlGoto
	Actions []*xmlAction

	Id int `xml:"id,attr"`
}
type xmlGoto struct {
	XMLName xml.Name `xml:"goto"`
	On      string   `xml:"on,attr"`
	State   int      `xml:"state,attr"`
}
type xmlAction struct {
	XMLName xml.Name `xml:"action"`
	On      string   `xml:"on,attr"`
	Do      string   `xml:"do,attr"`
}
