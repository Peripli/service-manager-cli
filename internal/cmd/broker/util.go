package broker

import "github.com/Peripli/service-manager-cli/pkg/types"

func getBrokerByName(brokers *types.Brokers, names []string) []types.Broker {
	result := make([]types.Broker, 0)
	for _, broker := range brokers.Brokers {
		for _, name := range names {
			if broker.Name == name {
				result = append(result, broker)
			}
		}
	}
	return result
}
