//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ualberta

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Type of AIS vessel, enum duplicated from AIS standard, https://gpsd.gitlab.io/gpsd/AIVDM.html
type AIS_TYPE = common.AIS_TYPE

const (
	// Not available (default).
	AIS_TYPE_UNKNOWN     AIS_TYPE = common.AIS_TYPE_UNKNOWN
	AIS_TYPE_RESERVED_1  AIS_TYPE = common.AIS_TYPE_RESERVED_1
	AIS_TYPE_RESERVED_2  AIS_TYPE = common.AIS_TYPE_RESERVED_2
	AIS_TYPE_RESERVED_3  AIS_TYPE = common.AIS_TYPE_RESERVED_3
	AIS_TYPE_RESERVED_4  AIS_TYPE = common.AIS_TYPE_RESERVED_4
	AIS_TYPE_RESERVED_5  AIS_TYPE = common.AIS_TYPE_RESERVED_5
	AIS_TYPE_RESERVED_6  AIS_TYPE = common.AIS_TYPE_RESERVED_6
	AIS_TYPE_RESERVED_7  AIS_TYPE = common.AIS_TYPE_RESERVED_7
	AIS_TYPE_RESERVED_8  AIS_TYPE = common.AIS_TYPE_RESERVED_8
	AIS_TYPE_RESERVED_9  AIS_TYPE = common.AIS_TYPE_RESERVED_9
	AIS_TYPE_RESERVED_10 AIS_TYPE = common.AIS_TYPE_RESERVED_10
	AIS_TYPE_RESERVED_11 AIS_TYPE = common.AIS_TYPE_RESERVED_11
	AIS_TYPE_RESERVED_12 AIS_TYPE = common.AIS_TYPE_RESERVED_12
	AIS_TYPE_RESERVED_13 AIS_TYPE = common.AIS_TYPE_RESERVED_13
	AIS_TYPE_RESERVED_14 AIS_TYPE = common.AIS_TYPE_RESERVED_14
	AIS_TYPE_RESERVED_15 AIS_TYPE = common.AIS_TYPE_RESERVED_15
	AIS_TYPE_RESERVED_16 AIS_TYPE = common.AIS_TYPE_RESERVED_16
	AIS_TYPE_RESERVED_17 AIS_TYPE = common.AIS_TYPE_RESERVED_17
	AIS_TYPE_RESERVED_18 AIS_TYPE = common.AIS_TYPE_RESERVED_18
	AIS_TYPE_RESERVED_19 AIS_TYPE = common.AIS_TYPE_RESERVED_19
	// Wing In Ground effect.
	AIS_TYPE_WIG             AIS_TYPE = common.AIS_TYPE_WIG
	AIS_TYPE_WIG_HAZARDOUS_A AIS_TYPE = common.AIS_TYPE_WIG_HAZARDOUS_A
	AIS_TYPE_WIG_HAZARDOUS_B AIS_TYPE = common.AIS_TYPE_WIG_HAZARDOUS_B
	AIS_TYPE_WIG_HAZARDOUS_C AIS_TYPE = common.AIS_TYPE_WIG_HAZARDOUS_C
	AIS_TYPE_WIG_HAZARDOUS_D AIS_TYPE = common.AIS_TYPE_WIG_HAZARDOUS_D
	AIS_TYPE_WIG_RESERVED_1  AIS_TYPE = common.AIS_TYPE_WIG_RESERVED_1
	AIS_TYPE_WIG_RESERVED_2  AIS_TYPE = common.AIS_TYPE_WIG_RESERVED_2
	AIS_TYPE_WIG_RESERVED_3  AIS_TYPE = common.AIS_TYPE_WIG_RESERVED_3
	AIS_TYPE_WIG_RESERVED_4  AIS_TYPE = common.AIS_TYPE_WIG_RESERVED_4
	AIS_TYPE_WIG_RESERVED_5  AIS_TYPE = common.AIS_TYPE_WIG_RESERVED_5
	AIS_TYPE_FISHING         AIS_TYPE = common.AIS_TYPE_FISHING
	AIS_TYPE_TOWING          AIS_TYPE = common.AIS_TYPE_TOWING
	// Towing: length exceeds 200m or breadth exceeds 25m.
	AIS_TYPE_TOWING_LARGE AIS_TYPE = common.AIS_TYPE_TOWING_LARGE
	// Dredging or other underwater ops.
	AIS_TYPE_DREDGING    AIS_TYPE = common.AIS_TYPE_DREDGING
	AIS_TYPE_DIVING      AIS_TYPE = common.AIS_TYPE_DIVING
	AIS_TYPE_MILITARY    AIS_TYPE = common.AIS_TYPE_MILITARY
	AIS_TYPE_SAILING     AIS_TYPE = common.AIS_TYPE_SAILING
	AIS_TYPE_PLEASURE    AIS_TYPE = common.AIS_TYPE_PLEASURE
	AIS_TYPE_RESERVED_20 AIS_TYPE = common.AIS_TYPE_RESERVED_20
	AIS_TYPE_RESERVED_21 AIS_TYPE = common.AIS_TYPE_RESERVED_21
	// High Speed Craft.
	AIS_TYPE_HSC             AIS_TYPE = common.AIS_TYPE_HSC
	AIS_TYPE_HSC_HAZARDOUS_A AIS_TYPE = common.AIS_TYPE_HSC_HAZARDOUS_A
	AIS_TYPE_HSC_HAZARDOUS_B AIS_TYPE = common.AIS_TYPE_HSC_HAZARDOUS_B
	AIS_TYPE_HSC_HAZARDOUS_C AIS_TYPE = common.AIS_TYPE_HSC_HAZARDOUS_C
	AIS_TYPE_HSC_HAZARDOUS_D AIS_TYPE = common.AIS_TYPE_HSC_HAZARDOUS_D
	AIS_TYPE_HSC_RESERVED_1  AIS_TYPE = common.AIS_TYPE_HSC_RESERVED_1
	AIS_TYPE_HSC_RESERVED_2  AIS_TYPE = common.AIS_TYPE_HSC_RESERVED_2
	AIS_TYPE_HSC_RESERVED_3  AIS_TYPE = common.AIS_TYPE_HSC_RESERVED_3
	AIS_TYPE_HSC_RESERVED_4  AIS_TYPE = common.AIS_TYPE_HSC_RESERVED_4
	AIS_TYPE_HSC_UNKNOWN     AIS_TYPE = common.AIS_TYPE_HSC_UNKNOWN
	AIS_TYPE_PILOT           AIS_TYPE = common.AIS_TYPE_PILOT
	// Search And Rescue vessel.
	AIS_TYPE_SAR         AIS_TYPE = common.AIS_TYPE_SAR
	AIS_TYPE_TUG         AIS_TYPE = common.AIS_TYPE_TUG
	AIS_TYPE_PORT_TENDER AIS_TYPE = common.AIS_TYPE_PORT_TENDER
	// Anti-pollution equipment.
	AIS_TYPE_ANTI_POLLUTION    AIS_TYPE = common.AIS_TYPE_ANTI_POLLUTION
	AIS_TYPE_LAW_ENFORCEMENT   AIS_TYPE = common.AIS_TYPE_LAW_ENFORCEMENT
	AIS_TYPE_SPARE_LOCAL_1     AIS_TYPE = common.AIS_TYPE_SPARE_LOCAL_1
	AIS_TYPE_SPARE_LOCAL_2     AIS_TYPE = common.AIS_TYPE_SPARE_LOCAL_2
	AIS_TYPE_MEDICAL_TRANSPORT AIS_TYPE = common.AIS_TYPE_MEDICAL_TRANSPORT
	// Noncombatant ship according to RR Resolution No. 18.
	AIS_TYPE_NONECOMBATANT         AIS_TYPE = common.AIS_TYPE_NONECOMBATANT
	AIS_TYPE_PASSENGER             AIS_TYPE = common.AIS_TYPE_PASSENGER
	AIS_TYPE_PASSENGER_HAZARDOUS_A AIS_TYPE = common.AIS_TYPE_PASSENGER_HAZARDOUS_A
	AIS_TYPE_PASSENGER_HAZARDOUS_B AIS_TYPE = common.AIS_TYPE_PASSENGER_HAZARDOUS_B
	AIS_TYPE_PASSENGER_HAZARDOUS_C AIS_TYPE = common.AIS_TYPE_PASSENGER_HAZARDOUS_C
	AIS_TYPE_PASSENGER_HAZARDOUS_D AIS_TYPE = common.AIS_TYPE_PASSENGER_HAZARDOUS_D
	AIS_TYPE_PASSENGER_RESERVED_1  AIS_TYPE = common.AIS_TYPE_PASSENGER_RESERVED_1
	AIS_TYPE_PASSENGER_RESERVED_2  AIS_TYPE = common.AIS_TYPE_PASSENGER_RESERVED_2
	AIS_TYPE_PASSENGER_RESERVED_3  AIS_TYPE = common.AIS_TYPE_PASSENGER_RESERVED_3
	AIS_TYPE_PASSENGER_RESERVED_4  AIS_TYPE = common.AIS_TYPE_PASSENGER_RESERVED_4
	AIS_TYPE_PASSENGER_UNKNOWN     AIS_TYPE = common.AIS_TYPE_PASSENGER_UNKNOWN
	AIS_TYPE_CARGO                 AIS_TYPE = common.AIS_TYPE_CARGO
	AIS_TYPE_CARGO_HAZARDOUS_A     AIS_TYPE = common.AIS_TYPE_CARGO_HAZARDOUS_A
	AIS_TYPE_CARGO_HAZARDOUS_B     AIS_TYPE = common.AIS_TYPE_CARGO_HAZARDOUS_B
	AIS_TYPE_CARGO_HAZARDOUS_C     AIS_TYPE = common.AIS_TYPE_CARGO_HAZARDOUS_C
	AIS_TYPE_CARGO_HAZARDOUS_D     AIS_TYPE = common.AIS_TYPE_CARGO_HAZARDOUS_D
	AIS_TYPE_CARGO_RESERVED_1      AIS_TYPE = common.AIS_TYPE_CARGO_RESERVED_1
	AIS_TYPE_CARGO_RESERVED_2      AIS_TYPE = common.AIS_TYPE_CARGO_RESERVED_2
	AIS_TYPE_CARGO_RESERVED_3      AIS_TYPE = common.AIS_TYPE_CARGO_RESERVED_3
	AIS_TYPE_CARGO_RESERVED_4      AIS_TYPE = common.AIS_TYPE_CARGO_RESERVED_4
	AIS_TYPE_CARGO_UNKNOWN         AIS_TYPE = common.AIS_TYPE_CARGO_UNKNOWN
	AIS_TYPE_TANKER                AIS_TYPE = common.AIS_TYPE_TANKER
	AIS_TYPE_TANKER_HAZARDOUS_A    AIS_TYPE = common.AIS_TYPE_TANKER_HAZARDOUS_A
	AIS_TYPE_TANKER_HAZARDOUS_B    AIS_TYPE = common.AIS_TYPE_TANKER_HAZARDOUS_B
	AIS_TYPE_TANKER_HAZARDOUS_C    AIS_TYPE = common.AIS_TYPE_TANKER_HAZARDOUS_C
	AIS_TYPE_TANKER_HAZARDOUS_D    AIS_TYPE = common.AIS_TYPE_TANKER_HAZARDOUS_D
	AIS_TYPE_TANKER_RESERVED_1     AIS_TYPE = common.AIS_TYPE_TANKER_RESERVED_1
	AIS_TYPE_TANKER_RESERVED_2     AIS_TYPE = common.AIS_TYPE_TANKER_RESERVED_2
	AIS_TYPE_TANKER_RESERVED_3     AIS_TYPE = common.AIS_TYPE_TANKER_RESERVED_3
	AIS_TYPE_TANKER_RESERVED_4     AIS_TYPE = common.AIS_TYPE_TANKER_RESERVED_4
	AIS_TYPE_TANKER_UNKNOWN        AIS_TYPE = common.AIS_TYPE_TANKER_UNKNOWN
	AIS_TYPE_OTHER                 AIS_TYPE = common.AIS_TYPE_OTHER
	AIS_TYPE_OTHER_HAZARDOUS_A     AIS_TYPE = common.AIS_TYPE_OTHER_HAZARDOUS_A
	AIS_TYPE_OTHER_HAZARDOUS_B     AIS_TYPE = common.AIS_TYPE_OTHER_HAZARDOUS_B
	AIS_TYPE_OTHER_HAZARDOUS_C     AIS_TYPE = common.AIS_TYPE_OTHER_HAZARDOUS_C
	AIS_TYPE_OTHER_HAZARDOUS_D     AIS_TYPE = common.AIS_TYPE_OTHER_HAZARDOUS_D
	AIS_TYPE_OTHER_RESERVED_1      AIS_TYPE = common.AIS_TYPE_OTHER_RESERVED_1
	AIS_TYPE_OTHER_RESERVED_2      AIS_TYPE = common.AIS_TYPE_OTHER_RESERVED_2
	AIS_TYPE_OTHER_RESERVED_3      AIS_TYPE = common.AIS_TYPE_OTHER_RESERVED_3
	AIS_TYPE_OTHER_RESERVED_4      AIS_TYPE = common.AIS_TYPE_OTHER_RESERVED_4
	AIS_TYPE_OTHER_UNKNOWN         AIS_TYPE = common.AIS_TYPE_OTHER_UNKNOWN
)
