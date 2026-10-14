package topics_test

import (
	"fmt"

	"github.com/mateusz-kacpura/integra-core/pkg/topics"
)

func ExampleTopicForAggregate() {
	fmt.Println("User aggregate topic:", topics.TopicForAggregate("user"))
	fmt.Println("Role aggregate topic:", topics.TopicForAggregate("role"))
	fmt.Println("License aggregate topic:", topics.TopicForAggregate("license"))
	// Output:
	// User aggregate topic: user-events
	// Role aggregate topic: user-roles-sync
	// License aggregate topic: license-events
}

func ExampleAllTopics() {
	all := topics.AllTopics()
	fmt.Println("Total registered topics:", len(all))
	// Output:
	// Total registered topics: 38
}
