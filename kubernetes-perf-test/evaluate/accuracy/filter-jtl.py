import csv
import sys

# first timestamp when `head *-measurement.jtl`
fileName = sys.argv[1]

csvFile = csv.reader(open(fileName),delimiter=',')
_ = next(csvFile)
firstRow = next(csvFile)
headTime = int(firstRow[0])

duration = 60000 # 1 min
startTime = ((headTime // duration) + 1) * duration # next min
startTime = startTime - 5000 # reduce 5 seconds
sampleCount = 3

startTime2 = startTime + duration
startTime3 = startTime2 + duration
startTime4 = startTime3 + duration

for i in range(1, sampleCount+1):
    endTime = startTime + duration
    csvFile = csv.reader(open(fileName),delimiter=',')
    filtered = filter(lambda row: row[0] == 'timeStamp' or startTime <= int(row[0]) < endTime, csvFile)
    csv.writer(open("filtered_" + str(i) + "_" + fileName,'w'),delimiter=',').writerows(filtered)
    startTime = endTime
