import requests
import pandas as pd
import sys
import os
from datetime import timezone
import datetime
import json
from time import perf_counter

# T = '2022-03-18T09:50:10.568317+00:00'
# first
# response = {'messages': [{'uuid': '0123456710', 'author': 'testname1', 'message': 'testmessage1', 'likes': 0},
#                          {'uuid': '0123456711', 'author': 'testname2', 'message': 'testmessage2',
#                           'likes': 1},
#                          {'uuid': '0123456712', 'author': 'testname3', 'message': 'testmessage3',
#                           'likes': 2},
#                          {'uuid': '0123456713', 'author': 'testname4', 'message': 'testmessage4',
#                           'likes': 3},
#                          {'uuid': '0123456714', 'author': 'testname5', 'message': 'testmessage5', 'likes': 4}],
#             'updates': []}

# only update
# response = {'messages': [],
#             'updates': [{'uuid': '0123456712', 'author': '', 'message': 'asdf3', 'likes': 3, 'is_deleted': 0},
#                         {'uuid': '0123456713', 'author': '', 'message': 'yqwef4', 'likes': 4, 'is_deleted': 0},
#                         {'uuid': '0123456711', 'author': '', 'message': '', 'likes': -1, 'is_deleted': 0}]}

# only delete
# response = {'messages': [],
#             'updates': [{'uuid': '0123456712', 'author': '', 'message': 'asdf3', 'likes': 3, 'is_deleted': 1},
#                         {'uuid': '0123456713', 'author': '', 'message': 'yqwef4',
#                          'likes': 4, 'is_deleted': 1},
#                         {'uuid': '0123456711', 'author': '', 'message': '', 'likes': -1, 'is_deleted': 1}]}

#create and update
# response = {'messages': [{'uuid': '0123456715', 'author': 'testname6', 'message': 'testmessage6', 'likes': 0}],
#             'updates': [{'uuid': '0123456714', 'author': '', 'message': 'krthdf5', 'likes': 10, 'is_deleted': 0}]}

#create and delete
# response = {'messages': [{'uuid': '0123456716', 'author': 'testname7', 'message': 'testmessage7', 'likes': 0}],
#             'updates': [{'uuid': '0123456714', 'author': '', 'message': 'krthdf5', 'likes': 10, 'is_deleted': 1}]}


def find(name):
    for root, dirs, files in os.walk(os.getcwd()):
        if name in files:
            return os.path.join(root, name)


def update_n_delete(df_saved, update):
    if not update:
        return df_saved

    df_update = pd.DataFrame(update)
    df_update.set_index('uuid', inplace=True)

    for index, row in df_update.iterrows():
        if row['is_deleted'] == 1:
            if index in df_saved.index:
                df_saved.drop(index, inplace=True)
            continue
        if len(str(row['message'])) > 0:
            df_saved.loc[index, 'message'] = row['message']
        if row['likes'] != -1:
            df_saved.loc[index, 'likes'] = row['likes']

    return df_saved


def add_new(df_saved, create):
    if not create:
        return df_saved

    df_create = pd.DataFrame(create)
    df_create.set_index('uuid', inplace=True)

    return pd.concat([df_saved, df_create])


def sync(endpoint):
    path = find('response.csv')
    data_saved = pd.DataFrame()
    timestamp = '0001-01-01T00:00:00'

    if path is not None:
        data_saved = pd.read_csv(path, dtype={'uuid': str},
                                 names=["uuid", "author", "message", "likes"])

        data_saved.set_index('uuid', inplace=True)
        timestamp = open("timestamp.txt", 'r').read()

    f = open('timestamp.txt', "w")
    f.write(datetime.datetime.now(timezone.utc).isoformat())
    f.close()

    # params = {'timestamp': timestamp}
    json_response = requests.get(url=endpoint+timestamp)  # send param
    response = json_response.json()

    data_saved = add_new(data_saved, response['messages'])
    data_saved = update_n_delete(data_saved, response['updates'])
    data_saved.to_csv('response.csv', header=False, index=True)


def main(url):
    t1_start = perf_counter()
    sync(url)
    t1_stop = perf_counter()
    print("Elapsed time during the whole program in seconds:",
          t1_stop-t1_start)


main(sys.argv[1])
