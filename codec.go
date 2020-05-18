package commercio

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
)

type typeMapping map[string]string

// cosmosType returns the Cosmos codec type associated with i.
func (tm typeMapping) cosmosType(i interface{}) string {
	if i == nil {
		return ""
	}

	t := strings.Split(reflect.TypeOf(i).String(), ".")
	if len(t) < 2 {
		return tm[t[0]]
	}

	return tm[t[1]]
}

// generateTypeMappings reads Cosmos type mappings from c and returns a map that associates each Go type to its
// Codec counterpart.
func generateTypeMappings(c *codec.Codec) typeMapping {
	tsb := bytes.Buffer{}
	if err := c.PrintTypes(&tsb); err != nil {
		panic(fmt.Errorf("cannot fetch cosmos codec types, %w", err))
	}

	ts := strings.NewReplacer("| ", "", " |", "").Replace(tsb.String())

	mappings := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(ts))
	skipped := 0
	for scanner.Scan() {
		text := scanner.Text()
		if skipped < 2 {
			skipped++
			continue
		}
		fields := strings.Split(text, " ")
		if len(fields) < 3 {
			panic("could not split type mappings, length less than 3")
		}

		// associate go type (field 0) to the cosmos type (field 1)
		mappings[fields[0]] = fields[1]
	}

	if scanner.Err() != nil {
		panic(fmt.Errorf("scanner gave error, %w", scanner.Err()))
	}

	return mappings
}
