package main

import (
	"bufio"
	gql "github.com/mattdamon108/gqlmerge/lib"
	"github.com/maximrub/thesis-diseases-risk-factors-server/internal/utils"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	// "  " is indent for the padding in generating schema
	// in case of using as go module, just " " would be fine
	//
	// paths should be a relative path
	schema := gql.Merge(" ", "types", "queries", "mutations")
	if schema != nil {
		f, err := os.Create("schema.graphqls")
		if err != nil {
			log.WithError(err).Fatal("error generating merged gql schema")
		}

		defer f.Close()
		writer := bufio.NewWriter(f)
		defer writer.Flush()
		_, err = writer.WriteString(utils.Deref(schema))
		if err != nil {
			panic(err)
		}
	}
}
