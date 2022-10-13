using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.AzureStackHCI.v2022_09_01.Cluster;

[ConstantType(ConstantTypeAttribute.ConstantType.String)]
internal enum SoftwareAssuranceIntentConstant
{
    [Description("Disable")]
    Disable,

    [Description("Enable")]
    Enable,
}
