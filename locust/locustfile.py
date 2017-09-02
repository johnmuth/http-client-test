from locust import HttpLocust, TaskSet

def api(l):
    l.client.get("/api")

class UserBehavior(TaskSet):
    tasks = {api: 5}

    def on_start(self):
        login(self)

class WebsiteUser(HttpLocust):
    task_set = UserBehavior
    min_wait = 5
    max_wait = 500