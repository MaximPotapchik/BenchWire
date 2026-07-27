import os
from datetime import datetime
from .stats import StatsResult
from .aggregator import Aggregate
from .presets.llvm_exegesis.markdown import ExegesisMarkdown
from .presets.llvm_exegesis.plots import ExegesisPlot
from .presets.llvm_exegesis.statpreset import GetStaticStats

def Pipeline(fullArgs, outputLocation):
    targetCnt = len(fullArgs["targets"][:2])
    labels = [target['label'] for target in fullArgs['targets']]

    aggregated = Aggregate(fullArgs["runs"], GetStaticStats(), labels) 

    # Output Directory
    plotDir = os.path.join(outputLocation, "results", "plots")

    # Timestamp
    now = datetime.now()
    timestamp = now.strftime("%y%m%d%H%M%S")

    if targetCnt > 1:
        statsA, statsB = StatsResult.FromAggregate(aggregated)
        plot = ExegesisPlot([statsA, statsB], fullArgs, targetCnt, plotDir, timestamp)
        ExegesisMarkdown([statsA, statsB], fullArgs, targetCnt, plotDir, timestamp, plot)
    else:
        stats = StatsResult.FromAggregate(aggregated)
        plot = ExegesisPlot(stats, fullArgs, targetCnt, plotDir, timestamp)
        ExegesisMarkdown(stats, fullArgs, targetCnt, plotDir, timestamp, plot)
   
    print(f"[BenchWire] Analysis completed for Batch #{timestamp}.")
