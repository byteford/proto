import click_pb2
import requests

def decodeUser(res):
    if res.status_code == 400:
        err = click_pb2.Error()
        err.ParseFromString(res.content)
        res.close()
        print(err)
    if res.status_code == 200:
        c = click_pb2.User()
        c.ParseFromString(res.content)
        res.close()
        return c

def NewUser(name):
    print("making user: ", name)
    c = click_pb2.User()
    c.name = name
    res = requests.post("http://127.0.0.1:3000/new", data=c.SerializeToString())
    return decodeUser(res)

def GetUserById(id):
    res = requests.get("http://127.0.0.1:3000/", params={"id":id})
    return decodeUser(res)

def GetUserByName(name):
    res = requests.get("http://127.0.0.1:3000/", params={"name":name})
    return decodeUser(res)

def GetUsers():
    res = requests.get("http://127.0.0.1:3000/")
    users = click_pb2.Users()
    users.ParseFromString(res.content)
    res.close()
    u = []
    for user in users.user:
        u.append(user)
    return u

print("what is your username:")
username = input()
user = GetUserByName(username)
if user is None:
    user = NewUser(username)
print("Connected as: ", user.name)
#users = GetUsers()
#GetUser(user.id)



