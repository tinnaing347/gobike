import plotly.express as px
import pandas as pd


def prepare_data_for_visualization(df):
    df = df.astype({"name": "string"})
    df['_time'] = pd.to_datetime(df["_time"])
    df.sort_values(by=["_time"], inplace=True)
    return standardize_data(df)

def fill_missing_values(df, missing_time_set):
    """fill ts df with place holder values as per missing_time_set"""
    df = df.reset_index()
    missing_rec = []
    for i in range(len(df)):
        start = df["_time"].iloc[i]
        end = df["_time"].iloc[i+1] if i < len(df)-1 else None
        time_filter = filter(lambda x: x > start and x < end, missing_time_set) if end else filter(lambda x: x > start, missing_time_set)
        default_dict = df[df["_time"] == start].to_dict(orient="records")[0]
        default_dict.pop("_time")
        missing_rec += [ {"_time": i, **default_dict}for i in time_filter]
    return missing_rec

def standardize_data(df):
    """standardize data. Since each bike station updates data at different times,
    the data need to be standardize (i.e fill inconsistent timeseries data with placeholder values) across stations.
    hoping this could be done in influxdb but I failed :("""
    time_set = set(df["_time"])
    results = []
    for _, g in df.groupby(by=["station_id"]):
        individual_time_set = set(g["_time"])
        missing_time_set = list(time_set - individual_time_set)
        results += fill_missing_values(g, missing_time_set)
    result =pd.concat([df, pd.DataFrame(results)])
    result.sort_values(by=["_time"],inplace=True)
    result.reset_index(inplace=True)
    result.drop(columns=["level_0"],inplace=True)
    return result

def visualize_data(df, mapbox_token_path, fig_output_path = "", return_figure = False):
    """visualize  e-bike data using Mapbox's NYC map. figure can be saved as html file.
    If return_figure is set to False, then fig_output_path musth be provided."""
    df = prepare_data_for_visualization(df)
    df = df.astype({"_time": "string"}) #need to be strings 'cause plotly animation_frame demands it
    df.sort_values(by=["_time"], inplace=True)

    px.set_mapbox_access_token(open(mapbox_token_path).read())
    fig = px.scatter_mapbox(df, lat="lat", lon="lon", size="num_ebikes_available", hover_name="name",
                  color_continuous_scale=px.colors.cyclical.IceFire,animation_frame="_time", zoom=10)

    assert fig_output_path == "" or return_figure, "If return_figure is set to False, then fig_output_path musth be provided."
    if fig_output_path:
        fig.write_html(fig_output_path)
    if return_figure:
        return fig