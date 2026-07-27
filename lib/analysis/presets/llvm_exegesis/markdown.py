from lib.analysis.reporting.markdown import MarkdownReporter

def ExegesisMarkdown(stats, fullArgs, targetCnt, outputDir, timestamp, plotFile):
    val = "value"
    
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
    
    opcode = ""
    mode = ""
    if targetCnt > 1:
        statsA, statsB = stats
        opcode = statsA.GetStat("instructions")[0][0].split()[0] 
        mode = statsA.GetStat("mode")

        r = MarkdownReporter(statsA, statsB)

        r.AddLine(f"**Methodology:** {fullArgs['methodology']}")
        r.AddLine(f"**Runs:** {statsA.GetRuns()}")
        r.AddLine(f"**CPU:** {statsA.GetStat('cpu_name')}")
        r.AddLine(f"**Triple:** {statsA.GetStat('llvm_triple')}")
        r.AddLine(f"**Min instructions:** {statsA.GetStat('min_instructions')}")

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
        opcode = stats.GetStat("instructions")[0][0].split()[0] 
        mode = stats.GetStat("mode")

        r = MarkdownReporter(stats)
        
        r.AddLine(f"**Methodology:** {fullArgs['methodology']}")
        r.AddLine(f"**Runs:** {stats.GetRuns()}")
        r.AddLine(f"**CPU:** {stats.GetStat('cpu_name')}")
        r.AddLine(f"**Triple:** {stats.GetStat('llvm_triple')}")
        r.AddLine(f"**Min instructions:** {stats.GetStat('min_instructions')}")

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
    
    r.SetPlotFile(plotFile)
    content = r.Render(opcode, mode, timestamp)
    name = f"report_{timestamp}.md"
    r.Save(outputDir, name, content)
    return content
