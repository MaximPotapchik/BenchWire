from enum import Enum
from .yaml_parser import YamlParser

# Add every new file format for fromats here.
class FORMATS(Enum):
    yaml = ".yaml"
    # html = ".html" next

# Dispatches to the right parser per-run for each requested file format.
def Selector(formats, run, collectedStats, labels):

    parsedRuns = {
        "formatTypes" : formats,
        "run" : run,
        "data" : {},
        "labels" : labels,
    }
    
    comparison = "|" in labels 
    
    # Add case per file format.
    for val in formats:
        match val:
            case FORMATS.yaml:            
                if comparison:
                    parsedRuns["data"]["yamlA"] = YamlParser("A", run, collectedStats)
                    parsedRuns["data"]["yamlB"] = YamlParser("B", run, collectedStats)
                else:
                    parsedRuns["data"]["yaml"] = YamlParser("", run, collectedStats)
            # Fallback.
            case _:
                raise NotImplementedError(f"No parser implemented for format: {val}")

    return parsedRuns

