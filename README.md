# roofwater
Run my sprinkler at regular intervals when it's too hot outside.

I put a simple sprinkler system on my roof with a cheap tasmota powered water valve controlling it. This software interfaces with the tasmota valve and runs it at a preprogrammed interval whenever the temperature outside is above a setpoint. This requires a [OpenWeather API Key](https://openweathermap.org/).

## Usage
```
This application is configured via the environment. The following environment
variables can be used:

KEY                      TYPE        DEFAULT    REQUIRED    DESCRIPTION
RW_VALVE_IP              String                 true
RW_PULSEWIDTH            Duration    15s
RW_PULSEINTERVAL         Duration    5m
RW_MINTEMP               Float       90
RW_OWM_API_KEY           String                 true
RW_OWM_ZIP               String                 true
RW_OWM_COUNTRY           String                 true
RW_OWM_CACHE_DURATION    Duration    1h
```
