import requests
import pandas as pd
import sys
import os
from datetime import timezone
import datetime

{
    "uuid": "28503f06-a608-11ec-b909-0242ac120702",
    "created_at": "2022-03-17T22:56:38.1985302+07:00",
    "updated_at": "2022-03-17T22:56:38.1985311+07:00",
    "message": "First Message",
    "author": "job",
    "likes": 0,
    "is_deleted": 0
}

# moc response first sync
response = [{'uuid': '0123456789', 'author': 'testname', 'message': 'testmessage', 'like': 0, 'time_stamp': 1647529287.822895},
            {'uuid': '0123456710', 'author': 'testname1',
                'message': 'testmessage1', 'like': 1, 'time_stamp': 1647529287.822895},
            {'uuid': '0123456711', 'author': 'testname2',
                'message': 'testmessage2', 'like': 2, 'time_stamp': 1647529287.822895},
            {'uuid': '0123456712', 'author': 'testname3', 'message': 'testmessage3', 'like': 3, 'time_stamp': 1647529287.822895}]

# moc response update sync
# response = [{'uuid': '0123456789', 'author': '', 'message': 'fsd', 'like': '', 'time_stamp': 1647530832.710686},
#             {'uuid': '0123456710', 'author': 'testname', 'message': 'gqs', 'like': '3', 'time_stamp': 1647530832.710686}]

endpoint = ''
headers = {'Content-type': 'application/json; charset=utf-8'}


def find(name):
    for root, dirs, files in os.walk(os.getcwd()):
        if name in files:
            return os.path.join(root, name)


def sync():
    path = find('response.csv')
    if path is not None:
        print("in")
        data = pd.read_csv(path, dtype={'uuid': str},
                           names=["uuid", "author", "message", "like", "time_stamp"])
        params = {'timestamp': data['time_stamp'].max()}
        # response = request.get(Endpoint,params=param,header)
        df = pd.DataFrame(response)
        for index, row in df.iterrows():
            # print(row['uuid'])
            if len(data[data['uuid'] == row['uuid']]):
                if len(str(row['message'])) > 0:
                    data.loc[data['uuid'] == row['uuid'],
                             'message'] = row['message']
                if len(str(row['like'])) > 0:
                    data.loc[data['uuid'] == row['uuid'], 'like'] = row['like']
                data.loc[data['uuid'] == row['uuid'],
                         'time_stamp'] = row['time_stamp']
            else:
                data.append(row, ignore_index=True)
        df = data
    else:
        params = {'timestamp': 0}
        # response = request.get(Endpoint,params=param,header)
        df = pd.DataFrame(response)
    df.to_csv(
        'response.csv', header=False, index=False)


def main(url):
    # endpoint = url
    # sync()
    print(datetime.datetime.now().isoformat())
    print(datetime.datetime.now(timezone.utc).isoformat())


main(sys.argv[1])
