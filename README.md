# roofwater
Run my sprinkler at regular intervals when it's too hot outside.

I put a simple sprinkler system on my roof with a cheap tasmota powered water valve controlling it. This software interfaces with the tasmota valve and runs it at a preprogrammed interval whenever the temperature outside is above a setpoint. This requires a [OpenWeather API Key](https://openweathermap.org/).

## Usage
```
This application is configured via the environment. The following environment
variables can be used:

KEY                        TYPE        DEFAULT       REQUIRED    DESCRIPTION
RW_PULSEWIDTH              Float       0.05                      Percent of the Interval expressed as a decimal to turn valve on
RW_PULSEPERIOD             Duration    5m                        Period of pulses
RW_MINTEMP                 Float       90                        Minimum temperature to run the valve
RW_OWM_API_KEY             String                    true        OpenWeatherMap API key
RW_OWM_ZIP                 String                    true        Zip code to get weather for
RW_OWM_COUNTRY             String                    true        Country to get weather for
RW_OWM_CACHE_DURATION      Duration    1h                        Duration to cache weather data
RW_OWM_UNIT                String      F                         Unit to use for weather data
RW_OWM_LANGUAGE            String      EN                        Language to use for weather data
RW_TRACING_SERVICE_NAME    String      roofwaterd                Service name to use for tracing
RW_VALVE_HTTP_ADDR         String                                HTTP Address of the valve
RW_VALVE_MQTT_TOPIC        String                                MQTT topic to publish to
RW_MQTT_URL                String                                MQTT Broker URL
RW_MQTT_USER               String                                MQTT username
RW_MQTT_PASS               String                                MQTT password

```
