using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;
using Pandora.Definitions.Attributes;
using Pandora.Definitions.Attributes.Validation;
using Pandora.Definitions.CustomTypes;


// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.


namespace Pandora.Definitions.ResourceManager.Kusto.v2023_08_15.Databases;


internal class ReadOnlyFollowingDatabasePropertiesModel
{
    [JsonPropertyName("attachedDatabaseConfigurationName")]
    public string? AttachedDatabaseConfigurationName { get; set; }

    [JsonPropertyName("databaseShareOrigin")]
    public DatabaseShareOriginConstant? DatabaseShareOrigin { get; set; }

    [JsonPropertyName("hotCachePeriod")]
    public string? HotCachePeriod { get; set; }

    [JsonPropertyName("leaderClusterResourceId")]
    public string? LeaderClusterResourceId { get; set; }

    [JsonPropertyName("originalDatabaseName")]
    public string? OriginalDatabaseName { get; set; }

    [JsonPropertyName("principalsModificationKind")]
    public PrincipalsModificationKindConstant? PrincipalsModificationKind { get; set; }

    [JsonPropertyName("provisioningState")]
    public ProvisioningStateConstant? ProvisioningState { get; set; }

    [JsonPropertyName("softDeletePeriod")]
    public string? SoftDeletePeriod { get; set; }

    [JsonPropertyName("statistics")]
    public DatabaseStatisticsModel? Statistics { get; set; }

    [JsonPropertyName("suspensionDetails")]
    public SuspensionDetailsModel? SuspensionDetails { get; set; }

    [JsonPropertyName("tableLevelSharingProperties")]
    public TableLevelSharingPropertiesModel? TableLevelSharingProperties { get; set; }
}
