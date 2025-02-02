using System.Collections.Generic;
using Pandora.Definitions.Interfaces;

namespace Pandora.Definitions.ResourceManager.Storage.v2023_01_01;

public partial class Definition : ApiVersionDefinition
{
    public string ApiVersion => "2023-01-01";
    public bool Preview => false;
    public Source Source => Source.ResourceManagerRestApiSpecs;

    public IEnumerable<ResourceDefinition> Resources => new List<ResourceDefinition>
    {
        new AccountMigrations.Definition(),
        new BlobContainers.Definition(),
        new BlobInventoryPolicies.Definition(),
        new BlobService.Definition(),
        new DeletedAccounts.Definition(),
        new EncryptionScopes.Definition(),
        new FileService.Definition(),
        new FileShares.Definition(),
        new LocalUsers.Definition(),
        new ManagementPolicies.Definition(),
        new ObjectReplicationPolicies.Definition(),
        new PrivateEndpointConnections.Definition(),
        new PrivateLinkResources.Definition(),
        new QueueService.Definition(),
        new QueueServiceProperties.Definition(),
        new Skus.Definition(),
        new StorageAccounts.Definition(),
        new TableService.Definition(),
        new TableServiceProperties.Definition(),
    };
}
