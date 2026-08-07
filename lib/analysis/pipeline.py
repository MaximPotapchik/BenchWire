import os
import time
from datetime import datetime
from .stats import StatsResult
from .aggregator import Aggregate
from .reporting.markdown import MarkdownReporter
from .presets.llvm_exegesis.markdown import ExegesisMarkdown
from .presets.llvm_exegesis.plots import ExegesisPlot
from .presets.llvm_exegesis.statpreset import GetStaticStats, GetMeasurementPreset

# These can be seperated from pipeline in the future if this gets too large.
def ExegesisPipe(specMatrix, plotDir, timestamp, itr, combine, outputPer, lastMatrix):

    # TODO: If checks per setting can turn into a function to reduce boilerplate.
    if outputPer == "single":
        if combine is None:
            combine = MarkdownReporter(None)
    elif outputPer == "perMatrix":
        combine = MarkdownReporter(None)
    elif outputPer == "perRun":
        combine = None  

    seq = specMatrix["sequence"]
    for i, labels in enumerate(specMatrix["sequence"]):
        
        labels = [labels] if isinstance(labels, str) else labels
        exegesisPreset = [GetStaticStats(), GetMeasurementPreset()]
        aggregated = Aggregate(specMatrix["runs"], exegesisPreset, specMatrix["name"], labels) 
        targetCnt = len(labels)
        
        lastStep = (i == len(seq) - 1)
        finalized= None
        if outputPer == "single" and lastMatrix and lastStep:
            finalized = ("BenchWire Batch", f"report_{timestamp}.md")
        elif outputPer == "perMatrix" and lastStep:
            finalized = (specMatrix["name"], f"{specMatrix['name']}.md")
        n = itr + 1
        if targetCnt > 1:
            statsA, statsB = StatsResult.FromAggregate(aggregated)
            plot = ExegesisPlot([statsA, statsB], specMatrix, targetCnt, plotDir, [n, timestamp])
            ExegesisMarkdown([statsA, statsB], specMatrix, targetCnt, plotDir, [n, timestamp], plot, combine=combine, finalized=finalized)

        else:
            stats = StatsResult.FromAggregate(aggregated)
            plot = ExegesisPlot(stats, specMatrix, targetCnt, plotDir, [n, timestamp])
            ExegesisMarkdown(stats, specMatrix, targetCnt, plotDir, [n, timestamp], plot, combine=combine, finalized=finalized)

        itr += 1

    return itr, (None if outputPer == "perMatrix" else combine)

def Pipeline(fullArgs, outputLocation):
    
    analysisOpts = fullArgs["analysis"]
    # Output Directory
    resultLocation = os.path.join(outputLocation, "results")

    # Timestamp for batch
    now = datetime.now()
    timestamp = now.strftime("%y%m%d%H%M%S")

    # TODO: For single file, something like if it exists with a space, it is a single file, is fine for now.
    if analysisOpts["output"].get("dirName"):
        resultLocation = os.path.join(resultLocation, analysisOpts["output"]["dirName"])
        os.makedirs(resultLocation, exist_ok=True)
        
    elif analysisOpts["output"]["storeIn"] == "directory" and analysisOpts["output"]["per"] == "suite":
        resultLocation = os.path.join(resultLocation, f"batch_{timestamp}")
        os.makedirs(resultLocation, exist_ok=True)

    elif analysisOpts["output"]["storeIn"] == "single file":
        resultLocation = os.path.join(resultLocation, "single_files")
        os.makedirs(resultLocation, exist_ok=True)
    
    storeIn, per = analysisOpts["output"]["storeIn"], analysisOpts["output"]["per"]
    # TODO: May need to chenge this. collides with methodology single.
    outputPer = "single" if (storeIn == "single file" or per == "suite") else ("perMatrix" if per == "specMatrix" else "perRun")

    runN = 0
    combine = None
    specMatrices = fullArgs["specMatrix"]

    print(f"[BenchWire] Analysis initiated for Batch #{timestamp}.")
    analysisTimeStart = time.perf_counter() 
    for idx, specMatrix in enumerate(specMatrices):
        if specMatrix["benchmarker"] == "llvm-exegesis":

            saveLocation = resultLocation
            if per == "specMatrix" and storeIn == "directory":
                saveLocation = os.path.join(resultLocation, specMatrix["name"])
                os.makedirs(saveLocation, exist_ok=True)

            isLastMatrix = (idx == len(specMatrices) - 1)
            runN, combine = ExegesisPipe(specMatrix, saveLocation, timestamp, runN, combine, outputPer, isLastMatrix)

    analysisTimeEnd = time.perf_counter()
    # Gets the time it took to do the anaylsis.
    analysisTime = (analysisTimeEnd - analysisTimeStart) * 1000

    print(f"[BenchWire] Analysis completed in {analysisTime:.3f}ms for Batch #{timestamp}.")
