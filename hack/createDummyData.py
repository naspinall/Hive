import requests
import random
import json

basePath = "http://localhost:3001/api/"
places = json.load(open("./places.json"))
industries = json.load(open("./types.json"))
types = json.load(open("./industries.json"))

for x in range(0, 25):
    place = random.choice(places)
    industry = random.choice(industries)
    kind = random.choice(types)
    for y in range(0, 10):
        r= requests.post(
            "http://localhost:3001/api/devices/",
            json={
                "name": " ".join([place["suburb"], industry, kind, str(y)]),
                "IMEI": str(random.randint(1, 999999999999999)),
                "latitude": -random.uniform(0.001, 0.000001) + place["latitude"],
                "longitude": random.uniform(0.001, 0.000001) + place["longitude"],
            },headers={"Authorization": "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImV4cCI6MTU5MDE0MjU0NSwiaXNzIjoiSGl2ZSIsInJvbGUiOnsiSUQiOjIsIkNyZWF0ZWRBdCI6IjIwMjAtMDUtMjFUMTM6MTA6NTMuNTA2MjM1WiIsIlVwZGF0ZWRBdCI6IjIwMjAtMDUtMjFUMTM6MTc6MTkuMjA0Mzc0WiIsIkRlbGV0ZWRBdCI6bnVsbCwiYWxhcm1zIjoxLCJ1c2VycyI6MCwibWVhc3VyZW1lbnRzIjowLCJkZXZpY2VzIjo0LCJzdWJzY3JpcHRpb25zIjowLCJ1c2VySUQiOjF9fQ.baOW3eGznzKj9VQw-Q5OMaoTO42dUtbLcJCyknU4VJDfh2yfqQuOJIR8Zg9XtIzAqyydzYCcFUebYUsC3rUrGg"}
        )

# for x in range(0, 100):
#     r = requests.post(
#         "/".join(["http://localhost:3001/api/devices",str(random.randint(1,50)),"alarms"]),
#         json={
#             "Type" : random.choice(["Security","Temperature","Flow","Humidity"]),
#             "Status" : random.choice(["ACTIVE", "ACKNOWLEDGED", "CLEARED"]),
#             "Severity" : random.choice(["MINOR","MAJOR","SEVERE"]),
#         }
#     )
#     print("/".join(["http://localhost:3001/api/devices",str(random.randint(1,50)),"alams"]))
#     print(r.status_code)
# -34.990608, 138.564841
# -34.990261, 138.564822
