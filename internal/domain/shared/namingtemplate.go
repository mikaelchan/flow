package shared

import (
	"errors"
	"strings"
)

var (
	variableNames     = map[string]struct{}{"title": {}, "year": {}, "season": {}, "episode": {}, "artist": {}, "album": {}, "track": {}, "bango": {}}
	invalidCharacters = map[byte]struct{}{'\\': {}, '/': {}, '\t': {}, '\n': {}, '\r': {}}
)

// NamingTemplate is a template for naming files in a library, formatted with the following variables:
// {title} - the title of the media
// {year} - the year of the media
// {season} - the season of the media
// {episode} - the episode of the media
// {artist} - the artist of the media
// {album} - the album of the media
// {track} - the track of the media
// {bango} - the bango of the media, for japanese media
type NamingTemplate string

func (nt NamingTemplate) String() string {
	return string(nt)
}

// Parser is a parser for NamingTemplate
type NamingTemplateParser struct {
	template     NamingTemplate // "x{title}y{year}z"
	placeholders map[int]string // {1: "title", 2: "year"}
	tokens       []string       // ["x", "y", "z"]
	currentIndex int
	tokenIndex   int
}

func NewNamingTemplateParser(template NamingTemplate) (*NamingTemplateParser, error) {
	parser := &NamingTemplateParser{template: template, placeholders: map[int]string{}, tokens: []string{}}
	if err := parser.parse(); err != nil {
		return nil, err
	}
	return parser, nil
}

func (p *NamingTemplateParser) parse() error {
	for p.currentIndex < len(p.template) {
		switch char := p.template[p.currentIndex]; char {
		case '{':
			p.currentIndex++
			placeholder, err := p.parsePlaceholder()
			if err != nil {
				return err
			}
			p.placeholders[p.tokenIndex] = placeholder
			p.tokenIndex++
			p.tokens = append(p.tokens, "")
		case '}':
			p.currentIndex++
		default:
			token, err := p.parseToken()
			if err != nil {
				return err
			}
			p.tokens = append(p.tokens, token)
			p.tokenIndex++
		}
	}
	return nil
}

func (p *NamingTemplateParser) parsePlaceholder() (string, error) {
	bytes := []byte{}
	for p.currentIndex < len(p.template) {
		if p.template[p.currentIndex] == '}' {
			break
		}
		bytes = append(bytes, p.template[p.currentIndex])
		p.currentIndex++
	}
	if p.currentIndex >= len(p.template) {
		return "", errors.New("unterminated placeholder")
	}
	placeholder := string(bytes)
	if _, ok := variableNames[placeholder]; !ok {
		return "", errors.New("invalid placeholder")
	}
	return placeholder, nil
}

func (p *NamingTemplateParser) parseToken() (string, error) {
	bytes := []byte{}
	for p.currentIndex < len(p.template) {
		if p.template[p.currentIndex] == '{' {
			break
		}
		if _, ok := invalidCharacters[p.template[p.currentIndex]]; ok {
			return "", errors.New("invalid character")
		}
		bytes = append(bytes, p.template[p.currentIndex])
		p.currentIndex++
	}
	return string(bytes), nil
}

func (p *NamingTemplateParser) Generate(variables map[string]string) string {
	var result []string
	for i := range p.tokens {
		if placeholder, ok := p.placeholders[i]; ok {
			result = append(result, variables[placeholder])
		} else {
			result = append(result, p.tokens[i])
		}
	}
	return strings.Join(result, "")
}
