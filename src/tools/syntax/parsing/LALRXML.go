package parsing

import "encoding/xml"

type xmlLALR struct {
	XMLName xml.Name `xml:"states"`
	States  []xmlState
}
type xmlSymbol struct {
}
type xmlProduction struct {
}
type xmlState struct {
	XMLName xml.Name `xml:"xmlState"`
	Gotos   []*xmlGoto
	Actions []*xmlAction
}
type xmlGoto struct {
	XMLName xml.Name `xml:"xmlGoto"`
	On      string   `xml:"on,attr"`
	State   int      `xml:"state,attr"`
}
type xmlAction struct {
	XMLName xml.Name `xml:"xmlAction"`
	On      string   `xml:"on,attr"`
	Do      string   `xml:"do,attr"`
}
