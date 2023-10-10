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
	game, err := f.CheckMainSteps()
	if err != nil {
		return err
	}
	err = f.CheckSecondSteps2(game)
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
	// game := chess.NewGame()
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
					// if oneStep.Step == ".." {
					// 	continue
					// } else {
					// err := game.MoveStr(oneStep.Step)
					// 	if err != nil {
					// 		err = ErrorSteps(err, j, v.NumberStep, v.BothString, value.StepString, v.Main, v.Paragraph, v.NumberInParagraph)
					// 		return err
					// 	}
					// }
					// fen := game.FEN()
					// position := game.Position()
					// oneStep.FEN = fen
					// oneStep.Position = position
					arOneStep = append(arOneStep, oneStep)
					aOneStep = append(aOneStep, oneStep)
					// }
				} else {
					oneStep.Step = j
					oneStep.Color = "b"
					// err := game.MoveStr(oneStep.Step)
					// if err != nil {
					// err = ErrorSteps(err, j, v.NumberStep, v.BothString, value.StepString, v.Main, v.Paragraph, v.NumberInParagraph)
					// return err
					// }
					// fen := game.FEN()
					// position := game.Position()
					// oneStep.FEN = fen
					// oneStep.Position = position
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

func (f *FileProcessingService) CheckMainSteps() (*chess.Game, error) {
	game := chess.NewGame()
	var arMetaStep []domain.MetaStep
	var mStep domain.MetaStep
	for _, value := range f.InfoStep.ArrayMetaStep {
		mStep = value
		var arrayMetaBoth []domain.MetaBoth
		var metaBoth domain.MetaBoth
		if value.Main {
			for _, v := range value.MetaBoth {
				v.Main = value.Main
				v.Paragraph = value.Paragraph
				v.NumberInParagraph = value.NumberInParagraph
				var oneSteps []domain.OneStep
				var oneStep domain.OneStep
				// var arrOneSteps []domain.OneStep
				metaBoth = v
				for _, j := range v.OneStep {
					if j.Color == "w" {
						oneStep = j
						if j.Step == ".." {
							continue
						}
						err := game.MoveStr(j.Step)
						if err != nil {
							return nil, errors.New(fmt.Sprintf("%s - %v", value.StepString, err))
						}
						oneStep.FEN = game.FEN()
						fmt.Println(game.Position().Board().Draw())
						fmt.Println(j.Step, " - ", j.Color)
					} else if j.Color == "b" {
						oneStep = j
						err := game.MoveStr(j.Step)
						if err != nil {
							return nil, errors.New(fmt.Sprintf("%s - %v", value.StepString, err))
						}
						oneStep.FEN = game.FEN()
						fmt.Println(game.Position().Board().Draw())
						fmt.Println(j.Step, " - ", j.Color)
					}
					oneSteps = append(oneSteps, oneStep)
				}
				metaBoth.OneStep = oneSteps
				// arrOneSteps = append(arrOneSteps, oneSteps)
				arrayMetaBoth = append(arrayMetaBoth, metaBoth)

			}
			// arrayMetaBoth = append(arrayMetaBoth, metaBoth)
			mStep.MetaBoth = arrayMetaBoth
		}
		arMetaStep = append(arMetaStep, mStep)
		f.InfoStep.ArrayMetaStep = arMetaStep
		f.getAllStep()
	}
	return game, nil
}

func (f *FileProcessingService) getAllStep() {
	var arAllStep []string
	for _, value := range f.InfoStep.ArrayMetaStep {
		if value.Main == true {
			for _, v := range value.MetaBoth {
				for _, va := range v.OneStep {
					arAllStep = append(arAllStep, va.Step)
				}
			}
		}
	}
	for _, value := range arAllStep {
		fmt.Printf("%s ", value)
	}
}

func (f *FileProcessingService) CheckSecondSteps(game *chess.Game) error {
	var control bool
	control = false
	for _, value := range f.InfoStep.ArrayMetaStep {
		if value.Main == false {
			for _, v := range value.MetaBoth {
				l := len(v.OneStep)
				if l <= 1 {
					for _, va := range f.InfoStep.ArrayMetaStep {
						if va.Main == true {
							for _, vi := range va.MetaBoth {
								if vi.NumberStep == v.NumberStep-1 {
									j := false
									if control == false {
										for _, vt := range v.OneStep {
											if j == false {
												fen, err := chess.FEN(vi.OneStep[1].FEN)
												if err != nil {
													return err
												}
												game = chess.NewGame(fen)
												fmt.Println(game.Position().Board().Draw())
												fmt.Println(game.Position())
												j = true
												control = true
											}
											fmt.Println(game.Position().Board().Draw())
											fmt.Println(game.Position())
											err := game.MoveStr(vt.Step)
											if err != nil {
												return errors.New(fmt.Sprintf("%s - %v", value.StepString, err))
											}
											fmt.Println(game.Position().Board().Draw())
											fmt.Println(game.FEN())
											fmt.Println()
											// i = i + 1
										}
									}

								}

							}
						}
					}
				} else {
					fmt.Println("next")
				}
			}
		}
	}
	return nil
}

func (f *FileProcessingService) CheckSecondSteps2(game *chess.Game) error {
	var arMetaStepFalse []domain.MetaStep
	for _, value := range f.InfoStep.ArrayMetaStep {

		if value.Main == false {
			arMetaStepFalse = append(arMetaStepFalse, value)
		}
	}

	for _, v := range arMetaStepFalse {
		b := v.MetaBoth[0].NumberStep
		g, err := f.CheckSecondSteps3(game, b)
		if err != nil {
			return err
		}
		var ar []domain.OneStep
		for _, val := range v.MetaBoth {
			for _, j := range val.OneStep {
				if j.Step == ".." {
					continue
				}
				var on domain.OneStep
				on = j
				err := g.MoveStr(j.Step)
				if err != nil {
					return err
				}
				on.FEN = g.FEN()
				ar = append(ar, on)
				fmt.Println(g.Position().Board().Draw())
				fmt.Println(g.FEN(), j.Step)
			}

		}

	}

	return nil
}

func (f *FileProcessingService) CheckSecondSteps3(game *chess.Game, numberStep int) (*chess.Game, error) {
	for _, value := range f.InfoStep.ArrayMetaStep {
		if value.Main == true {
			b := GetNumberStep(value, numberStep)
			// for _, v := range value.MetaBoth {
			// if v.NumberStep == numberStep {
			if len(b.OneStep) >= 2 {
				for _, q := range b.OneStep {
					fen, err := chess.FEN(b.OneStep[1].FEN)
					if err != nil {
						return nil, err
					}
					game = chess.NewGame(fen)
					fmt.Println(game.Position().Board().Draw())
					fmt.Println(game.Position(), q.Step)
				}
			} else {
				break
			}
		}
		break
	}
	return game, nil
}

func GetNumberStep(f domain.MetaStep, numberStep int) domain.MetaBoth {
	var b domain.MetaBoth
	for _, value := range f.MetaBoth {
		if value.NumberStep == numberStep-1 {
			b = value
		}
	}
	return b
}

func Replace(st string) string {
	st = strings.Replace(st, "Кр", "K", -1)
	st = strings.Replace(st, "К", "N", -1)
	st = strings.Replace(st, "С : ", "Bx", -1)
	st = strings.Replace(st, ".__", "", -1)
	st = strings.Replace(st, ". __", "", -1)
	st = strings.Replace(st, "Ф : ", "Qx", -1)
	st = strings.Replace(st, " : ", "x", -1)
	st = strings.Replace(st, "—", "-", -1)
	st = strings.Replace(st, "0-0", "O-O", -1)
	st = strings.Replace(st, "е", "e", -1)
	st = strings.Replace(st, " : ", "x", -1)
	st = strings.Replace(st, ":", "x", -1)
	st = strings.Replace(st, "Ф", "Q", -1)
	st = strings.Replace(st, "С", "B", -1)
	st = strings.Replace(st, "C", "B", -1)
	st = strings.Replace(st, "с", "c", -1)
	st = strings.Replace(st, "Л", "R", -1)
	st = strings.Replace(st, "0-0", "O-O", -1)
	return st
}
