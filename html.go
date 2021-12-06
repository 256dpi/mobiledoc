package mobiledoc

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"net/url"
)

// HTMLRenderer implements a basic HTML renderer.
type HTMLRenderer struct {
	Atoms map[string]func(*bufio.Writer, string, Map) error
	Cards map[string]func(*bufio.Writer, Map) error
}

// NewHTMLRenderer creates a new HTMLRenderer.
func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{
		Atoms: make(map[string]func(*bufio.Writer, string, Map) error),
		Cards: make(map[string]func(*bufio.Writer, Map) error),
	}
}

// Render will render the document to the provided writer.
func (r *HTMLRenderer) Render(w io.Writer, doc Document) error {
	// wrap writer
	bw := bufio.NewWriter(w)

	// render sections
	for _, section := range doc.Sections {
		err := r.renderSection(bw, section)
		if err != nil {
			return err
		}
	}

	// flush buffer
	err := bw.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (r *HTMLRenderer) renderSection(w *bufio.Writer, section Section) error {
	// select sub renderer based on type
	switch section.Type {
	case MarkupSection:
		return r.renderMarkupSection(w, section)
	case ImageSection:
		return r.renderImageSection(w, section)
	case ListSection:
		return r.renderListSection(w, section)
	case CardSection:
		return r.renderCardSection(w, section)
	}

	return nil
}

func (r *HTMLRenderer) renderMarkupSection(w *bufio.Writer, section Section) error {
	// write open tag
	_, err := w.WriteString(fmt.Sprintf("<%s>", section.Tag))
	if err != nil {
		return err
	}

	// render markers
	err = r.renderMarkers(w, section.Markers)
	if err != nil {
		return err
	}

	// write close tag
	_, err = w.WriteString(fmt.Sprintf("</%s>", section.Tag))
	if err != nil {
		return err
	}

	return nil
}

func (r *HTMLRenderer) renderImageSection(w *bufio.Writer, section Section) error {
	// parse url
	src, err := url.Parse(section.Source)
	if err != nil {
		return err
	}

	// write tag
	_, err = w.WriteString(fmt.Sprintf("<img src=\"%s\">", src.String()))
	if err != nil {
		return err
	}

	return nil
}

func (r *HTMLRenderer) renderListSection(w *bufio.Writer, section Section) error {
	// write open tag
	_, err := w.WriteString(fmt.Sprintf("<%s>", section.Tag))
	if err != nil {
		return err
	}

	// write all items
	for _, item := range section.Items {
		// write open tag
		_, err := w.WriteString("<li>")
		if err != nil {
			return err
		}

		// render markers
		err = r.renderMarkers(w, item)
		if err != nil {
			return err
		}

		// write close tag
		_, err = w.WriteString("</li>")
		if err != nil {
			return err
		}
	}

	// write close tag
	_, err = w.WriteString(fmt.Sprintf("</%s>", section.Tag))
	if err != nil {
		return err
	}

	return nil
}

func (r *HTMLRenderer) renderCardSection(w *bufio.Writer, section Section) error {
	// get card renderer
	renderer, ok := r.Cards[section.Card.Name]
	if !ok {
		return fmt.Errorf("missing card renderer")
	}

	// call renderer
	err := renderer(w, section.Card.Payload)
	if err != nil {
		return err
	}

	return nil
}

func (r *HTMLRenderer) renderMarkers(w *bufio.Writer, markers []Marker) error {
	// prepare stack
	stack := markupStack{}

	// write all markers
	for _, marker := range markers {
		// write opening markups
		for _, markup := range marker.OpenMarkups {
			// begin tag
			_, err := w.WriteString(fmt.Sprintf("<%s", markup.Tag))
			if err != nil {
				return err
			}

			// write attributes
			for key, value := range markup.Attributes {
				_, err = w.WriteString(fmt.Sprintf(" %s=\"%s\"", key, value))
				if err != nil {
					return err
				}
			}

			// close tag
			_, err = w.WriteString(">")
			if err != nil {
				return err
			}

			// push markup
			stack.push(markup)
		}

		// write marker
		switch marker.Type {
		case TextMarker:
			// write text
			_, err := w.WriteString(html.EscapeString(marker.Text))
			if err != nil {
				return err
			}
		case AtomMarker:
			// get renderer
			renderer, ok := r.Atoms[marker.Atom.Name]
			if !ok {
				return fmt.Errorf("missing atom renderer")
			}

			// call renderer
			err := renderer(w, marker.Atom.Text, marker.Atom.Payload)
			if err != nil {
				return err
			}
		}

		// close markups
		for i := 0; i < marker.ClosedMarkups; i++ {
			// get markup
			markup := stack.pop()

			// write closing tag
			_, err := w.WriteString(fmt.Sprintf("</%s>", markup.Tag))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type markupStack struct {
	list []*Markup
}

func (s *markupStack) push(m *Markup) {
	s.list = append(s.list, m)
}

func (s *markupStack) pop() *Markup {
	item := s.list[len(s.list)-1]
	s.list = s.list[0 : len(s.list)-1]
	return item
}
