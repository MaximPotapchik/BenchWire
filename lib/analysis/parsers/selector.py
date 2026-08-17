from enum import Enum
from .yaml_parser import YamlRunParser
from .json_parser import JsonRunParser

# Add every new file format here.
class FORMATS(Enum):
    yaml = ".yaml"
    html = ".html"
    json = ".json"
    custom = "custom"

# Dispatches to the right parser per-run for each needed file format.
def Selector(formats, run, matrixName, labels, customData=None):

    parsedRuns = {
        "formatTypes" : formats,
        "run" : run,
        "data" : {},
        "labels" : labels,
    }
    
    comparison =  len(labels) > 1
    
    # Add case per file format.
    for val in formats:
        match val:
            case FORMATS.custom:
                customRun = customData[run - 1]
                if comparison:
                    parsedRuns["data"]["customA"] = customRun[0]
                    parsedRuns["data"]["customB"] = customRun[1]
                else:
                    parsedRuns["data"]["custom"] = customRun

            case FORMATS.json:
                if comparison:
                    parsedRuns["data"]["jsonA"] = JsonRunParser(matrixName, labels[0], "A", run)
                    parsedRuns["data"]["jsonB"] = JsonRunParser(matrixName, labels[1], "B", run)
                else:
                    parsedRuns["data"]["json"] = JsonRunParser(matrixName, labels[0], "", run)

            case FORMATS.yaml:
                if comparison:
                    parsedRuns["data"]["yamlA"] = YamlRunParser(matrixName, labels[0], "A", run)
                    parsedRuns["data"]["yamlB"] = YamlRunParser(matrixName, labels[1], "B", run)
                else:
                    parsedRuns["data"]["yaml"] = YamlRunParser(matrixName, labels[0], "", run)

            case _:
                raise NotImplementedError(f"No parser implemented for format: {val}")

    return parsedRuns

