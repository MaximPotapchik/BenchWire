import os
import yaml

# TODO: Move this to its own file. root_finder.py or something.
# grabs yaml location 
def FindRootDirectory(rootDirName):
    currentDir = os.path.dirname(os.path.abspath(__file__))

    while True:
        if os.path.basename(currentDir) == rootDirName:
            return currentDir

        parentDir = os.path.dirname(currentDir)
        if parentDir == currentDir:
            raise FileNotFoundError("BenchWire directory not found.")

        currentDir = parentDir

# This is an example stat extractor. Not used anymore. 
"""
def ExtractStats(targetStats, data):
    extractedStats = {}

    def merge(new):
        for k, v in new.items():
            if k in extractedStats:
                raise ValueError(f"'{k}' found at multiple nesting levels, ambiguous extraction")
            extractedStats[k] = v

    if isinstance(data, dict):
        for key, val in data.items():
            if key in targetStats:
                merge({key: val})
            if isinstance(val, (dict, list)):
                merge(ExtractStats(targetStats, val))

    if isinstance(data, list):
        for item in data:
            if isinstance(item, (dict, list)):
                merge(ExtractStats(targetStats, item))

    return extractedStats
"""

def YamlRunParser(prefix, run):
    rootDir = FindRootDirectory("BenchWire") 
    yamlDir = os.path.join(rootDir, "results", "yaml")
    
    path = os.path.join(yamlDir, f"{prefix}run_{run}.yaml")
    with open(path) as f:
        stats = yaml.safe_load(f)
    
    return stats

def YamlConfigParser():
    rootDir = FindRootDirectory("BenchWire")
    configLoc = os.path.join(rootDir, "config.yaml")

    with open(configLoc) as f:
        config = yaml.safe_load(f)

    return config

