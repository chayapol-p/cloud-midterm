import requests
import pandas as pd
import sys
import os
from datetime import timezone
import datetime
import json
from time import perf_counter
import numpy as np


def find(name):
    for root, dirs, files in os.walk(os.getcwd()):
        if name in files:
            return os.path.join(root, name)


def update_value(df_merge, old_author, old_message, old_likes, new_author, new_message, new_likes):
    df_merge.loc[pd.isnull(new_author),
                 'author'] = old_author[pd.isnull(new_author)]
    df_merge.loc[pd.isnull(new_message),
                 'message'] = old_message[pd.isnull(new_message)]
    df_merge.loc[pd.isnull(new_likes),
                 'likes'] = old_likes[pd.isnull(new_likes)]


def sync(endpoint):
    t1_start = perf_counter()
    path = find('response.csv')
    data_saved = pd.DataFrame(columns=['uuid', 'author', 'message', 'likes'])
    timestamp = '0001-01-01T00:00:00'

    if path is not None:
        data_saved = pd.read_csv(path, dtype={'uuid': str},
                                 names=["uuid", "author", "message", "likes"])
        print("reading file...")
        # data_saved.set_index('uuid', inplace=True)
        timestamp = open("timestamp.txt", 'r').read()

    f = open('timestamp.txt', "w")
    f.write(datetime.datetime.now(timezone.utc).isoformat())
    f.close()

    i = 0
    limit = 70000

    updates, messages, deletes = [], [], []
    while True:
        print('Getting Messages...')
        query = {'offset': i, 'limit': limit, 'table': 'messages'}
        i += limit
        json_response = requests.get(
            url=endpoint+timestamp, params=query)  # send param

        response = json_response.json()
        messages += response['query_data']

        tc_stop = perf_counter()
        print("loop:", i, ":Elapsed time during the whole program in seconds:",
              tc_stop-t1_start)

        if len(response['query_data']) == 0:
            break

    while True:
        print('Getting Updates...')
        query = {'offset': i, 'limit': limit, 'table': 'updates'}
        i += limit
        json_response = requests.get(
            url=endpoint+timestamp, params=query)  # send param

        response = json_response.json()
        updates += response['query_data']

        tc_stop = perf_counter()
        print("loop:", i, ":Elapsed time during the whole program in seconds:",
              tc_stop-t1_start)

        if len(response['query_data']):
            break

    while True:
        print('Getting Deleted...')
        query = {'offset': i, 'limit': limit, 'table': 'deletes'}
        i += limit
        json_response = requests.get(
            url=endpoint+timestamp, params=query)  # send param

        response = json_response.json()
        deletes += response['query_data']

        tc_stop = perf_counter()
        print("loop:", i, ":Elapsed time during the whole program in seconds:",
              tc_stop-t1_start)

        if len(response['query_data']):
            break

    t2_stop = perf_counter()
    print("reciving:Elapsed time during the whole program in seconds:",
          t2_stop-t1_start)

    # add new message
    data_saved = pd.concat(
        [data_saved, pd.DataFrame(messages)])

    # deleted message
    data_saved.drop(
        data_saved[data_saved['uuid'].isin(deletes)].index, inplace=True)

    # update message
    df_update = pd.DataFrame(
        updates, columns=['uuid', 'author', 'message', 'likes'])

    data_saved = data_saved.merge(df_update[[
                                'uuid', 'author', 'message', 'likes']], on='uuid', how='left', suffixes=('_', ''))
    data_saved.replace('', np.nan, inplace=True)
    data_saved.replace(-1.0, np.nan, inplace=True)

    update_value(data_saved, data_saved['author_'].values, data_saved['message_'].values, data_saved['likes_'].values,
                 data_saved['author'].values, data_saved['message'].values, data_saved['likes'].values)
    data_saved.drop(columns=['author_', 'message_', 'likes_'], inplace=True)
    data_saved['likes'] = data_saved['likes'].astype(int)

    t3_stop = perf_counter()
    print("merging:Elapsed time during the whole program in seconds:",
          t3_stop-t1_start)

    data_saved.to_csv('response.csv', header=False, index=False)

    t1_stop = perf_counter()
    print("Elapsed time during the whole program in seconds:",
          t1_stop-t1_start)


def main(url):
    sync(url)


main(sys.argv[1])
