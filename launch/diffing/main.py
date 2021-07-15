"""
Diffing script.

This diff filters out people who have already been marketed to
and users who have already downloaded our APP and outputs a CSV(for now)
of new users we can market to.
"""
import csv
import json
import logging

from google.cloud import bigquery

logging.basicConfig(level=logging.INFO)

client = bigquery.Client()


def sladers_data_fieldnames():
    """Helper to store slader data headers."""
    return [
        "phone_contact",
        "active_card_types",
        "payer_slade_code",
        "beneficiary_code",
        "payer_name",
        "first_name",
        "last_name",
    ]


def get_bewellers_data():
    """Return data from bewellers dataset."""
    query = """
        SELECT * FROM `bewell-app.bewellers.bewell_user_profiles_raw_latest`
    """
    query_job = client.query(query)

    phone_numbers = []

    for row in query_job:
        data = json.loads(row[5])
        phone = data["primaryPhone"]
        phone_numbers.append(phone)

    logging.info(f"length of bewellers found: {len(phone_numbers)}\n")
    return phone_numbers


def get_bewell_marketing_data():
    """Return data from marketing dataset."""
    # Message sent should be true
    query = """
        SELECT * FROM `bewell-app.marketing_data.marketing_data_bewell_prod_raw_latest`
        WHERE JSON_VALUE(data, '$.message_sent')="TRUE"
    """
    query_job = client.query(query)
    phone_numbers = []

    for row in query_job:
        data = json.loads(row[5])
        if "phone" in data:
            phone = data["phone"]
            phone_numbers.append(phone)
        elif "Phone" in data["properties"]:
            phoneNumber = data["properties"]["Phone"]
            phone_numbers.append(phoneNumber) 
        else:
            print("missing phone number")

    logging.info(
        f"length of marketing data with sent messages: {len(phone_numbers)}\n"
    )
    return phone_numbers


def normalize_phone_number(phone):
    """Normalize a phone number."""
    if phone.startswith("0"):
        return f"+254{phone[1:]}"

    elif phone.startswith("254"):
        return f"+{phone}"

    else:
        return phone


def get_sladers_data():
    """Return data from APA dataset."""
    query = """
        SELECT * FROM `bewell-app.marketing_diffing_data.apa_segmentation_jul_07_11`
    """
    query_job = client.query(query)

    apa_data = []

    for row in query_job:
        data_we_are_interested_in = {
            "phone_contact": normalize_phone_number(str(row[0])),
            "active_card_types": row[23],
            "payer_slade_code": row[12],
            "beneficiary_code": row[11],
            "payer_name": row[16],
            "first_name": row[1],
            "last_name": row[2],
        }
        apa_data.append(data_we_are_interested_in)

    logging.info(f"length of APA data members: {len(apa_data)}\n")
    return apa_data


def diff_to_remove_existing_and_marketed_to_sladers():
    """
    Perform a simple diff on staged data.

    This diff filters out people who have already been marketed to
    and users who have already downloaded our APP.
    """
    blacklisted = get_bewellers_data() + get_bewell_marketing_data()
    logging.info(
        f"length of blacklisted (Phone numbers who exist on Be.Well or have sent messages): {len(blacklisted)}\n"
    )
    sladers = get_sladers_data()
    for slader in sladers:
        if slader["phone_contact"] in blacklisted:
            sladers.remove(slader)
    logging.info(
        f"length of diffed sladers (new people we are to target): {len(sladers)}\n"
    )
    return sladers


def write_unique_sladers_to_csv():
    """Create a csv for sladers we have not marketed to (for now)."""
    sladers = diff_to_remove_existing_and_marketed_to_sladers()
    with open(f"apa_diffed_segmented_data_jul_07.csv", mode="w") as csv_file:
        writer = csv.DictWriter(csv_file, fieldnames=sladers_data_fieldnames())
        writer.writeheader()
        for slader in sladers:
            writer.writerow(slader)


if __name__ == "__main__":
    write_unique_sladers_to_csv()
