import sys
from lib.analysis.pipeline import Pipeline
from lib.analysis.parsers.yaml_parser import YamlConfigParser, FindRootDirectory

# Directory name since Pipeline already handles it.
rootDir = FindRootDirectory("BenchWire")

# Check if any config was passed.
configPath = sys.argv[1] if len(sys.argv) > 1 else None

# Yaml config grabbing.
config = YamlConfigParser(configPath)

# Full analysis pipeline.
Pipeline(config, rootDir)
