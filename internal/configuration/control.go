package configuration

type ControlType DeviceType

const (
	CONTROL_TYPE_NORMALLY_OPEN_RELAY ControlType = "NORMALLY_OPEN_RELAY"
	CONTROL_TYPE_NORMALLY_CLOSED_RELAY ControlType = "NORMALLY_CLOSED_RELAY"
)

type ControlEffect string

const (
	CONTROL_EFFECT_TYPE_HUMIDITY ControlEffect = "HUMIDITY"
	CONTROL_EFFECT_TYPE_CO2 ControlEffect = "CO2"
	CONTROL_EFFECT_TYPE_TEMPERATURE ControlEffect = "TEMPERATURE"
)