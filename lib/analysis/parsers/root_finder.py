import os

# This grabs the location of the projects root directory.
def FindRootDirectory(rootDirName):
    currentDir = os.path.dirname(os.path.abspath(__file__))

    while True:
        if os.path.basename(currentDir) == rootDirName:
            return currentDir

        parentDir = os.path.dirname(currentDir)
        if parentDir == currentDir:
            raise FileNotFoundError("BenchWire directory not found.")

        currentDir = parentDir

