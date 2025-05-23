#!/bin/python3

import os
import requests
import sys
import zipfile

RPC_URL = "http://127.0.0.1:26657"

def get_current_block_height():
    """
    Fetch the current block height from the blockchain API.
    """
    url = "%s/block" % RPC_URL
    response = requests.get(url)
    if response.status_code == 200:
        return int(response.json()['result']['block']['header']['height'])
    else:
        raise Exception("Failed to fetch block height")
    
def align_to_bucket(bh, bucket_size=1000):
    """
    Align the block height to the nearest bucket size.
    """
    return (bh // bucket_size) * bucket_size

def zip_wasm_files(wasm_dir, output_zip):
    """
    Create a zip file containing all wasm files in the specified directory.
    """
    with zipfile.ZipFile(output_zip, 'w') as zipf:
        for root, _, files in os.walk(wasm_dir):
            for file in files:
                if file.endswith('.wasm'):
                    file_path = os.path.join(root, file)
                    zipf.write(file_path, os.path.relpath(file_path, wasm_dir))

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 wasm_packer.py <path_to_paxi>")
        sys.exit(1)

    paxi_path = sys.argv[1]
    wasm_snapshots_path = os.path.join(paxi_path, "wasm_snapshots")
    if not os.path.exists(wasm_snapshots_path):
        os.makedirs(wasm_snapshots_path)

    # Get the current block height
    block_height = get_current_block_height()
    aligned_height = align_to_bucket(block_height)

    # Create a zip file of the wasm files
    output_zip = os.path.join(wasm_snapshots_path, f"wasm_files_{aligned_height}.zip")
    if os.path.exists(output_zip):
        print(f"Zip file {output_zip} already exists. Skipping creation.")
        sys.exit(1)

    wasm_path = os.path.join(paxi_path, "wasm")
    zip_wasm_files(wasm_path, output_zip)
    print(f"Created zip file: {output_zip}")