package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

// Go has neither a native GUI, nor mature bindings to Qt or another similarly
// sophisticated library. So this program explores a way for Go to produce a
// locally running GUI app using an HTML5 web-app architecture, in which the
// content delivery and the dedicated server are compiled together into a
// single deployable executable. It additionally, compiles the html, css and
// template files required into the executable, so the executable has no
// runtime dependencies apart from a browser to display it. The auxiliary
// files are converted into compilable Go source code using the
// github.com/jteeuwen/go-bindata Go package. The example GUI is a loose copy
// of the Github GUI, and its controls, layout and style are all implemented with
// the Bootstrap CSS library. Go's native html templating is used.
func main() {
	// Unpack the compiled file resources into an in-memory virtual file system.
	virtualFs := &assetfs.AssetFS{
		Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}

	// Prepare an html template that will be combined with a data model to
	// serve html pages.
	// We need two html's one for the selectionPage (first screen) and another for the loadingPage (second screen)
	guiHTMLTemplateSelection = extractAndParseHTMLTemplate("selectionPage.html", templateSelectionName)

	// Route incoming web page requests for static URLs (like css files) to
	// the standard library's file server.
	http.Handle("/static/", http.FileServer(virtualFs))

	// Route incoming web page requests for the GUI home page to the dedicated
	// handler.
	http.HandleFunc("/medco-loader", guiSelectionPageHandler)
	http.HandleFunc("/medco-loader/lala", lala)

	fmt.Printf(
		"To see the GUI, visit this URL with your Web Browser:\n\n %s\n\n",
		"http://localhost:47066/medco-loader")

	// Spin-up the standard library's http server on a hard-coded port.
	http.ListenAndServe(":47066", nil)

}

// Provides a parsed html template, having first extracted the file
// representation of its text from a compiled resource.
func extractAndParseHTMLTemplate(htmlName string, templateName string) *template.Template {
	// Expose errors by permitting panic response.
	bytes, _ := Asset("templates/" + htmlName)
	parsedTemplate, _ := template.New(templateName).Parse(string(bytes))
	return parsedTemplate
}

// GuiDataModel is a data structure for the model part of the example GUI's model-view pattern.
type GuiDataModel struct {
	Title       string
	Unwatch     int
	Star        int
	Fork        int
	Commits     int
	Branch      int
	Release     int
	Contributor int
	RowsInTable []TableRow
}

// TableRow is a sub-model to the GuiDataModel that models a single row in the table
// displayed in the GUI.
type TableRow struct {
	File    string
	Comment string
	Ago     string
	Icon    string
}

// Sends the html required to render the GUI into the provided http
// response writer.
func guiSelectionPageHandler(w http.ResponseWriter, r *http.Request) {
	// Generate the html by plugging in data from the gui data model into the
	// prepared html template.
	err := guiHTMLTemplateSelection.ExecuteTemplate(w, templateSelectionName, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func lala(w http.ResponseWriter, r *http.Request) {
	fmt.Print("sdasd")
}

// Provides an illustrative, hard-coded instance of a GuiDataModel.
func guiData() *GuiDataModel {
	guiData := &GuiDataModel{
		Title:       "Golang Standalone GUI Example",
		Unwatch:     3,
		Star:        0,
		Fork:        2,
		Commits:     31,
		Release:     1,
		Contributor: 1,
		RowsInTable: []TableRow{},
	}
	guiData.RowsInTable = append(guiData.RowsInTable,
		TableRow{"do_this.go", "Initial commit", "1 month ago", "file"})
	guiData.RowsInTable = append(guiData.RowsInTable,
		TableRow{"do_that.go", "Initial commit", "1 month ago", "file"})
	guiData.RowsInTable = append(guiData.RowsInTable,
		TableRow{"index.go", "Initial commit", "1 month ago", "file"})
	guiData.RowsInTable = append(guiData.RowsInTable,
		TableRow{"resources", "Initial commit", "2 months ago", "folder-open"})
	guiData.RowsInTable = append(guiData.RowsInTable,
		TableRow{"docs", "Initial commit", "2 months ago", "folder-open"})
	return guiData
}

// Makes the the GUI templates available at module-scope.
var guiHTMLTemplateSelection *template.Template
var templateSelectionName = "selectionGUI"
