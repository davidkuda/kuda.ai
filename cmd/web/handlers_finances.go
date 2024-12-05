package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	sankeyNode = []opts.SankeyNode{
		{Name: "category1"},
		{Name: "category2"},
		{Name: "category3"},
		{Name: "category4"},
		{Name: "category5"},
		{Name: "category6"},
	}

	sankeyLink = []opts.SankeyLink{
		{Source: "category1", Target: "category2", Value: 10},
		{Source: "category2", Target: "category3", Value: 15},
		{Source: "category3", Target: "category4", Value: 20},
		{Source: "category5", Target: "category6", Value: 25},
	}
)

func (app *application) finances(w http.ResponseWriter, r *http.Request) {

	err := app.checkJWTCookie(r)
	if err != nil {
		log.Printf("could not authenticate client: %v", err)
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	} else {
		log.Println("Could authenticate client :)")
	}

	// TODO: htmlBytes := chart
	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Sankey-basic-example",
		}),
	)

	sankey.AddSeries("sankey", sankeyNode, sankeyLink, charts.WithLabelOpts(opts.Label{Show: opts.Bool(true)}))

	var buf bytes.Buffer
	err = sankey.Render(&buf)
	if err != nil {
		log.Printf("failed to render chart: %v\n", err)
	}

	pageData := Page{
		Title:   getTitleFromRequestPath(r),
		Content: template.HTML(buf.Bytes()),
	}

	renderPageSimple(w, &pageData)
}
