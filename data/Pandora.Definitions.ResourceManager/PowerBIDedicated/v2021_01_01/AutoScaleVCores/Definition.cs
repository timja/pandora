using System.Collections.Generic;
using Pandora.Definitions.Interfaces;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.PowerBIDedicated.v2021_01_01.AutoScaleVCores;

internal class Definition : ResourceDefinition
{
    public string Name => "AutoScaleVCores";
    public IEnumerable<Interfaces.ApiOperation> Operations => new List<Interfaces.ApiOperation>
    {
        new CreateOperation(),
        new DeleteOperation(),
        new GetOperation(),
        new ListByResourceGroupOperation(),
        new ListBySubscriptionOperation(),
        new UpdateOperation(),
    };
    public IEnumerable<System.Type> Constants => new List<System.Type>
    {
        typeof(VCoreProvisioningStateConstant),
        typeof(VCoreSkuTierConstant),
    };
    public IEnumerable<System.Type> Models => new List<System.Type>
    {
        typeof(AutoScaleVCoreModel),
        typeof(AutoScaleVCoreListResultModel),
        typeof(AutoScaleVCoreMutablePropertiesModel),
        typeof(AutoScaleVCorePropertiesModel),
        typeof(AutoScaleVCoreSkuModel),
        typeof(AutoScaleVCoreUpdateParametersModel),
    };
}
