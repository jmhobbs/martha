# Architecture

Martha is designed around control loops to regulate the state of the fruiting chambers.

The sensor systems are independent, feeding state values, which the controllers observe and act upon.

The state from the sensors is the current state of the chamber.  The configuration for the chamber is the desired state.

When the current state is outside the parameters defined by the desired state, the controller will take action upon the associated control to attempt to bring the current state into spec.

```
                       +--------------------------------------------------------------------------------------------+
                       |                                                                                            |
                       |                                   Configuration Provider                                   |
                       |                                                                                            |
                       +--------------------------------------------------------------------------------------------+

+------------------+   +-----------+   +--------------------------+   +---------------------------------------------+
|                  |   |           |   |                          |   |                                             |
|     tstorage     |<--|           |<--|     Sample Collector     |-->|              Hardware Manager               |
|                  |   |           |   |                          |   |                                             |
+------------------+   |           |   +-------------+------------+   +---------------------------------------------+
+------------------+   |           |                 |                      |                             |
|                  |   |  Storage  |                 v                      v                             v
|      Influx      |<--|  Engine   |   +--------------------------+   +-----------+   +---------+   +-----------+   +---------+
|                  |   |           |   |                          |   |           |   |         |   |           |   |         |
+------------------+   |           |   |   Chamber Controllers    |   |           |-->| BME280  |-->|           |-->|   I2C   |
                       |           |   |                          |   |           |   |         |   |           |   |         |
                       |           |   +--------------------------+   |           |   +---------+   |           |   +---------+
                       |           |                                  |           |   +---------+   |           |   +---------+
                       |           |                                  |           |   |         |   |           |   |         |
                       +-----------+                                  |           |-->|  DHT11  |-->|           |-->|   SPI   |
                                                                      |           |   |         |   |           |   |         |
                                                                      |  Device   |   +---------+   | periph.io |   +---------+
                                                                      |  Drivers  |   +---------+   |           |   +---------+
                                                                      |           |   |         |   |           |   |         |
                                                                      |           |-->|  Relay  |-->|           |-->|  GPIO   |
                                                                      |           |   |         |   |           |   |         |
                                                                      |           |   +---------+   |           |   +---------+
                                                                      |           |   +---------+   |           |
                                                                      |           |   |         |   |           |
                                                                      |           |-->| Max6675 |-->|           |
                                                                      |           |   |         |   |           |
                                                                      +-----------+   +---------+   +-----------+
```