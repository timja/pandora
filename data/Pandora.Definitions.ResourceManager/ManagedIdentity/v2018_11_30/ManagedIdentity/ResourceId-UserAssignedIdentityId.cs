using System.Collections.Generic;
using Pandora.Definitions.Interfaces;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.ManagedIdentity.v2018_11_30.ManagedIdentity;

internal class UserAssignedIdentityId : ResourceID
{
    public string? CommonAlias => "UserAssignedIdentity";

    public string ID => "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{resourceName}";

    public List<ResourceIDSegment> Segments => new List<ResourceIDSegment>
    {

    };
}
