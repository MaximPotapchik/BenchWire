def GetStaticStats():
    return ["mode", "cpu_name", "llvm_triple", "key", "info", "min_instructions"]

def GetMeasurementPreset():
    return {"collection": "measurements", "idField": "key", "subFields": ["value", "per_snippet_value"]}

def GetTrackedStats(stats):
    gotStats = stats.GetAllStats()
    trackedStats = []
    for key in gotStats.keys():
        if ":value" in key:
            trackedStats.append(key)
    return trackedStats
