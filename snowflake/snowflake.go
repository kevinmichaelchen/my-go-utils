package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

// InitSnowflakeNode initializes a Snowflake node.
//
// Example usage:
//
//   a.snowflakeNode = utils.InitSnowflakeNode(1)
func InitSnowflakeNode(nodeNumber int64) *snowflake.Node {
	node, err := snowflake.NewNode(nodeNumber)
	if err != nil {
		log.Fatalf("Could not generate Snowflake node for node number: %d", nodeNumber)
	}
	return node
}

// NewPrimaryKey generates a Snowflake ID and returns it as an int64.
func NewPrimaryKey(snowflakeNode *snowflake.Node) int64 {
	return snowflakeNode.Generate().Int64()
}