package gt

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Contra-Culture/report"
)

type (
	// styling
	SelectorInjection struct {
		Name string
	}
	StylingTemplate struct {
		selectorGenerators map[string]func(map[string]string) (string, error) // rule template name -> selector generator
		blocks             map[string]StyleBlock                              // rule template name ->  block of rule's style declarations
		name               string
	}
	StylingTemplateRule struct {
		stylingTemplate StylingTemplate
		selectors       map[string][]string // rule template name -> selectors
	}
	Stylesheet struct {
		rn                   report.Node
		stylingTemplateRules map[string]StylingTemplateRule
		predefined           string
	}
	StyleBlock [][]string // represents a multiple CSS declarations,
	// limbo templating
	LimboTemplate struct {
		name           string
		stylesheetName string
		content        interface{}
		rn             report.Node
	}
	Limbo struct {
		reportCreator    func(string, ...interface{}) report.Node
		rn               report.Node
		templates        []LimboTemplate
		stylesheets      map[string]Stylesheet
		stylingTemplates map[string]StylingTemplate
	}
	// universe templating
	Universe struct {
		reportCreator func(string, ...interface{}) report.Node
		templates     map[string]*Template
		stylesheets   map[string]string
	}
	Template struct {
		name      string
		fragments []interface{}
	}

	// trees traversing
	iterator struct {
		cursor     int
		path       []int
		label      string
		items      []interface{}
		paramsType string      // only for rendering,
		params     interface{} // []map[string]interface{} or map[string]interface{} // only for rendering,
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
	class struct {
		name                              string // class base name
		stylingTemplateName               string
		stylingTemplateSelectorInjections map[string]string
	}
	variant struct { // allows to place one or another of the predefined variants, for example: text or tag, depending on the key provided by params object on template rendering
		defaultTemplateName string // use __default key to provide data to the default rule. default rule is mandatory, but you can avoid rendering of anything with nothing rule.
		templates           map[string]string
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

const auto = "__auto__"

// returns rule template name and selector generator for that rule template
func ruleTemplateNameAndSelectorGenerator(template []interface{}) (string, func(map[string]string) (string, error)) {
	var sb strings.Builder
	for _, _f := range template {
		switch f := _f.(type) {
		case string:
			sb.WriteString(f)
		case SelectorInjection:
			sb.WriteString("{{")
			sb.WriteString(f.Name)
			sb.WriteString("}}")
		}
	}
	name := sb.String()
	fmt.Printf("debug ruleTemplateNameAndSelectorGenerator: name %s", name)
	selectorGenerator := func(injections map[string]string) (string, error) {
		_selector := []string{}
		for _, _f := range template {
			switch f := _f.(type) {
			case string:
				_selector = append(_selector, f)
			case SelectorInjection:
				inj, exists := injections[f.Name]
				if !exists {
					return "", fmt.Errorf("selector injection \"%s\" not provided", f.Name)
				}
				_selector = append(_selector, inj)
			}
		}
		selector:=strings.Join(_selector, " ")
		fmt.Printf("debug ruleTemplateNameAndSelectorGenerator: (%s) selector generator: %s", name, selector, " "))
		return selector , nil
	}
	return name, selectorGenerator
}
func Auto() string { // is for templates that have no key for params, but the params will be passed automatically (with Repeat() for example)
	return auto
}
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
func AttrInjection(name, key string) interface{} {
	return attributeInjection{
		name,
		key,
	}
}

const SELF_CLASS_PLACEMENT = "selfClass"

func Class(name string, stylingTemplateName string, styleTemplateSelectorInjections map[string]string) interface{} {
	if name[0] != '.' {
		name = "." + name
	}
	if styleTemplateSelectorInjections == nil {
		styleTemplateSelectorInjections = map[string]string{}
	}
	styleTemplateSelectorInjections[SELF_CLASS_PLACEMENT] = name
	return class{
		name,
		stylingTemplateName,
		styleTemplateSelectorInjections,
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
func TemplateInjection(k string) interface{} {
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
func Variant(dr string, variants map[string]string) interface{} {
	return variant{
		defaultTemplateName: dr,
		templates:           variants,
	}
}
func Nothing() interface{} {
	return nothing{nothing: true}
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
const paramsTypeMap = "map"
const paramsTypeSlice = "slice"

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
func newIteratorWithParamsMap(path []int, label string, items []interface{}, params map[string]interface{}) *iterator {
	return &iterator{
		cursor:     -1,
		path:       path,
		label:      label,
		items:      items,
		params:     params,
		paramsType: paramsTypeMap,
	}
}
func newIteratorWithParamsSlice(path []int, label string, items []interface{}, params []map[string]interface{}) *iterator {
	return &iterator{
		cursor:     -1,
		path:       path,
		label:      label,
		items:      items,
		params:     params,
		paramsType: paramsTypeSlice,
	}
}
func (iter *iterator) next() interface{} {
	iter.cursor = iter.cursor + 1
	if iter.cursor < len(iter.items) {
		return iter.items[iter.cursor]
	}
	return nil
}
func (iter *iterator) getParams() map[string]interface{} {
	switch iter.paramsType {
	case paramsTypeMap:
		return iter.params.(map[string]interface{})
	case paramsTypeSlice:
		return iter.params.([]map[string]interface{})[iter.cursor]
	default:
		panic("wrong params type") // can't occur
	}
}

// New() creates new Limbo object for dirty templates spec.
func New(rc func(string, ...interface{}) report.Node) *Limbo {
	return &Limbo{
		reportCreator:    rc,
		rn:               rc("limbo"),
		stylingTemplates: make(map[string]StylingTemplate),
		stylesheets:      make(map[string]Stylesheet),
	}
}
func WithLayout(content ...interface{}) func(*LimboTemplate) bool {
	return func(t *LimboTemplate) bool {
		if t.content != nil {
			t.rn.Error("templete content already specified")
			return false
		}
		t.content = documentContent(append(content, theEnd{theEnd: true}))
		return true
	}
}
func WithContent(content ...interface{}) func(*LimboTemplate) bool {
	return func(t *LimboTemplate) bool {
		if t.content != nil {
			t.rn.Error("templete content already specified")
			return false
		}
		t.content = TagContent(append(content, theEnd{theEnd: true}))
		return true
	}
}
func WithAttributes(attrs ...interface{}) func(*LimboTemplate) bool {
	return func(t *LimboTemplate) bool {
		if t.content != nil {
			t.rn.Error("templete content already specified")
			return false
		}
		t.content = TagAttributes(append(attrs, theEnd{theEnd: true}))
		return true
	}
}
func WithStylesheet(n string) func(*LimboTemplate) bool {
	return func(t *LimboTemplate) bool {
		if len(t.stylesheetName) > 0 {
			t.rn.Error("template stylesheet already specified")
			return false
		}
		t.stylesheetName = n
		return true
	}
}
func StylingRule(selectorTemplate []interface{}, block [][]string) func(*StylingTemplate) {
	return func(styleTemplate *StylingTemplate) {
		ruleTemplateName, selectorGenerator := ruleTemplateNameAndSelectorGenerator(selectorTemplate)
		styleTemplate.selectorGenerators[ruleTemplateName] = selectorGenerator
		styleTemplate.blocks[ruleTemplateName] = block
	}
}

// Defines a stylesheet
func (l *Limbo) Stylesheet(name string, opts ...func(*Stylesheet)) {
	s := Stylesheet{
		rn:                   l.rn.Structure("stylesheet \"%s\"", name),
		stylingTemplateRules: map[string]StylingTemplateRule{},
	}
	l.stylesheets[name] = s
	for _, opt := range opts {
		opt(&s)
	}
}
func CSSRule(selectors []string, block [][]string) func(*Stylesheet) {
	var sb strings.Builder
	sb.WriteString(strings.Join(selectors, ", "))
	sb.WriteString(" {\n")
	for _, declaration := range block {
		sb.WriteString(declaration[0])
		sb.WriteString(": ")
		sb.WriteString(strings.Join(declaration[1:], ", "))
		sb.WriteString(";\n")
	}
	sb.WriteString("}\n\n")
	css := sb.String()
	return func(s *Stylesheet) {
		fmt.Printf("debug CSSRule before: %s", s.predefined)
		s.predefined = s.predefined + css
		fmt.Printf("debug CSSRule after: %s", s.predefined)
	}
}
func Styling(name string, rules ...func(*StylingTemplate)) func(*Stylesheet) {
	t := StylingTemplate{
		name:               name,
		selectorGenerators: map[string]func(map[string]string) (string, error){},
		blocks:             map[string]StyleBlock{},
	}
	for _, rule := range rules {
		rule(&t)
	}
	return func(s *Stylesheet) {
		if _, exists := s.stylingTemplateRules[name]; exists {
			s.rn.Error("styling template \"%s\" already exists", name)
			return
		}
		s.stylingTemplateRules[name] = StylingTemplateRule{
			stylingTemplate: t,
			selectors:       map[string][]string{},
		}
	}
}
func (l *Limbo) Template(n string, opts ...func(*LimboTemplate) bool) {
	t := LimboTemplate{
		name: n,
		rn:   l.rn.Structure("template \"%s\"", n),
	}
	for _, opt := range opts {
		ok := opt(&t)
		if !ok {
			return
		}
	}
	if t.content == nil {
		t.rn.Error("template content should be specified")
		return
	}
	sname := t.stylesheetName
	if !(len(sname) > 0) {
		t.rn.Error("template stylesheet should be specified")
		return
	}
	_, exists := l.stylesheets[sname]
	if !exists {
		l.stylesheets[sname] = Stylesheet{
			stylingTemplateRules: map[string]StylingTemplateRule{},
		}
	}
	l.templates = append(l.templates, t)
}

// *Limbo.Universe() generates templating universe, which is the entity point to work with templates at the application runtime.
func (l *Limbo) Universe() (u *Universe, r report.Node) {
	r = l.rn
	u = &Universe{
		templates:     make(map[string]*Template),
		reportCreator: l.reportCreator,
		stylesheets:   map[string]string{},
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
		var topContent []interface{}
		switch rawTopContent := lt.content.(type) {
		case documentContent:
			topContent = []interface{}(rawTopContent)
			// todo: add check for doctype
		case TagContent:
			// todo: add check for only tag content (no doctype or attributes)
			topContent = []interface{}(rawTopContent)
		case TagAttributes:
			// todo: add check for only tag attribures (no doctype or tag content)
			topContent = []interface{}(rawTopContent)
		default:
			r.Error("wrong type of top content rule, expected: documentContent, TagContent, TagAttributes %#v", fragments)
			return nil, r
		}
		iter := newIterator([]int{}, "top", topContent)
		traverse := true
		for traverse {
			rule := iter.next()
			switch fragment := rule.(type) {
			case theEnd:
				traverse = false // stops the loop because rules tree traversing is finished
			case jump:
				iter = fragment.iterator
				continue
			case doctype:
				fragments = appendFragments(fragments, DOCTYPE)
			case tag:
				fragments = appendFragments(fragments, fmt.Sprintf("<%s", fragment.name))
				// for tag we flatten attributes and content rule into a single list of rules
				// because of that tagAttributes and tagContent rules are ignored, but not their content.
				// this allows to make less jumps and gets in theory some performance improvement.
				rules := []interface{}{}
				if len(fragment.attributesRule) > 0 {
					rules = append(rules, fragment.attributesRule...)
				}
				if selfClosingTag(fragment.name) {
					rules = append(rules, tagSelfClosing{selfClosing: true})
				} else {
					rules = append(rules, tagEnd{tagEnd: true})
					if len(fragment.contentRule) > 0 {
						rules = append(rules, fragment.contentRule...)
					}
					rules = append(rules, tagClosing{fragment.name})
				}
				rules = append(rules, jump{iterator: iter}) // allows to jump to the parrent's sibling at the end
				iter = newIterator(append(iter.path, iter.cursor), fmt.Sprintf("<%s>", fragment.name), rules)
			case tagEnd:
				fragments = appendFragments(fragments, ">")
			case tagClosing:
				fragments = appendFragments(fragments, fmt.Sprintf("</%s>", fragment.tag))
			case tagSelfClosing:
				fragments = appendFragments(fragments, "/>")
			case TagAttributes: // not achievable if tagAttributes is within tagRule because of flattening
				iter = newIterator(append(iter.path, iter.cursor), "attrs", []interface{}(fragment))
			case TagContent: // not achievable if tagContent is within tagRule because of flattening
				iter = newIterator(append(iter.path, iter.cursor), "content", []interface{}(fragment))
			case attribute:
				fragments = appendFragments(
					fragments,
					fmt.Sprintf(" %s=\"%s\"", fragment.name, fragment.value))
			case class:
				fragments = appendFragments(fragments, fmt.Sprintf(" class=\"%s\"", fragment.name))
				s := l.stylesheets[lt.stylesheetName]
				stylingTemplateRule, exists := s.stylingTemplateRules[fragment.stylingTemplateName]
				if !exists {
					r.Error("styling template \"%s\" does not exist", fragment.stylingTemplateName)
					return
				}
				for ruleTemplateName, selectorGenerator := range stylingTemplateRule.stylingTemplate.selectorGenerators {
					selector, err := selectorGenerator(fragment.stylingTemplateSelectorInjections)
					if err != nil {
						r.Error(err.Error())
					}
					stylingTemplateRule.selectors[ruleTemplateName] = append(stylingTemplateRule.selectors[ruleTemplateName], selector)
				}
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
				fragments = appendFragments(fragments, fragment)
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
				iter = newIterator(append(iter.path, iter.cursor), "document content", []interface{}(fragment))
			default:
				r.Error("wrong rule %#v, iterator: %#v", rule, iter)
				return
			}
		}
		fragments = appendFragments(fragments, theEnd{})
		t.fragments = fragments
		u.templates[t.name] = t
	}
	for n, stylesheet := range l.stylesheets {
		sr := r.Structure("stylesheet \"%s\" generation", n)
		var sb strings.Builder
		sb.WriteString(stylesheet.predefined)
		// ordering styling template rules by styling template name
		stylingTemplateNames := []string{}
		for stylingTemplateName := range stylesheet.stylingTemplateRules {
			stylingTemplateNames = append(stylingTemplateNames, stylingTemplateName)
		}
		sort.Strings(stylingTemplateNames)
		for _, stylingTemplateName := range stylingTemplateNames {
			stylingTemplateRule := stylesheet.stylingTemplateRules[stylingTemplateName]
			if len(stylingTemplateRule.selectors) < 1 {
				r.Warn("styling template \"has no use cases\"", stylingTemplateName)
				continue
			}
			sr.Info("styling template rule generation \"%s\"", stylingTemplateName)
			sb.WriteString("\n\n/* styling Template \"")
			sb.WriteString(stylingTemplateName)
			sb.WriteString("\" */\n")
			for ruleTemplateName, selectors := range stylingTemplateRule.selectors {
				sb.WriteString("/*   rule:")
				sb.WriteString(ruleTemplateName)
				sb.WriteString(" */\n")
				sb.WriteString(strings.Join(selectors, ", "))
				sb.WriteString(" {\n")
				block := stylingTemplateRule.stylingTemplate.blocks[ruleTemplateName]
				for _, declaration := range block {
					sb.WriteRune('\t')
					sb.WriteString(declaration[0])
					sb.WriteString(": ")
					sb.WriteString(strings.Join(declaration[1:], ", "))
					sb.WriteString(";\n")
				}
				sb.WriteString("}\n")
			}
		}
		u.stylesheets[n] = sb.String()
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
	r := u.reportCreator("rendering template \"%s\"", n)
	t, ok := u.templates[n]
	if !ok {
		r.Error("template \"%s\" not found", n)
		return "", r
	}
	var sb strings.Builder
	iter := newIteratorWithParamsMap([]int{}, "top", t.fragments, params)
	traverse := true
traverseLoop:
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
			tPl, exists := u.templates[f.name]
			if !exists {
				r.Error("template \"%s\" not found", f.name)
				return "", r
			}
			if f.key == auto { // when rendering within repeatable rule
				iter = newIteratorWithParamsMap(
					append(iter.path, iter.cursor),
					"template placement",
					append(tPl.fragments[:len(tPl.fragments)-1], jump{iterator: iter}),
					iter.getParams())
			} else {
				iter = newIteratorWithParamsMap(
					append(iter.path, iter.cursor),
					"template placement",
					append(tPl.fragments[:len(tPl.fragments)-1], jump{iterator: iter}),
					iter.getParams()[f.key].(map[string]interface{}))
			}
		case templateInjection:
			_data, exists := params[f.key]
			if !exists {
				r.Error("template injection \"%s\" not provided", f.key)
				return "", r
			}
			data, ok := _data.(map[string]interface{})
			if !ok {
				r.Error("wrong format of template injection \"%s\" data, expected map[string]interface{}, got: %#v", f.key, _data)
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
				r.Error("template \"%s\" for injection doesn't exist", tn)
				return "", r
			}
			iter = newIteratorWithParamsMap(
				append(iter.path, iter.cursor),
				"template injection",
				append(injT.fragments[:len(injT.fragments)-1], jump{iterator: iter}),
				injParams)
		case attributeInjection:
			_v, exists := iter.getParams()[f.key]
			if !exists {
				r.Error("attribute value injection \"%s\" not provided %#v", f.key)
				return "", r
			}
			v, ok := _v.(string)
			if !ok {
				r.Error("text injection \"%s\"should be a string", f.key)
				return "", r
			}
			sb.WriteString(v)
		case textInjection:
			_v, exists := iter.getParams()[f.key]
			if !exists {
				r.Error("text injection \"%s\" not provided (params: %#v)", f.key, iter.getParams())
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
		case repeatable:
			rawRepParams, ok := iter.getParams()[f.key]
			if !ok {
				r.Error("repeatable params \"%s\" are not provided", f.key)
				return "", r
			}
			repParams, ok := rawRepParams.([]map[string]interface{})
			if !ok {
				r.Error("repeatable params should be of type []map[string]interface{}")
				return "", nil
			}
			rules := []interface{}{}
			for range repParams {
				rules = append(rules, f.rule)
			}
			rules = append(rules, jump{iterator: iter})
			iter = newIteratorWithParamsSlice(
				append(iter.path, iter.cursor),
				"repeatable",
				rules,
				repParams)
		case variant:
			for k, n := range f.templates {
				if _, ok := iter.getParams()[k]; ok {
					iter = newIteratorWithParamsMap(
						append(iter.path, 0),
						"variant",
						[]interface{}{
							templatePlacement{name: n, key: k},
							jump{iterator: iter},
						},
						iter.getParams())
					continue traverseLoop
				}
			}
			iter = newIteratorWithParamsMap(
				append(iter.path, 0),
				"variant",
				[]interface{}{
					templatePlacement{name: f.defaultTemplateName, key: Auto()},
					jump{iterator: iter},
				},
				map[string]interface{}{})
		default:
			r.Error("wrong type of fragment %#v", f)
			return "", r
		}
	}
	return sb.String(), r
}
func (u *Universe) Stylesheets() map[string]string {
	return u.stylesheets
}
