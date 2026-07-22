import numpy as np
from parsers.selector import Selector, FORMATS

# This builds the stat value arrays and freezes non-numeric stats
# TODO: - Definitely modularizable further. 
def Aggregate(runs, collectedStats, labels):
    comparison = "|" in labels
    statList = []
    for opt in collectedStats:
        if comparison:
            statList.append({"statName": opt, "values": np.zeros(runs)})
            statList.append({"statName": opt, "values": np.zeros(runs)})
        else:
            statList.append({"statName": opt, "values": np.zeros(runs)})

    skippedStats = set()

    # TODO: This can be optimized and made into a seperate function
    for run in range(1, runs + 1):
        # TODO: de-hardcode this. Make it a function probably
        parsedStats = Selector([FORMATS.yaml], run, collectedStats, labels)
        parsedData = parsedStats.get("data", {})
        runIdx = run - 1

        if comparison:
            for i, opt in enumerate(collectedStats):
                if opt in skippedStats:
                    continue
                valA = parsedData.get("yamlA", {}).get(opt)
                valB = parsedData.get("yamlB", {}).get(opt)

                if run == 1 and (not isinstance(valA, (int, float)) or not isinstance(valB, (int, float))):
                    statList[2 * i]["values"] = str(valA)
                    statList[2 * i + 1]["values"] = str(valB)
                    skippedStats.add(opt)
                    continue
                statList[2 * i]["values"][runIdx] = valA
                statList[2 * i + 1]["values"][runIdx] = valB

        else: 
            for i, opt in enumerate(collectedStats):
                if opt in skippedStats:
                    continue                
                val = parsedData.get("yaml", {}).get(opt)
                if run == 1 and not isinstance(val, (int, float)):
                    statList[i]["values"] = str(val)
                    skippedStats.add(opt)
                    continue
                statList[i]["values"][runIdx] = val
    
    allStats = {}

    if comparison:
        labelA, labelB = labels.split("|")

        allStats[labelA] = {
            "runs" : runs,
            "stats" : statList[::2]
        }
        allStats[labelB] = {
            "runs" : runs,
            "stats" : statList[1::2]
        }

    else:
        allStats[labels] = {
            "runs" : runs,
            "stats" : statList
        }

    return allStats
