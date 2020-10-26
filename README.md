### gobike
- Pull real time data from Citibike (api?) and stored them in InfluxDB(v1.8.x).
- Visualize data with plotly.
- Except for python visualization code, the whole process can be dockerized.
- [Click here](https://www.citibikenyc.com/system-data) for information on data.

##### Why?
Had a conversation with a friend regarding e-bike availability in nyc.

### Setup
#### influxdb conifg
To use inlfuxdb client 2.0, flux need to be enabled; note that the latest version of official docker image for influxdb is v1.8.x. Influxdb config file `influxdb.conf` is located in `/etc/influxdb/`. Set `flux-enabled=true` in `[http]` section.

#### getting data
Set up .env file as in .env-example and build via docker-compose command.

#### visualization in python
Mapbox's NYC map is used to generate graph with plotly. See [here](https://www.mapbox.com/studio) to get mapboxtoken.
```
import client
import viz

query_api = client.make_query_api(url="http://localhost:8086", token="", org="") #create influxdb query client and 
bike = client.get_citibike_data(query_api, "2020-10-04T19:17:00+00:00", "2020-10-04T19:37:00+00:00") #get bike status data
station = client.get_station_data(query_api, "2020-10-03T18:17:00+00:00") # get station data
df = client.combine_bike_and_station_data(bike, station) #see sample.csv for data

#visualizing
fig = viz.visualize_data(df, ".mapbox_token", "citibike_map.html", return_figure=True)
fig.show()
```

![Citibike](/python/citibike_map.gif)


