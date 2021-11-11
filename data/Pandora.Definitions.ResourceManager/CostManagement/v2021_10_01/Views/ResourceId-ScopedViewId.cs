using System.Collections.Generic;
using Pandora.Definitions.Interfaces;

namespace Pandora.Definitions.ResourceManager.CostManagement.v2021_10_01.Views
{
    internal class ScopedViewId : ResourceID
    {
        public string? CommonAlias => null;

        public string ID => "/{scope}/providers/Microsoft.CostManagement/views/{viewName}";

        public List<ResourceIDSegment> Segments => new List<ResourceIDSegment>
        {
                new()
                {
                    Name = "scope",
                    Type = ResourceIDSegmentType.Scope
                },

                new()
                {
                    Name = "providers",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "providers"
                },

                new()
                {
                    Name = "microsoftCostManagement",
                    Type = ResourceIDSegmentType.ResourceProvider,
                    FixedValue = "Microsoft.CostManagement"
                },

                new()
                {
                    Name = "views",
                    Type = ResourceIDSegmentType.Static,
                    FixedValue = "views"
                },

                new()
                {
                    Name = "viewName",
                    Type = ResourceIDSegmentType.UserSpecified
                },

        };
    }
}
