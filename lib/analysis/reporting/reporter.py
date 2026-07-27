
class Reporter:
    def __init__(self, statsA, statsB=None):
        self.statsA, self.statsB = statsA, statsB

    def Render(self):
        raise NotImplementedError

    def Save(self, directory, name, content):
        with open(f"{directory}/{name}", "w") as f:
            f.write(content)

class DocumentReporter(Reporter):
    # StatsA also counts for single. Which is why B has a default value.
    def __init__(self, statsA, statsB=None):
        super().__init__(statsA, statsB)
        self.lines = []
        self.statIndex = {}
        self.plotFile = None

    def AddLine(self, content):
        self.lines.append(("text", content))
        return self

    def AddStat(self, label, valueA, valueB=None, diff=None):
        if label in self.statIndex:
            i = self.statIndex[label]
            _, _, oldA, oldB, oldDiff = self.lines[i]

            valueA = valueA if valueA is not None else oldA
            valueB = valueB if valueB is not None else oldB
            diff = diff if diff is not None else oldDiff
            self.lines[i] = ("stat", label, valueA, valueB, diff)

        else:
            self.statIndex[label] = len(self.lines)
            self.lines.append(("stat", label, valueA, valueB, diff))

        return self

    def AddTable(self, table):
        self.lines.append(("table", table))
        return self

    def SetPlotFile(self, path):
        self.plotFile = path
        return self

    def Render(self, title, mode, timestamp, **fields):
        output = self.RenderHeader(title, mode, timestamp, **fields) + self.RenderPlotRef()
        
        for line in self.lines:
            if line[0] == "stat":
                _, label, valueA, valueB, diff = line
                output += self.RenderStatLine(label, valueA, valueB, diff)
            elif line[0] == "table":
                _, table = line
                output += self.RenderTable(table)
            else:
                _, content = line
                output += self.RenderLines(content)

        return output

    def RenderHeader(self): raise NotImplementedError
    def RenderPlotRef(self): raise NotImplementedError
    def RenderStatLine(self, label, valueA, valueB, diff): raise NotImplementedError
    def RenderTable(self, table): raise NotImplementedError
    def RenderTableLine(self, label, valueA, valueB, diff, comparison, showDiff, unit): raise NotImplementedError
    def RenderLines(self, content): raise NotImplementedError

#TODO: This. For json, yaml, csv and more.
class DataReporter(Reporter):
    def Render(self):
        return self.Serialize(self.BuildData())

    def BuildData(self): raise NotImplementedError
    def Serialize(self, data): raise NotImplementedError
