// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

using Pandora.Definitions.CustomTypes;
using Pandora.Definitions.Interfaces;
using Pandora.Definitions.MicrosoftGraph.Beta.CommonTypes;
using System;

namespace Pandora.Definitions.MicrosoftGraph.Beta.Me.Beta.MeChatMessageReply;

internal class ListMeChatByIdMessageByIdRepliesOperation : Operations.ListOperation
{
    public override string? FieldContainingPaginationDetails() => "nextLink";
    public override ResourceID? ResourceId() => new MeChatIdMessageId();
    public override Type NestedItemType() => typeof(ChatMessageCollectionResponseModel);
    public override string? UriSuffix() => "/replies";
}
