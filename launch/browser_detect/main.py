"""
Phone browser OS detection.

This is a simple HTTP python function that detects a devices OS
through their browser and redirect them to Playstore if they are on android
or AppStore if the are on iOS.
"""
import flask
from user_agents import parse

app = flask.Flask(__name__)

IOS = "iOS"
ANDROID = "Android"

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

        if (RegExp('[?&]' + 'email' + '=([^&]*)').exec(window.location.search)) {
            identifyVisitor();
            _hsq.push(['trackPageView']);

            markAsBeWellAware({
                'email': atob(decodeURIComponent(getParameterByName("email"))),
            });
        }


        function getParameterByName(name) {
        var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);
        return match && decodeURIComponent(match[1].replace(/\+/g, ' '));
        }

        function identifyVisitor() {
        var _hsq = window._hsq = window._hsq || [];
        _hsq.push(["identify", {
            email: atob(decodeURIComponent(getParameterByName("email")))
        }]);
        }

        async function markAsBeWellAware(data) {
            let url = 'https://engagement-prod.healthcloud.co.ke/set_bewell_aware';
            const response = await fetch(url, {
                method: 'POST',
                mode: 'cors',
                cache: 'no-cache',
                headers: {
                'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });
            return response.json();
        }

        window.location.replace("https://play.google.com/store/apps/details?id=com.savannah.bewell");
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


        if (RegExp('[?&]' + 'email' + '=([^&]*)').exec(window.location.search)) {
            identifyVisitor();
            _hsq.push(['trackPageView']);

            markAsBeWellAware({
            'email': atob(decodeURIComponent(getParameterByName("email"))),
            });
        }


        function getParameterByName(name) {
        var match = RegExp('[?&]' + name + '=([^&]*)').exec(window.location.search);
        return match && decodeURIComponent(match[1].replace(/\+/g, ' '));
        }

        function identifyVisitor() {
            var _hsq = window._hsq = window._hsq || [];
            _hsq.push(["identify", {
            email: atob(decodeURIComponent(getParameterByName("email")))
            }]);
        }

        async function markAsBeWellAware(data) {
            let url = 'https://engagement-prod.healthcloud.co.ke/set_bewell_aware';
            const response = await fetch(url, {
            method: 'POST',
            mode: 'cors',
            cache: 'no-cache',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
            });
            return response.json();
        }

        window.location.replace("https://apps.apple.com/ke/app/be-well-by-slade360/id1496576692");
    </script>
</body>
</html>
"""


def detect_browser(request):
    """
    Detect a browser's user-agent.

    Given the family of OS we get, we redirect to either
    our Play store or App store.
    """
    user_agent = parse(flask.request.headers.get("User-Agent"))
    os_family = user_agent.os.family
    if os_family == IOS:
        return flask.render_template_string(IOS_DUMMY_TEMPLATE)

    if os_family == ANDROID:
        return flask.render_template_string(ANDROID_DUMMY_TEMPLATE)

    else:
        return "Run this on an Android or iOS phone and see the magic happen :). Be.Well By Slade360"


@app.route("/")
def index():
    """Flask app entrypoint."""
    return detect_browser(flask.request)
