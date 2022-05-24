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
											Content(Text("Test Header!")))))))))))
		// creates universe
		univ, r := limbo.Universe()
		Expect(report.ToString(r)).To(Equal("#[2022-05-02T10:11:12.0000001Z] universe\n"))
		Expect(r.HasErrors()).To(BeFalse())
		Expect(univ).NotTo(BeNil())
		rendered, r := univ.Render(
			"/layout/test",
			map[string]interface{}{
				"top-header-title": "Top Header Title",
			})
		Expect(report.ToString(r)).To(Equal("#[2022-05-02T10:11:12.0000002Z] rendering template \"/layout/test\"\n"))
		Expect(rendered).To(Equal("<!DOCTYPE html><html><head><title>::Test Template::</title><meta charset=\"utf-8\"/></head><body><header class=\"top-header\"><h1 class=\"top-header-title\">Test Header!</h1></header></body></html>"))
	})
})
