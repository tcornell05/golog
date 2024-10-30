package golog

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// appendStruct takes val and appends it to buf as a struct.
func (p *DevHandler) appendStruct(buf []byte, val reflect.Value, fgColor []byte, indent int) []byte {
	// Try to JSON marshal the struct with indentation
	jsonBytes, err := json.MarshalIndent(val.Interface(), "", "  ")
	if err == nil {
		// Successfully marshalled to JSON, use p.appendJSON to handle colors
		buf = p.appendJSON(buf, jsonBytes, indent)
		return buf
	}

	// If marshalling fails, proceed with the original method
	vType := val.Type()
	buf = fmt.Appendf(buf, "%s", vType.String())

	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name

		buf = append(buf, '\n')
		buf = fmt.Appendf(buf, "%*s", indent+2, "")
		buf = fmt.Appendf(buf, "%s|- %s%s : ", fgColor, fieldName, colorReset)
		buf = p.appendType(buf, val.Field(i), fgColor, indent+2)
	}

	return buf
}
