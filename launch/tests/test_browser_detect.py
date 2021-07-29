"""Browser detect main test file."""
import base64

import pytest
import requests
from launch.browser_detect.main import (
    APPLE_STORE_LINK,
    PLAY_STORE_LINK,
    app,
    events,
    htmlTemplate,
    mark_bewell_aware,
)

BASE_URL = "https://europe-west3-bewell-app.cloudfunctions.net/"


@pytest.fixture
def test_client():
    """Return a flask test client."""
    return app.test_client()


def _encode_email():
    """Base64 encodes a test email."""
    email = "something@users.bewell.co.ke"
    bs = bytes(email, "utf-8")
    return base64.b64encode(bs)


def _user_agent(os):
    """Return user agents header strings."""
    if os == "ANDROID":
        return (
            "Mozilla/5.0 (Linux; U; Android 4.0.4; en-gb; "
            "GT-I9300 Build/IMM76D) "
            "AppleWebKit/534.30 (KHTML, like Gecko) "
            "Version/4.0 Mobile Safari/534.30"
        )
    if os == "IOS":
        return (
            "Mozilla/5.0(iPad; U; CPU iPhone OS 3_2 like Mac OS X; en-us) "
            "AppleWebKit/531.21.10 (KHTML, like Gecko) "
            "Version/4.0.4 Mobile/7B314 Safari/531.21.10"
        )
    return (
        "Mozilla/5.0 (compatible; MSIE 10.0; "
        "Windows NT 6.2; Trident/6.0; Touch)"
    )


# todo restore this test after coomontools ref
# def test_mark_as_bewell_aware():
#     """
#     Test mark as bewell aware.

#     Scenario: Mark a user as Be.Well Aware
#         Given the request sent has a base64 encoded email param
#         Then the user with the email is marked as Be.Well Aware
#     """
#     resp = mark_bewell_aware(_encode_email())

#     assert resp.status_code == 200
#     assert resp.json() == {"status": "success"}


def test_detect_android_browser(test_client):
    """Test Android browser detection.

    Scenario: Detect a browser to get the device's os family
        Given the user-agent headers string is from an Android device
        When the request contains the email param, mark them as bewell aware
        Then render a page that redirects the user to playstore
    """
    headers = {"User-Agent": _user_agent("ANDROID")}
    resp = test_client.get(
        "/", query_string={"email": _encode_email()}, headers=headers
    )

    resp.status_code == 200
    resp.data.decode() in htmlTemplate(events["android"], PLAY_STORE_LINK)


def test_detect_ios_browser(test_client):
    """Test iOS browser detection.

    Scenario: Detect a browser to get the device's os family
        Given the user-agent headers string is from an iOS device
        When the request contains the email param, mark them as bewell aware
        Then render a page that redirects the user to apple store
    """
    headers = {"User-Agent": _user_agent("iOS")}
    resp = test_client.get(
        f"{BASE_URL}/detect_browser",
        query_string={"email": _encode_email()},
        headers=headers,
    )

    assert resp.status_code == 308
    assert (
        "pe-west3-bewell-app.cloudfunctions.net/?email=" in resp.data.decode()
    )


def test_detect_any_other_browser_with_email(test_client):
    """Test other browser detection.

    Scenario: Detect a browser to get the device's os family
        Given the user-agent headers string is not from an iOS or Android device
        When the request contains the email param, mark them as bewell aware
        Then redirect the user to our landing page A
    """
    headers = {"User-Agent": _user_agent("")}
    resp = test_client.get(
        f"{BASE_URL}/detect_browser",
        query_string={"email": _encode_email()},
        headers=headers,
    )

    assert resp.status_code == 308


def test_detect_android_browser_without_email(test_client):
    """Test Android browser detection without an email.

    Scenario: Detect a browser to get the device's os family
        Given the user-agent headers string is from Android device
        When the request does not contain the email param
        Then render a page that redirects the user to playstore
    """
    headers = {"User-Agent": _user_agent("ANDROID")}
    resp = test_client.get(f"{BASE_URL}/detect_browser", headers=headers)

    assert resp.status_code == 308
    resp.data.decode() in htmlTemplate(events["android"], PLAY_STORE_LINK)


def test_detect_ios_browser_without_email(test_client):
    """Test iOS browser detection without an email.

    Scenario: Detect a browser to get the device's os family
        Given the user-agent headers string is from an iOS device
        When the request does not contain the email param
        Then render a page that redirects the user to apple store
    """
    headers = {"User-Agent": _user_agent("iOS")}
    resp = test_client.get(f"{BASE_URL}/detect_browser", headers=headers)

    assert resp.status_code == 308
    resp.data.decode() in htmlTemplate(events["IOS"], APPLE_STORE_LINK)


def test_detect_other_browser_without_email(test_client):
    """Test other browser detection without an email.

    Scenario: Detect a browser to get the device's os family
        Given the user-agent headers string is is not from an iOS or Android device
        When the request does not contain the email param
        Then redirect the user to our landing page A
    """
    headers = {"User-Agent": _user_agent("iOS")}
    resp = test_client.get(f"{BASE_URL}/detect_browser", headers=headers)

    assert resp.status_code == 308
    resp.data.decode() in htmlTemplate(events["IOS"], APPLE_STORE_LINK)


def test_detect_browser_cloud_func_on_android_with_email():
    """Test android browser cloud function calls.

    Scenario: Detect a browser to get the device's os family
        Given the user-agent headers string is from an Android device
        When the request contains the email param, mark them as bewell aware
        Then render a page that redirects the user to playstore
    """
    headers = {"User-Agent": _user_agent("ANDROID")}
    args = {"email": _encode_email()}
    resp = requests.get(
        f"{BASE_URL}/detect_browser", headers=headers, params=args
    )

    assert resp.status_code == 200
    assert "redirected_to_android_playstore" in resp.text


def test_detect_browser_cloud_func_on_ios_with_email():
    """
    Test ios browser cloud function calls.

    Scenario: Detect a browser to get the device's os family
        Given the user-agent headers string is from an iOS device
        When the request contains the email param, mark them as bewell aware
        Then render a page that redirects the user to apple store
    """
    headers = {"User-Agent": _user_agent("iOS")}
    args = {"email": _encode_email()}
    resp = requests.get(
        f"{BASE_URL}/detect_browser", headers=headers, params=args
    )

    assert resp.status_code == 200
    assert "Be.Well by Slade360° - Simple. Caring. Trusted" in resp.text


def test_detect_browser_cloud_func_on_desktop():
    """
    Test other os families browser cloud function calls.

    Scenario: Detect a browser to get the device's os family
        Given the user-agent headers string is not from an iOS or Android device
        When the request contains the email param, mark them as bewell aware
        Then redirect the user to our landing page A
    """
    headers = {"User-Agent": _user_agent("")}
    args = {"email": _encode_email()}
    resp = requests.get(
        f"{BASE_URL}/detect_browser", headers=headers, params=args
    )

    assert resp.status_code == 200


def test_detect_browser_cloud_func_without_email():
    """
    Test cloud function call with email in the request params.

    Scenario: Detect a browser to get the device's os family
        Given the user-agent headers string is from Android device
        When the request does not contain the email param
        Then render a page that redirects the user to playstore
    """
    headers = {"User-Agent": _user_agent("ANDROID")}
    resp = requests.get(f"{BASE_URL}/detect_browser", headers=headers)

    assert resp.status_code == 200
