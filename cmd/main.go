package main

import (
	"chess-backend/internal/app"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoConnect = "mongodb://admin:S3cret@localhost:27017/admin"

const configsDir = "../configs"

func main() {

	// soundfile.Read("../001.m4a")

	// GetFEN()
	if err := app.Run(configsDir); err != nil {
	}

	// client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoConnect))
	// if err != nil {
	// 	panic(err)
	// }

	// defer func() {
	// 	if err := client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// Parse(client)
	// SaveDatabase(client)

	// router := gin.Default()
	// router.GET("/getallstep", getAllStep)

	// router.GET("/albums/:id", getAlbumByID)
	// router.POST("/albums", postAlbums)

	// router.Run("localhost:8080")
}

func SaveDatabase(client *mongo.Client) {
	// var infoSTtt []InfoStep

	// filter := bson.D{{"Main", true}}

	// coll := client.Database("ChessBoard").Collection("Parts")
	// cursor, err := coll.Find(context.TODO(), filter)
	// if err != nil {
	// 	log.Println(err)
	// }

	// for cursor.Next(context.Background()) {

	// 	var iifo InfoStep
	// 	if err = cursor.Decode(&iifo); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	infoSTtt = append(infoSTtt, iifo)
	// }
	// var stepMeta MetaStep
	// var steps []MetaStep
	// for _, va := range infoSTtt {
	// 	game := chess.NewGame()

	// 	zz := va.ArraySteps.Step

	// 	var lineStepsCheck []string

	// 	// var infoStep InfoStep

	// 	for _, value := range zz {

	// 		if err := game.MoveStr(value); err != nil {

	// 			log.Printf("неправильный ход - %v; line: %v ;   Error: - %v", value, va.Line, err)
	// 			fmt.Println(game.Position().Board().Draw())
	// 			break
	// 		} else {
	// 			stepMeta.StepString = value
	// 			stepMeta.FEN = game.FEN()
	// 			stepMeta.Position = game.Position()
	// 			lineStepsCheck = append(lineStepsCheck, value)
	// 			steps = append(steps, stepMeta)
	// 		}
	// 	}
	// 	fmt.Println(lineStepsCheck)
	// }
}

// func Parse(client *mongo.Client) {

// 	coll := client.Database("ChessBoard").Collection("Parts")

// 	fp := filepath.Clean("../example.docx")
// 	r, err := docc.NewReader(fp)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer r.Close()
// 	ps, _ := r.ReadAll()

// 	var re = regexp.MustCompile(`(?m)[(_]?\w.+?_[_)]`)

// 	var re3 = regexp.MustCompile(`\d+\.`)
// 	for kkk, value := range ps {
// 		stringSteps := re.FindAllString(value, -1)
// 		for _, val := range stringSteps {
// 			if len(val) != 0 {
// 				if len(val) != 4 {
// 					var infoStep InfoStep
// 					var stepAction []string
// 					infoStep.Line = val
// 					infoStep.NumberLine = kkk + 1
// 					ss := ChangeSymbol(val)

// 					jf := strings.Contains(val, "(")
// 					if jf {
// 						infoStep.Main = false
// 					} else {
// 						infoStep.Main = true
// 					}
// 					ss2 := strings.Split(ss, " ")

// 					red := make(map[string]string)
// 					for _, v := range ss2 {
// 						stSy := re3.FindAllString(v, -1)

// 						for _, vv := range stSy {
// 							red[vv] = vv
// 						}
// 					}

// 					for ke := range red {
// 						ss = strings.Replace(ss, ke, "", -1)
// 					}
// 					arrayStep := strings.Split(ss, " ")

// 					for _, value := range arrayStep {
// 						if len(value) != 0 {
// 							if value != "." {
// 								stepAction = append(stepAction, value)
// 							}
// 						}
// 					}

// 					oneLine := strings.Join(stepAction, " ")
// 					infoStep.OneLine = oneLine
// 					infoStep.ArraySteps.Step = stepAction
// 					game := chess.NewGame()
// 					for _, value := range stepAction {
// 						allInfo, err := checkStepInGame(value, game)
// 						if err != nil {
// 							log.Println("Обработать ...", value)
// 							log.Fatalln()
// 						}
// 						infoStep.AllInfoStep = append(infoStep.AllInfoStep, allInfo)
// 					}

// 					// infoStep.Fen =
// 					res, err := coll.InsertOne(context.TODO(), infoStep)
// 					if err != nil {
// 						log.Printf("don't write, %v", err)
// 					}
// 					fmt.Println(res)
// 				}
// 			}
// 		}

// 	}
// }

func ChangeSymbol(st string) string {
	var lt = map[string]string{
		"__":  "",
		".__": "",
		"(_":  "",
		"_)":  "",
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
		"...": "",
		"0":   "O",
		"—":   "-",
	}

	st = strings.Replace(st, "С : ", "Bx", -1)
	st = strings.Replace(st, ".__", "", -1)
	st = strings.Replace(st, "Ф : ", "Qx", -1)
	for key, value := range lt {
		st = strings.Replace(st, key, value, -1)
	}
	return st
}

// func GetFEN() {
// 	game := chess.NewGame()
// 	err := game.UnmarshalText([]byte("1.e4 e5 2.f4 exf4 3.Nf3 g5 4.Bc4 g4 5.O-O gxf3"))
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	fmt.Println(game.Position().Board().Draw())
// 	fmt.Printf("Game completed. %s by %s.\n", game.Outcome(), game.Method())
// 	fmt.Println(game.String())
// 	// fmt.Println(game.FEN())
// }

// func checkStepInGame(step string, game *chess.Game) (MetaStep, error) {
// 	var metaStep MetaStep

// 	if err := game.MoveStr(step); err != nil {
// 		fmt.Println(game.Position().Board().Draw())
// 		fmt.Println(game.Position().Board().String())
// 		fmt.Println(game.String())
// 		return metaStep, err
// 	} else {
// 		metaStep.FEN = game.FEN()
// 		metaStep.Position = game.Position().Board().String()
// 		metaStep.StepString = step
// 		metaStep.Draw = game.Position().Board().Draw()
// 		fmt.Println(metaStep.Draw)
// 		fmt.Println(game.String())
// 		return metaStep, nil
// 	}
// }

func getAllStep(c *gin.Context) {

}
