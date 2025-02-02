// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

using Pandora.Definitions.Interfaces;
using Pandora.Definitions.MicrosoftGraph.Beta.CommonTypes;
using System;

namespace Pandora.Definitions.MicrosoftGraph.Beta.Policies.Beta.PolicyClaimsMappingPolicy;

internal class Definition : ResourceDefinition
{
    public string Name => "PolicyClaimsMappingPolicy";

    public IEnumerable<Interfaces.ApiOperation> Operations => new List<Interfaces.ApiOperation>
    {
        new CreatePolicyClaimsMappingPolicyOperation(),
        new DeletePolicyClaimsMappingPolicyByIdOperation(),
        new GetPolicyClaimsMappingPolicyByIdOperation(),
        new GetPolicyClaimsMappingPolicyCountOperation(),
        new ListPolicyClaimsMappingPoliciesOperation(),
        new UpdatePolicyClaimsMappingPolicyByIdOperation()
    };

    public IEnumerable<System.Type> Constants => new List<System.Type>
    {

    };

    public IEnumerable<System.Type> Models => new List<System.Type>
    {

    };
}
