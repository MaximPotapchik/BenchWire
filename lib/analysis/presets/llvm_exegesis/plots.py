import math
from lib.analysis.plotting.plot import PlotGrid, BasePlot
from lib.analysis.plotting.theme import mocha
from .statpreset import GetTrackedStats

def ExegesisPlot(stats, fullArgs, targetCnt, outputDir, timestamp):
    firstStats = stats[0] if targetCnt > 1 else stats
    allTrackedStats = GetTrackedStats(firstStats)

    trackedStats = []
    for i in range(len(allTrackedStats)):
        val = allTrackedStats[i]
        if not firstStats.IsAllZero(val):
            trackedStats.append(val)

    n = len(trackedStats)
    name = f"plot_{timestamp}.png"
    
    def buildPlot(val):
        if targetCnt > 1:
            statsA, statsB = stats
            plot = BasePlot([statsA, statsB], mocha)
            mode = statsA.GetStat("mode")
        else:
            plot = BasePlot(stats, mocha)
            mode = stats.GetStat("mode")

        ylabel = "Uops" if mode == "uops" else "Cycles"
        plot.SetLabels(title=val.split(":")[0], xlabel="N of runs", ylabel=ylabel)
        plot.CleanXValues()
        return plot

    if n == 1:
        plot = buildPlot(trackedStats[0])
        plot.AddTimestamp(timestamp)
        plot.Render(trackedStats[0])
        plot.Save(outputDir, name)
        return name

    cols = min(3, n)
    rows = math.ceil(n / cols)
    grid = PlotGrid(mocha, rows, cols, figsize=(3.5 * cols, 2.8 * rows))

    for i in range(n):
        val = trackedStats[i]
        row = i // cols
        col = i % cols
        grid.Add(buildPlot(val), val, row, col)

    grid.AddTimestamp(timestamp)
    grid.SetPadding(0.7)
    grid.Render()
    grid.Save(outputDir, name)
    return name
