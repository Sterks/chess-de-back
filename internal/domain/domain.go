package domain

import "github.com/notnil/chess"

type InfoStep struct {
	Name           string     `bson:"Name"`
	AllStepsInPart [][]string `bson:"AllStepsInPart"`
	NumberParty    string     `bson:"NumberParty"`
	Party          string     `bson:"Party"`
	ArrayMetaStep  []MetaStep `bson:"ArrayMetaStep"`
}

type MetaStep struct {
	Main              bool       `bson:"Main"`
	Steps             string     `bson:"Steps"`
	Color             string     `bson:"Color"`
	StepString        string     `bson:"StepString"`
	CleanStepString   string     `bson:"CleanStepString"`
	Draw              string     `bson:"Draw"`
	Paragraph         int        `bson:"Paragraph"`
	NumberInParagraph int        `bson:"NumberInParagraph"`
	MetaBoth          []MetaBoth `bson:"MetaBoth"`
}

type MetaBoth struct {
	Main              bool
	Paragraph         int       `bson:"Paragraph"`
	NumberInParagraph int       `bson:"NumberInParagraph"`
	StepString        string    `bson:"StepString"`
	Both              []string  `bson:"Both"`
	BothString        string    `bson:"BothString"`
	NumberStep        int       `bson:"NumberStep"`
	OneStep           []OneStep `bson:"OneStep"`
}

type OneStep struct {
	Step     string          `bson:"Step"`
	Color    string          `bson:"Color"`
	FEN      string          `bson:"FEN"`
	Position *chess.Position `bson:"Position"`
}
