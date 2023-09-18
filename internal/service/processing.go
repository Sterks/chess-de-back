package service

import (
	"chess-backend/internal/domain"
	"chess-backend/internal/repository"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/notnil/chess"
	"github.com/tenkoh/go-docc"
)

// ErrorSteps(err, j, v.NumberStep, v.BothString, value.StepString, v.Main, v.Paragraph, v.NumberInParagraph)
func ErrorSteps(err error,
	step string,
	numberStep int,
	bothString string,
	stepString string,
	main bool,
	paragraph int,
	numberInParagraph int,
) error {
	er := fmt.Sprintf("error occurred at the step - %s, number - %d, both - %s, string steps - %s, main - %v, paragraph - %d, number in paragraph - %d", step, numberStep, bothString, stepString, main, paragraph, numberInParagraph)
	e := errors.New(er)
	return e
}

type FileProcessingService struct {
	NameBook string
	File     string
	Repo     repository.IRepositories
	Source   []string
	InfoStep domain.InfoStep
}

func NewFileProcessingService(repo repository.IRepositories) *FileProcessingService {
	return &FileProcessingService{
		Repo: repo,
	}
}

func (f *FileProcessingService) AddNamesBook(str string) {
	f.NameBook = str
}

func (f *FileProcessingService) GetSource(source []string) {
	f.Source = source
}

func (f *FileProcessingService) GetPage() {
	var info domain.InfoStep
	re := regexp.MustCompile(`(?ms)____\d\d\d\d\d`)
	re2 := regexp.MustCompile(`(?m)__имя.*`)
	for _, value := range f.Source {
		page := re.FindAllString(value, -1)
		page2 := re2.FindAllString(value, -1)
		if len(page) != 0 {
			info.Party = page[0]
		}
		if len(page2) != 0 {
			info.NumberParty = page2[0]
		}
	}
	info.NumberParty = strings.ReplaceAll(info.NumberParty, "\"", "")
	info.NumberParty = strings.ReplaceAll(info.NumberParty, "”", "")
	info.NumberParty = strings.Trim(info.NumberParty, " ")

	f.InfoStep = info
}

func (f *FileProcessingService) GetAllSteps() {
	var stringSteps []string
	var str [][]string
	re := regexp.MustCompile(`(?m)[(_]?\w.+?_[_)]`)

	for _, v := range f.Source {

		stringSteps = re.FindAllString(v, -1)
		if len(stringSteps) != 0 && stringSteps[0] != "____" {
			str = append(str, stringSteps)
		}
	}
	f.InfoStep.AllStepsInPart = str

	var meta domain.MetaStep
	var metaA []domain.MetaStep
	for key, value := range f.InfoStep.AllStepsInPart {
		if len(value) > 1 {
			for k, v := range value {
				meta.Paragraph = key + 1
				meta.NumberInParagraph = k + 1
				meta.StepString = v
				metaA = append(metaA, meta)
			}
		} else {
			meta.Paragraph = key + 1
			meta.NumberInParagraph = 0
			meta.StepString = value[0]
			metaA = append(metaA, meta)
		}
	}

	var metaArrr []domain.MetaStep
	var metaa domain.MetaStep

	for _, value := range metaA {
		metaa = value
		metaa.StepString = Replace(value.StepString)
		metaArrr = append(metaArrr, metaa)
	}

	f.InfoStep.ArrayMetaStep = metaArrr
}

func (f *FileProcessingService) ReadProcessing(file string) error {
	f.File = file

	r, err := docc.NewReader("upload/" + file)
	if err != nil {
		return err
	}

	ps, err := r.ReadAll()
	if err != nil {
		return err
	}
	f.GetSource(ps)
	f.GetPage()
	f.GetAllSteps()
	f.ParseMainAndSlaveSteps()
	err = f.GetInfoBoth()
	if err != nil {
		return err
	}
	r.Close()
	return nil
}

func (f *FileProcessingService) ParseMainAndSlaveSteps() {
	var mm domain.MetaStep
	var arrayTypeStep []domain.MetaStep

	for _, value := range f.InfoStep.ArrayMetaStep {

		jf := strings.Contains(value.StepString, "(")
		if !jf {
			mm = value
			mm.Main = true
			arrayTypeStep = append(arrayTypeStep, mm)
		} else {
			mm = value
			mm.Main = false
			arrayTypeStep = append(arrayTypeStep, mm)
		}
	}
	f.InfoStep.ArrayMetaStep = arrayTypeStep
}

func (f *FileProcessingService) GetInfoBoth() error {
	var metaA []domain.MetaStep
	var meta domain.MetaStep
	var re = regexp.MustCompile(`(?m)[0-9]{1,2}\.`)
	game := chess.NewGame()
	for _, value := range f.InfoStep.ArrayMetaStep {

		value.StepString = strings.Replace(value.StepString, "(_", "", -1)
		value.StepString = strings.Replace(value.StepString, "_)", "", -1)
		value.StepString = strings.Replace(value.StepString, ".__", "", -1)
		value.StepString = strings.Replace(value.StepString, "__", "", -1)
		meta = value

		arIntIndexString := re.FindAllString(value.StepString, -1)
		var arIntIndex []int
		for _, value := range arIntIndexString {
			value = strings.Replace(value, ".", "", -1)
			i, _ := strconv.Atoi(value)
			arIntIndex = append(arIntIndex, i)
		}

		arrayStep := re.Split(value.StepString, -1)

		var metaBoth domain.MetaBoth
		var arMetaBothTrim []domain.MetaBoth
		for k, v := range arrayStep {
			v1 := strings.Trim(v, " ")
			if len(v) > 0 {
				metaBoth.BothString = v1
				metaBoth.NumberStep = arIntIndex[k-1]
				arMetaBothTrim = append(arMetaBothTrim, metaBoth)
			}

			for k, v := range arMetaBothTrim {
				v.NumberStep = k
			}

			meta.MetaBoth = arMetaBothTrim
		}

		metaA = append(metaA, meta)
	}
	f.InfoStep.ArrayMetaStep = metaA

	var oneStep domain.OneStep
	var arOneStep []domain.OneStep
	var arMetaStep []domain.MetaStep
	var mStep domain.MetaStep
	var mMBoth domain.MetaBoth

	for _, value := range f.InfoStep.ArrayMetaStep {
		var arMetaBoth []domain.MetaBoth
		mStep = value
		for _, v := range value.MetaBoth {
			mMBoth = v
			arSteps := strings.Split(v.BothString, " ")
			var aOneStep []domain.OneStep
			for i, j := range arSteps {
				if i == 0 {
					oneStep.Step = j
					oneStep.Color = "w"
					if oneStep.Step == ".." {
						continue
					} else {
						err := game.MoveStr(oneStep.Step)
						if err != nil {
							err = ErrorSteps(err, j, v.NumberStep, v.BothString, value.StepString, v.Main, v.Paragraph, v.NumberInParagraph)
							return err
						}
					}
					fen := game.FEN()
					position := game.Position()
					oneStep.FEN = fen
					oneStep.Position = position
					arOneStep = append(arOneStep, oneStep)
					aOneStep = append(aOneStep, oneStep)
				} else {
					oneStep.Step = j
					oneStep.Color = "b"
					err := game.MoveStr(oneStep.Step)
					if err != nil {
						err = ErrorSteps(err, j, v.NumberStep, v.BothString, value.StepString, v.Main, v.Paragraph, v.NumberInParagraph)
						return err
					}
					fen := game.FEN()
					position := game.Position()
					oneStep.FEN = fen
					oneStep.Position = position
					arOneStep = append(arOneStep, oneStep)
					aOneStep = append(aOneStep, oneStep)
				}
				mMBoth.OneStep = aOneStep
			}
			arMetaBoth = append(arMetaBoth, mMBoth)
			mStep.MetaBoth = arMetaBoth
		}
		arMetaStep = append(arMetaStep, mStep)
		f.InfoStep.ArrayMetaStep = arMetaStep

		var book domain.InfoStep
		book = f.InfoStep
		book.Name = f.NameBook
		err := f.Repo.StepsSave(book)
		if err != nil {
			return err
		}
	}
	return nil
}

// func (f *FileProcessingService) CheckSteps() error {
// 	game := chess.NewGame()
// 	n := true
// 	for _, value := range f.InfoStep.ArrayMetaStep {

// 		if !n {
// 			break
// 		}
// 		if value.Main {
// 			for _, v := range value.MetaBoth {
// 				for _, j := range v.OneStep {
// 					if err := game.MoveStr(j.Step); err != nil {
// 						if j.Step == ".." {
// 							continue
// 						}
// 						log.Println(err, value)
// 						n = false
// 						return err
// 					}
// 					fmt.Println(game.Position().Board().Draw())
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

func Replace(st string) string {
	var lt = map[string]string{
		"Л":   "R",
		"К":   "N",
		"с":   "c",
		"C":   "B",
		"С":   "B",
		"Ф":   "Q",
		":":   "x",
		" : ": "x",
		"е":   "e",
		// "...": "",
		"0-0": "O-O",
		// "0": "O",
		"—": "-",
	}
	st = strings.Replace(st, "Кр", "K", -1)
	st = strings.Replace(st, "С : ", "Bx", -1)
	st = strings.Replace(st, ".__", "", -1)
	st = strings.Replace(st, ". __", "", -1)
	st = strings.Replace(st, "Ф : ", "Qx", -1)
	st = strings.Replace(st, " : ", "x", -1)
	for key, value := range lt {
		st = strings.Replace(st, key, value, -1)
	}
	st = strings.Replace(st, "0-0", "O-O", -1)
	return st
}
