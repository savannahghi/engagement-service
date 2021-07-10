"""
Phone browser OS detection.

This is a simple HTTP python function that detects a devices OS
through their browser and redirect them to Playstore if they are on android
or AppStore if the are on iOS.
"""
import base64

import flask
import requests
from user_agents import parse

app = flask.Flask(__name__)

IOS = "iOS"
ANDROID = "Android"
BASE_URL = "https://engagement-prod-uyajqt434q-ew.a.run.app/"

ANDROID_DUMMY_TEMPLATE = """
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Be.Well By Slade360</title>
</head>
<body>
    <!-- Start of HubSpot Embed Code -->
    <script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/20198195.js"></script>
    <!-- End of HubSpot Embed Code -->

    <!-- The core Firebase JS SDK  -->
    <script src="https://www.gstatic.com/firebasejs/8.7.0/firebase-app.js"></script>
    <script src="https://www.gstatic.com/firebasejs/8.7.0/firebase-analytics.js"></script>
    <script>
        var firebaseConfig = {
            apiKey: "AIzaSyAv2aRsSSHkOR6xGwwaw6-UTkvED3RNlBQ",
            authDomain: "bewell-app.firebaseapp.com",
            databaseURL: "https://bewell-app.firebaseio.com",
            projectId: "bewell-app",
            storageBucket: "bewell-app.appspot.com",
            messagingSenderId: "841947754847",
            appId: "1:841947754847:web:034e338de70038796686ea",
            measurementId: "G-WR8JPLG8ZH"
        };

        firebase.initializeApp(firebaseConfig);
        var analytics = firebase.analytics();

        analytics.logEvent('redirected_to_android_playstore');    

        window.location.replace("https://play.google.com/store/apps/details?id=com.savannah.bewell");
    </script>

    <!-- Facebook Pixel Code -->
    <script>
        !function (f, b, e, v, n, t, s) {
            if (f.fbq) return; n = f.fbq = function () {
                n.callMethod ?
                n.callMethod.apply(n, arguments) : n.queue.push(arguments)
            };
            if (!f._fbq) f._fbq = n; n.push = n; n.loaded = !0; n.version = '2.0';
            n.queue = []; t = b.createElement(e); t.async = !0;
            t.src = v; s = b.getElementsByTagName(e)[0];
            s.parentNode.insertBefore(t, s)
        }(window, document, 'script',
        'https://connect.facebook.net/en_US/fbevents.js');
        fbq('init', '400335678066977');
        fbq('track', 'PageView');
    </script>
    <noscript>
        <img height="1" width="1" style="display:none"
        src="https://www.facebook.com/tr?id=400335678066977&ev=PageView&noscript=1" />
    </noscript>
    <!-- End Facebook Pixel Code -->

    <!-- Global site tag (gtag.js) - Google Ads: 1025904802 -->
    <script async src="https://www.googletagmanager.com/gtag/js?id=AW-1025904802"></script>
    <script>
        window.dataLayer = window.dataLayer || [];
        function gtag() { dataLayer.push(arguments); }
        gtag('js', new Date());

        gtag('config', 'AW-1025904802');
    </script>

    <!-- Event snippet for BeWell Well Campaign conversion page -->
    <script>
        gtag('event', 'conversion', { 'send_to': 'AW-1025904802/KKiKCOq1_dECEKKhmOkD' });
    </script>
</body>
</html>
"""

IOS_DUMMY_TEMPLATE = """
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Be.Well By Slade360</title>
</head>
<body>
    <!-- Start of HubSpot Embed Code -->
    <script type="text/javascript" id="hs-script-loader" async defer src="//js.hs-scripts.com/20198195.js"></script>
    <!-- End of HubSpot Embed Code -->

    <!-- The core Firebase JS SDK  -->
    <script src="https://www.gstatic.com/firebasejs/8.7.0/firebase-app.js"></script>
    <script src="https://www.gstatic.com/firebasejs/8.7.0/firebase-analytics.js"></script>
    <script>
        var firebaseConfig = {
            apiKey: "AIzaSyAv2aRsSSHkOR6xGwwaw6-UTkvED3RNlBQ",
            authDomain: "bewell-app.firebaseapp.com",
            databaseURL: "https://bewell-app.firebaseio.com",
            projectId: "bewell-app",
            storageBucket: "bewell-app.appspot.com",
            messagingSenderId: "841947754847",
            appId: "1:841947754847:web:034e338de70038796686ea",
            measurementId: "G-WR8JPLG8ZH"
        };

        firebase.initializeApp(firebaseConfig);
        var analytics = firebase.analytics();

        analytics.logEvent('redirected_to_iOS_appstore');

        window.location.replace("https://apps.apple.com/ke/app/be-well-by-slade360/id1496576692");
    </script>

    <!-- Facebook Pixel Code -->
    <script>
        !function (f, b, e, v, n, t, s) {
            if (f.fbq) return; n = f.fbq = function () {
                n.callMethod ?
                n.callMethod.apply(n, arguments) : n.queue.push(arguments)
            };
            if (!f._fbq) f._fbq = n; n.push = n; n.loaded = !0; n.version = '2.0';
            n.queue = []; t = b.createElement(e); t.async = !0;
            t.src = v; s = b.getElementsByTagName(e)[0];
            s.parentNode.insertBefore(t, s)
        }(window, document, 'script',
        'https://connect.facebook.net/en_US/fbevents.js');
        fbq('init', '400335678066977');
        fbq('track', 'PageView');
    </script>
    <noscript>
        <img height="1" width="1" style="display:none"
        src="https://www.facebook.com/tr?id=400335678066977&ev=PageView&noscript=1" />
    </noscript>
    <!-- End Facebook Pixel Code -->

    <!-- Global site tag (gtag.js) - Google Ads: 1025904802 -->
    <script async src="https://www.googletagmanager.com/gtag/js?id=AW-1025904802"></script>
    <script>
        window.dataLayer = window.dataLayer || [];
        function gtag() { dataLayer.push(arguments); }
        gtag('js', new Date());

        gtag('config', 'AW-1025904802');
    </script>

    <!-- Event snippet for BeWell Well Campaign conversion page -->
    <script>
        gtag('event', 'conversion', { 'send_to': 'AW-1025904802/KKiKCOq1_dECEKKhmOkD' });
    </script>
</body>
</html>
"""


def mark_bewell_aware(email):
    """Marks a user as bewell aware."""
    url = BASE_URL + "set_bewell_aware"
    decoded_email = base64.b64decode(email)
    response = requests.post(
        url=url, json={"email": decoded_email.decode("utf-8")}
    )
    result = response.json()
    return result


def detect_browser(request):
    """
    Detect a browser's user-agent.

    Given the family of OS we get, we redirect to either
    our Play store or App store.
    """
    user_agent = parse(flask.request.headers.get("User-Agent"))
    os_family = user_agent.os.family
    if os_family == IOS:
        email = request.args.get("email")
        mark_bewell_aware(email)
        return flask.render_template_string(IOS_DUMMY_TEMPLATE)

    if os_family == ANDROID:
        email = request.args.get("email")
        mark_bewell_aware(email)
        return flask.render_template_string(ANDROID_DUMMY_TEMPLATE)

    else:
        return "Run this on an Android or iOS phone and see the magic happen :). Be.Well By Slade360"


@app.route("/")
def index():
    """Flask app entrypoint."""
    return detect_browser(flask.request)
