package service

import (
	"chess-backend/internal/domain"
	"chess-backend/internal/repository"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/notnil/chess"
	"github.com/tenkoh/go-docc"
)

type FileProcessingService struct {
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

	// var info domain.InfoStep
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

// func (f *FileProcessingService) GetAll() {
// 	// var infoA []domain.InfoStep
// 	// var info domain.InfoStep
// 	for key, value := range f.InfoStep.AllStepsInPart {
// 		if len(value) > 1 {
// 			for k, v := range value {
// 				Replace(v)
// 				fmt.Println(key, value, k, v)
// 			}
// 		} else {
// 			Replace(value[0])
// 		}
// 	}
// 	fmt.Println(f.InfoStep.AllStepsInPart)
// }

// func (f *FileProcessingService) GetParagraph() {
// 	// var arMeta []domain.InfoStep
// 	// var meta domain.InfoStep
// 	for key, value := range f.InfoStep.AllStepsInPart {
// 		// meta.
// 		fmt.Println(key, value)
// 	}
// }

// func (f *FileProcessingService) GetMainSteps() {
// 	// var infoArr []domain.InfoStep
// 	var info domain.InfoStep
// 	var arrayOneRange []string
// 	for key, value := range f.InfoStep.AllStepsInPart {
// 		info.Paragraph = key
// 		if len(value) > 0 {
// 			for k, v := range value {
// 				info.NumberInParagraph = k
// 				arrayOneRange = append(arrayOneRange, v)
// 			}
// 		}
// 	}
// 	// f.InfoStep.ArrayOneRange = arrayOneRange
// 	var arrayChangeOneRange []string
// 	for _, value := range arrayOneRange {
// 		res := Replace(value)
// 		arrayChangeOneRange = append(arrayChangeOneRange, res)
// 	}
// 	f.InfoStep.ArrayOneRange = arrayChangeOneRange
// }

// func (f *FileProcessingService) ReadFF() {
// 	for key, value := range f.InfoStep.AllStepsInPart {
// 		fmt.Println(key, value)
// 	}
// }

func (f *FileProcessingService) ReadProcessing(file string) error {
	f.File = file

	r, err := docc.NewReader("upload/" + file)
	if err != nil {
		return err
	}

	// f.InfoStep.FileName = file

	ps, err := r.ReadAll()
	if err != nil {
		return err
	}
	f.GetSource(ps)
	f.GetPage()
	f.GetAllSteps()
	f.ParseMainAndSlaveSteps()
	f.GetInfoBoth()
	f.CheckSteps()
	r.Close()
	return nil
}

// func ReplaceFixStep(st string) string {
// 	// f.InfoStep.ArrayReplaces = make([]string, len(f.InfoStep.ArrayOneRange))
// 	// copy(f.InfoStep.ArrayReplaces, f.InfoStep.ArrayOneRange)
// 	// var doneReplace []string
// 	var lt = map[string]string{
// 		"Л":   "R",
// 		"К":   "N",
// 		"с":   "c",
// 		"C":   "B",
// 		"С":   "B",
// 		"Ф":   "Q",
// 		"Кр":  "K",
// 		":":   "x",
// 		" : ": "x",
// 		"е":   "e",
// 		"...": "",
// 		"0":   "O",
// 		"—":   "-",
// 	}
// 	// for _, st := range f.InfoStep.ArrayOneRange {
// 	// st = strings.Replace(st, "(_", "", -1)
// 	// st = strings.Replace(st, "_)", "", -1)
// 	// st = strings.Replace(st, ".__", "", -1)
// 	// st = strings.Replace(st, "__", "", -1)
// 	st = strings.Replace(st, "С : ", "Bx", -1)
// 	st = strings.Replace(st, ".__", "", -1)
// 	st = strings.Replace(st, ". __", "", -1)
// 	st = strings.Replace(st, "Ф : ", "Qx", -1)
// 	st = strings.Replace(st, " : ", "x", -1)
// 	for key, value := range lt {
// 		st = strings.Replace(st, key, value, -1)
// 	}
// 	return st
// 	// doneReplace = append(doneReplace, st)
// 	// }
// 	// copy(f.InfoStep.ArrayReplaces, doneReplace)
// }

// func (f *FileProcessingService) CleanLastDota() {
// 	var arrayCleanLastDora []string

// 	var re = regexp.MustCompile(`(?m)\.$`)
// 	for _, value := range f.InfoStep.ArrayReplaces {
// 		s := strings.TrimRight(value, " ")
// 		result := re.ReplaceAllString(s, "")
// 		arrayCleanLastDora = append(arrayCleanLastDora, result)
// 	}
// 	f.InfoStep.ArrayReplaces = arrayCleanLastDora
// }

func (f *FileProcessingService) ParseMainAndSlaveSteps() {
	var mm domain.MetaStep
	var arrayTypeStep []domain.MetaStep

	for _, value := range f.InfoStep.ArrayMetaStep {
		jf := strings.Contains(value.StepString, "(")
		if jf {
			mm = value
			mm.Main = false
			// mm.StepString = value
			arrayTypeStep = append(arrayTypeStep, mm)
		} else {
			mm = value
			mm.Main = true
			// mm.StepString = value
			arrayTypeStep = append(arrayTypeStep, mm)
		}
	}
	f.InfoStep.ArrayMetaStep = arrayTypeStep
}

func (f *FileProcessingService) GetInfoBoth() {
	var metaA []domain.MetaStep
	var meta domain.MetaStep
	var re = regexp.MustCompile(`(?m)[0-9]{1,2}\.`)
	for _, value := range f.InfoStep.ArrayMetaStep {
		var arIntIndex []int
		value.StepString = strings.Replace(value.StepString, "(_", "", -1)
		value.StepString = strings.Replace(value.StepString, "_)", "", -1)
		value.StepString = strings.Replace(value.StepString, ".__", "", -1)
		value.StepString = strings.Replace(value.StepString, "__", "", -1)
		meta = value
		arNumber := re.FindAllString(value.StepString, -1)
		for _, val := range arNumber {
			v1 := strings.Trim(val, "")
			v2 := strings.ReplaceAll(v1, ".", "")
			v3 := strings.ReplaceAll(v2, "", "")
			i, err := strconv.Atoi(v3)
			if err != nil {
				log.Fatal(i)
			}
			arIntIndex = append(arIntIndex, i)
		}

		arrayStep := re.Split(value.StepString, -1)
		var arMetaBothTrim []domain.MetaBoth
		var metaBoth domain.MetaBoth
		for _, v := range arrayStep {
			for _, va := range arIntIndex {
				v = strings.Trim(v, " ")
				if len(v) > 0 {
					metaBoth.BothString = v
					metaBoth.NumberStep = va
					arMetaBothTrim = append(arMetaBothTrim, metaBoth)
				}
			}
		}
		meta.MetaBoth = arMetaBothTrim
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
					arOneStep = append(arOneStep, oneStep)
					aOneStep = append(aOneStep, oneStep)
				} else {
					oneStep.Step = j
					oneStep.Color = "b"
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
	}
}

func (f *FileProcessingService) CheckSteps() {
	game := chess.NewGame()
	for _, value := range f.InfoStep.ArrayMetaStep {
		if value.Main {
			for _, v := range value.MetaBoth {
				// var araiOneStep []domain.OneStep
				for _, j := range v.OneStep {
					// araiOneStep = append(araiOneStep, j)
					if err := game.MoveStr(j.Step); err != nil {
						log.Println(j, " - ", err)
						break
					}
					j.FEN = game.FEN()
					fmt.Println(game.Position().Board().Draw())
				}
			}
		}
	}
}

func Replace(st string) string {
	var lt = map[string]string{
		"Л":   "R",
		"К":   "N",
		"с":   "c",
		"C":   "B",
		"С":   "B",
		"Ф":   "Q",
		"Кр":  "K",
		":":   "x",
		" : ": "x",
		"е":   "e",
		// "...": "",
		"0-0": "O-O",
		// "0": "O",
		"—": "-",
	}

	st = strings.Replace(st, "С : ", "Bx", -1)
	st = strings.Replace(st, ".__", "", -1)
	st = strings.Replace(st, ". __", "", -1)
	st = strings.Replace(st, "Ф : ", "Qx", -1)
	st = strings.Replace(st, " : ", "x", -1)
	for key, value := range lt {
		st = strings.Replace(st, key, value, -1)
	}
	return st
}

// func (f *FileProcessingService) EveryBothSteps() {
// 	var re = regexp.MustCompile(`(?m)[0-9]{1,2}\.`)
// 	for _, value := range f.InfoStep.ArrayTypeSteps {
// 		value.StepString = strings.Replace(value.StepString, "(_", "", -1)
// 		value.StepString = strings.Replace(value.StepString, "_)", "", -1)
// 		value.StepString = strings.Replace(value.StepString, ".__", "", -1)
// 		value.StepString = strings.Replace(value.StepString, "__", "", -1)
// 		// fmt.Println(value.StepString)
// 		arrayIndex := re.FindAllString(value.StepString, -1)
// 		arrayStep := re.Split(value.StepString, -1)

// 		var arrIndex []int
// 		var arrStep []string
// 		var st domain.MetaStep
// 		var stBoth []domain.MetaStep

// 		// re2 := regexp.MustCompile(`(?m)\.`)
// 		for _, value := range arrayIndex {
// 			v1 := strings.Trim(value, "")
// 			v2 := strings.ReplaceAll(v1, ".", "")
// 			v3 := strings.ReplaceAll(v2, "", "")
// 			i, err := strconv.Atoi(v3)
// 			if err != nil {
// 				log.Fatal(i)
// 			}
// 			// fmt.Println(v3)
// 			arrIndex = append(arrIndex, i)
// 		}
// 		// fmt.Println(arrIndex, " - ", arrayStep)

// 		for _, value := range arrayStep {
// 			v1 := strings.Trim(value, " ")
// 			if value != "" {
// 				arrStep = append(arrStep, v1)
// 			}
// 		}

// 		for i := 0; i < len(arrIndex); i++ {
// 			st.Number = arrIndex[i]
// 			st.BothString = arrStep[i]
// 			stBoth = append(stBoth, st)
// 		}

// 		f.InfoStep.AllInfoStep = stBoth

// 		var oneStep domain.OneStep
// 		var arrayOneStep []domain.OneStep

// 		var va []domain.MetaStep
// 		// var strA []string
// 		for _, value = range f.InfoStep.AllInfoStep {
// 			v1 := strings.Split(value.StepString, " ")
// 			value.Both = v1
// 			for i, _ := range value.Both {
// 				if i != 0 {
// 					oneStep.Color = "w"
// 					oneStep.Step = v1[i]
// 					oneStep.Number = value.Number
// 					arrayOneStep = append(arrayOneStep, oneStep)
// 				} else {
// 					oneStep.Color = "b"
// 					oneStep.Step = v1[i]
// 					oneStep.Number = value.Number
// 					arrayOneStep = append(arrayOneStep, oneStep)
// 				}
// 			}
// 			value.EveryStep = arrayOneStep
// 			va = append(va, value)
// 		}
// 		f.InfoStep.AllInfoStep = va

// 	}
// }

// func (f *FileProcessingService) InfoStepsN() {

// 	var arraySteps []domain.MetaStep
// 	for _, v := range f.InfoStep.ArrayMapSteps {
// 		for k, j := range v {
// 			var step domain.MetaStep
// 			res, err := strconv.ParseBool(k)
// 			if err != nil {
// 				step.Main = res
// 			}

// 			step.StepString = j
// 			arraySteps = append(arraySteps, step)
// 		}
// 	}
// 	// fmt.Println(arraySteps)
// }

// func (f *FileProcessingService) ParseSteps() {

// 	var re = regexp.MustCompile(`(?m)\d{0,2}.\.`)
// 	var splitArray [][]string
// 	var ll []string

// 	for _, value := range f.InfoStep.ArrayReplaces {
// 		number := re.FindAllString(value, -1)
// 		fmt.Println(number)
// 		sa := re.Split(value, -1)
// 		for _, v := range sa {
// 			if len(v) != 0 {
// 				v = strings.Trim(v, " ")
// 				ll = append(ll, v)
// 			}
// 		}
// 		splitArray = append(splitArray, ll)
// 	}
// 	f.InfoStep.ArraySplitSteps = splitArray
// }
