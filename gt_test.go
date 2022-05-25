package gt_test

import (
	"time"

	. "github.com/Contra-Culture/gt"
	"github.com/Contra-Culture/report"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("gt", func() {
	It("creates limbo, than universe and then renders templates", func() {
		// limbo creation
		now, err := time.Parse(time.RFC3339Nano, "2022-05-02T10:11:12.000000000Z")
		Expect(err).NotTo(HaveOccurred())
		limbo := New(report.ReportCreator(report.DumbTimer(now)))
		Expect(limbo).NotTo(BeNil())
		// limbo.Trait(
		// 	"text/top-header",
		// 	Prop("color", "red"),
		// 	Prop("font-family", "Helvetica Neue", "Helvetica", "Arial", "sans-serif"),
		// 	Prop("font-weight", "600"),
		// )
		// limbo.Trait(
		// 	"text/main",
		// 	Prop("color", "#000000"),
		// 	Prop("font-size", "1em"),
		// 	Prop("line-height", "1.6em"),
		// )
		// limbo.Trait(
		// 	"/block/header",
		// 	Prop("background", "#f5f6f7"),
		// 	Prop("margin", "20px", "0", "20px", "0"),
		// 	Prop("padding", "20px"),
		// )
		limbo.Template(
			"/layout/test",
			DocumentContent(
				Doctype(),
				Tag("html",
					Attributes(),
					Content(
						Tag("head",
							Attributes(),
							Content(
								Tag(
									"title",
									Attributes(),
									Content(Text("::Test Template::"))),
								Tag(
									"meta",
									Attributes(
										Attr("charset", "utf-8")),
									Content()))),
						Tag(
							"body",
							Attributes(),
							Content(
								Tag(
									"header",
									Attributes(
										Attr("class", "top-header"),
									),
									Content(
										Tag(
											"h1",
											Attributes(
												Attr("class", "top-header-title")),
											Content(Text("Test Header!"))))),
								Repeat(
									"articles",
									TemplatePlacement("/card/article", Auto()),
								),
								TemplateInjection("bottom")))))))
		limbo.Template(
			"/card/article",
			Content(
				Tag(
					"div",
					Attributes(
						Attr("class", "article-card"),
					),
					Content(
						Tag(
							"h1",
							Attributes(
								Attr("class", "article-card-title")),
							Content(
								TextInj("article-title"))),
						Tag(
							"span",
							Attributes(
								Attr("class", "article-card-preview")),
							Content(
								TextInj("article-preview"))),
						Tag(
							"a",
							Attributes(
								Attr("class", "article-card-link"),
								AttrInjection("href", "article-link")),
							Content(
								TextInj("article-link-anchor")))))))
		limbo.Template(
			"/btn/mailme",
			Content(
				Tag(
					"a",
					Attributes(
						Attr("class", "mailme-btn"),
						AttrInjection("href", "mailme-mailto"),
					),
					Content(TextInj("mailme-text")))))
		// creates universe
		univ, r := limbo.Universe()
		Expect(report.ToString(r)).To(Equal("#[2022-05-02T10:11:12.0000001Z] universe\n"))
		Expect(r.HasErrors()).To(BeFalse())
		Expect(univ).NotTo(BeNil())
		rendered, r := univ.Render(
			"/layout/test",
			map[string]interface{}{
				"bottom": map[string]interface{}{
					"name": "/btn/mailme",
					"params": map[string]interface{}{
						"mailme-mailto": "mailto:egotraumatic@example.com",
						"mailme-text":   "Mail me",
					},
				},
				"articles": []map[string]interface{}{
					{
						"article-title":       "Article 1",
						"article-preview":     "Preview for article 1.",
						"article-link":        "http://google.com",
						"article-link-anchor": "google",
					},
					{
						"article-title":       "Article 2",
						"article-preview":     "Preview for article 2.",
						"article-link":        "http://yahoo.com",
						"article-link-anchor": "yahoo!",
					},
				}})
		Expect(report.ToString(r)).To(Equal("#[2022-05-02T10:11:12.0000002Z] rendering template \"/layout/test\"\n"))
		Expect(rendered).To(Equal("<!DOCTYPE html><html><head><title>::Test Template::</title><meta charset=\"utf-8\"/></head><body><header class=\"top-header\"><h1 class=\"top-header-title\">Test Header!</h1></header><div class=\"article-card\"><h1 class=\"article-card-title\">Article 1</h1><span class=\"article-card-preview\">Preview for article 1.</span><a class=\"article-card-link\" href=\"http://google.com\">google</a></div><div class=\"article-card\"><h1 class=\"article-card-title\">Article 2</h1><span class=\"article-card-preview\">Preview for article 2.</span><a class=\"article-card-link\" href=\"http://yahoo.com\">yahoo!</a></div><a class=\"mailme-btn\" href=\"mailto:egotraumatic@example.com\">Mail me</a></body></html>"))
	})
})
