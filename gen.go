// The following directive is necessary to make the package coherent:
//go:build ignore

// This script generates new module in internal/modules/. It can be invoked by running go run gen.go
package main

import (
	_ "embed"

	"fmt"
	// "io/ioutil"
	"log"
	"os"
	// "os/exec"
	"strings"
	"sync"
	"text/template"

	// "golang.org/x/text/cases"
	// "golang.org/x/text/language"
	"github.com/iancoleman/strcase"
)

type entity struct {
	ModuleName string
	StructName string // just Uppercase 1st char
	Plural     string // the plural of the new module
	Initial    *string
	RouteName  *string
	TableName  *string
	Path       string
}

var wg sync.WaitGroup

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("error: missing arg[1] arg[2], try go run gen.go document doc")
		return
	}
	if len(os.Args) == 2 {
		fmt.Println("error: missing new module name")
		fmt.Println("try: go run gen.go <module-name-in-singular-lower-case> <initial e.g: u (for User)>")
		return
	}

	newModule := getNewModuleStruct(os.Args[1])

	fmt.Printf("newModule %+v,\n initial: %s,\n route: %s,\n tableName: %s\n", newModule, *newModule.Initial, *newModule.RouteName, *newModule.TableName)

	/* create directory */
	if err := os.Mkdir(newModule.Path, 0755); err != nil {
		fmt.Println("err:", err)
	}
	fmt.Printf("created %s\n\n", newModule.Path)
	// wg.Add(7) // number of go routines
	/* generate all related files route, controller, service, repo, interface, model, migration */
	// modelsDirectory := "internal/database/models"
	newModule.createFile("type", typeTemplate)
	newModule.createFile("controller", controllerTemplate)
	// go generateRoute(newModule, newDirectory)
	// go generateController(newModule, newDirectory)
	// go generateService(newModule, newDirectory)
	// go generateRepository(newModule, newDirectory)
	// go generateInterfaces(newModule, newDirectory)
	// go generateModel(newModule, modelsDirectory)
	// go generateMigration(&newModule)
	// wg.Wait()
	//
	// /* re-generate server.go */
	// reGenerateServerFile(newModule)
}

func getNewModuleStruct(inputName string) *entity {
	var (
		// inputName string = os.Args[1]
		// structName string = fmt.Sprintf("%s", cases.Title(language.English, cases.Compact).String(inputName))
		structName string = strcase.ToCamel(inputName)
		plural     string = Pluralfy(structName)
		routeName  string = strcase.ToKebab(plural)
		tableName  string = strings.ToLower(Pluralfy(strcase.ToSnake(inputName)))
		initial    string = os.Args[2]
	)

	newDirectory := fmt.Sprintf("internal/modules/%s", inputName)
	newModule := &entity{inputName, structName, plural, &initial, &routeName, &tableName, newDirectory}
	return newModule
}

func Pluralfy(word string) (plural string) {
	if word[len(word)-1:] == "y" { // handle the word ends with y --> ies
		plural = word[0:len(word)-1] + "ies"
	} else if word[len(word)-1:] == "s" {
		plural = word + "es" // handle the word ends with s --> es
	} else {
		plural = word + "s"
	}
	return plural
}

func (e *entity) createFile(fileName, templateFile string) {
	filePath := fmt.Sprintf("%s/%s.go", e.Path, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("create %s failed: %s\n", filePath, err)
		return
	}

	t := template.Must(template.New(filePath).Parse(templateFile))
	t.Execute(file, e)
}

func (e *entity) generateType() {
	// wg.Done()
}

// func generateRoute(newModule entity, dirPath string) {
// 	filePath := fmt.Sprintf("%s/route.go", dirPath)
// 	createFile(newModule, filePath, "route", routeTemplate)
// 	wg.Done()
// }

// func generateController(newModule entity, dirPath string) {
// 	filePath := fmt.Sprintf("%s/controller.go", dirPath)
// 	createFile(newModule, filePath, "controller", controllerTemplate)
// 	wg.Done()
// }

// func generateService(newModule entity, dirPath string) {
// 	filePath := fmt.Sprintf("%s/service.go", dirPath)
// 	createFile(newModule, filePath, "service", serviceTemplate)
// 	wg.Done()
// }

// func generateRepository(newModule entity, dirPath string) {
// 	filePath := fmt.Sprintf("%s/repository.go", dirPath)
// 	createFile(newModule, filePath, "repository", repositoryTemplate)
// 	wg.Done()
// }

// func generateInterfaces(newModule entity, dirPath string) {
// 	filePath := fmt.Sprintf("%s/interfaces.go", dirPath)
// 	createFile(newModule, filePath, "interfaces", interfacesTemplate)
// 	wg.Done()
// }

// func generateModel(newModule entity, dirPath string) {
// 	filePath := fmt.Sprintf("%s/%s.go", dirPath, newModule.ModuleName)
// 	createFile(newModule, filePath, "model", modelTemplate)
// 	fmt.Printf("created %s/%s.go, \nplease add the fields(columns) in this file\n\n", dirPath, newModule.ModuleName)
// 	wg.Done()
// }

// func generateMigration(newModule *entity) {
// 	newModule.Plural = strings.ToLower(newModule.Plural)
// 	whichOutput, _ := exec.Command("which", "migrate").Output()
// 	// fmt.Println("whichOutput",string(whichOutput))
// 	migrateBinPath := strings.Fields(string(whichOutput))
// 	if len(migrateBinPath) == 0 {
// 		log.Fatal("migrate command not found")
// 	}
// 	// fmt.Println("migrateBinPath?",migrateBinPath[1])
//
// 	migrationName := fmt.Sprintf("create_%s", *newModule.TableName)
// 	argstr := []string{"create", "-ext", "sql", "-dir", "migrations", "-seq", migrationName}
// 	out, err := exec.Command(migrateBinPath[0], argstr...).CombinedOutput()
// 	migrationOut := strings.Split(strings.ReplaceAll(string(out), "\r\n", "\n"), "\n")
// 	// fmt.Println("up", migrationOut[0])
// 	// fmt.Println("down", migrationOut[1])
//
// 	/* create migrateion-up.sql */
// 	createFile(*newModule, migrationOut[0], "migrate-up", migrateUpTemplate)
// 	fmt.Printf("created %s\n", migrationOut[0])
//
// 	/* create migrateion-down.sql */
// 	createFile(*newModule, migrationOut[1], "migrate-down", migrateDownTemplate)
// 	fmt.Printf("created %s\n\n", migrationOut[1])
//
// 	if err != nil {
// 		log.Fatal(err)
// 		os.Exit(1)
// 	}
//
// 	fmt.Printf("DB migration files for %s created in ./migrations, \nplease go to add the SQL statements in up+down files, and then run: make migrate-up \n\n", newModule.ModuleName)
// 	wg.Done()
// }

// func reGenerateServerFile(newModule entity) {
// 	var allModules []entity
// 	moduleDirs, err := ioutil.ReadDir("internal/modules/")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	for _, dir := range moduleDirs {
// 		// structName := fmt.Sprintf("%s", cases.Title(language.English, cases.Compact).String(dir.Name()))
// 		// structName := strcase.ToCamel(dir.Name())
// 		// routeName := strcase.ToKebab(dir.Name())
// 		// tableName := strings.ToLower(Pluralfy(routeName))
// 		module := getNewModuleStruct(dir.Name())
// 		module.Initial = nil
//
// 		// module := &entity{ModuleName: dir.Name(), StructName: structName, Plural: strings.ToLower(Pluralfy(structName)), RouteName: &routeName, TableName:&tableName }
// 		allModules = append(allModules, module)
// 	}
//
// 	// fmt.Println(allModules)
// 	filePath := fmt.Sprintf("cmd/server/server.go")
// 	tmplData := map[string][]entity{"Modules": allModules}
//
// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		log.Println("re-generate server.go failed: ", err)
// 		return
// 	}
//
// 	t := template.Must(template.New("server").Parse(serverTemplate))
// 	t.Execute(file, tmplData)
// }

//go:embed skel/route.tmpl
var routeTemplate string

//go:embed skel/controller.tmpl
var controllerTemplate string

//go:embed skel/service.tmpl
var serviceTemplate string

//go:embed skel/repository.tmpl
var repositoryTemplate string

//go:embed skel/type.tmpl
var typeTemplate string

//go:embed skel/model.tmpl
var modelTemplate string

//go:embed skel/migrate-up.tmpl
var migrateUpTemplate string

//go:embed skel/migrate-down.tmpl
var migrateDownTemplate string

//go:embed skel/server.tmpl
var serverTemplate string
