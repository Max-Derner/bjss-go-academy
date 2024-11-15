import requests
from uuid import uuid4
import json

HOME = "http://localhost:8080/"
CREATE = HOME + "create"
READ = HOME + "read"
UPDATE = HOME + "update"
DELETE = HOME + "delete"


def prompt_data() -> dict:
    id = uuid4()
    title = input("Give title: ")
    prio = input("Give priority: ")
    comp = False
    return {
        "id": str(id),
        "title": title,
        "priority": prio,
        "complete": comp,
    }


def main():
    print("What do?")
    while True:
        print("\t1: Create")
        print("\t2: Read")
        print("\t3: Update")
        print("\t4: Delete")
        print("\t0: Exit")
        match input("#: "):
            case "1":
                data = prompt_data()
                from pprint import pprint
                pprint(data)
                r = requests.post(
                    url=CREATE,
                    data=json.dumps(data),
                )
                print(r.content.decode("utf-8"))
            case "2":
                r = requests.get(
                    url=READ,
                )
                print(r.content.decode("utf-8"))
            case "3":
                requests.post(
                    url=UPDATE,
                )
            case "4":
                requests.post(
                    url=DELETE,
                )
            case "0":
                exit(0)


if __name__ == "__main__":
    main()
