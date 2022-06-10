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
		limbo.Trait(
			"text/top-header",
			Normal(
				[]string{"color", "red"},
				[]string{"font-family", "Helvetica Neue", "Helvetica", "Arial", "sans-serif"},
				[]string{"font-weight", "600"}))
		limbo.Trait(
			"text/main",
			Normal(
				[]string{"color", "#000000"},
				[]string{"font-size", "1em"},
				[]string{"line-height", "1.6em"}))
		limbo.Trait(
			"block/header",
			Normal(
				[]string{"background", "#f5f6f7"},
				[]string{"margin", "20px", "0", "20px", "0"},
				[]string{"padding", "20px"}),
			Pseudo(
				"hover",
				[]string{"background", "#e5e6e7"}))

		limbo.Template(
			"/layout/test",
			WithStylesheet("main"),
			WithLayout(
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
										SemClass("top-header", UsesTraits("block/header"))),
									Content(
										Tag(
											"h1",
											Attributes(
												SemClass("top-header-title", UsesTraits("text/top-header"))),
											Content(Text("Test Header!"))))),
								Repeat(
									"articles",
									TemplatePlacement("/card/article", Auto()),
								),
								TemplateInjection("bottom")))))))

		limbo.Template(
			"/card/article",
			WithStylesheet("main"),
			WithContent(
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
								TextInj("article-link-anchor"))),
						Variant("/comments/empty", map[string]string{"top-comments": "/comments/top"})))))

		limbo.Template(
			"/comments/empty",
			WithStylesheet("main"),
			WithContent(Text("no comments")))

		limbo.Template(
			"/comments/top",
			WithStylesheet("main"),
			WithContent(
				Repeat("comments", TemplatePlacement("/card/comment", Auto()))))

		limbo.Template(
			"/card/comment",
			WithStylesheet("main"),
			WithContent(
				Tag(
					"div",
					Attributes(
						Attr("class", "comment-card")),
					Content(
						Tag(
							"span",
							Attributes(
								Attr("class", "comment-card-author"),
							),
							Content(
								TextInj("comment-author"))),
						Tag(
							"p",
							Attributes(
								Attr("class", "comment-card-text"),
							),
							Content(
								TextInj("comment-text")))))))

		limbo.Template(
			"/btn/mailme",
			WithStylesheet("main"),
			WithContent(
				Tag(
					"a",
					Attributes(
						Attr("class", "mailme-btn"),
						AttrInjection("href", "mailme-mailto"),
					),
					Content(TextInj("mailme-text")))))
		// creates universe
		univ, r := limbo.Universe()
		Expect(report.ToString(r)).To(Equal("#[2022-05-02T10:11:12.0000001Z] limbo\n\t#[2022-05-02T10:11:12.0000002Z] trait \"text/top-header\"\n\t#[2022-05-02T10:11:12.0000003Z] trait \"text/main\"\n\t#[2022-05-02T10:11:12.0000004Z] trait \"block/header\"\n\t#[2022-05-02T10:11:12.0000005Z] template \"/layout/test\"\n\t#[2022-05-02T10:11:12.0000006Z] template \"/card/article\"\n\t#[2022-05-02T10:11:12.0000007Z] template \"/comments/empty\"\n\t#[2022-05-02T10:11:12.0000008Z] template \"/comments/top\"\n\t#[2022-05-02T10:11:12.0000009Z] template \"/card/comment\"\n\t#[2022-05-02T10:11:12.000001Z] template \"/btn/mailme\"\n\t#[2022-05-02T10:11:12.0000011Z] stylesheet \"main\" generation\n\t\t<info>[2022-05-02T10:11:12.0000012Z] trait generation \"block/header\"\n\t\t<info>[2022-05-02T10:11:12.0000013Z] trait generation \"text/top-header\"\n"))

		Expect(r.HasErrors()).To(BeFalse())
		Expect(univ).NotTo(BeNil())
		Expect(univ.Stylesheets()).NotTo(BeNil())
		Expect(len(univ.Stylesheets())).To(Equal(1))
		Expect(univ.Stylesheets()["main"]).To(Equal("\n\n/* trait \"block/header\" */.top-header{\nbackground: #f5f6f7;\nmargin: 20px 0 20px 0;\npadding: 20px;\n}\n\n\n/* trait \"block/header\":hover */\n.top-header:hover{\nbackground: #e5e6e7;\n}\n\n\n/* trait \"text/top-header\" */.top-header-title{\ncolor: red;\nfont-family: Helvetica Neue Helvetica Arial sans-serif;\nfont-weight: 600;\n}\n"))
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
						"top-comments": map[string]interface{}{
							"comments": []map[string]interface{}{
								{
									"comment-author": "Sam",
									"comment-text":   "good article",
								},
								{
									"comment-author": "John",
									"comment-text":   "bullshit article",
								},
							},
						},
					},
					{
						"article-title":       "Article 2",
						"article-preview":     "Preview for article 2.",
						"article-link":        "http://yahoo.com",
						"article-link-anchor": "yahoo!",
					},
				}})
		Expect(report.ToString(r)).To(Equal("#[2022-05-02T10:11:12.0000012Z] rendering template \"/layout/test\"\n"))
		Expect(rendered).To(Equal("<!DOCTYPE html><html><head><title>::Test Template::</title><meta charset=\"utf-8\"/></head><body><header class=\"top-header\"><h1 class=\"top-header-title\">Test Header!</h1></header><div class=\"article-card\"><h1 class=\"article-card-title\">Article 1</h1><span class=\"article-card-preview\">Preview for article 1.</span><a class=\"article-card-link\" href=\"http://google.com\">google</a><div class=\"comment-card\"><span class=\"comment-card-author\">Sam</span><p class=\"comment-card-text\">good article</p></div><div class=\"comment-card\"><span class=\"comment-card-author\">John</span><p class=\"comment-card-text\">bullshit article</p></div></div><div class=\"article-card\"><h1 class=\"article-card-title\">Article 2</h1><span class=\"article-card-preview\">Preview for article 2.</span><a class=\"article-card-link\" href=\"http://yahoo.com\">yahoo!</a>no comments</div><a class=\"mailme-btn\" href=\"mailto:egotraumatic@example.com\">Mail me</a></body></html>"))
	})
})
