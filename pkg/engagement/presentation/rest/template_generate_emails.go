package rest

// ColLectCMREmailTemplate generates an email template
const ColLectCMREmailTemplate = `

<!DOCTYPE html>
<html>
  <head>
    <title></title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <style type="text/css">
      @media screen {
        @font-face {
          font-family: "Lato";
          font-style: normal;
          font-weight: 400;
          src: local("Lato Regular"), local("Lato-Regular"),
            url(https://fonts.gstatic.com/s/lato/v11/qIIYRU-oROkIk8vfvxw6QvesZW2xOQ-xsNqO47m55DA.woff)
              format("woff");
        }

        @font-face {
          font-family: "Lato";
          font-style: normal;
          font-weight: 700;
          src: local("Lato Bold"), local("Lato-Bold"),
            url(https://fonts.gstatic.com/s/lato/v11/qdgUG4U09HnJwhYI-uK18wLUuEpTyoUstqEm5AMlJo4.woff)
              format("woff");
        }

        @font-face {
          font-family: "Lato";
          font-style: italic;
          font-weight: 700;
          src: local("Lato Bold Italic"), local("Lato-BoldItalic"),
            url(https://fonts.gstatic.com/s/lato/v11/HkF_qI1x_noxlxhrhMQYELO3LdcAZYWl9Si6vvxL-qU.woff)
              format("woff");
        }
      }

      body,
      table,
      td {
        -webkit-text-size-adjust: 100%;
        -ms-text-size-adjust: 100%;
      }

      img {
        -ms-interpolation-mode: bicubic;
      }

      img {
        border: 0;
        height: auto;
        line-height: 100%;
        outline: none;
        text-decoration: none;
      }

      table {
        border-collapse: collapse !important;
      }

      body {
        height: 100% !important;
        margin: 0 !important;
        padding: 0 !important;
        width: 100% !important;
      }

      /* MOBILE STYLES */
      @media screen and (max-width: 600px) {
        h1 {
          font-size: 32px !important;
          line-height: 32px !important;
        }
      }

      /* ANDROID CENTER FIX */
      div[style*="margin: 16px 0;"] {
        margin: 0 !important;
      }
    </style>
  </head>

  <body
    style="
      background-color: #f4f4f4;
      margin: 0 !important;
      padding: 0 !important;
    "
  >
    <table
      border="0"
      cellpadding="0"
      cellspacing="0"
      width="100%"
      height="100%"
    >
      <tr>
        <td bgcolor="#7B54C4" align="center">
          <table
            border="0"
            cellpadding="0"
            cellspacing="0"
            width="100%"
            style="max-width: 600px"
          >
            <tr>
              <td
                align="center"
                valign="top"
                style="padding: 40px 10px 40px 10px"
              ></td>
            </tr>
          </table>
        </td>
      </tr>
      <tr>
        <td bgcolor="#7B54C4" align="center" style="padding: 0px 10px 0px 10px">
          <table
            border="0"
            cellpadding="0"
            cellspacing="0"
            width="85%"
            style="max-width: 600px"
          >
            <tr>
              <td
                bgcolor="#ffffff"
                valign="top"
                style="
                  padding: 40px 20px 10px 24px;
                  border-radius: 4px 4px 0px 0px;
                  color: #111111;
                  font-family: 'Lato', Helvetica, Arial, sans-serif;
                  font-size: 48px;
                  font-weight: 400;
                  line-height: 48px;
                "
              >
                <img
                  src="https://lh3.googleusercontent.com/pw/ACtC-3fN_p8U8EZgmtQymnwrhr_-5Go6Kw5e5U7lkjyk1jjMIEwSs6rDNELplpgVk2IciMfw5AbnphxJYwdocnsE6Y88xyKGlNXm1E1x3Sm9uxeMHhwjf8YgNwo622G8cb-d7ntlbNl7-uPCEylu5O_KzZY=s638-no"
                  width="125"
                  height="120"
                  style="display: block; border: 0px; margin-bottom: 0"
                />
              </td>
            </tr>
          </table>
        </td>
      </tr>
      <tr>
        <td bgcolor="#f4f4f4" align="center" style="padding: 0px 10px 0px 10px">
          <table
            border="0"
            cellpadding="0"
            cellspacing="0"
            width="85%"
            style="max-width: 600px"
          >
            <tr>
              <td
                bgcolor="#ffffff"
                align="left"
                style="
                  padding: 20px 30px 20px 30px;
                  color: #666666;
                  font-family: 'Lato', Helvetica, Arial, sans-serif;
                  font-size: 18px;
                  font-weight: 400;
                  line-height: 25px;
                "
              >
                <p style="margin: 0;font-size: 25px;font-weight: bold; color: black;">Join Be.Well today</p>
                <p></p>
                
                <p style="margin: 0">
                  <span style="color: black;">Hi Vincent Michuki</span>, click the link below to download the Be.Well App on your Android or iOS device
               </p>
               <br>
               <p style="margin-bottom: 10px;font-size: 16px;font-weight: bold; color: black;">The Be.Well app lets you:</p>
               <ul style="margin: 0;">
                 <li>Add your cover</li>
                 <li>View your benefits</li>
                 <li>See how you have used your benefits, etc</li>
               </ul>
               <p style="display: flex;">
                 <a href="https://play.google.com/store/apps/details?id=com.savannah.bewell" style="margin-right: 20px;">
                  <img
                  src="https://lh3.googleusercontent.com/AX2HV2ghizQTq_7Y4NnTpSy8It07rQwitQMp3J3U8JNIbe5XFnJ6wppwDTX_An7yN1FZ-XFnqhFOQol6enNd9LcyFnWNedKZOw4_vnvVKOBDXfiPX55DLHgqezcKLxBWdAWLIi5utMB1NvD5FXiseUcGkOAf8QvaeM7OeNZt2qdBE315oiP4Q2yFA9Cj9qyHDUb9PNGNfiC5jyEXfFw6OGEJU2ZMk5FsnL44Pv2bwF7k7W2DT6E1x_frR3OML1RrzJDLZlUC_WAn5G1P0Kbz49T_nhw_9fvLdW2lHjk2EB0R3GcsgAbVnhTa1HLhGq6A4UlouQI6aGXxDxDW7-ET9MYbQgoRTct94dQDuN3iy1X2sn1AnGQ11zIYPpAwGaEu7bZg6_pUtPE-RN074LWRMczCVi7zhul4bSEexZOAsc7zrYjk37qoKMyJs8QcG_odYheAvGkJhhOp_0KG8C2sXjVfyjtEolDVAyrsLfUAYZwq3SucWwzFpg2rq9fBdqs99rAvbsgorRutUPFHAcB5JX8g1E6wrz4w8s16EJZ0SsJBtTtgIuRxhxJyOmXWDdV-tOJSp1EIJVqKBHQYFoTZ6EEpG4Ft-LcaHOJul_7y15Zh7ORXEAqPoRVGlHvLLV0knDPG6olF-vQrN2VzSLIA889hx1M9pzcA3JOhhVEXC16eiZwOPxiW9F5OS9tjy1FA53XcxkMSmzfyhPAO-P3A_BI=w440-h130-no?authuser=1"
                  width="150"
                  style="display: block; border: 0px; margin-bottom: 0"
                />
                </a>

                <a href="https://play.google.com/store/apps/details?id=com.savannah.bewell">
                  <img
                  src="https://lh3.googleusercontent.com/nS3ZVoQs6NJ3eFA4Gx-wJJ5gemF10OFYK7pZBbcHqWizXFLmh4NruQHywPD4_z7n0SQychIO_gOQMi80_qpu2tpsNqu7wXvdSRZhDwYXkvM3qd0vV0sK4H90Vtq13Pyhh4WAnSRV7AisfHEu6xgAp_gDe8ScRQdt3UvdD34vK7q6GjEQz_E3w_Hq1RcNssmkqus0emXQAd-PMv2CCqKdoe0jeVdwUk_261StQCq-SndTNI2R0J3vQXmujzEklgWqSV20gLr1ihu2y1zqblXK8VUVtS3Cvwk8EjuCU80QjNp52NYLrAGAVqEP0PibVmt7lDVQ0ri0UD-CxFo-q4c9zP4O18lbz3PVG4KA3bvqss9WcO6w3xhR4QkwmkULfaPAgi-qrgSXEROWQ5o4PcTnl6xWG-XWcxVUnCYbwjWhqwH7I4xzFBYReZdU1YZmdKiuGb4Juq-Bs2qQkiWnBD8bMeS1pX9XxX2xI3towtIKc-_zsXO806ncfJjIKFrD5cVEUpBDRdbuc7S7NsVWtqpXxItKd-ye9dbXkaI-aSkN-59mHvdN-YMhixU0ovD3GG55rhpEgAilrdSeIPDkJChCdmxB3ORwOWWvsZlznW-Zu29rAm5Fy7y0O1N98f8Q3fl7oMASUqPRCl70fYJKcKVHHqJUNJJR9L7KBn4iVf75M3OtUJ5qmBaizDjtR6dck4D7ZdsAHsSvGyYbQQafIodRpco=w440-h129-no?authuser=1"
                  width="150"
                  style="display: block; border: 0px; margin-bottom: 0"
                />
                </a>
               </p>
               <br>
               <p style="margin: 0">
               Thank you for using Be.Well.
               </p>
              </td>
            </tr>

            <tr>
              <td
                bgcolor="#ffffff"
                align="center"
                style="
                  color: #666666;
                  font-family: 'Lato', Helvetica, Arial, sans-serif;
                  font-size: 40px;
                  font-weight: 900;
                  line-height: 40px;
                "
              ></td>
            </tr>

            <tr>
              <td
                bgcolor="#ffffff"
                align="left"
                style="
                  padding: 0px 30px 40px 30px;
                  border-radius: 0px 0px 4px 4px;
                  color: #666666;
                  font-family: 'Lato', Helvetica, Arial, sans-serif;
                  font-size: 18px;
                  font-weight: 400;
                  line-height: 25px;
                "
              >
                <p style="margin: 0">
                  Regards,<br />
                  The Be.Well Team
                </p>
              </td>
            </tr>
          </table>
        </td>
      </tr>
      <table
        border="0"
        cellpadding="0"
        cellspacing="0"
        width="100%"
        style="padding-top: 40px"
      >
        <tr>
          <td
            bgcolor="#f4f4f4"
            align="center"
            style="
              padding: 40px 30px 40px 30px;
              border-radius: 0px 0px 4px 4px;
              color: #666666;
              font-family: 'Lato', Helvetica, Arial, sans-serif;
              font-size: 18px;
              font-weight: 400;
              line-height: 25px;
            "
          >
            <p style="margin: 0">
              For more information or queries, contact us at
              <a href="mailto:feedback@bewell.co.ke">feedback@bewell.co.ke</a>
            </p>
          </td>
        </tr>
      </table>
    </table>
    <script
      src="https://cdn.jsdelivr.net/npm/publicalbum@latest/embed-ui.min.js"
      async
    ></script>
  </body>
</html>

`
