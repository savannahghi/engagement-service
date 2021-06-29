"""Send bulk SMS to a campaign."""
import base64
import enum
import json
import logging
import os
import time
from datetime import datetime

import requests

logger = logging.getLogger(__name__)

BASE_URL = os.getenv("BASE_URL")
TRACKING_URL = os.getenv("TRACKING_URL")
FIREBASE_WEB_API_KEY = os.getenv("FIREBASE_WEB_API_KEY")
ANDROID_PACKAGE_NAME = os.getenv("ANDROID_PACKAGE_NAME")
IOS_BUNDLE_ID = os.getenv("IOS_BUNDLE_ID")
DOMAIN_URI_PREFIX = os.getenv("DOMAIN_URI_PREFIX")
FIREBASE_DYNAMIC_LINK_URL = os.getenv("FIREBASE_DYNAMIC_LINK_URL")

# todo change this message for the different segments
MESSAGE = ""

# todo change the segment name after first run
SEGMENT_NAME = ""


class SenderID(enum.Enum):
    """SenderID enum values."""

    BeWell = "BEWELL"
    Slade360 = "SLADE360"


def generate_shortened_dynamic_links(long_link):
    """Generate a shortened Firebase Dynamic Link from the tracking URL."""
    headers = {"Content-Type": "application/json"}
    params = {
        "dynamicLinkInfo": {
            "domainUriPrefix": DOMAIN_URI_PREFIX,
            "link": long_link,
            "androidInfo": {
                "androidPackageName": ANDROID_PACKAGE_NAME,
                "androidFallbackLink": long_link,
            },
            "iosInfo": {
                "iosBundleId": IOS_BUNDLE_ID,
                "iosFallbackLink": long_link,
            },
        },
        "suffix": {"option": "SHORT"},
    }

    resp = requests.post(
        FIREBASE_DYNAMIC_LINK_URL + FIREBASE_WEB_API_KEY,
        data=json.dumps(params),
        headers=headers,
    )
    if resp.status_code != 200:
        raise Exception(
            "unable to shorten link with status code "
            f"{resp.status_code} and data {resp.content}"
        )

    result = resp.json()
    time.sleep(2)
    return result["shortLink"]


def generate_marketing_url(identifier):
    """Generate tracking URL."""
    request = requests.models.PreparedRequest()
    bs = bytes(identifier, "utf-8")
    encoded_identifier = base64.b64encode(bs)
    params = {"email": encoded_identifier}
    request.prepare_url(TRACKING_URL, params)
    return request.url


def send_sms(payload):
    """Helper function to send the actual SMS."""
    url = BASE_URL + "send_marketing_sms"
    response = requests.post(url=url, json=payload)
    if response.status_code > 299:
        raise Exception(
            "unable to send SMS with status code "
            f"{response.status_code} and data {response.content}"
        )
    return


def current_time():
    """Return the current time"""
    return datetime.now()


def convert_datetime_to_hours(date_time):
    """Convert a date time to hours"""
    return date_time / 3600


def get_segmented_contacts(segment):
    """
    Get segmented contacts details from a data store.
    """
    headers = {"Content-Type": "application/json"}
    url = BASE_URL + "marketing_data"
    payload = {
        "wing": segment,
    }
    response = requests.post(url=url, json=payload, headers=headers)
    if response.status_code > 299:
        raise Exception(
            "unable to get marketing data with status code "
            f"{response.status_code} and data {response.content}"
        )
    return response.json()


def send_marketing_bulk_sms(request):
    """
    Send bulk SMS.

    The call is made to our engagement service to send bulk SMS
    to our segments using either our BeWell or Slade360 sender
    """
    contacts = get_segmented_contacts(SEGMENT_NAME)
    if contacts is None:
        raise Exception("No contacts found")

    phone_message_list = []
    for contact in contacts:
        phone = contact["phone"]
        first_name = contact["firstname"]
        payer_name = contact["lastname"]
        email = contact["email"]

        phone_message_dict = {
            "phone_number": phone,
            "message": MESSAGE.format(
                first_name,
                payer_name,
                generate_shortened_dynamic_links(
                    generate_marketing_url(email)
                ),
            ),
        }
        phone_message_list.append(phone_message_dict)

    contact_count = 0

    for data in phone_message_list:
        start_time = current_time()
        payload = {
            "to": [data["phone_number"]],
            "message": data["message"],
            "sender": SenderID.BeWell.value,
        }

        send_sms_start_time = current_time()
        send_sms(payload)
        send_sms_end_time = current_time()

        send_sms_total_time = send_sms_end_time - send_sms_start_time
        sms_rate = f"{send_sms_total_time.total_seconds()}s/message"

        contact_count += 1

        if contact_count % 100 == 0:
            t = (current_time() - start_time).total_seconds()
            time_taken_so_far = convert_datetime_to_hours(t)
            time_left = (
                len(phone_message_list) - contact_count
            ) * send_sms_total_time

            time_left_in_hr = convert_datetime_to_hours(
                time_left.total_seconds()
            )
            logger.warning(
                f"{contact_count} contacts marketed to, "
                f"{time_taken_so_far} hours taken so far, "
                f"{sms_rate}, {time_left_in_hr} hours left"
            )

    if contact_count == len(phone_message_list):
        logger.warning(
            f"{len(phone_message_list)} contacts engaged successfully!"
        )


def main():
    send_marketing_bulk_sms("")


if __name__ == "__main__":
    main()
