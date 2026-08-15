import os
import orjson
from .root_finder import FindRootDirectory

def JsonRunParser(specMatrix, label, prefix, run):
    rootDir = FindRootDirectory("BenchWire")
    jsonDir = os.path.join(rootDir, "results", "json")

    path = os.path.join(jsonDir, f"{specMatrix}_{label}_{prefix}run{run}.json")
    with open(path, "rb") as f:
        return orjson.loads(f.read())
