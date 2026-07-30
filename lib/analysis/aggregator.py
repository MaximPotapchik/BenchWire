import numpy as np
from .parsers.selector import Selector, FORMATS

def GetAllKeys(d, skipCollection=None):
    for key, value in d.items():
        if key == skipCollection:
            continue
        yield key

        if isinstance(value, dict):
            yield from GetAllKeys(value, skipCollection)

        elif isinstance(value, list):
            for item in value:
                if isinstance(item, dict):
                    yield from GetAllKeys(item, skipCollection)
            
def ExtractMeasurementStats(yamlOutput, measurementPreset):
    collection = yamlOutput.get(measurementPreset["collection"], [])
    idField = measurementPreset["idField"]
    stats = {}
    for entry in collection:
        for field in measurementPreset["subFields"]:
            stats[f"{entry[idField]}:{field}"] = entry[field]
    return stats 

def FormatStats(yamlOutput, preset, measurementPreset=None):
    collectedStats = {}
    skipCollection = measurementPreset["collection"] if measurementPreset else None

    for key in GetAllKeys(yamlOutput, skipCollection):
        collectedStats[key] = "static" if key in preset else ""

    if measurementPreset:
        for statName in ExtractMeasurementStats(yamlOutput, measurementPreset):
            collectedStats[statName] = ""

    return collectedStats

# TODO: Refactor this. It's annoyingly nested.
def FindValue(data, targetKey, measurementPreset=None):
    if measurementPreset and ":" in targetKey:
        idValue, field = targetKey.split(":", 1)

        for entry in data.get(measurementPreset["collection"], []):
            if entry.get(measurementPreset["idField"]) == idValue:
                return entry.get(field)
        return None

    if isinstance(data, dict):
        for key, value in data.items():
            if key == targetKey:
                return value

            if isinstance(value, dict):
                found = FindValue(value, targetKey, measurementPreset)
                if found is not None:
                    return found

            elif isinstance(value, list):
                for item in value:
                    if isinstance(item, dict):
                        found = FindValue(item, targetKey, measurementPreset)
                        if found is not None:
                            return found

    return None

def FillStatArrays(statList, collectedStats, parsedData, run, labelCnt, isInitialized, measurementPreset=None):
    runIdx = run - 1

    for i, opt in enumerate(collectedStats):
        isStatic = collectedStats[opt] == "static"
        if isStatic and isInitialized:
            continue 

        if labelCnt:
            valA = FindValue(parsedData.get("yamlA", {}), opt, measurementPreset)
            valB = FindValue(parsedData.get("yamlB", {}), opt, measurementPreset)
            pairs = [(2 * i, valA), (2 * i + 1, valB)]
        else:
            val = FindValue(parsedData.get("yaml", {}), opt, measurementPreset)
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
def FillStats(preset, measurementPreset, runs, labels, labelCnt):
    checkParsedStats = Selector([FORMATS.yaml], 1, labels)
    checkParsedData = checkParsedStats.get("data", {})
    checkData = checkParsedData["yamlA"] if labelCnt else checkParsedData["yaml"]

    collectedStats = FormatStats(checkData, preset, measurementPreset)
    statList = []

    for key in collectedStats:
        isStatic = collectedStats[key] == "static"
        sample = FindValue(checkData, key, measurementPreset)
        isNumeric = isinstance(sample, (int, float))
        for _ in range(2 if labelCnt else 1):
            values = None if isStatic else (np.zeros(runs) if isNumeric else [])
            statList.append({"statName": key, "values": values})
    isInitialized = False

    for run in range(1, runs + 1):
        parsedStats = Selector([FORMATS.yaml], run, labels)
        parsedData = parsedStats.get("data", {})
        ran = FillStatArrays(statList, collectedStats, parsedData, run, labelCnt, isInitialized, measurementPreset)
        isInitialized = ran

    return statList

# This builds the stat value arrays.
def Aggregate(runs, preset, labels):
    staticPreset, measurementPreset = preset
    labelCnt = len(labels) > 1
    statList = FillStats(staticPreset, measurementPreset, runs, labels, labelCnt)
        
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
