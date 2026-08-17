import os
import yaml
from .root_finder import FindRootDirectory

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

def YamlRunParser(specMatrix, label, prefix, run):
    rootDir = FindRootDirectory("BenchWire") 
    yamlDir = os.path.join(rootDir, "results", "yaml")
    
    path = os.path.join(yamlDir, f"{specMatrix}_{label}_{prefix}run{run}.yaml")
    with open(path) as f:
        stats = yaml.safe_load(f)
    
    return stats

def YamlConfigParser(configPath=None):
    rootDir = FindRootDirectory("BenchWire")

    if configPath is None:
        configPath = os.path.join(rootDir, "config.yaml")

    with open(configPath) as f:
        config = yaml.safe_load(f)

    return config
