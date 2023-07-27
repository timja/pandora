using Pandora.Definitions.Attributes;
using Pandora.Definitions.CustomTypes;
using Pandora.Definitions.Interfaces;
using System;
using System.Collections.Generic;
using System.Net;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.ApiManagement.v2022_08_01.TenantConfigurationSyncState;

internal class TenantConfigurationGetSyncStateOperation : Pandora.Definitions.Operations.GetOperation
{
    public override ResourceID? ResourceId() => new ServiceId();

    public override Type? ResponseObject() => typeof(TenantConfigurationSyncStateContractModel);

    public override string? UriSuffix() => "/tenant/configuration/syncState";


}
