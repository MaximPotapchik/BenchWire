from lib.analysis.pipeline import Pipeline
from lib.analysis.parsers.yaml_parser import YamlConfigParser, FindRootDirectory

# Directory name since Pipeline already handles it
rootDir = FindRootDirectory("BenchWire")

# Yaml config grabbing
config = YamlConfigParser()

# Full analysis pipeline
Pipeline(config, rootDir)
