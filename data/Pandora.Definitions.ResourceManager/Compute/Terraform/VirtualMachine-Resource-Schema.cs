using System.Collections.Generic;
using Pandora.Definitions.Attributes;

namespace Pandora.Definitions.ResourceManager.Compute.Terraform;

public class VirtualMachineResourceSchema
{
    [HclName("additional_capabilities")]
    [ForceNew]
    [Optional]
    public AdditionalCapabilitiesSchema AdditionalCapabilities { get; set; }

    [HclName("application_profile")]
    [ForceNew]
    [Optional]
    public ApplicationProfileSchema ApplicationProfile { get; set; }

    [HclName("availability_set_id")]
    [ForceNew]
    [Optional]
    public string AvailabilitySetId { get; set; }

    [HclName("billing_profile")]
    [ForceNew]
    [Optional]
    public BillingProfileSchema BillingProfile { get; set; }

    [HclName("capacity_reservation")]
    [ForceNew]
    [Optional]
    public CapacityReservationProfileSchema CapacityReservationProfile { get; set; }

    [HclName("diagnostics_profile")]
    [ForceNew]
    [Optional]
    public DiagnosticsProfileSchema DiagnosticsProfile { get; set; }

    [HclName("eviction_policy")]
    [ForceNew]
    [Optional]
    public VirtualMachineEvictionPolicyTypesConstant? VirtualMachineEvictionPolicyTypes { get; set; }

    [HclName("extensions_time_budget")]
    [ForceNew]
    [Optional]
    public string ExtensionsTimeBudget { get; set; }

    [HclName("hardware_profile")]
    [ForceNew]
    [Optional]
    public HardwareProfileSchema HardwareProfile { get; set; }

    [HclName("host_group_id")]
    [ForceNew]
    [Optional]
    public string HostGroupId { get; set; }

    [HclName("host_id")]
    [ForceNew]
    [Optional]
    public string HostId { get; set; }

    [HclName("identity")]
    [Optional]
    public CustomTypes.LegacySystemAndUserAssignedIdentityMap Identity { get; set; }

    [HclName("instance_view")]
    [ForceNew]
    [Optional]
    public VirtualMachineInstanceViewSchema VirtualMachineInstanceView { get; set; }

    [HclName("location")]
    [ForceNew]
    [Required]
    public CustomTypes.Location Location { get; set; }

    [HclName("name")]
    [ForceNew]
    [Required]
    public string Name { get; set; }

    [HclName("network_profile")]
    [ForceNew]
    [Optional]
    public NetworkProfileSchema NetworkProfile { get; set; }

    [HclName("os_profile")]
    [ForceNew]
    [Optional]
    public OSProfileSchema OSProfile { get; set; }

    [HclName("platform_fault_domain")]
    [ForceNew]
    [Optional]
    public int PlatformFaultDomain { get; set; }

    [HclName("priority")]
    [ForceNew]
    [Optional]
    public VirtualMachinePriorityTypesConstant? VirtualMachinePriorityTypes { get; set; }

    [HclName("proximity_placement_group_id")]
    [ForceNew]
    [Optional]
    public string ProximityPlacementGroupId { get; set; }

    [HclName("scheduled_events_profile")]
    [ForceNew]
    [Optional]
    public ScheduledEventsProfileSchema ScheduledEventsProfile { get; set; }

    [HclName("security_profile")]
    [ForceNew]
    [Optional]
    public SecurityProfileSchema SecurityProfile { get; set; }

    [HclName("storage_profile")]
    [ForceNew]
    [Optional]
    public StorageProfileSchema StorageProfile { get; set; }

    [HclName("tags")]
    [Optional]
    public CustomTypes.Tags Tags { get; set; }

    [HclName("time_created")]
    [ForceNew]
    [Optional]
    public string TimeCreated { get; set; }

    [HclName("user_data")]
    [ForceNew]
    [Optional]
    public string UserData { get; set; }

    [HclName("virtual_machine_scale_set_id")]
    [ForceNew]
    [Optional]
    public string VirtualMachineScaleSetId { get; set; }

    [HclName("vm_id")]
    [ForceNew]
    [Optional]
    public string VmId { get; set; }


}
