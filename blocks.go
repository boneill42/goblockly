package goblockly

import (
	"encoding/xml"
)

// Top-level blockly XML container
type BlockXml struct {
	XMLName xml.Name `xml:"xml"`
	Blocks  []Block  `xml:"block"`
}

// A single Blockly block
type Block struct {
	XMLName    xml.Name         `xml:"block"`
	Type       string           `xml:"type,attr"`
	X          string           `xml:"x,attr"`
	Y          string           `xml:"y,attr"`
	Values     []BlockValue     `xml:"value"`
	Fields     []BlockField     `xml:"field"`
	Statements []BlockStatement `xml:"statement"`
	Next       *Block           `xml:"next>block"`
	Mutation   *BlockMutation   `xml:"mutation"`
}

// A value in a Blockly block. Values are blocks that will evaluate to a value.
type BlockValue struct {
	Name   string  `xml:"name,attr"`
	Blocks []Block `xml:"block"`
}

// A field attached to a Blockly block
type BlockField struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// A statement in a Blockly block. Statements are (usually stacks of) Blockly
// blocks with an ignorable return value.
type BlockStatement struct {
	XMLName xml.Name `xml:"statement"`
	Name    string   `xml:"name,attr"`
	Blocks  []Block  `xml:"block"`
}

// Modifiers on Blockly blocks. Indicates if a block has special encoding (such
// as the "elseif" / "elses" mutations on if blocks).
type BlockMutation struct {
	At        bool   `xml:"at,attr"`
	At1       bool   `xml:"at1,attr"`
	At2       bool   `xml:"at2,attr"`
	ElseIf    int    `xml:"elseif,attr"`
	Else      int    `xml:"else,attr"`
	Items     int    `xml:"items,attr"`
	Mode      string `xml:"mode,attr"`
	Statement bool   `xml:"statement,attr"`
}

// FieldWithName fetches the field with a given name, or returns nil if the
// field doesn't exist.
func (b *Block) FieldWithName(name string) *BlockField {
	for _, v := range b.Fields {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

// SingleFieldWithName fetches the value of the field with the given name, or
// "". Fails interpretation if the field doesn't exist.
func (b *Block) SingleFieldWithName(i *Interpreter, name string) string {
	fv := b.FieldWithName(name)
	if fv == nil {
		i.Fail("No field named " + name)
		return ""
	}
	return fv.Value

}

// BlockValueWithName retrieves the block value with the specified name, or
// returns nil if the block value doesn't exist.
func (b *Block) BlockValueWithName(name string) *BlockValue {
	for _, v := range b.Values {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

// SingleBlockValueWithName retrieves the single block in the block value with
// the specified name, or nil. Fails interpretation if there is not exactly one
// single block for the specified value.
func (b *Block) SingleBlockValueWithName(i *Interpreter, name string) *Block {
	bv := b.BlockValueWithName(name)
	if bv == nil {
		i.Fail("No block with value " + name)
		return nil
	}
	if len(bv.Blocks) != 1 {
		i.Fail("Block socket does not have exactly one block attached to it.")
		return nil
	}
	return &bv.Blocks[0]
}

// SingleBlockStatementWIthName retrieves the single block in the block
// statement with the specified name, or nil. Fails interpretation if there is
// not a single block in the specified statement.
func (b *Block) SingleBlockStatementWithName(i *Interpreter, name string) *Block {
	for _, v := range b.Statements {
		if v.Name == name {
			if len(v.Blocks) != 1 {
				i.Fail("Block socket does not have exactly one block attached to it.")
				return nil
			}
			return &v.Blocks[0]
		}
	}
	i.Fail("No statement with name " + name)
	return nil
}
