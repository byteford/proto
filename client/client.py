import proto.click_pb2
import requests

c = proto.click_pb2.Click()
c.id = 1
c.name = "James"

print(c.SerializeToString())
f = open ("./temp.txt", "wb")
f.write(c.SerializeToString())
f.close()

res = requests.get("http://127.0.0.1:3000/", params=c.SerializeToString())
print(res.content)