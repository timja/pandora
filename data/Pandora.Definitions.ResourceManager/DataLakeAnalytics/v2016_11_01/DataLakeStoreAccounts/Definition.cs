using System.Collections.Generic;
using Pandora.Definitions.Interfaces;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.DataLakeAnalytics.v2016_11_01.DataLakeStoreAccounts;

internal class Definition : ResourceDefinition
{
    public string Name => "DataLakeStoreAccounts";
    public IEnumerable<Interfaces.ApiOperation> Operations => new List<Interfaces.ApiOperation>
    {
        new AddOperation(),
        new DeleteOperation(),
        new GetOperation(),
        new ListByAccountOperation(),
    };
    public IEnumerable<System.Type> Constants => new List<System.Type>
    {

    };
    public IEnumerable<System.Type> Models => new List<System.Type>
    {
        typeof(AddDataLakeStoreParametersModel),
        typeof(AddDataLakeStorePropertiesModel),
        typeof(DataLakeStoreAccountInformationModel),
        typeof(DataLakeStoreAccountInformationPropertiesModel),
    };
}
