using System.Collections.Generic;
using Pandora.Definitions.Interfaces;

namespace Pandora.Definitions.ResourceManager.AppConfiguration.v2020_06_01.PrivateEndpointConnections
{
    internal class Definition : ApiDefinition
    {
        // Generated from Swagger revision "28e60e3f539b44b60e7b4d6fa2cf4476519bcf93" 

        public string ApiVersion => "2020-06-01";
        public string Name => "PrivateEndpointConnections";
        public IEnumerable<Interfaces.ApiOperation> Operations => new List<Interfaces.ApiOperation>
        {
            new CreateOrUpdateOperation(),
            new DeleteOperation(),
            new GetOperation(),
            new ListByConfigurationStoreOperation(),
        };
    }
}
