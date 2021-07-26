package helper

import (
	"bestprice/bestprice-api/internal/data/model"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Checks if a string consists of Latin letters, digits and (_) in order to avoid SQL injections
// This technique is not optimal but in most cases functional. Guideline: https://stackoverflow.com/questions/30867337/golang-order-by-issue-with-mysql
// TODO : Consider to use an ORM like GORM to handle the structure of an SQL statement
func SupportedSqlString(input string) bool {
	return regexp.MustCompile("^[A-Za-z0-9_]+$").MatchString(input)
}

// Transforms a value of 'sort' query string parameter to SQL 'ORDER BY' keyword
// Not equal number of sort keys and directions have to be specified, example: sort=key1:asc,key2,key3
// The default direction is considered the ascending order and empty imnput is not handled as an error
// The function doesn't check for the validity of the final SQL statement - tha validity is up to the DB
// The pagination, filter and sort query convention is based on OpenStack API guideline https://specs.openstack.org/openstack/api-wg/guidelines/pagination_filter_sort.html
func ExportSqlSorting(sortInput string) (string, error) {
	var output string
	var ordering [][2]string // Slice(not a Map): Keeps the order of user input
	const (
		ORDER_BY            = "ORDER BY"
		ASC                 = "ASC"
		PAIR_SEPARATOR      = ","
		KEY_VALUE_SEPARATOR = ":"
		COLUMN_SEPARATOR    = ","
	)

	if sortInput == "" {
		return output, nil
	}

	pairs := strings.Split(sortInput, PAIR_SEPARATOR)
	for i := range pairs {
		pair := strings.Split(pairs[i], KEY_VALUE_SEPARATOR)
		for j := range pair {
			if !SupportedSqlString(pair[j]) {
				return output, errors.New("Not supported characters in 'sort' query string")
			}
		}

		switch len(pair) {
		case 1:
			ordering = append(ordering, [2]string{pair[0], ASC})
		case 2:
			ordering = append(ordering, [2]string{pair[0], strings.ToUpper(pair[1])})
		default:
			return output, errors.New("Sort value is not a key:value pair")
		}
	}

	orderBy := ORDER_BY
	for i := range ordering {
		orderBy = orderBy + " " + ordering[i][0] + " " + ordering[i][1] + COLUMN_SEPARATOR
	}
	output = strings.TrimRight(orderBy, COLUMN_SEPARATOR)

	return output, nil
}

// Checks compliance of 'limit' and 'offset' query string parameters with SQL 'LIMIT' and 'OFFSET' keywords
// If 'limit' is empty, is not handled as an error but offset is ignored
// The values validity (non negative integer etc.) is up to DB
// The pagination, filter and sort query convention is based on OpenStack API guideline https://specs.openstack.org/openstack/api-wg/guidelines/pagination_filter_sort.html
// TODO : For purfomance purposes consider to alter the pagination from limit/offset to a mechanism with auto icrement ID as a cursor
func ExportSqlPagination(limitInput string, offsetInput string) (string, string, error) {
	var limitOut, offsetOut string

	if limitInput == "" {
		return limitOut, offsetOut, nil
	}

	if !SupportedSqlString(limitInput) {
		return limitOut, offsetOut, errors.New("Not supported characters in 'limit' query string")
	}
	if offsetInput != "" && !SupportedSqlString(offsetInput) {
		return limitOut, offsetOut, errors.New("Not supported characters in 'offset' query string")
	}

	limitOut = limitInput
	offsetOut = offsetInput
	return limitOut, offsetOut, nil
}

// Reads and HTTP request and returns an array of bytes
// At any case, closes the body in order to avoid resource leak and the connection to be re-used
func ReadHttpRequestAndClose(r *http.Request) ([]byte, error) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Retrieves the value of the environment variable named by the key
// If the variable is present in the environment the value is returned, otherwise the fallback value is returned
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// Redis keys are binary safe, this means that we can use any binary sequence as a key
// The empty string is also a valid key as also as the space ( ) character
// So this function prevents Redis to set empty Keys and Values or Keys that contain the space (_) character
// Also prevents to call the Redis DB with key that doesn't follow the rules above
func ValidRedisKeyValuePair(k string, v string) bool {
	if k == "" || v == "" || strings.Contains(k, " ") {
		return false
	}

	return true
}

// Transforms a slice of model.Category to a slice of model.CategoryExposed
func ExposeCategories(c []model.Category) []model.CategoryExposed {
	var output []model.CategoryExposed
	for i := range c {
		output = append(output, model.CategoryExposed(c[i]))
	}
	return output
}

// Transforms a slice of model.Product to a slice of model.ProductExposed
func ExposeProducts(p []model.Product) []model.ProductExposed {
	var output []model.ProductExposed
	for i := range p {
		output = append(output, model.ProductExposed(p[i]))
	}

	return output
}
