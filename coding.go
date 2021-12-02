package mobiledoc

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
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

// UnmarshalBSON implements the bson.Unmarshaler interface.
func (d *Document) UnmarshalBSON(bytes []byte) error {
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

// MarshalBSON implements the bson.Marshaler interface.
func (d Document) MarshalBSON() ([]byte, error) {
	// compile document
	raw, err := Compile(d)
	if err != nil {
		return nil, err
	}

	// marshal map
	bytes, err := bson.Marshal(raw)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
