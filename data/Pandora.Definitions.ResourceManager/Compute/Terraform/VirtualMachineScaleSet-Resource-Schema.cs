using Pandora.Definitions.Attributes;

namespace Pandora.Definitions.ResourceManager.Compute.Terraform;

public class VirtualMachineScaleSetResourceSchema
{

    [HclName("location")]
    [ForceNew]
    [Required]
    public CustomTypes.Location Location { get; set; }

    [HclName("scale_in_policy")]
    [ForceNew]
    [Optional]
    public ScaleInPolicySchema ScaleInPolicy { get; set; }

    [HclName("spot_restore_policy")]
    [Optional]
    public SpotRestorePolicySchema SpotRestorePolicy { get; set; }

    [HclName("do_not_run_extensions_on_overprovisioned_v_ms")]
    [ForceNew]
    [Optional]
    public bool DoNotRunExtensionsOnOverprovisionedVMs { get; set; }

    [HclName("unique_id")]
    [Optional]
    public string UniqueId { get; set; }

    [HclName("overprovision")]
    [ForceNew]
    [Optional]
    public bool Overprovision { get; set; }

    [HclName("tags")]
    [Optional]
    public CustomTypes.Tags Tags { get; set; }

    [HclName("name")]
    [ForceNew]
    [Required]
    public string Name { get; set; }

    [HclName("platform_fault_domain_count")]
    [Optional]
    public int PlatformFaultDomainCount { get; set; }

    [HclName("additional_capabilities")]
    [ForceNew]
    [Optional]
    public AdditionalCapabilitiesSchema AdditionalCapabilities { get; set; }

    [HclName("time_created")]
    [Optional]
    public string TimeCreated { get; set; }

    [HclName("identity")]
    [Optional]
    public CustomTypes.LegacySystemAndUserAssignedIdentityMap Identity { get; set; }

    [HclName("single_placement_group")]
    [ForceNew]
    [Optional]
    public bool SinglePlacementGroup { get; set; }

    [HclName("automatic_repairs_policy")]
    [ForceNew]
    [Optional]
    public AutomaticRepairsPolicySchema AutomaticRepairsPolicy { get; set; }

    [HclName("upgrade_policy")]
    [ForceNew]
    [Optional]
    public UpgradePolicySchema UpgradePolicy { get; set; }

    [HclName("orchestration_mode")]
    [Optional]
    public OrchestrationModeSchema OrchestrationMode { get; set; }

    [HclName("zone_balance")]
    [Optional]
    public bool ZoneBalance { get; set; }

    [HclName("host_group_id")]
    [Optional]
    public HostGroupIdSchema HostGroupId { get; set; }

    [HclName("proximity_placement_group_id")]
    [ForceNew]
    [Optional]
    public ProximityPlacementGroupIdSchema ProximityPlacementGroupId { get; set; }

    [HclName("virtual_machine_profile")]
    [ForceNew]
    [Optional]
    public VirtualMachineScaleSetVMProfileSchema VirtualMachineScaleSetVMProfile { get; set; }



}
