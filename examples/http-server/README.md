# Prerequisites
- Make sure `locust` is installed. See https://locust.io/
- Commands run assume you are in the same local dir as this file.
- Has been tested on docker-desktop
- Install traceable-agent or make sure it's somehow running. This is important so that you can inject the proxy, but if you have a different way to inject the proxy it's fine.

# Instructions on how to use this
- If you will use k8s make sure you are on the `docker-desktop` k8s context.
- Generate some request bodies by running `./generate_request_bodies.sh`
- Build the go app: `GOOS=linux go build`
- Build the docker image: `docker build -t tmwangi/hackgoapp:0.1.0 .`
- To deploy on k8s run `kubectl apply -f hackgoapp_k8s.yaml`. Exposes 
- To run some load using locust, run `locust -f bigfoorequest_locust.py -u 20 -r 100 -H http://localhost:8081 -t 180s --headless`. You can change the concurrent users by changing the `-u` flag and the number of users spawned at a time by changing `-r`. I've set 100 for that so that I get all the users spawned at once.
