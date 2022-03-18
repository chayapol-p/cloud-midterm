from urllib import response
import requests
import pandas as pd
import sys
import os
from datetime import timezone
import datetime
import json


# T = '2022-03-18T09:50:10.568317+00:00'
# # first
# response = {'create': [{'uuid': '0123456710', 'author': 'testname1', 'message': 'testmessage1', 'like': 0, 'timestamp': T},
#                        {'uuid': '0123456711', 'author': 'testname2', 'message': 'testmessage2',
#                            'like': 1, 'timestamp': T},
#                        {'uuid': '0123456712', 'author': 'testname3', 'message': 'testmessage3',
#                            'like': 2, 'timestamp': T},
#                        {'uuid': '0123456713', 'author': 'testname4', 'message': 'testmessage4',
#                            'like': 3, 'timestamp': T},
#                        {'uuid': '0123456714', 'author': 'testname5', 'message': 'testmessage5', 'like': 4, 'timestamp': T}],
#             'update': []}

# only update
# response = {'create': [],
#             'update': [{'uuid': '0123456712', 'author': '', 'message': 'asdf3', 'like': 3, 'timestamp': T, 'is_deleted': 0},
#                        {'uuid': '0123456713', 'author': '', 'message': 'yqwef4', 'like': 4,
#                            'timestamp': T, 'is_deleted': 0},
#                        {'uuid': '0123456711', 'author': '', 'message': '', 'like': -1, 'timestamp': T, 'is_deleted': 0}]}

# only delete
# response = {'create': [],
#             'update': [{'uuid': '0123456712', 'author': '', 'message': 'asdf3', 'like': 3, 'timestamp': T, 'is_deleted': 1},
#                        {'uuid': '0123456713', 'author': '', 'message': 'yqwef4',
#                            'like': 4, 'timestamp': T, 'is_deleted': 1},
#                        {'uuid': '0123456711', 'author': '', 'message': '', 'like': -1, 'timestamp': T, 'is_deleted': 1}]}

#create and update
# response = {'create': [{'uuid': '0123456715', 'author': 'testname6', 'message': 'testmessage6', 'like': 0, 'timestamp': T}],
#             'update': [{'uuid': '0123456714', 'author': '', 'message': 'krthdf5', 'like': 10, 'timestamp': T, 'is_deleted': 0}]}

#create and delete
# response = {'create': [{'uuid': '0123456716', 'author': 'testname7', 'message': 'testmessage7', 'like': 0, 'timestamp': T}],
#             'update': [{'uuid': '0123456714', 'author': '', 'message': 'krthdf5', 'like': 10, 'timestamp': T, 'is_deleted': 1}]}

endpoint = ''
headers = {'Content-type': 'application/json; charset=utf-8'}


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
        # print(index, row['uuid'], row['is_deleted'])
        if row['is_deleted'] == 1:
            df_saved.drop(index, inplace=True)
            continue
        if len(str(row['message'])) > 0:
            df_saved.loc[index, 'message'] = row['message']
        if row['like'] != -1:
            df_saved.loc[index, 'like'] = row['like']
        df_saved.loc[index, 'timestamp'] = row['timestamp']
    return df_saved


def add_new(df_saved, create):
    if not create:
        return df_saved
    df_create = pd.DataFrame(create)
    df_create.set_index('uuid', inplace=True)
    return pd.concat([df_saved, df_create])


def sync():
    path = find('response.csv')
    data_saved = pd.DataFrame()
    timestamp = '0'
    if path is not None:
        data_saved = pd.read_csv(path, dtype={'uuid': str},
                                 names=["uuid", "author", "message", "like", "timestamp", "is_deleted"])
        data_saved.set_index('uuid', inplace=True)
        timestamp = open("timestamp.txt", 'r').read()
    f = open('timestamp.txt', "w")
    f.write(datetime.datetime.now(timezone.utc).isoformat())
    f.close()
    params = {'timestamp': timestamp}
    # json_response = requests.get(url=endpoint,params=params, headers=headers)
    # response = json.loads(json_response)
    data_saved = add_new(data_saved, response['create'])
    data_saved = update_n_delete(data_saved, response['update'])
    data_saved.to_csv(
        'response.csv', header=False, index=True)


def main(url):
    endpoint = url
    sync()


main(sys.argv[1])
