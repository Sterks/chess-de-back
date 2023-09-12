package domain

type InfoStep struct {
	AllStepsInPart [][]string
	NumberParty    string
	Party          string
	ArrayMetaStep  []MetaStep
	AllStep        []OneStep
}
type MetaStep struct {
	Main              bool   `bson:"Main"`
	Steps             string `bson:"Steps"`
	Color             string `bson:"Color"`
	StepString        string `bson:"StepString"`
	CleanStepString   string
	Draw              string     `bson:"Draw"`
	Paragraph         int        `bson:"Paragraph"`
	NumberInParagraph int        `bson:"NumberInParagraph"`
	MetaBoth          []MetaBoth `bson:"MetaBoth"`
}

type MetaBoth struct {
	Main              bool
	Paragraph         int `bson:"Paragraph"`
	NumberInParagraph int
	StepString        string
	Both              []string  `bson:"Both"`
	BothString        string    `bson:"BothString"`
	OneStep           []OneStep `bson:"OneStep"`
	NumberStep        int
}

type OneStep struct {
	Step  string `bson:"Step"`
	Color string `bson:"Color"`
	FEN   string `bson:"Color"`
}
