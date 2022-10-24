using Pandora.Definitions.Attributes;
using Pandora.Definitions.CustomTypes;
using Pandora.Definitions.Interfaces;
using Pandora.Definitions.Operations;
using System;
using System.Collections.Generic;
using System.Net;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.Batch.v2022_01_01.Location;

internal class ListSupportedCloudServiceSkusOperation : Operations.ListOperation
{
    public override string? FieldContainingPaginationDetails() => "nextLink";

    public override ResourceID? ResourceId() => new LocationId();

    public override Type NestedItemType() => typeof(SupportedSkuModel);

    public override Type? OptionsObject() => typeof(ListSupportedCloudServiceSkusOperation.ListSupportedCloudServiceSkusOptions);

    public override string? UriSuffix() => "/cloudServiceSkus";

    internal class ListSupportedCloudServiceSkusOptions
    {
        [QueryStringName("$filter")]
        [Optional]
        public string Filter { get; set; }

        [QueryStringName("maxresults")]
        [Optional]
        public int Maxresults { get; set; }
    }
}
