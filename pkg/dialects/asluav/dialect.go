// Package asluav contains the asluav dialect.
//
//autogenerated:yes
package asluav

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

// Dialect contains the dialect definition.
var Dialect = dial

// dial is not exposed directly in order not to display it in godoc.
var dial = &dialect.Dialect{
	Version: 3,
	Messages: []message.Message{
		// minimal
		&MessageHeartbeat{},
		&MessageProtocolVersion{},
		// standard
		// common
		&MessageSysStatus{},
		&MessageSystemTime{},
		&MessagePing{},
		&MessageChangeOperatorControl{},
		&MessageChangeOperatorControlAck{},
		&MessageAuthKey{},
		&MessageLinkNodeStatus{},
		&MessageSetMode{},
		&MessageParamRequestRead{},
		&MessageParamRequestList{},
		&MessageParamValue{},
		&MessageParamSet{},
		&MessageGpsRawInt{},
		&MessageGpsStatus{},
		&MessageScaledImu{},
		&MessageRawImu{},
		&MessageRawPressure{},
		&MessageScaledPressure{},
		&MessageAttitude{},
		&MessageAttitudeQuaternion{},
		&MessageLocalPositionNed{},
		&MessageGlobalPositionInt{},
		&MessageRcChannelsScaled{},
		&MessageRcChannelsRaw{},
		&MessageServoOutputRaw{},
		&MessageMissionRequestPartialList{},
		&MessageMissionWritePartialList{},
		&MessageMissionItem{},
		&MessageMissionRequest{},
		&MessageMissionSetCurrent{},
		&MessageMissionCurrent{},
		&MessageMissionRequestList{},
		&MessageMissionCount{},
		&MessageMissionClearAll{},
		&MessageMissionItemReached{},
		&MessageMissionAck{},
		&MessageSetGpsGlobalOrigin{},
		&MessageGpsGlobalOrigin{},
		&MessageParamMapRc{},
		&MessageMissionRequestInt{},
		&MessageSafetySetAllowedArea{},
		&MessageSafetyAllowedArea{},
		&MessageAttitudeQuaternionCov{},
		&MessageNavControllerOutput{},
		&MessageGlobalPositionIntCov{},
		&MessageLocalPositionNedCov{},
		&MessageRcChannels{},
		&MessageRequestDataStream{},
		&MessageDataStream{},
		&MessageManualControl{},
		&MessageRcChannelsOverride{},
		&MessageMissionItemInt{},
		&MessageVfrHud{},
		&MessageCommandInt{},
		&MessageCommandLong{},
		&MessageCommandAck{},
		&MessageCommandCancel{},
		&MessageManualSetpoint{},
		&MessageSetAttitudeTarget{},
		&MessageAttitudeTarget{},
		&MessageSetPositionTargetLocalNed{},
		&MessagePositionTargetLocalNed{},
		&MessageSetPositionTargetGlobalInt{},
		&MessagePositionTargetGlobalInt{},
		&MessageLocalPositionNedSystemGlobalOffset{},
		&MessageHilState{},
		&MessageHilControls{},
		&MessageHilRcInputsRaw{},
		&MessageHilActuatorControls{},
		&MessageOpticalFlow{},
		&MessageGlobalVisionPositionEstimate{},
		&MessageVisionPositionEstimate{},
		&MessageVisionSpeedEstimate{},
		&MessageViconPositionEstimate{},
		&MessageHighresImu{},
		&MessageOpticalFlowRad{},
		&MessageHilSensor{},
		&MessageSimState{},
		&MessageRadioStatus{},
		&MessageFileTransferProtocol{},
		&MessageTimesync{},
		&MessageCameraTrigger{},
		&MessageHilGps{},
		&MessageHilOpticalFlow{},
		&MessageHilStateQuaternion{},
		&MessageScaledImu2{},
		&MessageLogRequestList{},
		&MessageLogEntry{},
		&MessageLogRequestData{},
		&MessageLogData{},
		&MessageLogErase{},
		&MessageLogRequestEnd{},
		&MessageGpsInjectData{},
		&MessageGps2Raw{},
		&MessagePowerStatus{},
		&MessageSerialControl{},
		&MessageGpsRtk{},
		&MessageGps2Rtk{},
		&MessageScaledImu3{},
		&MessageDataTransmissionHandshake{},
		&MessageEncapsulatedData{},
		&MessageDistanceSensor{},
		&MessageTerrainRequest{},
		&MessageTerrainData{},
		&MessageTerrainCheck{},
		&MessageTerrainReport{},
		&MessageScaledPressure2{},
		&MessageAttPosMocap{},
		&MessageSetActuatorControlTarget{},
		&MessageActuatorControlTarget{},
		&MessageAltitude{},
		&MessageResourceRequest{},
		&MessageScaledPressure3{},
		&MessageFollowTarget{},
		&MessageControlSystemState{},
		&MessageBatteryStatus{},
		&MessageAutopilotVersion{},
		&MessageLandingTarget{},
		&MessageFenceStatus{},
		&MessageMagCalReport{},
		&MessageEfiStatus{},
		&MessageEstimatorStatus{},
		&MessageWindCov{},
		&MessageGpsInput{},
		&MessageGpsRtcmData{},
		&MessageHighLatency{},
		&MessageHighLatency2{},
		&MessageVibration{},
		&MessageHomePosition{},
		&MessageSetHomePosition{},
		&MessageMessageInterval{},
		&MessageExtendedSysState{},
		&MessageAdsbVehicle{},
		&MessageCollision{},
		&MessageV2Extension{},
		&MessageMemoryVect{},
		&MessageDebugVect{},
		&MessageNamedValueFloat{},
		&MessageNamedValueInt{},
		&MessageStatustext{},
		&MessageDebug{},
		&MessageSetupSigning{},
		&MessageButtonChange{},
		&MessagePlayTune{},
		&MessageCameraInformation{},
		&MessageCameraSettings{},
		&MessageStorageInformation{},
		&MessageCameraCaptureStatus{},
		&MessageCameraImageCaptured{},
		&MessageFlightInformation{},
		&MessageMountOrientation{},
		&MessageLoggingData{},
		&MessageLoggingDataAcked{},
		&MessageLoggingAck{},
		&MessageVideoStreamInformation{},
		&MessageVideoStreamStatus{},
		&MessageCameraFovStatus{},
		&MessageCameraTrackingImageStatus{},
		&MessageCameraTrackingGeoStatus{},
		&MessageCameraThermalRange{},
		&MessageGimbalManagerInformation{},
		&MessageGimbalManagerStatus{},
		&MessageGimbalManagerSetAttitude{},
		&MessageGimbalDeviceInformation{},
		&MessageGimbalDeviceSetAttitude{},
		&MessageGimbalDeviceAttitudeStatus{},
		&MessageAutopilotStateForGimbalDevice{},
		&MessageGimbalManagerSetPitchyaw{},
		&MessageGimbalManagerSetManualControl{},
		&MessageEscInfo{},
		&MessageEscStatus{},
		&MessageWifiConfigAp{},
		&MessageAisVessel{},
		&MessageUavcanNodeStatus{},
		&MessageUavcanNodeInfo{},
		&MessageParamExtRequestRead{},
		&MessageParamExtRequestList{},
		&MessageParamExtValue{},
		&MessageParamExtSet{},
		&MessageParamExtAck{},
		&MessageObstacleDistance{},
		&MessageOdometry{},
		&MessageTrajectoryRepresentationWaypoints{},
		&MessageTrajectoryRepresentationBezier{},
		&MessageCellularStatus{},
		&MessageIsbdLinkStatus{},
		&MessageCellularConfig{},
		&MessageRawRpm{},
		&MessageUtmGlobalPosition{},
		&MessageDebugFloatArray{},
		&MessageOrbitExecutionStatus{},
		&MessageSmartBatteryInfo{},
		&MessageFuelStatus{},
		&MessageBatteryInfo{},
		&MessageGeneratorStatus{},
		&MessageActuatorOutputStatus{},
		&MessageTimeEstimateToTarget{},
		&MessageTunnel{},
		&MessageCanFrame{},
		&MessageOnboardComputerStatus{},
		&MessageComponentInformation{},
		&MessageComponentInformationBasic{},
		&MessageComponentMetadata{},
		&MessagePlayTuneV2{},
		&MessageSupportedTunes{},
		&MessageEvent{},
		&MessageCurrentEventSequence{},
		&MessageRequestEvent{},
		&MessageResponseEventError{},
		&MessageIlluminatorStatus{},
		&MessageCanfdFrame{},
		&MessageCanFilterModify{},
		&MessageWheelDistance{},
		&MessageWinchStatus{},
		&MessageOpenDroneIdBasicId{},
		&MessageOpenDroneIdLocation{},
		&MessageOpenDroneIdAuthentication{},
		&MessageOpenDroneIdSelfId{},
		&MessageOpenDroneIdSystem{},
		&MessageOpenDroneIdOperatorId{},
		&MessageOpenDroneIdMessagePack{},
		&MessageOpenDroneIdArmStatus{},
		&MessageOpenDroneIdSystemUpdate{},
		&MessageHygrometerSensor{},
		// asluav
		&MessageCommandIntStamped{},
		&MessageCommandLongStamped{},
		&MessageSensPower{},
		&MessageSensMppt{},
		&MessageAslctrlData{},
		&MessageAslctrlDebug{},
		&MessageAsluavStatus{},
		&MessageEkfExt{},
		&MessageAslObctrl{},
		&MessageSensAtmos{},
		&MessageSensBatmon{},
		&MessageFwSoaringData{},
		&MessageSensorpodStatus{},
		&MessageSensPowerBoard{},
		&MessageGsmLinkStatus{},
		&MessageSatcomLinkStatus{},
		&MessageSensorAirflowAngles{},
	},
}
