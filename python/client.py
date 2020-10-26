import pandas as pd
from influxdb_client import InfluxDBClient, Point
from influxdb_client.client.write_api import SYNCHRONOUS
import datetime

DATETIME_NOW = datetime.datetime.now(tz=datetime.timezone.utc).isoformat()

def make_query_api(url="http://localhost:8086", token="", org=""):
    client = InfluxDBClient(url=url, token=token, org=org)
    query_api = client.query_api()
    return query_api

def get_citibike_data(query_api, start, end=DATETIME_NOW):
    """pull bike status data from database"""
    bike_df = query_api.query_data_frame('from(bucket:"bikedb") '
                                        '|> range(start: {}, stop: {}) '
                                        '|> filter(fn: (r) => r._measurement == "bike_status_at_station" and (r._field == "num_ebikes_available" or r._field == "station_id"))'
                                        '|> pivot(rowKey:["_time", "station_id"], columnKey: ["_field"], valueColumn: "_value") '
                                        '|> keep(columns: ["num_ebikes_available", "num_docks_available", "_time", "station_id"])'.format(start,end))
    bike_df.drop(columns=["result","table"],inplace=True,  errors="ignore")
    return bike_df

def get_station_data(query_api, start, end=DATETIME_NOW):
    """pull station data from database"""
    stations = query_api.query_data_frame('from(bucket:"bikedb") '
                                        '|> range(start: {}, stop: {}) '
                                        '|> filter(fn: (r) => r._measurement == "stations")'
                                        '|> pivot(rowKey:["_time", "station_id"], columnKey: ["_field"], valueColumn: "_value") '
                                        '|> keep(columns: ["capacity","name","_time", "station_id", "lat", "lon"])'.format(start,end))
    stations.drop_duplicates(subset=["name", "station_id"], inplace=True)
    stations.drop(columns=["_time", "result", "table"], inplace=True, errors="ignore")
    return stations

def combine_bike_and_station_data(bike_df, stations, localize=True):
    """combine and clean bike status data and station data.
    if localize is set to True, timestamp will be converted to EST from UTC."""
    df = pd.merge(left=bike_df, right=stations, how='left', left_on='station_id', right_on='station_id')
    df = df.astype({"num_ebikes_available": "int32", "name": "string"})
    if localize:
        df["_time"] = df["_time"].dt.tz_convert(tz="EST")
    df["_time"]=df["_time"].dt.strftime("%Y-%m-%d %H:%M:%S")
    return df