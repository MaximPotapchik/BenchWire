import pytest
from stats import StatsResult

def makeStats(label, values):
    return StatsResult({label: {"runs": len(values), "stats": [{"statName": "value", "values": values}]}})

@pytest.mark.parametrize("valueA,valueB,expectedPct", [
    (1.0093, 1.0085, 0.0793),
    (1.0000, 1.0000, 0.0),
])
def testDiff(valueA, valueB, expectedPct):
    a = makeStats("A", [valueA])
    result = a.Diff(valueA, valueB)
    assert round(result, 4) == expectedPct
