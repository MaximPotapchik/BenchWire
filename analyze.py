import os
import yaml
import numpy as np
import matplotlib.pyplot as plt
import argparse
from enum import Enum
from datetime import datetime


class MODE(Enum):
    SINGLE = "single"
    COMPARE = "compare"

# cli parser
parser = argparse.ArgumentParser(description="Python CLI parser")
parser.add_argument("benchmode", type=str, choices=[m.value for m in MODE])
parser.add_argument("runs", type=int)
parser.add_argument("labels", type=str)
parser.add_argument("methodology", type=str)
parser.add_argument("cooldowntimer", type=int)

args = parser.parse_args()
BenchMode = args.benchmode
Runs = args.runs
Methodology = args.methodology
CooldownTimer = args.cooldowntimer

if "|" in args.labels:
    Labels = args.labels.split("|")
elif args.labels != "":
    Labels = args.labels
else:
    Labels = ""

if BenchMode == MODE.COMPARE.value and (not isinstance(Labels, list) or len(Labels) != 2):
    parser.error("compare mode needs two labels separated by '|', e.g. \"Raw|Libpfm\"")

ExegesisMode = "Latency"
YAML_DIR = os.path.join("results", "yaml")

# YAML parsing functions. Currently only supports latency.
def LoadRun(Prefix, Run, Mode="latency"):
    Path = os.path.join(YAML_DIR, f"{Prefix}run_{Run}.yaml")
    with open(Path) as f:
        data = yaml.safe_load(f)
    Measurements = data.get("measurements")
    if not Measurements:
        raise ValueError(f"No measurements found in {Path}")
    return Measurements[0]

def CollectValues(Prefix, Runs):
    Arr = np.empty(Runs)
    for Run in range(1, Runs + 1):
        Arr[Run - 1] = LoadRun(Prefix, Run).get('value')
    return Arr

# Aesthetics
mocha = {
    "blue": "#89b4fa",
    "red": "#f38ba8",
    "green": "#a6e3a1",
    "text": "#cdd6f4",
    "base": "#1e1e2e",
}

# plot - Yes it is a lot of styling
if BenchMode == "single":
    RunCount = np.arange(1, Runs + 1)
    ValueArr = CollectValues("", Runs)
    plt.figure(facecolor=mocha["base"])
    ax = plt.gca()
    ax.set_facecolor(mocha["base"])
    ax.tick_params(colors=mocha["text"])
    ax.xaxis.label.set_color(mocha["text"])
    ax.yaxis.label.set_color(mocha["text"])
    ax.title.set_color(mocha["text"])
    plt.plot(RunCount, ValueArr, color=mocha["blue"], label=Labels)
    plt.title("Exegesis Runs")
    plt.xlabel("N of runs")
    plt.ylabel("Cycles")
    plt.legend()

if BenchMode == "compare":
    RunCount = np.arange(1, Runs + 1)
    ValueArrA = CollectValues("A", Runs)
    ValueArrB = CollectValues("B", Runs)
    plt.figure(facecolor=mocha["base"])
    ax = plt.gca()
    ax.set_facecolor(mocha["base"])
    ax.tick_params(colors=mocha["text"])
    ax.xaxis.label.set_color(mocha["text"])
    ax.yaxis.label.set_color(mocha["text"])
    ax.title.set_color(mocha["text"])
    plt.plot(RunCount, ValueArrA, color=mocha["blue"], label=Labels[0])
    plt.plot(RunCount, ValueArrB, color=mocha["red"], label=Labels[1])
    plt.title("Exegesis runs")
    plt.xlabel("N of runs")
    plt.ylabel("Cycles")
    plt.legend()

# output
OutputDir = "results/plots"
Time = datetime.now()
CurTime = Time.strftime("%d_%H%M%S")
plt.savefig(f"{OutputDir}/plot{CurTime}.png", dpi=300)

with open(f"{OutputDir}/result{CurTime}.md", "w") as f:
    f.write(f"## {ExegesisMode} | Mode - {BenchMode} \n\n")
    f.write(f"[plt](plot{CurTime}.png)\n\n")
    if BenchMode == "single":
        f.write(f"## Statistics\n\n")
        f.write(f"{Runs} Runs | Cooldown time - {CooldownTimer}ms\n")
        f.write(f"Mean: {ValueArr.mean():.4f}\n")
        f.write(f"Median: {np.median(ValueArr):.4f}\n")
        f.write(f"Standard Deviation: {ValueArr.std():.4f}\n")
        f.write(f"Coefficient of Variation: {ValueArr.std() / ValueArr.mean() * 100:.4f}%\n")
        f.write(f"Min: {ValueArr.min():.4f}\n")
        f.write(f"Max: {ValueArr.max():.4f}\n\n")
        f.write(f"Percentile statistics\n")
        f.write(f"P50: {np.percentile(ValueArr, 50):.4f}\n")
        f.write(f"P75: {np.percentile(ValueArr, 75):.4f}\n")
        f.write(f"P90: {np.percentile(ValueArr, 90):.4f}\n")
        f.write(f"P99: {np.percentile(ValueArr, 99):.4f}\n")
        f.write(f"P99.9: {np.percentile(ValueArr, 99.9):.4f}\n")
    if BenchMode == "compare":
        f.write(f"## Statistics\n\n")
        f.write(f"A. {Labels[0]} B. {Labels[1]}\n\n")
        f.write(f"Methodology - {Methodology} | {Runs} Runs | Cooldown time - {CooldownTimer}ms\n\n")
        f.write(f"Mean | A: {ValueArrA.mean():.4f} B: {ValueArrB.mean():.4f}\n")
        MeanDiff = ValueArrB.mean() - ValueArrA.mean()
        MeanDiffPct = abs(MeanDiff) / ((ValueArrA.mean() + ValueArrB.mean()) / 2) * 100
        f.write(f"Mean Difference: {MeanDiffPct:.4f}% ({Labels[1] if MeanDiff > 0 else Labels[0]} higher)\n\n")
        f.write(f"Median | A: {np.median(ValueArrA):.4f} B: {np.median(ValueArrB):.4f}\n")
        MedianDiff = np.median(ValueArrB) - np.median(ValueArrA)
        MedianDiffPct = abs(MedianDiff) / ((np.median(ValueArrA) + np.median(ValueArrB)) / 2) * 100
        f.write(f"Median Difference: {MedianDiffPct:.4f}% ({Labels[1] if MedianDiff > 0 else Labels[0]} higher)\n\n")
        f.write(f"Standard Deviation | A: {ValueArrA.std():.4f} B: {ValueArrB.std():.4f}\n")
        StdDiff = ValueArrB.std() - ValueArrA.std()
        StdDiffPct = abs(StdDiff) / ((ValueArrA.std() + ValueArrB.std()) / 2) * 100
        f.write(f"StdDev Difference: {StdDiffPct:.4f}% ({Labels[1] if StdDiff > 0 else Labels[0]} higher)\n\n")
        CvA = ValueArrA.std() / ValueArrA.mean() * 100
        CvB = ValueArrB.std() / ValueArrB.mean() * 100
        f.write(f"Coefficient of Variation | A: {CvA:.4f}% B: {CvB:.4f}%\n")
        CvDiff = CvB - CvA
        CvDiffPct = abs(CvDiff) / ((CvA + CvB) / 2) * 100
        f.write(f"CoV Difference: {CvDiffPct:.4f}% ({Labels[1] if CvDiff > 0 else Labels[0]} higher)\n\n")
        f.write(f"Min | A: {ValueArrA.min():.4f} B: {ValueArrB.min():.4f}\n")
        f.write(f"Max | A: {ValueArrA.max():.4f} B: {ValueArrB.max():.4f}\n\n")
        f.write(f"Percentile statistics\n")
        f.write(f"P50  | A: {np.percentile(ValueArrA, 50):.4f}  B: {np.percentile(ValueArrB, 50):.4f}\n")
        f.write(f"P75  | A: {np.percentile(ValueArrA, 75):.4f}  B: {np.percentile(ValueArrB, 75):.4f}\n")
        f.write(f"P90  | A: {np.percentile(ValueArrA, 90):.4f}  B: {np.percentile(ValueArrB, 90):.4f}\n")
        f.write(f"P99  | A: {np.percentile(ValueArrA, 99):.4f}  B: {np.percentile(ValueArrB, 99):.4f}\n")
        f.write(f"P99.9| A: {np.percentile(ValueArrA, 99.9):.4f} B: {np.percentile(ValueArrB, 99.9):.4f}\n")
