package gomavlib

// this file contains a test dialect used in other tests.
// it's better not to import real dialects, but to use a separate one

var testDialect = MustDialect(3, []Message{
	&MessageTest5{},
	&MessageTest6{},
	&MessageTest8{},
	&MessageHeartbeat{},
	&MessageOpticalFlow{},
})

type MAV_TYPE int
type MAV_AUTOPILOT int
type MAV_MODE_FLAG int
type MAV_STATE int
type MAV_SYS_STATUS_SENSOR int
type MAV_CMD int

type MessageTest5 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest5) GetId() uint32 {
	return 5
}

type MessageTest6 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest6) GetId() uint32 {
	return 0x0607
}

type MessageTest8 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest8) GetId() uint32 {
	return 8
}

type MessageHeartbeat struct {
	Type           MAV_TYPE      `mavenum:"uint8"`
	Autopilot      MAV_AUTOPILOT `mavenum:"uint8"`
	BaseMode       MAV_MODE_FLAG `mavenum:"uint8"`
	CustomMode     uint32
	SystemStatus   MAV_STATE `mavenum:"uint8"`
	MavlinkVersion uint8
}

func (*MessageHeartbeat) GetId() uint32 {
	return 0
}

type MessageRequestDataStream struct {
	TargetSystem    uint8
	TargetComponent uint8
	ReqStreamId     uint8
	ReqMessageRate  uint16
	StartStop       uint8
}

func (*MessageRequestDataStream) GetId() uint32 {
	return 66
}

type MessageSysStatus struct {
	OnboardControlSensorsPresent MAV_SYS_STATUS_SENSOR `mavenum:"uint32"`
	OnboardControlSensorsEnabled MAV_SYS_STATUS_SENSOR `mavenum:"uint32"`
	OnboardControlSensorsHealth  MAV_SYS_STATUS_SENSOR `mavenum:"uint32"`
	Load                         uint16
	VoltageBattery               uint16
	CurrentBattery               int16
	BatteryRemaining             int8
	DropRateComm                 uint16
	ErrorsComm                   uint16
	ErrorsCount1                 uint16
	ErrorsCount2                 uint16
	ErrorsCount3                 uint16
	ErrorsCount4                 uint16
}

func (m *MessageSysStatus) GetId() uint32 {
	return 1
}

type MessageChangeOperatorControl struct {
	TargetSystem   uint8
	ControlRequest uint8
	Version        uint8
	Passkey        string `mavlen:"25"`
}

func (m *MessageChangeOperatorControl) GetId() uint32 {
	return 5
}

type MessageAttitudeQuaternionCov struct {
	TimeUsec   uint64
	Q          [4]float32
	Rollspeed  float32
	Pitchspeed float32
	Yawspeed   float32
	Covariance [9]float32
}

func (m *MessageAttitudeQuaternionCov) GetId() uint32 {
	return 61
}

type MessageOpticalFlow struct {
	TimeUsec       uint64
	SensorId       uint8
	FlowX          int16
	FlowY          int16
	FlowCompMX     float32
	FlowCompMY     float32
	Quality        uint8
	GroundDistance float32
	FlowRateX      float32 `mavext:"true"`
	FlowRateY      float32 `mavext:"true"`
}

func (*MessageOpticalFlow) GetId() uint32 {
	return 100
}

type MessagePlayTune struct {
	TargetSystem    uint8
	TargetComponent uint8
	Tune            string `mavlen:"30"`
	Tune2           string `mavext:"true" mavlen:"200"`
}

func (*MessagePlayTune) GetId() uint32 {
	return 258
}

type MessageAhrs struct {
	OmegaIx     float32 `mavname:"omegaIx"`
	OmegaIy     float32 `mavname:"omegaIy"`
	OmegaIz     float32 `mavname:"omegaIz"`
	AccelWeight float32
	RenormVal   float32
	ErrorRp     float32
	ErrorYaw    float32
}

func (*MessageAhrs) GetId() uint32 {
	return 163
}

type MessageTrajectoryRepresentationWaypoints struct {
	TimeUsec    uint64
	ValidPoints uint8
	PosX        [5]float32
	PosY        [5]float32
	PosZ        [5]float32
	VelX        [5]float32
	VelY        [5]float32
	VelZ        [5]float32
	AccX        [5]float32
	AccY        [5]float32
	AccZ        [5]float32
	PosYaw      [5]float32
	VelYaw      [5]float32
	Command     [5]MAV_CMD `mavenum:"uint16"`
}

func (*MessageTrajectoryRepresentationWaypoints) GetId() uint32 {
	return 332
}
