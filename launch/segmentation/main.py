"""Segment CSV data into the different segments."""
import csv

import click
import numpy as np

WINGS = 2


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
    ]


def randomize_data(data):
    """Helper to shuffle data randomly."""
    np.random.shuffle(data)
    return data


def normalize_phone_number(phone):
    """Normalize a phne number."""
    if phone.startswith("0"):
        return f"+254{phone[1:]}"

    elif phone.startswith("254"):
        return f"+{phone}"

    else:
        return phone


def create_custom_properties_from_slade_data(path_to_csv, segment_name):
    """Infer Hubspot data from the slade CSV data."""
    list_of_properties_we_want = []
    with open(path_to_csv) as csv_file:
        csv_reader = csv.DictReader(csv_file)
        for row in csv_reader:
            has_virtual_card = (
                "YES" if "VIRTUAL" in row["active_card_types"] else "NO"
            )
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
            }
            list_of_properties_we_want.append(custom_properties)
    return list_of_properties_we_want


def split_segment_data_into_wings(path_to_csv, segment_name):
    """
    Split the data we have into wings.

    The output of this will be a tuple of equally split data
    based on the number of wings defined
    """
    random_data = randomize_data(
        create_custom_properties_from_slade_data(path_to_csv, segment_name),
    )
    wing_1_data, wing_2_data = np.array_split(
        random_data,
        WINGS,
    )
    return wing_1_data, wing_2_data


def write_wing_data_to_csv(segment_name, path_to_csv):
    """
    Write the extracted data to a CSV.

    This helper write the first wing data (data in wing one)
    to a CSV file
    """
    wing_1_data, wing_2_data = split_segment_data_into_wings(
        path_to_csv, segment_name
    )
    with open(f"{segment_name}_wing_A.csv", mode="w") as wing_A_csv:
        fieldnames = custom_hubspot_properties()
        writer = csv.DictWriter(wing_A_csv, fieldnames=fieldnames)
        writer.writeheader()
        for dataset in wing_1_data:
            writer.writerow(dataset)

    with open(f"{segment_name}_wing_B.csv", mode="w") as wing_B_csv:
        fieldnames = custom_hubspot_properties()
        writer = csv.DictWriter(wing_B_csv, fieldnames=fieldnames)
        writer.writeheader()
        for dataset in wing_2_data:
            writer.writerow(dataset)


@click.command()
@click.argument("path_to_csv")
@click.argument("segment_name")
def generate_children_csv(path_to_csv, segment_name):
    """Entry point to our script."""
    write_wing_data_to_csv(segment_name, path_to_csv)


def main():
    generate_children_csv()


if __name__ == "__main__":
    main()
