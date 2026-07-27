from enum import Enum
from .yaml_parser import YamlRunParser

# Add every new file format for fromats here.
class FORMATS(Enum):
    yaml = ".yaml"
    # html = ".html" next

# Dispatches to the right parser per-run for each needed file format.
def Selector(formats, run, labels):

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
            case FORMATS.yaml:            
                if comparison:
                    parsedRuns["data"]["yamlA"] = YamlRunParser("A", run)
                    parsedRuns["data"]["yamlB"] = YamlRunParser("B", run)
                else:
                    parsedRuns["data"]["yaml"] = YamlRunParser("", run)
            # Fallback.
            case _:
                raise NotImplementedError(f"No parser implemented for format: {val}")

    return parsedRuns

