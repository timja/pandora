using System.Collections.Generic;
using Pandora.Definitions.Interfaces;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.ConfidentialLedger.v2022_05_13.ConfidentialLedger;

internal class Definition : ResourceDefinition
{
    public string Name => "ConfidentialLedger";
    public IEnumerable<Interfaces.ApiOperation> Operations => new List<Interfaces.ApiOperation>
    {
        new LedgerCreateOperation(),
        new LedgerDeleteOperation(),
        new LedgerGetOperation(),
        new LedgerListByResourceGroupOperation(),
        new LedgerListBySubscriptionOperation(),
        new LedgerUpdateOperation(),
    };
}
