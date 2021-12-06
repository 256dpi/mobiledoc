package mobiledoc

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Document) UnmarshalJSON(bytes []byte) error {
	// unmarshal to map
	var m Map
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}

	// parse map
	doc, err := Parse(m)
	if err != nil {
		return err
	}

	// validate document
	err = formatValidator.Validate(doc)
	if err != nil {
		return err
	}

	// set document
	*d = doc

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (d Document) MarshalJSON() ([]byte, error) {
	// compile document
	m, err := Compile(d)
	if err != nil {
		return nil, err
	}

	// marshal map
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (d *Document) UnmarshalBSONValue(typ bsontype.Type, bytes []byte) error {
	// check type
	if typ == bson.TypeNull {
		*d = Document{}
		return nil
	} else if typ != 0 && typ != bson.TypeEmbeddedDocument {
		return fmt.Errorf("unexpected type: %s", typ.String())
	}

	// unmarshal to map
	var m Map
	err := bson.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}

	// parse map
	doc, err := Parse(m)
	if err != nil {
		return err
	}

	// validate document
	err = formatValidator.Validate(doc)
	if err != nil {
		return err
	}

	// set document
	*d = doc

	return nil
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (d *Document) MarshalBSONValue() (bsontype.Type, []byte, error) {
	// handle nil
	if d == nil {
		return bsontype.Null, nil, nil
	}

	// handle zero
	if isZero(*d) {
		return bson.MarshalValue(Map{})
	}

	// compile document
	raw, err := Compile(*d)
	if err != nil {
		return 0, nil, err
	}

	// marshal map
	val, bytes, err := bson.MarshalValue(raw)
	if err != nil {
		return 0, nil, err
	}

	return val, bytes, nil
}

func isZero(d Document) bool {
	return d.Version == "" &&
		len(d.Markups) == 0 &&
		len(d.Atoms) == 0 &&
		len(d.Cards) == 0 &&
		len(d.Sections) == 0
}
