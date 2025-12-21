package gomind

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TOON (Token-Oriented Object Notation) encoder for LLM context serialization.
// Provides ~40% token reduction compared to JSON for uniform arrays.

// Encode converts any Go value to TOON format string.
func Encode(v any) string {
	return encodeValue(reflect.ValueOf(v), 0)
}

// EncodeTabular encodes a slice of maps as TOON tabular format.
// Format: name[N]{field1,field2,...}:\n  val1,val2,...\n  ...
func EncodeTabular(name string, rows []map[string]string, fields ...string) string {
	if len(rows) == 0 {
		return fmt.Sprintf("%s[0]{%s}:", name, strings.Join(fields, ","))
	}

	var sb strings.Builder

	// Write header: name[count]{fields}:
	sb.WriteString(fmt.Sprintf("%s[%d]{%s}:\n", name, len(rows), strings.Join(fields, ",")))

	// Write rows
	for _, row := range rows {
		sb.WriteString("  ")
		for i, field := range fields {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(escapeValue(row[field]))
		}
		sb.WriteString("\n")
	}

	return strings.TrimSuffix(sb.String(), "\n")
}

// EncodeTabularAuto encodes a slice of structs as TOON tabular format,
// automatically extracting fields based on JSON tags.
// If fields is empty, all exported fields with json tags are included.
func EncodeTabularAuto[T any](name string, items []T, fields ...string) string {
	if len(items) == 0 {
		if len(fields) == 0 {
			return fmt.Sprintf("%s[0]{}:", name)
		}
		return fmt.Sprintf("%s[0]{%s}:", name, strings.Join(fields, ","))
	}

	// Get field mappings from struct type
	var sample T
	fieldMap := getJSONFieldMap(reflect.TypeOf(sample))

	// If no fields specified, use all JSON-tagged fields
	if len(fields) == 0 {
		fields = getOrderedJSONFields(reflect.TypeOf(sample))
	}

	// Build rows
	rows := make([]map[string]string, len(items))
	for i, item := range items {
		rows[i] = extractFieldValues(reflect.ValueOf(item), fieldMap, fields)
	}

	return EncodeTabular(name, rows, fields...)
}

// getJSONFieldMap returns a map of JSON tag name -> struct field index.
func getJSONFieldMap(t reflect.Type) map[string]int {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fieldMap := make(map[string]int)
	if t.Kind() != reflect.Struct {
		return fieldMap
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse json tag (handle "name,omitempty" format)
		tagName := strings.Split(jsonTag, ",")[0]
		if tagName != "" {
			fieldMap[tagName] = i
		}
	}

	return fieldMap
}

// getOrderedJSONFields returns JSON field names in struct declaration order.
func getOrderedJSONFields(t reflect.Type) []string {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		tagName := strings.Split(jsonTag, ",")[0]
		if tagName != "" {
			fields = append(fields, tagName)
		}
	}

	return fields
}

// extractFieldValues extracts string values for specified fields from a struct.
func extractFieldValues(v reflect.Value, fieldMap map[string]int, fields []string) map[string]string {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return make(map[string]string)
		}
		v = v.Elem()
	}

	result := make(map[string]string)
	for _, fieldName := range fields {
		idx, ok := fieldMap[fieldName]
		if !ok {
			result[fieldName] = ""
			continue
		}

		fieldVal := v.Field(idx)
		result[fieldName] = valueToString(fieldVal)
	}

	return result
}

// valueToString converts a reflect.Value to its string representation.
func valueToString(v reflect.Value) string {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Struct:
		// Handle nested structs - try to get Name field (common pattern)
		if nameField := v.FieldByName("Name"); nameField.IsValid() && nameField.Kind() == reflect.String {
			return nameField.String()
		}
		return fmt.Sprintf("%v", v.Interface())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// encodeValue recursively encodes a reflect.Value to TOON format.
func encodeValue(v reflect.Value, indent int) string {
	if !v.IsValid() {
		return "null"
	}

	// Handle pointers and interfaces
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return "null"
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return escapeValue(v.String())

	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)

	case reflect.Slice, reflect.Array:
		return encodeSlice(v, indent)

	case reflect.Map:
		return encodeMap(v, indent)

	case reflect.Struct:
		return encodeStruct(v, indent)

	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// encodeSlice encodes a slice/array to TOON format.
func encodeSlice(v reflect.Value, indent int) string {
	if v.Len() == 0 {
		return "[]"
	}

	// Check if all elements are primitives (simple array)
	if isPrimitiveSlice(v) {
		return encodePrimitiveSlice(v)
	}

	// YAML-like list format for non-uniform or complex arrays
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(indentStr)
		sb.WriteString("- ")
		elem := encodeValue(v.Index(i), indent+1)
		// Handle multi-line values
		if strings.Contains(elem, "\n") {
			sb.WriteString("\n")
			sb.WriteString(indentLines(elem, indent+1))
		} else {
			sb.WriteString(elem)
		}
	}

	return sb.String()
}

// encodePrimitiveSlice encodes a slice of primitives as comma-separated values.
func encodePrimitiveSlice(v reflect.Value) string {
	parts := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		parts[i] = encodeValue(v.Index(i), 0)
	}
	return fmt.Sprintf("[%d]: %s", v.Len(), strings.Join(parts, ","))
}

// encodeMap encodes a map to TOON format (YAML-like).
func encodeMap(v reflect.Value, indent int) string {
	if v.Len() == 0 {
		return "{}"
	}

	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)
	keys := v.MapKeys()

	for i, key := range keys {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(indentStr)
		sb.WriteString(fmt.Sprintf("%v: ", key.Interface()))

		val := encodeValue(v.MapIndex(key), indent+1)
		if strings.Contains(val, "\n") {
			sb.WriteString("\n")
			sb.WriteString(indentLines(val, indent+1))
		} else {
			sb.WriteString(val)
		}
	}

	return sb.String()
}

// encodeStruct encodes a struct to TOON format (YAML-like).
func encodeStruct(v reflect.Value, indent int) string {
	t := v.Type()
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)
	first := true

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldVal := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get field name from json tag or use field name
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				fieldName = parts[0]
			}
			// Skip if json:"-"
			if parts[0] == "-" {
				continue
			}
			// Handle omitempty
			if len(parts) > 1 && parts[1] == "omitempty" && isZeroValue(fieldVal) {
				continue
			}
		}

		if !first {
			sb.WriteString("\n")
		}
		first = false

		sb.WriteString(indentStr)
		sb.WriteString(fieldName)
		sb.WriteString(": ")

		val := encodeValue(fieldVal, indent+1)
		if strings.Contains(val, "\n") {
			sb.WriteString("\n")
			sb.WriteString(indentLines(val, indent+1))
		} else {
			sb.WriteString(val)
		}
	}

	return sb.String()
}

// escapeValue escapes a string value for TOON format.
// Quotes the value if it contains special characters.
func escapeValue(s string) string {
	if s == "" {
		return ""
	}

	// Check if quoting is needed
	needsQuote := false
	for _, c := range s {
		if c == ',' || c == '"' || c == '\n' || c == '\r' || c == ':' || c == '[' || c == ']' || c == '{' || c == '}' {
			needsQuote = true
			break
		}
	}

	if !needsQuote {
		return s
	}

	// Escape quotes and wrap in quotes
	escaped := strings.ReplaceAll(s, `"`, `""`)
	return `"` + escaped + `"`
}

// isPrimitiveSlice checks if all elements in a slice are primitive types.
func isPrimitiveSlice(v reflect.Value) bool {
	if v.Len() == 0 {
		return true
	}

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		for elem.Kind() == reflect.Ptr || elem.Kind() == reflect.Interface {
			if elem.IsNil() {
				continue
			}
			elem = elem.Elem()
		}

		switch elem.Kind() {
		case reflect.String, reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			continue
		default:
			return false
		}
	}
	return true
}

// isZeroValue checks if a value is the zero value for its type.
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}

// indentLines adds indentation to each line of a string.
func indentLines(s string, indent int) string {
	indentStr := strings.Repeat("  ", indent)
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indentStr + line
		}
	}
	return strings.Join(lines, "\n")
}
