import math
import matplotlib.pyplot as plt
import numpy as np 

class BasePlot:
    def __init__(self, statsResult, theme, ax=None):
        self.stats = statsResult if isinstance(statsResult, list) else [statsResult]
        self.theme = theme 
        self.arguments = []

        if ax is None:
            self.fig, self.ax = plt.subplots()
        else:
            self.fig, self.ax = ax.figure, ax

    def Draw(self, statName):
        for stats in self.stats:
            plotRuns = np.arange(1, stats.GetRuns() + 1)
            self.ax.plot(plotRuns, stats.GetStat(statName), label=stats.label)
    
    def ApplyTheme(self):
        self.fig.patch.set_facecolor(self.theme["base"])
        self.ax.set_facecolor(self.theme["base"])
        self.ax.tick_params(colors=self.theme["text"])
        self.ax.title.set_color(self.theme["text"])
        self.ax.xaxis.label.set_color(self.theme["text"])
        self.ax.yaxis.label.set_color(self.theme["text"])

    def AddArg(self, *args, **kwargs):
        self.arguments.append((args, kwargs))
        return self
    
    # Built-in args.
    def SetLabels(self, title=None, xlabel=None, ylabel=None, legend=True):
        if title:
            self.AddArg("set_title", title)
        if xlabel:
            self.AddArg("set_xlabel", xlabel)
        if ylabel:
            self.AddArg("set_ylabel", ylabel)
        if legend:
            self.AddArg("legend", fontsize=plt.rcParams["font.size"] * 0.75)
        return self

    # Expects same amount of runs in comparison which is hardcoded anyways.
    def CleanXValues(self):
        runs = self.stats[0].GetRuns()
        rough = max(runs / 10, 1)
        magnitude = 10 ** math.floor(math.log10(rough))
        step = next((s * magnitude for s in [1, 2, 5, 10] if s * magnitude >= rough), 10 * magnitude)
        step = max(1, int(step))
        ticks = sorted(set([1] + list(range(step, runs, step)) + [runs]))
        self.AddArg("set_xticks", ticks)
        self.AddArg("xaxis.set_major_formatter", "{x:.0f}")
        self.AddArg("ticklabel_format", axis="y", useOffset=False, style="plain")

    def ThemeLegend(self):
        legend = self.ax.get_legend()
        if legend:
            legend.get_frame().set_facecolor(self.theme["base"])
            legend.get_frame().set_edgecolor(self.theme["text"])
            for text in legend.get_texts():
                text.set_color(self.theme["text"])
    
    def AddTimestamp(self, timestamp):
        self.fig.text(0.98, 0.98, timestamp, ha="right", va="top", fontsize=8, color=self.theme["text"])
        return self

    def Render(self, statName):
        self.ApplyTheme()
        self.Draw(statName)
        for args, kwargs in self.arguments:
            target = self.ax
            *path, methodName = args[0].split(".")
            for attr in path:
                target = getattr(target, attr)
            getattr(target, methodName)(*args[1:], **kwargs)
        self.ThemeLegend()
        return self.fig

    def Save(self, path, name):
        self.fig.savefig(f"{path}/{name}", dpi=200)

# This class holds instances for BasePlots. 
class PlotGrid:
    def __init__(self, theme, rows, cols, figsize):
        self.fig, self.axes = plt.subplots(rows, cols, figsize=figsize, squeeze=False, layout="constrained")
        self.theme = theme
        self.plots = []
        self.origWidth = figsize[0]
        self.origHeight = figsize[1]

    # TODO: Auto-layout the plots instead of hardcoding the row/cols.
    def Add(self, plot, statName, row, col):
        cellAx = self.axes[row][col]
        plt.close(plot.fig)
        plot.fig, plot.ax = cellAx.figure, cellAx
        self.plots.append((plot, statName))
        return plot 
    
    def AddTimestamp(self, timestamp):
        self.fig.text(0.99, 0.98, timestamp, ha="right", va="top", fontsize=8, color=self.theme["text"])
        return self

    def SetPadding(self, extraRight=1.5):
        newWidth = self.origWidth + extraRight
        self.fig.set_size_inches(newWidth, self.origHeight)
        self.fig.get_layout_engine().set(rect=(0, 0, self.origWidth / newWidth, 1))
        return self
    
    def Render(self):
        for plot, statName in self.plots:
            plot.Render(statName)
        return self.fig

    def Save(self, path, name):
        self.fig.savefig(f"{path}/{name}", dpi=175)
