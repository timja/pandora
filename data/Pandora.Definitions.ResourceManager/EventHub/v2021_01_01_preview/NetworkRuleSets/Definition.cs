using System.Collections.Generic;
using Pandora.Definitions.Interfaces;

namespace Pandora.Definitions.ResourceManager.EventHub.v2021_01_01_preview.NetworkRuleSets
{
    internal class Definition : ApiDefinition
    {
        // Generated from Swagger revision "3865f04d22e82db481be0727b406021d29cd2b70" 

        public string ApiVersion => "2021-01-01-preview";
        public string Name => "NetworkRuleSets";
        public IEnumerable<Interfaces.ApiOperation> Operations => new List<Interfaces.ApiOperation>
        {
            new NamespacesCreateOrUpdateNetworkRuleSetOperation(),
            new NamespacesGetNetworkRuleSetOperation(),
        };
    }
}
