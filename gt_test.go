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
		limbo.Stylesheet(
			"main",
			CSSRule(
				[]string{"*"},
				[][]string{
					{"margin", "0"},
					{"padding", "0"},
					{"font-size", "16px"}}),
			Styling(
				"layout/header",
				StylingRule(
					[]interface{}{SelectorInjection{Name: SELF_CLASS_PLACEMENT}},
					[][]string{
						{"border", "1px solid black"},
						{"padding", "1rem"}}),
				StylingRule(
					[]interface{}{SelectorInjection{Name: SELF_CLASS_PLACEMENT}, "> h1"},
					[][]string{
						{"font-size", "2rem"},
						{"font-weight", "600"},
						{"color", "#454647"}}),
			))

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
										Class("top-header", "layout/header", nil)),
									Content(
										Tag(
											"h1",
											Attributes(),
											Content(Text("Test Header!"))))),
								Repeat(
									"articles",
									TemplatePlacement("/card/article", Auto())),
								TemplateInjection("bottom")))))))

		limbo.Template(
			"/card/article",
			WithStylesheet("main"),
			WithContent(
				Tag(
					"div",
					Attributes(
						Attr("class", "article-card")),
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
		Expect(report.ToString(r)).To(Equal("#[2022-05-02T10:11:12.0000001Z] limbo\n\t#[2022-05-02T10:11:12.0000002Z] stylesheet \"main\"\n\t#[2022-05-02T10:11:12.0000003Z] template \"/layout/test\"\n\t#[2022-05-02T10:11:12.0000004Z] template \"/card/article\"\n\t#[2022-05-02T10:11:12.0000005Z] template \"/comments/empty\"\n\t#[2022-05-02T10:11:12.0000006Z] template \"/comments/top\"\n\t#[2022-05-02T10:11:12.0000007Z] template \"/card/comment\"\n\t#[2022-05-02T10:11:12.0000008Z] template \"/btn/mailme\"\n\t#[2022-05-02T10:11:12.0000009Z] stylesheet \"main\" generation\n\t\t<info>[2022-05-02T10:11:12.000001Z] styling template rule generation \"layout/header\"\n"))

		Expect(r.HasErrors()).To(BeFalse())
		Expect(univ).NotTo(BeNil())
		Expect(univ.Stylesheets()).NotTo(BeNil())
		Expect(len(univ.Stylesheets())).To(Equal(1))
		Expect(univ.Stylesheets()["main"]).To(Equal(""))
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
