"""Segment CSV data into the different segments."""
import csv

import click
import numpy as np

WINGS = 2
CHUNKS = 5


def custom_hubspot_properties():
    """
    Define a single source of HubSpot custom properties.

    These are the custom fields we want in a new CSV that will
    be inferred from the CSV and imported to the CRM."""
    return [
        "be_well_enrolled",
        "opt_out",
        "be_well_aware",
        "be_well_persona",
        "has_wellness_card",
        "has_cover",
        "payor",
        "first_channel_of_contact",
        "initial_segment",
        "has_virtual_card",
        "email",
        "phone_number",
        "firstname",
        "lastname",
        "wing",
        "message_sent",
        "payer_slade_code",
        "member_number",
    ]


def randomize_data(data):
    """Helper to shuffle data randomly."""
    np.random.shuffle(data)
    return data


def normalize_phone_number(phone):
    """Normalize a phone number."""
    if phone.startswith("0"):
        return f"+254{phone[1:]}"

    elif phone.startswith("254"):
        return f"+{phone}"

    else:
        return phone


def create_custom_properties_from_slade_data(csv_file, segment_name):
    """Infer Hubspot data from the slade CSV data."""
    list_of_properties_we_want = []
    with open(f"{csv_file}.csv") as csv_file:
        csv_reader = csv.DictReader(csv_file)
        for row in csv_reader:
            has_virtual_card = (
                "YES" if "VIRTUAL" in row["active_card_types"] else "NO"
            )
            try:
                payer_slade_code = row["payer_slade_code"]
                member_number = row["beneficiary_code"]
            except KeyError:
                payer_slade_code = ""
                member_number = ""

            phone_number = normalize_phone_number(row["phone_contact"])
            email = f"{phone_number}@users.bewell.co.ke"
            custom_properties = {
                "email": email,
                "be_well_enrolled": "NO",
                "opt_out": "NO",
                "be_well_aware": "NO",
                "be_well_persona": "SLADER",
                "has_wellness_card": "YES",
                "has_cover": "YES",
                "payor": row["payer_name"],
                "first_channel_of_contact": "SMS",
                "initial_segment": segment_name,
                "has_virtual_card": has_virtual_card,
                "phone_number": phone_number,
                "firstname": row["first_name"],
                "lastname": row["last_name"],
                "payer_slade_code": payer_slade_code,
                "member_number": member_number,
            }
            list_of_properties_we_want.append(custom_properties)
    return list_of_properties_we_want


def split_segment_data_into_wings(csv_file, segment_name):
    """
    Split the data we have into wings.

    The output of this will be a tuple of equally split data
    based on the number of wings defined
    """
    random_data = randomize_data(
        create_custom_properties_from_slade_data(csv_file, segment_name),
    )
    return np.array_split(
        random_data,
        WINGS,
    )


def chunk_winged_data(csv_file, segment_name):
    """For each wing split the data further into chunks."""
    winged_data = split_segment_data_into_wings(csv_file, segment_name)
    chunked_data = []
    for data in winged_data:
        chunked_data.append(np.array_split(data, CHUNKS))
    return chunked_data


def write_wing_data_to_csv(segment_name, csv_file):
    """
    Write the extracted data to a CSV.

    This helper write the first wing data (data in wing one)
    to a CSV file
    """
    chunked_data = chunk_winged_data(csv_file, segment_name)
    for i in range(0, WINGS):
        with open(f"{csv_file}_wing_{i}.csv", mode="w") as wing_A_csv:
            fieldnames = custom_hubspot_properties()
            writer = csv.DictWriter(wing_A_csv, fieldnames=fieldnames)
            writer.writeheader()

            for index, chunk_list in enumerate(chunked_data[i]):
                for data in chunk_list:
                    data["wing"] = f"WING {chr(ord('@')+(i + 1))} 0{index + 1}"
                    data["message_sent"] = "FALSE"
                    writer.writerow(data)


@click.command()
@click.argument("csv_file")
@click.argument("segment_name")
def generate_children_csv(csv_file, segment_name):
    """Entry point to our script."""
    write_wing_data_to_csv(segment_name, csv_file)


def main():
    generate_children_csv()


if __name__ == "__main__":
    main()
