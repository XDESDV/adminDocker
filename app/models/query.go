package models

import (
	"adminDocker/app/functions"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// QueryParams : Parametres requetes HTTP
type QueryParams struct {
	ID               string
	Path             string
	Body             []byte
	View             string
	FilterClause     []string
	FilterLikeClause []string
	SortClause       []string
	Offset           int
	Count            int
	Export           string
	GroupBy          string
	Columns          []string
	SearchClause     []string
	Collection       string
	TestDeleted      bool
}

// Parse : QueryParams parser
func (q *QueryParams) Parse(c *gin.Context) {

	q.Count, _ = strconv.Atoi(c.Query("count"))
	q.Offset, _ = strconv.Atoi(c.Query("offset"))
	q.View = c.Query("view")
	q.GroupBy = c.Query("col")

	q.Path = c.Request.URL.Path

	if search := c.Query("search"); len(search) > 0 {
		tabKeyword := strings.Split(search, " ")
		for i, keyword := range tabKeyword {
			// Double apostrophes
			tabKeyword[i] = strings.Replace(keyword, "'", "''", -1)
		}
		q.SearchClause = tabKeyword
	}

	if sort := c.Query("sort"); len(sort) > 0 {
		q.SortClause = strings.Split(sort, ",")
	}

	// For POST or PUT requests, reading body JSON
	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			log.Fatal().Err(err).Msg("")
		}
		// Convert map to JSON bytes
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		q.Body = bodyBytes
	}

	q.FilterClause = c.QueryArray("filter")
	q.FilterLikeClause = c.QueryArray("filter_like")

	// Assuming functions.RemoveDuplicate is your custom function to remove duplicates
	functions.RemoveDuplicate(&q.FilterClause)
	functions.RemoveDuplicate(&q.FilterLikeClause)
	functions.RemoveDuplicate(&q.SearchClause)
}
