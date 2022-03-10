import click_pb2
import requests

c = click_pb2.Click()
c.id = 1
c.name = "James"

print(c.SerializeToString())

res = requests.post("http://127.0.0.1:3000/new", data=c.SerializeToString())
print(res)
print(res.content)
if res.status_code == 400:
    err = click_pb2.Error()
    err.ParseFromString(res.content)
    res.close()
    print(err)
if res.status_code == 200:
    c = click_pb2.Click()
    c.ParseFromString(res.content)
    res.close()
    print(c)