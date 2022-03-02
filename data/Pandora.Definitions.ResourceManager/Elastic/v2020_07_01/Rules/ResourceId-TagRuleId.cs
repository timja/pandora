using System.Collections.Generic;
using Pandora.Definitions.Interfaces;

namespace Pandora.Definitions.ResourceManager.Elastic.v2020_07_01.Rules;

internal class TagRuleId : ResourceID
{
    public string? CommonAlias => null;

    public string ID => "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Elastic/monitors/{monitorName}/tagRules/{ruleSetName}";

    public List<ResourceIDSegment> Segments => new List<ResourceIDSegment>
    {
                new()
                {
                    Name = "staticSubscriptions",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "subscriptions"
                },

                new()
                {
                    Name = "subscriptionId",
                    Type = ResourceIDSegmentType.SubscriptionId
                },

                new()
                {
                    Name = "staticResourceGroups",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "resourceGroups"
                },

                new()
                {
                    Name = "resourceGroupName",
                    Type = ResourceIDSegmentType.ResourceGroup
                },

                new()
                {
                    Name = "staticProviders",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "providers"
                },

                new()
                {
                    Name = "staticMicrosoftElastic",
                    Type = ResourceIDSegmentType.ResourceProvider,
                    FixedValue = "Microsoft.Elastic"
                },

                new()
                {
                    Name = "staticMonitors",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "monitors"
                },

                new()
                {
                    Name = "monitorName",
                    Type = ResourceIDSegmentType.UserSpecified
                },

                new()
                {
                    Name = "staticTagRules",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "tagRules"
                },

                new()
                {
                    Name = "ruleSetName",
                    Type = ResourceIDSegmentType.UserSpecified
                },

    };
}
