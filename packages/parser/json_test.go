package parser

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type JSONParserTestSuite struct {
	suite.Suite
}

func (suite *JSONParserTestSuite) SetupTest() {
}

func (t *testing.T) TestJSONParserTestSuite() {
	suite.Run(t, new(JSONParserTestSuite))
}
