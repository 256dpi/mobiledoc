package mobiledoc

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Document) UnmarshalJSON(bytes []byte) error {
	var m Map
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}
	doc, err := Parse(m)
	if err != nil {
		return err
	}
	err = formatValidator.Validate(doc)
	if err != nil {
		return err
	}
	*d = doc
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (d Document) MarshalJSON() ([]byte, error) {
	raw, err := Compile(d)
	if err != nil {
		return nil, err
	}
	return json.Marshal(raw)
}

// UnmarshalBSON implements the bson.Unmarshaler interface.
func (d *Document) UnmarshalBSON(bytes []byte) error {
	var m Map
	err := bson.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}
	doc, err := Parse(m)
	if err != nil {
		return err
	}
	err = formatValidator.Validate(doc)
	if err != nil {
		return err
	}
	*d = doc
	return nil
}

// MarshalBSON implements the bson.Marshaler interface.
func (d Document) MarshalBSON() ([]byte, error) {
	raw, err := Compile(d)
	if err != nil {
		return nil, err
	}
	return bson.Marshal(raw)
}
