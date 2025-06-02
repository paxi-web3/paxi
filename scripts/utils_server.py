# utils_server.py

"""FastAPI server for utility functions like snapshot download."""

from fastapi import FastAPI
import os

app = FastAPI()

@app.get("/latest_wasm_snapshot")
def get_snapshot():
    wasm_snapshots_path = os.path.join(os.path.expanduser("~"), "wasm_snapshots")
    
    # Get the latest snapshot file
    file_list = []
    for file in os.listdir(wasm_snapshots_path):
        if file.endswith('.zip'):
            file_height = int(file.split('.')[0])
            file_list.append((file, file_height))
    file_list.sort(key=lambda x: x[1], reverse=True)

    return {"url": f"http://testnet-snapshot.paxinet.io/wasm_snapshots/{file_list[0][0]}"}