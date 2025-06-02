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

def zip_wasm_files(wasm_dir, output_tmp, output_zip):
    """
    Create a zip file containing all wasm files in the specified directory.
    """
    with zipfile.ZipFile(output_tmp, 'w') as zipf:
        for root, _, files in os.walk(wasm_dir):
            for file in files:
                if file.endswith('.wasm'):
                    file_path = os.path.join(root, file)
                    zipf.write(file_path, os.path.relpath(file_path, wasm_dir))
    os.rename(output_tmp, output_tmp.replace('.tmp', '.zip'))

def clear_old_snapshots(wasm_snapshots_path):
    file_list = []
    for file in os.listdir(wasm_snapshots_path):
        if file.endswith('.zip'):
            file_height = int(file.split('.')[0])
            file_list.append((file, file_height))
    file_list.sort(key=lambda x: x[1], reverse=True)

    keep_top = 5
    for file, _ in file_list[keep_top:]:
        os.remove(os.path.join(wasm_snapshots_path, file))

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 wasm_packer.py <path_to_paxi>")
        sys.exit(1)

    paxi_path = sys.argv[1]
    wasm_snapshots_path = os.path.join(os.path.expanduser("~"), "wasm_snapshots")
    if not os.path.exists(wasm_snapshots_path):
        os.makedirs(wasm_snapshots_path)

    # Get the current block height
    block_height = get_current_block_height()

    # Create a zip file of the wasm files
    output_tmp = os.path.join(wasm_snapshots_path, f"{block_height}.tmp")
    output_zip = os.path.join(wasm_snapshots_path, f"{block_height}.zip")
    if os.path.exists(output_tmp):
        print(f"Zip file {output_tmp} already exists. Skipping creation.")
        sys.exit(1)

    wasm_path = os.path.join(paxi_path, "wasm/wasm/state/wasm/")
    has_wasm = any(f.endswith('.wasm') for f in os.listdir(wasm_path) if os.path.isfile(os.path.join(wasm_path, f)))
    if not has_wasm:
        print(f"No wasm files found in {wasm_path}. Please ensure the path is correct.")
        sys.exit(1)
    
    zip_wasm_files(wasm_path, output_tmp, output_zip)
    print(f"Created zip file: {output_zip}")
    
    # Clear old snapshots
    clear_old_snapshots(wasm_snapshots_path)
    