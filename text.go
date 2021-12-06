package mobiledoc

import (
	"bufio"
	"fmt"
	"io"
)

// TextRenderer implements a basic text renderer.
type TextRenderer struct {
	Atoms map[string]func(*bufio.Writer, string, Map) error
	Cards map[string]func(*bufio.Writer, Map) error
}

// NewTextRenderer creates a new TextRenderer.
func NewTextRenderer() *TextRenderer {
	return &TextRenderer{
		Atoms: make(map[string]func(*bufio.Writer, string, Map) error),
		Cards: make(map[string]func(*bufio.Writer, Map) error),
	}
}

// Render will render the document to the provided writer.
func (r *TextRenderer) Render(w io.Writer, doc Document) error {
	// wrap writer
	bw := bufio.NewWriter(w)

	// render sections
	for i, section := range doc.Sections {
		err := r.renderSection(bw, section)
		if err != nil {
			return err
		}

		// write newline
		if i < len(doc.Sections)-1 {
			_, err = bw.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}

	// flush buffer
	err := bw.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (r *TextRenderer) renderSection(w *bufio.Writer, section Section) error {
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

func (r *TextRenderer) renderMarkupSection(w *bufio.Writer, section Section) error {
	// render markers
	err := r.renderMarkers(w, section.Markers)
	if err != nil {
		return err
	}

	return nil
}

func (r *TextRenderer) renderImageSection(w *bufio.Writer, section Section) error {
	// write source
	_, err := w.WriteString(fmt.Sprintf("[%s]", section.Source))
	if err != nil {
		return err
	}

	return nil
}

func (r *TextRenderer) renderListSection(w *bufio.Writer, section Section) error {
	// write all items
	for i, item := range section.Items {
		// write dash
		if section.Tag == "ol" {
			_, err := w.WriteString(fmt.Sprintf("%d. ", i+1))
			if err != nil {
				return err
			}
		} else {
			_, err := w.WriteString("- ")
			if err != nil {
				return err
			}
		}

		// render markers
		err := r.renderMarkers(w, item)
		if err != nil {
			return err
		}

		if i < len(section.Items)-1 {
			_, err := w.WriteString("\n")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *TextRenderer) renderCardSection(w *bufio.Writer, section Section) error {
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

func (r *TextRenderer) renderMarkers(w *bufio.Writer, markers []Marker) error {
	// write all markers
	for i, marker := range markers {
		// write marker
		switch marker.Type {
		case TextMarker:
			// write text
			_, err := w.WriteString(marker.Text)
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

		// write space
		if i < len(markers)-1 {
			_, err := w.WriteString(" ")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
