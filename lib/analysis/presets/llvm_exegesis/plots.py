from lib.analysis.plotting.plot import PlotGrid, BasePlot
from lib.analysis.plotting.theme import mocha

def ExegesisPlot(stats, fullArgs, targetCnt, outputDir, timestamp):
    if targetCnt > 1:
        statsA, statsB = stats
        plot = BasePlot([statsA, statsB], mocha)
        plot.SetLabels(title="Exegesis Runs, Comparison", xlabel="N of runs", ylabel="Cycles")
    else:
        plot = BasePlot(stats, mocha)
        plot.SetLabels(title="Exegesis Runs", xlabel="N of runs", ylabel="Cycles")

    plot.CleanXValues()
    plot.AddTimestamp(timestamp)
    plot.Render("value")

    name = f"plot_{timestamp}.png"
    plot.Save(outputDir, name)

    return name
