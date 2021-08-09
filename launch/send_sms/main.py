"""Send bulk SMS to a campaign."""
import base64
import enum
import json
from launch.browser_detect.main import get_authorization_token
import os
import time
from datetime import datetime

import click
import requests


def get_env(var_name):
    """Get an environment variable."""
    var = os.getenv(var_name)
    if var is None:
        click.secho(
            f"Can't find environment {var_name}. "
            "Confirm if you have sourced your environment.",
            fg="red",
            bold=True,
        )
        return
    return var


BASE_URL = get_env("BASE_URL")
EDI_BASE_URL = get_env("EDI_BASE_URL")
FIREBASE_WEB_API_KEY = get_env("FIREBASE_WEB_API_KEY")
ANDROID_PACKAGE_NAME = get_env("ANDROID_PACKAGE_NAME")
IOS_BUNDLE_ID = get_env("IOS_BUNDLE_ID")
DOMAIN_URI_PREFIX = get_env("DOMAIN_URI_PREFIX")
FIREBASE_DYNAMIC_LINK_URL = get_env("FIREBASE_DYNAMIC_LINK_URL")
TRACKING_URL_B = get_env("TRACKING_URL_B")


class SenderID(enum.Enum):
    """SenderID enum values."""

    BeWell = "BEWELL"
    Slade360 = "SLADE360"


class Wings(enum.Enum):
    """Wing names."""

    WingA = "WING A"
    WingB = "WING B"


MESSAGE_A = {
    "message": (
        "Did you know you can now link your {} medical cover "
        "on the Be.Well App and view your benefits? "
        "To get started, Download Now {}. "
        "For more information, call 0790 360 360. "
        "To opt-out dial *384*600# Be.Well by Slade360"
    ),
    "tracking_url": TRACKING_URL_B,
}


MESSAGE_B = {
    "message": (
        "Did you know you can now link your {} medical cover on "
        "the Be.Well App and view your benefits? To get started, "
        "Download Now {}. "
        "For more information, call 0790 360 360. "
        "To opt-out dial *384*600# Be.Well by Slade360"
    ),
    "tracking_url": TRACKING_URL_B,
}


def current_time():
    """Return the current time"""
    return datetime.now()


def convert_datetime_to_hours(date_time):
    """Convert a date time to hours"""
    return date_time / 3600


def send_sms(payload):
    """Helper function to send the actual SMS."""
    url = BASE_URL + "send_marketing_sms"
    token = get_authorization_token()
    headers = {'Authorization': token.strip('"')}
    response = requests.post(url=url, json=payload, headers=headers)
    result = response.json()
    if response.status_code > 299:
        click.secho(
            "\n\nunable to send SMS with status code "
            f"{response.status_code} and data {result}",
            fg="red",
            bold=True,
        )
        raise Exception("unable to send SMS")


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
    result = resp.json()
    if resp.status_code != 200:
        click.secho(
            "\n\nunable to shorten link with status code "
            f"{resp.status_code} and data {result}",
            fg="red",
            bold=True,
        )

    time.sleep(2)
    return result["shortLink"]


def generate_marketing_url(identifier, tracking_url):
    """Generate tracking URL."""
    request = requests.models.PreparedRequest()
    bs = bytes(identifier, "utf-8")
    encoded_identifier = base64.b64encode(bs)
    params = {"email": encoded_identifier}
    request.prepare_url(tracking_url, params)
    return request.url


def get_segmented_contacts(wing, segment):
    """
    Get segmented contacts details from a data store.
    """
    headers = {"Content-Type": "application/json"}
    url = EDI_BASE_URL + "marketing_data"
    params = {
        "wing": wing,
        "segment": segment,
    }
    click.secho(
        "We are about to fetch contacts from our table ... \n", fg="green"
    )
    response = requests.get(url=url, params=params, headers=headers)

    result = response.json()
    if response.status_code > 299:
        click.secho(
            "\n\nunable to get marketing data with status code "
            f"{response.status_code} and data {result}",
            fg="red",
            bold=True,
        )
        return

    return result


def send_marketing_bulk_sms(segment, wing, message_data):
    """
    Send bulk SMS.

    The call is made to our engagement service to send bulk SMS
    to our segments using either our BeWell or Slade360 sender
    """
    click.secho("Launch campaign starts now ...", fg="green")
    start_time = current_time()
    contacts = get_segmented_contacts(wing, segment)
    if contacts is None:
        click.secho(
            "No contacts have been found from your segment.",
            fg="red",
            bold=True,
        )
        click.secho(
            "This can be caused by on of the follwing: \n"
            "\t 1. All messages have been sent to your contacts (If you are re-running with the same arguments).\n"
            "\t 2. You have provided a non existent segment name.\n",
            bold=True,
        )
        return

    contacts_found = len(contacts)
    click.secho(
        f"A total of {contacts_found} contacts have been found ..", bold=True
    )

    contact_count = 0
    try:
        with click.progressbar(contacts, length=contacts_found) as contacts:
            for contact in contacts:
                phone = contact["phone"]
                payer_name = contact["payor"]
                email = contact["email"]

                phone_message_dict = {
                    "phone_number": phone,
                    "message": message_data["message"].format(
                        payer_name,
                        generate_shortened_dynamic_links(
                            generate_marketing_url(
                                email, message_data["tracking_url"]
                            )
                        ),
                    ),
                }

                to = [phone_message_dict["phone_number"]]
                payload = {
                    "to": to,
                    "message": phone_message_dict["message"],
                    "sender": SenderID.BeWell.value,
                    "segment": segment,
                }

                try:
                    send_sms_start_time = current_time()
                    send_sms(payload)
                except:
                    continue

                send_sms_end_time = current_time()
                send_sms_total_time = send_sms_end_time - send_sms_start_time
                sms_time_in_secs = send_sms_total_time.total_seconds()
                sms_rate = f"{sms_time_in_secs}s/message"

                click.secho(
                    f"\nMessage has been sent to phone number {to}. "
                    f"Message count {contact_count + 1} out of {contacts_found} "
                    f"with {sms_time_in_secs} seconds"
                )

                contact_count += 1

                if contact_count % 100 == 0:
                    t = (current_time() - start_time).total_seconds()
                    time_taken_so_far = convert_datetime_to_hours(t)
                    time_left = (
                        contacts_found - contact_count
                    ) * send_sms_total_time

                    time_left_in_hr = convert_datetime_to_hours(
                        time_left.total_seconds()
                    )
                    click.secho(
                        f"\n{contact_count} contacts marketed to, "
                        f"{time_taken_so_far} hours taken so far, "
                        f"{sms_rate}, {time_left_in_hr} hours left\n",
                        fg="blue",
                        blink=True,
                        bold=True,
                    )
    except:
        click.secho(f"\nExiting gracefully!", fg="red")

    if contact_count == contacts_found:
        click.secho(
            f"\n{contacts_found} contacts engaged successfully!", fg="green"
        )
    else:
        click.secho(
            f"\n{contact_count} contacts engaged successfully!", fg="green"
        )


@click.command()
@click.argument("segment")
@click.argument("wing")
def run_campaign(segment, wing):
    """Run the campaign."""
    if "WING A" in wing:
        send_marketing_bulk_sms(segment, wing, MESSAGE_A)
        return

    if "WING B" in wing:
        send_marketing_bulk_sms(segment, wing, MESSAGE_B)
        return

    else:
        click.secho(
            f"{wing} is an unidentified wing name",
            fg="red",
            blink=True,
            bold=True,
        )
        return


if __name__ == "__main__":
    run_campaign()
