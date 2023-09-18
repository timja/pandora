package resourceids

import (
	"github.com/hashicorp/pandora/tools/importer-rest-api-specs/models"
	"github.com/hashicorp/pandora/tools/sdk/resourcemanager"
)

var _ commonIdMatcher = commonIdBotServiceChannel{}

type commonIdBotServiceChannel struct{}

func (c commonIdBotServiceChannel) id() models.ParsedResourceId {
	name := "BotServiceChannel"
	return models.ParsedResourceId{
		CommonAlias: &name,
		Constants:   map[string]resourcemanager.ConstantDetails{},
		Segments: []resourcemanager.ResourceIdSegment{
			models.StaticResourceIDSegment("subscriptions", "subscriptions"),
			models.SubscriptionIDResourceIDSegment("subscriptionId"),
			models.StaticResourceIDSegment("resourceGroups", "resourceGroups"),
			models.ResourceGroupResourceIDSegment("resourceGroupName"),
			models.StaticResourceIDSegment("providers", "providers"),
			models.ResourceProviderResourceIDSegment("resourceProvider", "Microsoft.BotService"),
			models.StaticResourceIDSegment("botServices", "botServices"),
			models.UserSpecifiedResourceIDSegment("botServiceName"),
			models.StaticResourceIDSegment("channels", "channels"),
			models.ConstantResourceIDSegment("channelType", "ChannelType"),
		},
	}
}
