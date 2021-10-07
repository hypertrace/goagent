# change -u to change number of concurrent users, -r to change how many users are spawned off at a time and -t to
# change how long the test runs for.
#
# locust -f bigfoorequest_locust.py -u 20 -r 100 -H http://localhost:8081 -t 180s --headless
import time, json, random
from locust import HttpUser, task, between

# header = {'Content-type':'application/json, */*;q=0.5', 'Accept':'application/json'}
header = {'Content-type':'application/json', 'Accept':'application/json'}

# JSON file
f = open('request16K.json', "r")
bigfoo_request_json_str = f.read()
#print(bigfoo_request_json_str)
# bigfoo_request_data = json.dumps(f.read())


# curl --silent -XPOST http://localhost:8081/bigfoorequest -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json' \
#    -d @request.json > /dev/null
# curl --silent -XPOST http://localhost:8081/bigfoo -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'  -d '{"name": "Anomandaris"}' > /dev/null
class QuickstartUser(HttpUser):

    @task
    def on_start(self):
        # time.sleep(0.005)
        self.client.post("/bigfoorequest", data=bigfoo_request_json_str, headers=header)

