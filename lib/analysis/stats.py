import numpy as np

class StatsResult:
    def __init__(self, data: dict):
        self.label, body = next(iter(data.items()))
        self.runs = body["runs"]
        self.index = {s["statName"]: s["values"] for s in body["stats"]}
    
    # TODO: Tighten up for construction
    @classmethod # Real constructor
    def FromAggregate(cls, aggregated: dict):
        results = tuple(cls({label: body}) for label, body in aggregated.items())
        return results[0] if len(results) == 1 else results

    def GetStat(self, stat):
        return self.index[stat]

    def GetAllStats(self):
        return self.index
    
    def GetStatNames(self):
        return list(self.index.keys())

    def GetRuns(self):
        return self.runs
    
    def Compare(self, other, stat):
        return self.GetStat(stat), other.GetStat(stat)

    # Add new statistical calculations here.
    def Mean(self, stat):
        values = self.GetStat(stat)
        return values.mean() if isinstance(values, np.ndarray) else None
    
    def Median(self, stat):
        values = self.GetStat(stat)
        return np.median(values) if isinstance(values, np.ndarray) else None
    
    def StandardDeviation(self, stat):
        values = self.GetStat(stat)
        return values.std() if isinstance(values, np.ndarray) else None

    def Variance(self, stat):
        values = self.GetStat(stat)
        return values.var() if isinstance(values, np.ndarray) else None
   
    def CoefficientofV(self, stat):
        mean = self.Mean(stat)
        return self.StandardDeviation(stat) / mean * 100 if mean else None

    def Range(self, stat):
        minV, maxV = self.Min(stat), self.Max(stat)
        return maxV - minV if minV is not None and maxV is not None else None

    def InterquartileRange(self, stat):
        p75, p25 = self.P(stat, 75), self.P(stat, 25)
        return p75 - p25 if p75 is not None and p25 is not None else None

    # Percentiles/Ordering
    def P(self, stat, percent):
        values = self.GetStat(stat)
        return np.percentile(values, percent) if isinstance(values, np.ndarray) else None

    def Min(self, stat):
        values = self.GetStat(stat)
        return np.min(values) if isinstance(values, np.ndarray) else None

    def Max(self, stat):
        values = self.GetStat(stat)
        return np.max(values) if isinstance(values, np.ndarray) else None
