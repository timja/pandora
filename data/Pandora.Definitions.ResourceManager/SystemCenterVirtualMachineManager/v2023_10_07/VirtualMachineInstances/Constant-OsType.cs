using Pandora.Definitions.Attributes;
using System.ComponentModel;

namespace Pandora.Definitions.ResourceManager.SystemCenterVirtualMachineManager.v2023_10_07.VirtualMachineInstances;

[ConstantType(ConstantTypeAttribute.ConstantType.String)]
internal enum OsTypeConstant
{
    [Description("Linux")]
    Linux,

    [Description("Other")]
    Other,

    [Description("Windows")]
    Windows,
}
