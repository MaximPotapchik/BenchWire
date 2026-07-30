from .reporter import DocumentReporter

class MarkdownReporter(DocumentReporter):

    def RenderHeader(self, title, mode, timestamp, **fields):
        header = f"# {title} {mode} Batch #{timestamp}\n\n"
        for key, value in fields.items():
            if value:
                header += f"{key}: {value}\n"
        return header + "\n"

    def RenderPlotRef(self):
        return f'<img src="{self.plotFile}" width="100%">\n\n' if self.plotFile else ""

    def RenderStatLine(self, label, valueA, valueB, diff):
        if valueB is None:
            return f"{label}: {valueA:.4f}\n\n"

        higher = self.statsB.label if valueB > valueA else self.statsA.label
        line = f"{label} | A: {valueA:.4f} B: {valueB:.4f}\n"

        if diff is not None:
            line += f"{label} Difference: {diff:.4f}% ({higher} higher)\n"

        return line + "\n"
    
    def StartTable(self, comparison=False, showDiff=False):
        self.currentTable = {"comparison": comparison, "showDiff": showDiff, "columns": [], "rows": []}
        self.AddTable(self.currentTable)
        return self
 
    def SetTableColumns(self, *columns):
        self.currentTable["columns"] = columns
        return self
 
    def AddTableRow(self, label, valueA, valueB=None, diff=None, unit=""):
        self.currentTable["rows"].append((label, valueA, valueB, diff, unit))
        return self
    
    def EndTable(self):
        self.currentTable = None
        return self
 
    def RenderTable(self, table):
        columns = list(table["columns"])
        if table["showDiff"]:
            columns.append("Diff")
        header = "| " + " | ".join(columns) + " |\n"
        sep = "|" + "|".join(["---"] * len(columns)) + "|\n"
        body = ""

        for label, valueA, valueB, diff, unit in table["rows"]:
            higher = self.statsB.label if valueB is not None and valueB > valueA else self.statsA.label
            body += self.RenderTableLine(label, valueA, valueB, diff, table["comparison"],
                                        table["showDiff"], unit, higher)
        return header + sep + body + "\n"
 
    def RenderTableLine(self, label, valueA, valueB, diff, comparison, showDiff, unit, higher=None):
        if unit == "%":
            fmt = "{:.2f}"
        elif unit == "sci":
            fmt = "{:.4e}"
            unit = ""
        else:
            fmt = "{:.4f}"
 
        def cell(v):
            return fmt.format(v) + unit
 
        if not comparison:
            return f"| {label} | {cell(valueA)} |\n"
        cells = [label, cell(valueA), cell(valueB)]
        if showDiff and diff is not None:
            cells.append(f"{diff:.2f}% ({higher})")
        return "| " + " | ".join(cells) + " |\n"

    def RenderLines(self, content):
        return f"{content}\n\n"
