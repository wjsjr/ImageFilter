import matplotlib
matplotlib.use("agg")
import matplotlib.pyplot as plt
import sys
import subprocess
import os

def buildGraphs(speedups, sizes, modes, numThreads):
    for i in range(len(modes)):
        mode = modes[i]
        graph = plt.subplot()
        graph.set_xlabel("Num Threads")
        graph.set_ylabel("Speedup")
        graph.set_title(mode)
        plt.xticks(numThreads)
        for j in range(len(sizes)):
            graph.scatter(numThreads, speedups[i][j])
        
        path = f"graphs/{modes[i]}.png"
        if os.path.isfile(path):
            os.remove(path)
        plt.savefig(path, facecolor='w')
        
        plt.legend(sizes)
        plt.clf()



def runTrial(mode, threads, size):
    args = ["go", "run", "../editor/editor.go", size]
    if mode != "s":
        args.append(mode)
        args.append(str(threads))
    job = subprocess.Popen(args, stdout=subprocess.PIPE)
    job.wait()
    time = float(job.stdout.read().decode("utf-8")[:-1])
    return time

def calculateSpeedsup(results, benchmarks):
    for i in range(len(results)):
        resultsPerMode = results[i] 
        for j in range(len(resultsPerMode)):
            resultsPerSize = resultsPerMode[j]
            benchmarkPerSize = benchmarks[j]
            for k in range(len(resultsPerSize)):
                results[i][j][k] = benchmarkPerSize / resultsPerSize[k]
    
                

def runTests(modes, numThreads, sizes):
    benchmark = []
    for size in sizes:
        avgSerial = 0
        for i in range(5):
            avgSerial += runTrial("s", -1, size)
        avgSerial /= 5
        benchmark.append(avgSerial)
    
    resultsPerMode = []
    for mode in modes:
        resultsPerSize = []
        for size in sizes:
            resultsPerThreadCount = []
            for threadCount in numThreads:
                avgResults = 0.0
                for i in range(5):
                    avgResults += runTrial(mode, threadCount, size)
                avgResults /= 5
                resultsPerThreadCount.append(round(avgResults, 3))
            resultsPerSize.append(resultsPerThreadCount)
        resultsPerMode.append(resultsPerSize)
    
    calculateSpeedsup(resultsPerMode, benchmark)
    return resultsPerMode


def main(): 
    numThreads = [2, 4, 6, 8, 12]
    modes = ["pipeline", "bsp"]
    sizes = ["small", "mixture", "big"]
    if len(sys.argv) == 2:
        if sys.argv[1] == "small":
            sizes = [sizes[0]]
            numThreads = [numThreads[0]]
            modes = [modes[0]]
        else:
            print("Invalid arg")
            exit()
    speedups = runTests(modes, numThreads, sizes)
    buildGraphs(speedups, sizes, modes, numThreads)

if __name__ == '__main__':
    main()
