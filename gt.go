package gt

import (
	"fmt"
	"strings"

	"github.com/Contra-Culture/report"
)

type (
	// styling
	StyleTrait struct {
		name         string
		declarations []StyleDeclaration
	}
	StyleDeclaration []string
	StyleRule        struct {
		selectors    []string
		declarations []StyleDeclaration
	}
	Stylesheet struct {
		name  string
		rules []StyleRule
	}

	// limbo templating
	LimboTemplate struct {
		name    string
		content interface{}
	}
	Limbo struct {
		reportCreator func(string, ...interface{}) report.Node
		templates     []LimboTemplate
	}

	// universe templating
	Universe struct {
		reportCreator func(string, ...interface{}) report.Node
		templates     map[string]*Template
	}
	Template struct {
		name      string
		fragments []interface{}
	}

	// trees traversing
	iterator struct {
		cursor int
		path   []int
		label  string
		items  []interface{}
		params map[string]interface{} // only for rendering,
	}

	// rules
	doctype   string   // allows to place HTML5 doctype: <!DOCTYPE html>
	attribute struct { // allows to place attribute, works only as a child of tagAttributes rule
		name  string
		value string
	}
	attributeInjection struct { // allows to inject attribute value, works only as a child of tagAttributes rule
		name string
		key  string
	}
	tag struct { // allows to place an HTML tag
		name           string
		attributesRule TagAttributes
		contentRule    TagContent
	}
	TagAttributes  []interface{} // attribute or attributeInjection rules
	tagSelfClosing struct {      // nil, for "/>""
		selfClosing interface{}
	}
	tagClosing struct { // tag name for closing "</tag-name>"
		tag string
	}
	tagEnd struct { // nil, for ">"
		tagEnd interface{}
	}
	TagContent      []interface{} // text, textInjection, tag, templatePlacement, templateInjection, repeatable, variant rules
	documentContent []interface{} // same as tagContent, but allows doctype rule
	text            struct {      // text node represents exact text placement
		unsafe bool // text could be safe (HTML escape will be applied) or unsafe (text will be placed as is)
		text   string
	}
	textInjection struct { // allows to inject safe or unsafe text by the key on template rendering
		unsafe bool // text could be safe (HTML escape will be applied) or unsafe (text will be placed as is)
		key    string
	}
	templatePlacement struct { // allows to use other templates within the current one
		name string // template name
		key  string // params key for namespaceing to avoid naming conflicts for injections
	}
	templateInjection struct { // allows to inject template (place other template content on template rendering)
		key string
	}
	repeatable struct { // allows to repeat given rule rendering for N times, when N is a len() of a slice, provided through params on template rendering
		key  string
		rule interface{}
	}
	variant struct { // allows to place one or another of the predefined variants, for example: text or tag, depending on the key provided by params object on template rendering
		defaultRule interface{} // use __default key to provide data to the default rule. default rule is mandatory, but you can avoid rendering of anything with nothing rule.
		rules       map[string]interface{}
	}
	nothing struct { // allows to place nothing, makes sense only as a direct child of variant rule.
		nothing interface{}
	}
	jump struct { // allows to jump back to the next parrent's sibling, we use theEnd rule for both: template preparation and rendering
		iterator *iterator
	}
	theEnd struct { // signalizes about the end of rules tree traversing, we use theEnd rule for both: template preparation and rendering
		theEnd interface{}
	}
)

func Tag(n string, attrs TagAttributes, content TagContent) interface{} {
	return tag{
		name:           n,
		attributesRule: attrs,
		contentRule:    content,
	}
}
func Doctype() interface{} {
	return doctype(DOCTYPE)
}
func Attributes(attrs ...interface{}) TagAttributes {
	return TagAttributes(attrs)
}
func Content(content ...interface{}) TagContent {
	return TagContent(content)
}
func DocumentContent(content ...interface{}) interface{} {
	return documentContent(content)
}
func Attr(n, v string) interface{} {
	return attribute{
		name:  n,
		value: v,
	}
}
func AttrInjection(n, k string) interface{} {
	return attributeInjection{
		name: n,
		key:  k,
	}
}
func Text(t string) interface{} {
	return text{
		unsafe: false,
		text:   t,
	}
}
func UnsafeText(t string) interface{} {
	return text{
		unsafe: true,
		text:   t,
	}
}
func TextInj(k string) interface{} {
	return textInjection{
		unsafe: false,
		key:    k,
	}
}
func UnsafeTextInj(k string) interface{} {
	return textInjection{
		unsafe: true,
		key:    k,
	}
}
func TemplatePlacement(n, k string) interface{} {
	return templatePlacement{
		name: n,
		key:  k,
	}
}
func TemplateInjection(n, k string) interface{} {
	return templateInjection{
		key: k,
	}
}
func Repeat(k string, r interface{}) interface{} {
	return repeatable{
		key:  k,
		rule: r,
	}
}
func Variant(variants map[string]interface{}, dr interface{}) interface{} {
	return variant{
		defaultRule: dr,
		rules:       variants,
	}
}
func Nothing() interface{} {
	return nothing{}
}

// void elements
var selfClosingTags = []string{
	"area",
	"base",
	"br",
	"col",
	"embed",
	"hr",
	"img",
	"input",
	"link",
	"meta",
	"param",
	"source",
	"track",
	"wbr",
}

const DOCTYPE = "<!DOCTYPE html>"

func selfClosingTag(n string) bool {
	for _, t := range selfClosingTags {
		if t == n {
			return true
		}
	}
	return false
}
func newIterator(path []int, label string, items []interface{}) *iterator {
	return &iterator{
		cursor: -1,
		path:   path,
		label:  label,
		items:  items,
	}
}
func newIteratorWithParams(path []int, label string, items []interface{}, params map[string]interface{}) *iterator {
	return &iterator{
		cursor: -1,
		path:   path,
		label:  label,
		items:  items,
		params: params,
	}
}
func (iter *iterator) next() interface{} {
	iter.cursor = iter.cursor + 1
	fmt.Printf("\n\nnext(): %#v : %d/%d : %#v\n\n", iter.path, iter.cursor, len(iter.items), iter.items[iter.cursor])
	if iter.cursor < len(iter.items) {
		return iter.items[iter.cursor]
	}
	return nil
}

// New() creates new Limbo object for dirty templates spec.
func New(rc func(string, ...interface{}) report.Node) *Limbo {
	return &Limbo{reportCreator: rc}
}

func (l *Limbo) Template(n string, rule interface{}) {
	switch rule.(type) {
	case TagAttributes, TagContent, documentContent:
		l.templates = append(l.templates, LimboTemplate{
			name:    n,
			content: rule,
		})
	default:
		panic(fmt.Sprintf("wrong rule for template content: `%#v`", rule))
	}
}

// *Limbo.Universe() generates templating universe, which is the entity point to work with templates at the application runtime.
func (l *Limbo) Universe() (u *Universe, r report.Node) {
	r = l.reportCreator("universe")
	u = &Universe{
		templates:     make(map[string]*Template),
		reportCreator: l.reportCreator,
	}
	// go through limbo template to prepare final (universe) templates
	for _, lt := range l.templates {
		if _, exists := u.templates[lt.name]; exists {
			r.Error("template \"%s\" already specified", lt.name)
			continue
		}
		t := &Template{
			name: lt.name,
		}
		fragments := []interface{}{}
		ri := newIterator([]int{}, "top", []interface{}(lt.content.(documentContent)))
		traverse := true
		for traverse {
			rule := ri.next()
			//fragments = appendFragments(fragments, rule)
			switch fragment := rule.(type) {
			case theEnd:
				traverse = false // stops the loop because rules tree traversing is finished
			case jump:
				ri = fragment.iterator
				continue
			case doctype:
				fragments = appendFragments(fragments, DOCTYPE)
			case tag:
				fmt.Printf("\ntag: %s\n", fragment.name)
				fragments = appendFragments(fragments, fmt.Sprintf("<%s", fragment.name))
				// for tag we flatten attributes and content rule into a single list of rules
				// because of that tagAttributes and tagContent rules are ignored, but not their content.
				// this allows to make less jumps and gets in theory some performance improvement.
				rules := []interface{}{}
				if len(fragment.attributesRule) > 0 {
					rules = append(rules, fragment.attributesRule...)
				}
				if selfClosingTag(fragment.name) {
					rules = append(rules, tagSelfClosing{})
				} else {
					rules = append(rules, tagEnd{})
					if len(fragment.contentRule) > 0 {
						rules = append(rules, fragment.contentRule...)
					}
					rules = append(rules, tagClosing{fragment.name})
				}
				if len(ri.path) > 0 {
					rules = append(rules, jump{iterator: ri}) // allows to jump to the parrent's sibling at the end
				} else {
					rules = append(rules, theEnd{})
				}
				ri = newIterator(append(ri.path, ri.cursor), fmt.Sprintf("<%s>", fragment.name), rules)
			case tagEnd:
				fragments = appendFragments(fragments, ">")
			case tagClosing:
				fragments = appendFragments(fragments, fmt.Sprintf("</%s>", fragment.tag))
			case tagSelfClosing:
				fragments = appendFragments(fragments, "/>")
			case TagAttributes: // not achievable if tagAttributes is within tagRule because of flattening
				ri = newIterator(append(ri.path, ri.cursor), "attrs", append([]interface{}(fragment), jump{iterator: ri}))
			case TagContent: // not achievable if tagContent is within tagRule because of flattening
				ri = newIterator(append(ri.path, ri.cursor), "content", append([]interface{}(fragment), jump{iterator: ri}))
			case attribute:
				fragments = appendFragments(
					fragments,
					fmt.Sprintf(" %s=\"%s\"", fragment.name, fragment.value))
			case attributeInjection:
				fragments = appendFragments(
					fragments,
					fmt.Sprintf(" %s=\"", fragment.name),
					fragment, // attribute value injection
					"\"",
				)
			case text:
				var text = fragment.text
				if !fragment.unsafe {
					text = safeTextReplacer.Replace(text)
				}
				fragments = appendFragments(fragments, text)
			case textInjection:

			case templatePlacement:
				exists := false
				for _, lt := range l.templates {
					if lt.name == fragment.name {
						exists = true
						break
					}
				}
				if !exists {
					r.Error("template \"%s\" not defined", fragment.name)
					return nil, r
				}
				fragments = appendFragments(fragments, fragment)
			case templateInjection:
				fragments = appendFragments(fragments, fragment)
			case repeatable:
				fragments = appendFragments(fragments, fragment)
			case variant:
				fragments = appendFragments(fragments, fragment)
			case documentContent:
				ri = newIterator(append(ri.path, ri.cursor), "document content", []interface{}(fragment))
			default:
				r.Error("wrong rule %#v, iterator: %#v", rule, ri)
				return
			}
		}
		fragments = appendFragments(fragments, theEnd{})
		fmt.Printf("\n\n\n::fragments: %#v", fragments)
		t.fragments = fragments
		u.templates[t.name] = t
	}
	return
}

var safeTextReplacer = strings.NewReplacer("<", "&lt;", ">", "&gt;", "\"", "&quot", "'", "&quot")

func appendFragments(fragments []interface{}, newRawFragments ...interface{}) []interface{} {
	if len(fragments) == 0 {
		fragments = append(fragments, newRawFragments[0])
		newRawFragments = newRawFragments[1:]
	}
	for _, newRawFragment := range newRawFragments {
		lastIdx := len(fragments) - 1
		prevFragment := fragments[lastIdx]
		switch newFragment := newRawFragment.(type) {
		case string:
			switch prevValue := prevFragment.(type) {
			case string:
				fragments[lastIdx] = prevValue + newFragment
			default:
				fragments = append(fragments, newFragment)
			}
		default:
			fragments = append(fragments, newFragment)
		}
	}
	return fragments
}

func (u *Universe) Render(n string, params map[string]interface{}) (string, report.Node) {
	fmt.Printf("\n\n\n\nrendering template %s\n", n)

	r := u.reportCreator("rendering template \"%s\"", n)
	t, ok := u.templates[n]
	if !ok {
		r.Error("template not found")
		return "", r
	}
	var sb strings.Builder
	iter := newIteratorWithParams([]int{}, "top", t.fragments, params)
	traverse := true
	for traverse {
		rawFragment := iter.next()
		switch f := rawFragment.(type) {
		case theEnd:
			traverse = false
		case jump:
			iter = f.iterator
		case string:
			sb.WriteString(f)
		case templatePlacement:
			tPl := u.templates[f.name]
			iter = newIteratorWithParams(append(iter.path, iter.cursor), "template placement", append(tPl.fragments, jump{iterator: iter}), params)
		case templateInjection:
			_data, exists := params[f.key]
			if !exists {
				r.Error("template injection \"%s\" not provided", f.key)
				return "", r
			}
			data, ok := _data.(map[string]interface{})
			if !ok {
				r.Error("wrong format of template injection \"%s\" data", f.key)
				return "", r
			}
			_tn, ok := data["name"]
			if !ok {
				r.Error("template name for injection \"%s\" not provided", f.key)
				return "", r
			}
			tn, ok := _tn.(string)
			if !ok {
				r.Error("template name for injection \"%s\"should be string", f.key)
				return "", r
			}
			_injParams, ok := data["params"]
			if !ok {
				r.Error("template params for injection \"%s\" not provided", f.key)
				return "", r
			}
			injParams, ok := _injParams.(map[string]interface{})
			if !ok {
				r.Error("template params for injection \"%s\" should be map[string]interface{}", f.key)
			}
			injT, ok := u.templates[tn]
			if !ok {
				r.Error("template \"%s\" for injection doesn't exist", t.name)
				return "", r
			}
			iter = newIteratorWithParams(append(iter.path, iter.cursor), "template injection", append(injT.fragments, jump{iterator: iter}), injParams)
		case attributeInjection:
			_v, exists := params[f.key]
			if !exists {
				r.Error("attribute value injection \"%s\" not provided", f.key)
				return "", r
			}
			v, ok := _v.(string)
			if !ok {
				r.Error("text injection \"%s\"should be a string", f.key)
				return "", r
			}
			sb.WriteString(v)
		case textInjection:
			_v, exists := params[f.key]
			if !exists {
				r.Error("text injection \"%s\" not provided", f.key)
				return "", r
			}
			v, ok := _v.(string)
			if !ok {
				r.Error("text injection \"%s\"should be a string", f.key)
				return "", r
			}
			if !f.unsafe {
				v = safeTextReplacer.Replace(v)
			}
			sb.WriteString(v)
		}
	}
	return sb.String(), r
}
func (u *Universe) Stylesheets(n string) map[string]Stylesheet {
	return nil
}
