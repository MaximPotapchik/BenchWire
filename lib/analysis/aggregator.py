import numpy as np
from .parsers.selector import Selector, FORMATS

def GetAllKeys(d):
    for key, value in d.items():
        yield key

        if isinstance(value, dict):
            yield from GetAllKeys(value)

        elif isinstance(value, list):
            for item in value:
                if isinstance(item, dict):
                    yield from GetAllKeys(item)
            
def FormatStats(yamlOutput, preset):
    collectedStats = {}
    foundStats = GetAllKeys(yamlOutput)

    for key in foundStats:
        if key in preset:
            collectedStats[key] = "static"
        else:
            collectedStats[key] = ""

    return collectedStats

# TODO: Refactor this. It's annoyingly nested.
def FindValue(data, targetKey):
    if isinstance(data, dict):
        for key, value in data.items():
            if key == targetKey:
                return value

            if isinstance(value, dict):
                found = FindValue(value, targetKey)
                if found is not None:
                    return found

            elif isinstance(value, list):
                for item in value:
                    if isinstance(item, dict):
                        found = FindValue(item, targetKey)
                        if found is not None:
                            return found

    return None

def FillStatArrays(statList, collectedStats, parsedData, run, labelCnt, isInitialized):
    runIdx = run - 1

    for i, opt in enumerate(collectedStats):
        isStatic = collectedStats[opt] == "static"
        if isStatic and isInitialized:
            continue 

        if labelCnt:
            valA = FindValue(parsedData.get("yamlA", {}), opt)
            valB = FindValue(parsedData.get("yamlB", {}), opt)
            pairs = [(2 * i, valA), (2 * i + 1, valB)]
        else:
            val = FindValue(parsedData.get("yaml", {}), opt)
            pairs = [(i, val)]

        for idx, val in pairs:
            entry = statList[idx]
            # TODO: Make this safer.
            if isStatic:
                entry["values"] = val
            elif isinstance(entry["values"], list):
                entry["values"].append(val)
            else:
                entry["values"][runIdx] = val

    return True

# TODO: - Algorithm can be made more efficient.
# - Needs to sample both sides.
def FillStats(preset, runs, labels, labelCnt):
    checkParsedStats = Selector([FORMATS.yaml], 1, labels)
    checkParsedData = checkParsedStats.get("data", {})
    collectedStats = FormatStats(checkParsedData["yamlA"], preset) if labelCnt else FormatStats(checkParsedData["yaml"], preset)
    statList = []

    for key in collectedStats:
        isStatic = collectedStats[key] == "static"
        sample = FindValue(checkParsedData["yamlA" if labelCnt else "yaml"], key)
        isNumeric = isinstance(sample, (int, float))

        for _ in range(2 if labelCnt else 1):
            if isStatic:
                values = None
            elif isNumeric:
                values = np.zeros(runs)
            else:
                values = []
            statList.append({"statName": key, "values": values})

    isInitialized = False

    for run in range(1, runs + 1):
        parsedStats = Selector([FORMATS.yaml], run, labels)
        parsedData = parsedStats.get("data", {})
        ran = FillStatArrays(statList, collectedStats, parsedData, run, labelCnt, isInitialized)
        isInitialized = ran

    return statList

# This builds the stat value arrays.
def Aggregate(runs, statsPreset, labels):
    labelCnt = len(labels) > 1
    statList = FillStats(statsPreset, runs, labels, labelCnt)
        
    allStats = {}

    if labelCnt:
        labelA, labelB = labels[:2]

        allStats[labelA] = {
            "runs" : runs,
            "stats" : statList[::2]
        }
        allStats[labelB] = {
            "runs" : runs,
            "stats" : statList[1::2]
        }

    else:
        label = labels[0]

        allStats[label] = {
            "runs": runs,
            "stats": statList
        }

    return allStats
