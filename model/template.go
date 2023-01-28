package model

import (
	"blog/common"
	"blog/config"
	"html/template"
	"log"
	"net/http"
	"time"
)

type TemplateBlog struct {
	*template.Template
}

type HtmlTemplate struct {
	Index      TemplateBlog
	Category   TemplateBlog
	Custom     TemplateBlog
	Detail     TemplateBlog
	Login      TemplateBlog
	Pigeonhole TemplateBlog
	Writing    TemplateBlog
	Error      TemplateBlog
}

var Template HtmlTemplate

func init() {
	templateDir := config.Cfg.System.CurrentDir + "/template/"
	tp := readTemplate(
		[]string{"index", "category", "custom", "detail", "login", "pigeonhole", "writing", "error"},
		templateDir,
	)
	htp := HtmlTemplate{
		Index:      tp[0],
		Category:   tp[1],
		Custom:     tp[2],
		Detail:     tp[3],
		Login:      tp[4],
		Pigeonhole: tp[5],
		Writing:    tp[6],
		Error:      tp[7],
	}
	Template = htp
}

func readTemplate(templates []string, templateDir string) (tbs []TemplateBlog) {
	home := templateDir + "home.html"
	header := templateDir + "layout/header.html"
	footer := templateDir + "layout/footer.html"
	personal := templateDir + "layout/personal.html"
	post := templateDir + "layout/post-list.html"
	pagination := templateDir + "layout/pagination.html"

	for _, view := range templates {
		viewName := view + ".html"
		t := template.New(viewName)

		t.Funcs(template.FuncMap{"isODD": isODD, "getNextName": getNextName, "date": date, "dateDay": common.DateDay})

		t, err := t.ParseFiles(templateDir+viewName, home, header, footer, personal, post, pagination)
		if err != nil {
			panic("加载模板失败，" + err.Error())
		}
		tbs = append(tbs, TemplateBlog{t})
	}
	return
}

func (t *TemplateBlog) WriteData(w http.ResponseWriter, data any) {
	err := t.Execute(w, data)
	if err != nil {
		WriteError(w, err, "加载页面失败")
	}
}

func WriteError(w http.ResponseWriter, err error, msgAndCode ...interface{}) {
	if err != nil {
		log.Println(err)
	}

	code, msg := 500, "未知错误"
	if len(msgAndCode) >= 1 {
		msg = msgAndCode[0].(string)
	}
	if len(msgAndCode) >= 2 {
		code = msgAndCode[1].(int)
	}

	w.WriteHeader(code)
	t := Template.Error
	t.WriteData(w, struct {
		config.Viewer
		Code int
		Msg  string
	}{config.Cfg.Viewer, code, msg})
}

func isODD(num int) bool {
	return num%2 == 0
}

func getNextName(navs []string, index int) string {
	return navs[index+1]
}

func date(layout string) string {
	return time.Now().Format(layout)
}
