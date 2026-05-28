#!/usr/bin/env python3

import os
import re
import sys
from html import unescape
from html.parser import HTMLParser

import requests

if not os.environ.get("TAURON_USER"):
    print("Error: TAURON_USER is not set", file=sys.stderr)
    sys.exit(1)
if not os.environ.get("TAURON_PASS"):
    print("Error: TAURON_PASS is not set", file=sys.stderr)
    sys.exit(1)

USERNAME = os.environ["TAURON_USER"]
PASSWORD = os.environ["TAURON_PASS"]

LOGIN_URL = "https://logowanie.tauron-dystrybucja.pl/login?service=https://elicznik.tauron-dystrybucja.pl"
DATA_URL = "https://elicznik.tauron-dystrybucja.pl/energia/do/dane"
REF = "https://elicznik.tauron-dystrybucja.pl/energia"
UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"


class FormActionParser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.action_url = None

    def handle_starttag(self, tag, attrs):
        if tag == "form" and self.action_url is None:
            attrs = dict(attrs)
            if attrs.get("id") == "kc-form-login" or attrs.get("method", "").lower() == "post":
                self.action_url = unescape(attrs.get("action", ""))


session = requests.Session()
session.headers.update({"User-Agent": UA})

# Step 1: GET login page, extract form action URL
response = session.get(LOGIN_URL)
parser = FormActionParser()
parser.feed(response.text)
action_url = parser.action_url

# Step 2: POST credentials
session.post(
    action_url,
    data={"username": USERNAME, "password": PASSWORD, "credentialId": "", "login": ""},
)


def fetch_csv(from_date, to_date, output_file):
    params = {
        "form[from]": from_date,
        "form[to]": to_date,
        "form[type]": "godzin",
        "form[energy][consum]": "1",
        "form[energy][oze]": "1",
        "form[fileType]": "CSV",
    }
    r = session.get(DATA_URL, params=params, headers={"Referer": REF})
    with open(output_file, "wb") as f:
        f.write(r.content)
    print(f"Saved {output_file}")


fetch_csv("15.10.2021", "21.10.2021", "sample.csv")
fetch_csv("30.10.2022", "30.10.2022", "sample-CEST-to-CET.csv")
fetch_csv("27.03.2022", "27.03.2022", "sample-CET-to-CEST.csv")


def filter_lines(input_file, output_file, exclude):
    with open(input_file) as f:
        lines = f.readlines()
    with open(output_file, "w") as f:
        f.writelines(line for line in lines if exclude not in line)


def filter_lines_regex(input_file, output_file, pattern):
    compiled = re.compile(pattern)
    with open(input_file) as f:
        lines = f.readlines()
    with open(output_file, "w") as f:
        f.writelines(line for line in lines if not compiled.search(line))


filter_lines("sample.csv", "sample-missing-data-on-day1.csv", "2021-10-15 2:00")
filter_lines("sample.csv", "sample-missing-data-on-day4.csv", "2021-10-18 14:00")
filter_lines("sample.csv", "sample-missing-data-on-day7.csv", "2021-10-21 20:00")
filter_lines("sample.csv", "sample-missing-day1.csv", "2021-10-15")
filter_lines("sample.csv", "sample-missing-day4.csv", "2021-10-18")
filter_lines("sample.csv", "sample-missing-day7.csv", "2021-10-21")
filter_lines_regex("sample.csv", "sample-different-array-sizes.csv", r"2021-10-16 2:00;.*;pobór;")
