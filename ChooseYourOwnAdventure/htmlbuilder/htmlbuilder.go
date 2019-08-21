package htmlbuilder

import (
	"bytes"
	"html/template"

	"github.com/mishuk-sk/gopher/ChooseYourOwnAdventure/htmlbuilder/storyparser"
)

const (
	htmlTemplate = `<html><head><meta charset="utf-8"><title>Choose Your Own Adventure</title></head>{{template "body" .}}{{template "styles"}}</html>`
	bodyTemplate = `{{define "body"}}<body><section class="page"><h1>{{.Title}}</h1>{{range .Paragraphs}}<p>{{.}}</p>{{end}}
	{{template "options" .}}
	</section></body>{{end}}`
	optionsTemplate = `{{define "options"}}<ul>{{range .Options}}<li><a href="{{.Link}}">{{.Text}}</a></li>{{end}}</ul>{{end}}`
	stylesTemplate  = `{{define "styles"}}<style>
	body {
	  font-family: helvetica, arial;
	}
	h1 {
	  text-align:center;
	  position:relative;
	}
	.page {
	  width: 80%;
	  max-width: 500px;
	  margin: auto;
	  margin-top: 40px;
	  margin-bottom: 40px;
	  padding: 80px;
	  background: #FFFCF6;
	  border: 1px solid #eee;
	  box-shadow: 0 10px 6px -6px #777;
	}
	ul {
	  border-top: 1px dotted #ccc;
	  padding: 10px 0 0 0;
	  -webkit-padding-start: 0;
	}
	li {
	  padding-top: 10px;
	}
	a,
	a:visited {
	  text-decoration: none;
	  color: #6295b5;
	}
	a:active,
	a:hover {
	  color: #7792a2;
	}
	p {
	  text-indent: 1em;
	}
  </style>{{end}}`
)

// GetPage returns slice of bytes, representing rendered html page
func GetPage(arc storyparser.Arc) []byte {
	html := template.Must(template.New("html").Parse(htmlTemplate))
	template.Must(html.Parse(bodyTemplate))
	template.Must(html.Parse(optionsTemplate))
	tmplt := template.Must(html.Parse(stylesTemplate))
	buf := bytes.NewBuffer([]byte{})
	tmplt.Execute(buf, arc)
	return buf.Bytes()
}
