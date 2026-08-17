from lib.analysis.reporting.markdown import MarkdownReporter
from .statpreset import GetTrackedStats

# Uses the markdown.py api from lib/analysis/reporting/markdown.py and reporter.py

# Default mode. Creates the markdown report per measurement.
def BuildMarkdownReport(stats, methodology, summary, targetCnt, trackedStats):
    
    r = None
    opcode = ""
    mode = ""

    if targetCnt > 1:
        statsA, statsB = stats

        try:
            opcode = statsA.GetStat("instructions")[0][0].split()[0]
        except (KeyError, IndexError, TypeError, AttributeError):
            opcode = None 
        
        try:
            mode = statsA.GetStat("mode")
        except (KeyError, IndexError, TypeError, AttributeError):
            mode = None

        r = MarkdownReporter(statsA, statsB)    

        r.AddLine(f"**Methodology:** {methodology}")
        r.AddLine(f"**Runs:** {statsA.GetRuns()}")
        r.AddLine(f"**CPU:** {statsA.GetStat('cpu_name')}")
        r.AddLine(f"**Triple:** {statsA.GetStat('llvm_triple')}")
        r.AddLine(f"**Min instructions:** {statsA.GetStat('min_instructions')}")

        for val in trackedStats:
            if not statsA.IsAllZero(val) and not statsB.IsAllZero(val):
                r.AddLine(f"## Measurement - {val.split(':')[0]}")

                r.AddLine("### Summary Statistics")
                r.StartTable(True, True)
                r.SetTableColumns("Metric", statsA.label, statsB.label)
                for label, method, unit in summary:
                    a, b = getattr(statsA, method)(val), getattr(statsB, method)(val)
                    r.AddTableRow(label, a, b, statsA.Diff(a, b), unit=unit)

                r.AddLine("### Percentile Statistics")
                r.StartTable(True, True)
                r.SetTableColumns("Percentile", statsA.label, statsB.label)
                for p in [50, 75, 90, 99, 99.9, 99.99]:
                    pA, pB = statsA.P(val, p), statsB.P(val, p)
                    r.AddTableRow(f"P{p}", pA, pB, statsA.Diff(pA, pB))
            else:
                r.AddLine(f"## Measurement - {val.split(':')[0]} = 0")

    else:
        try:
            opcode = stats.GetStat("instructions")[0][0].split()[0]
        except (KeyError, IndexError, TypeError, AttributeError):
            opcode = None 
        
        try:
            mode = stats.GetStat("mode")
        except (KeyError, IndexError, TypeError, AttributeError):
            mode = None

        r = MarkdownReporter(stats)

        r.AddLine(f"**Methodology:** {methodology}")
        r.AddLine(f"**Runs:** {stats.GetRuns()}")
        r.AddLine(f"**CPU:** {stats.GetStat('cpu_name')}")
        r.AddLine(f"**Triple:** {stats.GetStat('llvm_triple')}")
        r.AddLine(f"**Min instructions:** {stats.GetStat('min_instructions')}")
       
        for val in trackedStats:
            if not stats.IsAllZero(val):
                r.AddLine(f"## Measurement - {val.split(':')[0]}")

                r.AddLine("### Summary Statistics")
                r.StartTable(False)
                r.SetTableColumns("Metric", stats.label)
                for label, method, unit in summary:
                    r.AddTableRow(label, getattr(stats, method)(val), unit=unit)

                r.AddLine("### Percentile statistics")
                r.StartTable(False)
                r.SetTableColumns("Percentile", stats.label)
                for p in [50, 75, 90, 99, 99.9, 99.99]:
                    r.AddTableRow(f"P{p}", stats.P(val, p))
            else:
                r.AddLine(f"## Measurement - {val.split(':')[0]} = 0")

    return r, opcode, mode 

# If the opcode failed in LLVM-exegesis, it shows tthe error message.
def BuildFailureReport(stats, targetCnt):
    if targetCnt > 1:
        statsA, statsB = stats
        firstStats = statsA
        r = MarkdownReporter(statsA, statsB)
    else:
        firstStats = stats
        r = MarkdownReporter(stats)

    try:
        opcode = firstStats.GetStat("instructions")[0][0]
    except (KeyError, IndexError, TypeError, AttributeError):
        opcode = "unknown opcode"

    errorMsg = firstStats.GetStat("error")[0]

    r.AddLine("**Status:** failed")
    r.AddLine(f"**Error:** {errorMsg}")

    return r, opcode, "failed"

# Used for opcode sweep.
def BuildMeasurementSection(stats, val, summary):
    r = MarkdownReporter(stats)
    r.AddLine(f"## Measurement - {val.split(':')[0]}")
    r.AddLine("### Summary Statistics")
    r.StartTable(False)
    r.SetTableColumns("Metric", stats.label)

    for label, method, unit in summary:
        r.AddTableRow(label, getattr(stats, method)(val), unit=unit)

    r.AddLine("### Percentile statistics")
    r.StartTable(False)
    r.SetTableColumns("Percentile", stats.label)

    for p in [50, 75, 90, 99, 99.9, 99.99]:
        r.AddTableRow(f"P{p}", stats.P(val, p))

    return r.RenderBody()

# sweep.yaml presets.
MODE_NAMES = {"Lat": "latency", "Uop": "uops", "Thru": "inverse_throughput"}

# TODO: Make this a method probably.
# Indexes the opcodes together, and builds the proper markdown report for sweep.
def AccumulateSweepStep(combine, stats, trackedStats, summary):
    label = stats.label
    mode = label.split(" ")[-1]
    opcode = label.rsplit(" ", 1)[0]

    if not hasattr(combine, "bodies"):
        combine.bodies = []
        combine.currentIndex = 0
        combine.lastMode = None
        combine.header = None

    if mode != combine.lastMode:
        combine.currentIndex = 0
        combine.lastMode = mode

    if combine.currentIndex >= len(combine.bodies):
        combine.bodies.append({"opcode": opcode, "content": "", "failures": {}})

    entry = combine.bodies[combine.currentIndex]

    # Error messages remain the same if it failed once, so don't repeat them.
    if not trackedStats:
        errorMsg = stats.GetStat("error")[0]
        entry["failures"].setdefault(errorMsg, []).append(mode)

    else:
        for val in trackedStats:
            if not stats.IsAllZero(val):
                entry["content"] += BuildMeasurementSection(stats, val, summary)
            else:
                entry["content"] += f"## Measurement - {val.split(':')[0]} = 0\n\n"
        
        # These are the same every sweep, so we hold the first success.
        if combine.header is None:
            combine.header = {
                "cpu": stats.GetStat("cpu_name"),
                "triple": stats.GetStat("llvm_triple"),
                "minInstructions": stats.GetStat("min_instructions"),
            }

    combine.currentIndex += 1


def RenderIndexed(combine, title, timestamp, methodology):
    headr = combine.header or {}
    out = f"# {title} Batch #{timestamp}\n\n"
    out += f"**Methodology:** {methodology}\n\n**CPU:** {headr.get('cpu','')}\n\n**Triple:** {headr.get('triple','')}\n\n**Min instructions:** {headr.get('minInstructions','')}\n\n"

    for entry in combine.bodies:
        out += f"# {entry['opcode']}\n\n"

        for errorMsg, modes in entry["failures"].items():
            modeNames = "/".join(MODE_NAMES.get(m, m) for m in modes)
            out += f"## Measurement - failed {modeNames}\n\n**Error:** {errorMsg}\n\n"

        out += entry["content"]

    return out

# Interface for every BenchWire mode for LLVM-exegesis.
def ExegesisMarkdown(stats, fullArgs, targetCnt, outputDir, timestamp, plotFile, combine=None, finalized=None):
    
    summary = [
        ("Mean", "Mean", ""),
        ("Median", "Median", ""),
        ("Standard Deviation","StandardDeviation", ""),
        ("Variance", "Variance", "sci"),
        ("Coefficient of Variation", "CoefficientofV", "%"),
        ("Range", "Range", ""),
        ("IQR", "InterquartileRange", ""),
        ("Min", "Min", ""),
        ("Max", "Max", ""),
    ]
    
    if targetCnt > 1:
        trackedStats = GetTrackedStats(stats[0])
    else:
        trackedStats = GetTrackedStats(stats)

    # Calls the sweep report building functions and saves them once finalized is ready.
    if fullArgs["name"] == "Opcode Sweep" and combine is not None:
        AccumulateSweepStep(combine, stats, trackedStats, summary)
        if finalized is not None:
            title, name = finalized
            content = RenderIndexed(combine, title, timestamp[1], fullArgs["methodology"])
            combine.Save(outputDir, name, content)

        return None

    if not trackedStats:
        r, opcode, mode = BuildFailureReport(stats, targetCnt)
    else:
        r, opcode, mode = BuildMarkdownReport(stats, fullArgs["methodology"], summary, targetCnt, trackedStats)

    r.SetPlotFile(plotFile)
    content = r.Render(opcode, mode, timestamp[1])
    name = f"report{timestamp[0]}_{timestamp[1]}.md"
    
    if combine is not None:
        combine.AppendReport(content)
        if finalized is not None:
            title, name = finalized
            combine.Save(outputDir, name, combine.Render(title, "", timestamp[1]))

        return content

    r.Save(outputDir, name, content)

    return content
