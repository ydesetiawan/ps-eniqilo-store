package helper

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ConvertSliceToPostgresArray(slice []string) string {
	arrayString := "{"
	for i, value := range slice {
		// Escape any double quotes in the string value
		value = strings.Replace(value, `"`, `\"`, -1)
		// Append each string value to the arrayString
		arrayString += `"` + value + `"`
		// Add a comma separator except for the last element
		if i < len(slice)-1 {
			arrayString += ","
		}
	}
	arrayString += "}"
	return arrayString
}

func ParsePostgresArray(src string) []string {
	// Remove curly braces from the string
	src = strings.Trim(src, "{}")
	// Split the string by comma to get individual values
	values := strings.Split(src, ",")
	// Trim whitespace from each value
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}
	return values
}

func IsStructEmpty(s interface{}) bool {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.String:
			if field.String() != "" {
				return false
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() != 0 {
				return false
			}
		case reflect.Bool:
			if field.Bool() != false {
				return false
			}
		// Add cases for other types as needed
		default:
			// Handle other types if necessary
		}
	}
	return true
}

func IntToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func CombineAndUniqueWithExclusion(a, b []int64, exclude ...int64) []int64 {
	// Create a map to store unique elements
	unique := make(map[int64]struct{})

	// Add elements of slice a to the map
	for _, val := range a {
		if !contains(exclude, val) {
			unique[val] = struct{}{}
		}
	}

	// Add elements of slice b to the map
	for _, val := range b {
		if !contains(exclude, val) {
			unique[val] = struct{}{}
		}
	}

	// Extract keys from the map to form the result slice
	result := make([]int64, 0, len(unique))
	for key := range unique {
		result = append(result, key)
	}

	return result
}

// Helper function to check if a slice contains a specific value
func contains(slice []int64, val int64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func IsOrderValueValid(val string) bool {
	orders := []string{
		"asc",
		"desc",
	}

	for _, c := range orders {
		if c == val {
			return true
		}
	}
	return false
}

func PlaceholdersString(n int) string {
	placeholders := make([]string, n)
	for i := range placeholders {
		placeholders[i] = "$" + strconv.Itoa(i+1) + ""
	}
	return strings.Join(placeholders, ", ")
}

func ValidateStruct(req interface{}) error {
	validate := validator.New()
	return validate.Struct(req)
}

func IdIsInteger(id string) bool {
	var idPattern = regexp.MustCompile(`^[0-9]+$`)
	return idPattern.MatchString(id)
}
