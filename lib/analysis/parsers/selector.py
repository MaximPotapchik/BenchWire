from enum import Enum
from .yaml_parser import YamlRunParser

# Add every new file format here.
class FORMATS(Enum):
    yaml = ".yaml"
    # html = ".html" next

# Dispatches to the right parser per-run for each needed file format.
def Selector(formats, run, matrixName, labels):

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
                    parsedRuns["data"]["yamlA"] = YamlRunParser(matrixName, labels[0], "A", run)
                    parsedRuns["data"]["yamlB"] = YamlRunParser(matrixName, labels[1], "B", run)
                else:
                    parsedRuns["data"]["yaml"] = YamlRunParser(matrixName, labels[0], "", run)
            # Fallback.
            case _:
                raise NotImplementedError(f"No parser implemented for format: {val}")

    return parsedRuns

