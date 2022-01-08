package configuration

type DeviceType string

type SensorType DeviceType

const (
	SENSOR_TYPE_SYNTHETIC SensorType = "SYNTHETIC"
	SENSOR_TYPE_EXTERNAL SensorType = "EXTERNAL"
	SENSOR_TYPE_BME280 SensorType = "BME280"
	SENSOR_TYPE_MAX6675 SensorType = "MAX6675"
)
